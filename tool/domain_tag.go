package tool

import (
	"fengqi/qbittorrent-tool/config"
	"fengqi/qbittorrent-tool/qbittorrent"
	"fmt"
	"log"
	"strings"
)

// DomainTag 根据域名设置tag, 主要是给webui用
// 等webui可以和桌面端一样自动合并tracker的时候, 可以放弃使用
func DomainTag(c *config.Config, torrent *qbittorrent.Torrent) {
	if !c.DomainTag.Enable || c.DomainTag.MapConfig == nil {
		log.Printf("[DEBUG] 跳过域名标签处理 %s: 启用=%v, 映射配置为空=%v\n", torrent.Name, c.DomainTag.Enable, c.DomainTag.MapConfig == nil)
		return
	}

	if torrent.Tracker == "" {
		trackerList, err := qbittorrent.Api.GetTorrentTrackers(torrent.Hash)
		if err == nil && len(trackerList) > 0 {
			torrent.Tracker = trackerList[0].Url // 暂时默认用第一个
			log.Printf("[DEBUG] 获取到种子tracker %s: %s\n", torrent.Name, torrent.Tracker)
		} else {
			log.Printf("[DEBUG] 无法获取种子tracker %s: 错误=%v, tracker列表长度=%d\n", torrent.Name, err, len(trackerList))
		}
	}

	tag, err := torrent.GetTrackerHost()
	if err != nil {
		fmt.Printf("[ERR] 获取种子 %s 标签错误: %v\n", torrent.Name, err)
		return
	}
	log.Printf("[DEBUG] 提取标签 %s: %s\n", torrent.Name, tag)

	if custom, ok := c.DomainTag.MapConfig[tag]; ok {
		tag = custom
		log.Printf("[DEBUG] 应用自定义标签映射 %s: %s -> %s\n", torrent.Name, tag, custom)
	}

	if strings.Contains(torrent.Tags, tag) {
		log.Printf("[DEBUG] 标签 %s 已存在于种子 %s 中，跳过\n", tag, torrent.Name)
		return
	}

	err = qbittorrent.Api.AddTags(torrent.Hash, tag)
	if err != nil {
		fmt.Printf("[ERR] 添加标签 %s 到种子 %s 错误: %v\n", tag, torrent.Name, err)
		return
	}

	fmt.Printf("[INFO] 添加标签 %s 到种子 %s\n", tag, torrent.Name)
}