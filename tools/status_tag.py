#!/usr/bin/env python3
# -*- coding: utf-8 -*-

from config.config import get_logger
from qbittorrent.api import QBittorrentAPI, Torrent, Tracker
from config.config import StatusTagConfig

logger = get_logger()

class StatusTag:
    def __init__(self, qb_api: QBittorrentAPI, config: StatusTagConfig):
        self.qb_api = qb_api
        self.config = config
        
    def process(self, torrent: Torrent):
        """根据tracker状态设置标签"""
        if not self.config.enable or not self.config.map_config:
            logger.debug(f"跳过状态标签处理 {torrent.name}: 启用={self.config.enable}, 映射配置为空={not bool(self.config.map_config)}")
            return
            
        try:
            trackers = self.qb_api.get_torrent_trackers(torrent.hash)
        except Exception as e:
            logger.error(f"获取种子 {torrent.name} 的 trackers 错误: {e}")
            return
            
        for tracker in trackers:
            # 检查tracker状态消息是否匹配配置的映射
            for status_msg, tag in self.config.map_config.items():
                if status_msg in tracker.msg:
                    # 如果配置中的标签为空字符串，则移除对应标签
                    if tag == "":
                        # 移除标签的实现
                        logger.debug(f"应该移除种子 {torrent.name} 的标签，基于状态: {status_msg}")
                    else:
                        # 检查标签是否已存在
                        existing_tags = [t.strip() for t in torrent.tags.split(',')] if torrent.tags else []
                        if tag in existing_tags:
                            logger.debug(f"状态标签 {tag} 已存在于种子 {torrent.name} 中，跳过")
                            continue
                            
                        # 添加标签
                        try:
                            self.qb_api.add_tags(torrent.hash, tag)
                            logger.debug(f"成功添加状态标签 {tag} 到种子 {torrent.name}")
                        except Exception as e:
                            logger.error(f"添加状态标签 {tag} 到种子 {torrent.name} 错误: {e}")
                    break  # 找到匹配就退出循环