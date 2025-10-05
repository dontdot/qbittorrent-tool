#!/usr/bin/env python3
# -*- coding: utf-8 -*-

from config.config import get_logger
from qbittorrent.api import QBittorrentAPI, Torrent
from config.config import AutoCategoryConfig

logger = get_logger()

class AutoCategory:
    def __init__(self, qb_api: QBittorrentAPI, config: AutoCategoryConfig):
        self.qb_api = qb_api
        self.config = config
        
    def process(self, torrent: Torrent):
        """根据保存路径自动设置分类"""
        if not self.config.enable or not self.config.map_config:
            logger.debug(f"跳过自动分类处理 {torrent.name}: 启用={self.config.enable}, 映射配置为空={not bool(self.config.map_config)}")
            return
            
        # 使用 torrent 对象自带的 save_path 属性
        torrent_path = torrent.save_path
        if not torrent_path:
            logger.debug(f"无法获取种子 {torrent.name} 的保存路径")
            return

        # 查找匹配的分类
        category = None
        for path_prefix, cat in self.config.map_config.items():
            if torrent_path.startswith(path_prefix):
                category = cat
                break
                
        if not category:
            logger.debug(f"种子 {torrent.name} 的路径 {torrent_path} 未匹配到任何分类规则")
            return
            
        # 直接访问 category 属性（qBittorrent-API 支持）
        if torrent.category == category:
            logger.debug(f"种子 {torrent.name} 已经设置为分类 {category}，跳过")
            return
            
        # 设置分类
        try:
            self.qb_api.set_category(torrent.hash, category)
            logger.debug(f"成功为种子 {torrent.name} 设置分类 {category}")
        except Exception as e:
            logger.error(f"为种子 {torrent.name} 设置分类 {category} 错误: {e}")