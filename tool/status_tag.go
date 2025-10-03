package tool

import (
	"fengqi/qbittorrent-tool/config"
	"fengqi/qbittorrent-tool/qbittorrent"
	"fmt"
	"log"
	"strings"
)

func StatusTag(c *config.Config, torrent *qbittorrent.Torrent) {
	if !c.StatusTag.Enable || c.StatusTag.MapConfig == nil {
		log.Printf("[DEBUG] 跳过状态标签 %s: 启用=%v, 映射配置为空=%v\n", torrent.Name, c.StatusTag.Enable, c.StatusTag.MapConfig == nil)
		return
	}

	trackerList, err := qbittorrent.Api.GetTorrentTrackers(torrent.Hash)
	if err != nil || len(trackerList) == 0 {
		fmt.Printf("[ERR] 获取种子 %s tracker列表错误: %v, 数量: %d\n", torrent.Name, err, len(trackerList))
		return
	}
	log.Printf("[DEBUG] 获取到 %d 个trackers %s\n", len(trackerList), torrent.Name)

	tag := ""
	miss := make(map[string]int, 0)
	for i, tracker := range trackerList {
		log.Printf("[DEBUG] 检查tracker #%d %s: 状态=%d, 消息=\"%s\"\n", i, torrent.Name, tracker.Status, tracker.Msg)
		if tracker.Status == 2 || tracker.Msg == "" {
			log.Printf("[DEBUG] 跳过tracker #%d %s: 状态=%d 或 消息为空=%v\n", i, torrent.Name, tracker.Status, tracker.Msg == "")
			return
		}

		if custom, ok := c.StatusTag.MapConfig[tracker.Msg]; ok {
			tag = custom
			log.Printf("[DEBUG] 找到自定义标签映射 %s: \"%s\" -> %s\n", torrent.Name, tracker.Msg, custom)
		} else {
			miss[tracker.Msg] += 1
			log.Printf("[DEBUG] 未找到tracker消息映射 \"%s\" %s\n", tracker.Msg, torrent.Name)
		}
	}

	if len(miss) > 0 {
		for item, _ := range miss {
			fmt.Printf("错误: \"%s: %s\" 未配置映射\n", torrent.Name, item)
		}
	}

	if tag == "" || strings.Contains(torrent.Tags, tag) {
		log.Printf("[DEBUG] 无需添加标签 %s: 标签=\"%s\", 已存在=%v\n", torrent.Name, tag, strings.Contains(torrent.Tags, tag))
		return
	}

	err = qbittorrent.Api.AddTags(torrent.Hash, tag)
	if err != nil {
		fmt.Printf("[ERR] 添加标签 %s 到种子 %s 错误: %v\n", tag, torrent.Name, err)
		return
	}

	fmt.Printf("[INFO] 添加标签 %s 到种子 %s\n", tag, torrent.Name)
}