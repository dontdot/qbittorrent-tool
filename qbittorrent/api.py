#!/usr/bin/env python3
# -*- coding: utf-8 -*-

import requests
import logging
from typing import Dict, List, Any, Optional
from urllib.parse import urljoin
from config.config import get_logger

logger = get_logger()

class Torrent:
    def __init__(self, data: Dict[str, Any]):
        self.hash = data.get('hash', '')
        self.name = data.get('name', '')
        self.tracker = data.get('tracker', '')
        self.tags = data.get('tags', '')
        self.category = data.get('category', '')
        self.seeds = data.get('nb_seeders', 0)
        self.ratio = data.get('ratio', 0.0)
        self.seeding_time = data.get('seeding_time', 0)  # 秒
        self.activity_time = data.get('last_activity', 0)  # 时间戳
        self.save_path = data.get('save_path', '')
        # 可以根据需要添加更多字段
        
    def get_tracker_host(self) -> str:
        """从 tracker URL 提取主机名"""
        if self.tracker:
            # 简单实现，实际可能需要更复杂的 URL 解析
            if '://' in self.tracker:
                host = self.tracker.split('://')[1].split('/')[0]
                if ':' in host:
                    return host.split(':')[0]
                return host
        return ''

class Tracker:
    def __init__(self, data: Dict[str, Any]):
        self.url = data.get('url', '')
        self.status = data.get('status', 0)
        self.msg = data.get('msg', '')
        # 可以根据需要添加更多字段

class QBittorrentAPI:
    def __init__(self, host: str, username: str, password: str):
        self.host = host.rstrip('/')
        self.username = username
        self.password = password
        self.session = requests.Session()
        
    def login(self) -> bool:
        """登录到 qBittorrent Web UI"""
        url = urljoin(self.host, '/api/v2/auth/login')
        data = {
            'username': self.username,
            'password': self.password
        }
        
        try:
            response = self.session.post(url, data=data)
            response.raise_for_status()
            return True
        except Exception as e:
            logger.error(f"登录 qBittorrent 错误: {e}")
            raise
            
    def get_torrent_list(self, params: Dict[str, str]) -> List[Torrent]:
        """获取种子列表"""
        url = urljoin(self.host, '/api/v2/torrents/info')
        
        try:
            response = self.session.get(url, params=params)
            response.raise_for_status()
            torrents_data = response.json()
            return [Torrent(torrent_data) for torrent_data in torrents_data]
        except Exception as e:
            logger.error(f"获取种子列表错误: {e}")
            raise
            
    def get_torrent_trackers(self, torrent_hash: str) -> List[Tracker]:
        """获取种子的 trackers"""
        url = urljoin(self.host, '/api/v2/torrents/trackers')
        params = {'hash': torrent_hash}
        
        try:
            response = self.session.get(url, params=params)
            response.raise_for_status()
            trackers_data = response.json()
            return [Tracker(tracker_data) for tracker_data in trackers_data]
        except Exception as e:
            logger.error(f"获取种子 trackers 错误: {e}")
            raise
            
    def add_tags(self, torrent_hash: str, tags: str) -> bool:
        """给种子添加标签"""
        url = urljoin(self.host, '/api/v2/torrents/addTags')
        data = {
            'hashes': torrent_hash,
            'tags': tags
        }
        
        try:
            response = self.session.post(url, data=data)
            response.raise_for_status()
            return True
        except Exception as e:
            logger.error(f"添加标签错误: {e}")
            raise
            
    def set_category(self, torrent_hash: str, category: str) -> bool:
        """设置种子分类"""
        url = urljoin(self.host, '/api/v2/torrents/setCategory')
        data = {
            'hashes': torrent_hash,
            'category': category
        }
        
        try:
            response = self.session.post(url, data=data)
            response.raise_for_status()
            return True
        except Exception as e:
            logger.error(f"设置分类错误: {e}")
            raise
            
    def set_torrent_limits(self, torrent_hash: str, limits: Dict[str, Any]) -> bool:
        """设置种子限制"""
        url = urljoin(self.host, '/api/v2/torrents/setShareLimits')
        data = {
            'hashes': torrent_hash,
            **limits
        }
        
        try:
            response = self.session.post(url, data=data)
            response.raise_for_status()
            return True
        except Exception as e:
            logger.error(f"设置种子限制错误: {e}")
            raise