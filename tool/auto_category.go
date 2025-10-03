package tool

import (
	"fengqi/qbittorrent-tool/config"
	"fengqi/qbittorrent-tool/qbittorrent"
	"fmt"
	"log"
)

// AutoCategory 根据保存目录设置分类
func AutoCategory(c *config.Config, torrent *qbittorrent.Torrent) {
	if torrent.Category != "" || !c.AutoCategory.Enable || c.AutoCategory.MapConfig == nil {
		log.Printf("[DEBUG] 跳过自动分类 %s: 已有分类=%v, 启用=%v, 映射配置为空=%v\n", torrent.Name, torrent.Category != "", c.AutoCategory.Enable, c.AutoCategory.MapConfig == nil)
		return
	}

	category, ok := c.AutoCategory.MapConfig[torrent.SavePath]
	if !ok {
		log.Printf("[WARN] 种子保存路径 %s 未找到对应分类 %s\n", torrent.SavePath, torrent.Name)
		return
	}
	log.Printf("[DEBUG] 找到分类映射 %s: %s -> %s\n", torrent.Name, torrent.SavePath, category)

	err := qbittorrent.Api.SetCategory(torrent.Hash, category)
	if err != nil {
		log.Printf("[ERR] 设置分类错误 %s 到种子: %s 错误: %v\n", category, torrent.Name, err)
		return
	}

	fmt.Printf("[INFO] 设置分类: %s 到种子: %s\n", category, torrent.Name)
}