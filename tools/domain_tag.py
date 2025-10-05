#!/usr/bin/env python3
# -*- coding: utf-8 -*-

import logging
from typing import List
from qbittorrent.api import QBittorrentAPI, Torrent
from config.config import DomainTagConfig

logger = logging.getLogger(__name__)

class DomainTag:
    def __init__(self, qb_api: QBittorrentAPI, config: DomainTagConfig):
        self.qb_api = qb_api
        self.config = config
        
    def process(self, torrent: Torrent):
        """根据域名设置标签"""
        if not self.config.enable or not self.config.map_config:
            logger.debug(f"跳过域名标签处理 {torrent.name}: 启用={self.config.enable}, 映射配置为空={not bool(self.config.map_config)}")
            return
            
        if not torrent.tracker:
            try:
                tracker_list = self.qb_api.get_torrent_trackers(torrent.hash)
                if tracker_list:
                    torrent.tracker = tracker_list[0].url  # 暂时默认用第一个
                    logger.debug(f"获取到种子tracker {torrent.name}: {torrent.tracker}")
                else:
                    logger.debug(f"无法获取种子tracker {torrent.name}: tracker列表为空")
            except Exception as e:
                logger.debug(f"无法获取种子tracker {torrent.name}: 错误={e}")
                
        try:
            tracker_host = torrent.get_tracker_host()
        except Exception as e:
            logger.error(f"获取种子 {torrent.name} 标签错误: {e}")
            return
            
        logger.debug(f"提取标签 {torrent.name}: {tracker_host}")
        
        # 应用自定义标签映射
        if tracker_host in self.config.map_config:
            custom_tag = self.config.map_config[tracker_host]
            tag = custom_tag
            logger.debug(f"应用自定义标签映射 {torrent.name}: {tracker_host} -> {custom_tag}")
        else:
            tag = tracker_host
            
        # 检查标签是否已存在
        if tag in torrent.tags:
            logger.debug(f"标签 {tag} 已存在于种子 {torrent.name} 中，跳过")
            return
            
        # 添加标签
        try:
            self.qb_api.add_tags(torrent.hash, tag)
            logger.debug(f"成功添加标签 {tag} 到种子 {torrent.name}")
        except Exception as e:
            logger.error(f"添加标签 {tag} 到种子 {torrent.name} 错误: {e}")