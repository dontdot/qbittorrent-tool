#!/usr/bin/env python3
# -*- coding: utf-8 -*-

from typing import Dict, List, Optional, Any

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