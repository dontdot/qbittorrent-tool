#!/usr/bin/env python3
# -*- coding: utf-8 -*-

import logging
from qbittorrent.api import QBittorrentAPI, Torrent
from config.config import AutoCategoryConfig

logger = logging.getLogger(__name__)

class AutoCategory:
    def __init__(self, qb_api: QBittorrentAPI, config: AutoCategoryConfig):
        self.qb_api = qb_api
        self.config = config
        
    def process(self, torrent: Torrent):
        """根据保存路径自动设置分类"""
        if not self.config.enable or not self.config.map_config:
            logger.debug(f"跳过自动分类处理 {torrent.name}: 启用={self.config.enable}, 映射配置为空={not bool(self.config.map_config)}")
            return
            
        # 注意：由于API限制，我们可能无法直接获取种子的保存路径
        # 在这个简化版本中，我们假设可以通过某种方式获取路径信息
        # 实际实现中你可能需要根据 qBittorrent API 文档进行调整
        torrent_path = self._get_torrent_save_path(torrent)
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
            
        # 检查是否已经设置了正确的分类
        # 注意：这里假设 torrent 对象有 category 属性，实际情况可能需要根据 API 调整
        if hasattr(torrent, 'category') and torrent.category == category:
            logger.debug(f"种子 {torrent.name} 已经设置为分类 {category}，跳过")
            return
            
        # 设置分类
        try:
            self.qb_api.set_category(torrent.hash, category)
            logger.debug(f"成功为种子 {torrent.name} 设置分类 {category}")
        except Exception as e:
            logger.error(f"为种子 {torrent.name} 设置分类 {category} 错误: {e}")
            
    def _get_torrent_save_path(self, torrent: Torrent) -> str:
        """
        获取种子保存路径
        注意：这只是一个占位实现，实际需要根据 qBittorrent API 提供的方法来获取
        """
        # 在真实的实现中，你可能需要通过额外的 API 调用来获取种子的保存路径
        # 这里暂时返回空字符串表示无法获取
        return ""