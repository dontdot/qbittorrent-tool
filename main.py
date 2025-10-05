#!/usr/bin/env python3
# -*- coding: utf-8 -*-

import argparse
import json
import logging
import sys
import time
import os
from typing import Dict, List, Any

from config.config import Config, setup_logger
from qbittorrent.api import QBittorrentAPI
from tools.auto_category import AutoCategory
from tools.domain_tag import DomainTag
from tools.seeding_limits import SeedingLimits
from tools.status_tag import StatusTag

def main():
    parser = argparse.ArgumentParser(description='qBittorrent Tool')
    parser.add_argument('-c', '--config', default='./config.json', help='配置文件路径')
    parser.add_argument('-log', '--log-file', default='./qbittorrent-tool.log', help='日志文件路径')
    
    args = parser.parse_args()
    
    # 设置全局日志记录器
    logger = setup_logger(args.log_file)
    
    try:
        # 加载配置
        config = load_config(args.config)
        
        # 初始化 qBittorrent API
        qb_api = QBittorrentAPI(config.host, config.username, config.password)
        qb_api.login()
        
        logger.info(f"程序开始执行 {time.strftime('%Y-%m-%d %H:%M:%S')}")
        
        # 初始化工具
        auto_category = AutoCategory(qb_api, config.auto_category)
        domain_tag = DomainTag(qb_api, config.domain_tag)
        seeding_limits = SeedingLimits(qb_api, config.seeding_limits)
        status_tag = StatusTag(qb_api, config.status_tag)
        
        # 分批处理种子
        offset = 0
        limit = 1000
        total_processed = 0
        
        while True:
            params = {
                'filter': 'all',
                'sort': 'added_on',
                'limit': limit,
                'offset': offset
            }
            
            torrent_list = qb_api.get_torrent_list(params)
            logger.info(f"获取种子列表数量: {len(torrent_list)}")
            
            if len(torrent_list) == 0:
                logger.info("没有更多种子需要处理，退出循环")
                break
            
            for i, torrent in enumerate(torrent_list):
                logger.debug(f"处理第 {i+1} 个种子: {torrent.name}")
                auto_category.process(torrent)
                domain_tag.process(torrent)
                seeding_limits.process(torrent)
                status_tag.process(torrent)
                total_processed += 1
            
            if len(torrent_list) < limit:
                logger.info(f"处理完最后一批种子，总共处理: {total_processed}")
                break
            
            offset += limit
            logger.info(f"移动到下一批，偏移量: {offset}")
        
        logger.info(f"完成所有种子处理，总计: {total_processed}")
        logger.info(f"程序执行结束 {time.strftime('%Y-%m-%d %H:%M:%S')}")
        
    except Exception as e:
        logger.error(f"程序执行出错: {e}")
        sys.exit(1)

def load_config(config_file: str) -> Config:
    """加载配置文件"""
    try:
        with open(config_file, 'r', encoding='utf-8') as f:
            config_data = json.load(f)
        return Config(config_data)
    except Exception as e:
        # 在logger初始化前，使用print输出错误
        print(f"加载配置错误: {e}")
        raise

if __name__ == "__main__":
    main()