#!/usr/bin/env python3
# -*- coding: utf-8 -*-

import logging
from qbittorrent.api import QBittorrentAPI, Torrent
from config.config import SeedingLimitsConfig, SeedingLimitsRule

logger = logging.getLogger(__name__)

class SeedingLimits:
    def __init__(self, qb_api: QBittorrentAPI, config: SeedingLimitsConfig):
        self.qb_api = qb_api
        self.config = config
        
    def process(self, torrent: Torrent):
        """根据规则设置做种限制"""
        if not self.config.enable or not self.config.rules:
            logger.debug(f"跳过做种限制处理 {torrent.name}: 启用={self.config.enable}, 规则数量={len(self.config.rules) if self.config.rules else 0}")
            return
            
        matched_rule = None
        # 按顺序检查规则，找到最后一个匹配的规则
        for rule in self.config.rules:
            if self._match_rule(rule, torrent):
                matched_rule = rule
                logger.debug(f"种子 {torrent.name} 匹配规则: {rule}")
                
        if not matched_rule:
            logger.debug(f"种子 {torrent.name} 未匹配任何规则")
            return
            
        # 根据规则的动作类型执行相应的操作
        if matched_rule.action == 0:  # 继续做种/应用限制
            self._apply_limits(torrent, matched_rule)
        elif matched_rule.action == 1:  # 暂停做种
            # 暂停种子的实现
            pass
        elif matched_rule.action in [2, 3]:  # 删除种子/删除种子及文件
            # 删除种子的实现
            pass
        elif matched_rule.action == 4:  # 启动超级做种
            # 启动超级做种的实现
            pass
            
    def _match_rule(self, rule: SeedingLimitsRule, torrent: Torrent) -> bool:
        """检查种子是否匹配规则"""
        # 检查标签匹配
        if rule.tag:
            tag_matched = any(tag in torrent.tags for tag in rule.tag)
            if not tag_matched:
                return False
                
        # 检查分类匹配
        if rule.category and hasattr(torrent, 'category'):
            category_matched = torrent.category in rule.category
            if not category_matched:
                return False
                
        # 检查关键词匹配
        if rule.keyword:
            keyword_matched = any(keyword in torrent.name for keyword in rule.keyword)
            if not keyword_matched:
                return False
                
        # 其他匹配条件可以根据需要添加
        # 这里为了简化只实现了标签、分类和关键词匹配
        
        return True
        
    def _apply_limits(self, torrent: Torrent, rule: SeedingLimitsRule):
        """应用做种限制"""
        if not rule.limits:
            logger.debug(f"规则没有设置限制参数，跳过应用限制 {torrent.name}")
            return
            
        limits_data = {}
        
        # 设置下载限速
        if rule.limits.download > 0:
            limits_data['download_limit'] = rule.limits.download
            
        # 设置上传限速
        if rule.limits.upload > 0:
            limits_data['upload_limit'] = rule.limits.upload
            
        # 如果有待设置的限制参数，则应用它们
        if limits_data:
            try:
                # 这里应该调用适当的API方法来应用限制
                # 由于API细节可能有所不同，这里只是示意
                logger.debug(f"为种子 {torrent.name} 应用限制: {limits_data}")
                # self.qb_api.set_torrent_limits(torrent.hash, limits_data)
            except Exception as e:
                logger.error(f"为种子 {torrent.name} 应用限制错误: {e}")
        else:
            logger.debug(f"没有需要应用的限制参数 {torrent.name}")