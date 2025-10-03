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
		log.Printf("[DEBUG] AutoCategory skipped for %s: hasCategory=%v, enable=%v, map_config_nil=%v\n", torrent.Name, torrent.Category != "", c.AutoCategory.Enable, c.AutoCategory.MapConfig == nil)
		return
	}

	category, ok := c.AutoCategory.MapConfig[torrent.SavePath]
	if !ok {
		log.Printf("[WARN] get path %s category empty for torrent %s\n", torrent.SavePath, torrent.Name)
		return
	}
	log.Printf("[DEBUG] Found category mapping for %s: %s -> %s\n", torrent.Name, torrent.SavePath, category)

	err := qbittorrent.Api.SetCategory(torrent.Hash, category)
	if err != nil {
		log.Printf("[ERR] set category: %s \tto: %s err: %v\n", category, torrent.Name, err)
		return
	}

	fmt.Printf("[INFO] set category: %s \tto: %s\n", category, torrent.Name)
}