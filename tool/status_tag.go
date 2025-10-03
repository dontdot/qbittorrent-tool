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
		log.Printf("[DEBUG] StatusTag skipped for %s: enable=%v, map_config_nil=%v\n", torrent.Name, c.StatusTag.Enable, c.StatusTag.MapConfig == nil)
		return
	}

	trackerList, err := qbittorrent.Api.GetTorrentTrackers(torrent.Hash)
	if err != nil || len(trackerList) == 0 {
		fmt.Printf("[ERR] get %s tracker list err: %v, count: %d\n", torrent.Name, err, len(trackerList))
		return
	}
	log.Printf("[DEBUG] Got %d trackers for %s\n", len(trackerList), torrent.Name)

	tag := ""
	miss := make(map[string]int, 0)
	for i, tracker := range trackerList {
		log.Printf("[DEBUG] Checking tracker #%d for %s: status=%d, msg=\"%s\"\n", i, torrent.Name, tracker.Status, tracker.Msg)
		if tracker.Status == 2 || tracker.Msg == "" {
			log.Printf("[DEBUG] Skipping tracker #%d for %s: status=%d or msg_empty=%v\n", i, torrent.Name, tracker.Status, tracker.Msg == "")
			return
		}

		if custom, ok := c.StatusTag.MapConfig[tracker.Msg]; ok {
			tag = custom
			log.Printf("[DEBUG] Found custom tag mapping for %s: \"%s\" -> %s\n", torrent.Name, tracker.Msg, custom)
		} else {
			miss[tracker.Msg] += 1
			log.Printf("[DEBUG] No mapping found for tracker message \"%s\" in %s\n", tracker.Msg, torrent.Name)
		}
	}

	if len(miss) > 0 {
		for item, _ := range miss {
			fmt.Printf("err: \"%s: %s\" not map config\n", torrent.Name, item)
		}
	}

	if tag == "" || strings.Contains(torrent.Tags, tag) {
		log.Printf("[DEBUG] No tag to add for %s: tag=\"%s\", already_exists=%v\n", torrent.Name, tag, strings.Contains(torrent.Tags, tag))
		return
	}

	err = qbittorrent.Api.AddTags(torrent.Hash, tag)
	if err != nil {
		fmt.Printf("[ERR] add tag %s to %s err: %v\n", tag, torrent.Name, err)
		return
	}

	fmt.Printf("[INFO] add tag %s to %s\n", tag, torrent.Name)
}