#!/usr/bin/env python3
# -*- coding: utf-8 -*-

import logging
import sys
from typing import List, Dict, Any

# 全局logger实例
logger = None

def setup_logger(log_file: str = './qbittorrent-tool.log') -> logging.Logger:
    """设置并返回全局logger实例"""
    global logger
    if logger is not None:
        return logger
    
    # 创建logger
    logger = logging.getLogger('qbittorrent_tool')
    logger.setLevel(logging.DEBUG)
    
    # 避免重复添加处理器
    if not logger.handlers:
        # 创建日志格式
        formatter = logging.Formatter('%(asctime)s [%(levelname)s] %(filename)s:%(lineno)d - %(message)s')
        
        # 创建控制台处理器
        console_handler = logging.StreamHandler(sys.stdout)
        console_handler.setLevel(logging.DEBUG)
        console_handler.setFormatter(formatter)
        logger.addHandler(console_handler)
        
        # 创建文件处理器
        try:
            file_handler = logging.FileHandler(log_file, encoding='utf-8')
            file_handler.setLevel(logging.DEBUG)
            file_handler.setFormatter(formatter)
            logger.addHandler(file_handler)
        except Exception as e:
            # 如果无法创建文件处理器，至少保持控制台输出
            print(f"警告: 无法创建日志文件 {log_file}: {e}")
    
    return logger

def get_logger() -> logging.Logger:
    """获取全局logger实例"""
    global logger
    if logger is None:
        logger = setup_logger()
    return logger

class AutoCategoryConfig:
    def __init__(self, data: Dict[str, Any]):
        self.enable = data.get('enable', False)
        self.map_config = data.get('map_config', {})

class DomainTagConfig:
    def __init__(self, data: Dict[str, Any]):
        self.enable = data.get('enable', False)
        self.map_config = data.get('map_config', {})

class Limits:
    def __init__(self, data: Dict[str, Any]):
        self.download = data.get('download', 0)
        self.upload = data.get('upload', 0)
        self.ratio = data.get('ratio', 0.0)
        self.seeding_time = data.get('seeding_time', 0)
        self.inactive_seeding_time = data.get('inactive_seeding_time', 0)

class SeedingLimitsRule:
    def __init__(self, data: Dict[str, Any]):
        self.ratio = data.get('ratio', 0.0)
        self.seeding_time = data.get('seeding_time', 0)
        self.activity_time = data.get('activity_time', 0)
        self.tag = data.get('tag', [])
        self.category = data.get('category', [])
        self.tracker = data.get('tracker', [])
        self.seeds_gt = data.get('seeds_gt', 0)
        self.seeds_lt = data.get('seeds_lt', 0)
        self.keyword = data.get('keyword', [])
        self.action = data.get('action', 0)
        # 处理 limits 字段，如果存在则创建 Limits 对象
        self.limits = Limits(data['limits']) if 'limits' in data and data['limits'] else None

class SeedingLimitsConfig:
    def __init__(self, data: Dict[str, Any]):
        self.enable = data.get('enable', False)
        self.resume = data.get('resume', False)
        # 处理规则列表
        self.rules = [SeedingLimitsRule(rule) for rule in data.get('rules', [])]

class StatusTagConfig:
    def __init__(self, data: Dict[str, Any]):
        self.enable = data.get('enable', False)
        self.map_config = data.get('map_config', {})

class Config:
    def __init__(self, data: Dict[str, Any]):
        self.host = data.get('host', '')
        self.username = data.get('username', '')
        self.password = data.get('password', '')
        
        # 处理各模块配置
        self.auto_category = AutoCategoryConfig(data.get('auto_category', {}))
        self.domain_tag = DomainTagConfig(data.get('domain_tag', {}))
        self.seeding_limits = SeedingLimitsConfig(data.get('seeding_limits', {}))
        self.status_tag = StatusTagConfig(data.get('status_tag', {}))