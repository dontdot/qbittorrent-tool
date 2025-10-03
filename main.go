package main

import (
	"fengqi/qbittorrent-tool/config"
	"fengqi/qbittorrent-tool/qbittorrent"
	"fengqi/qbittorrent-tool/tool"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
)

func main() {
	configFile := flag.String("c", "./config.json", "config file path")
	debug := flag.Bool("debug", false, "enable debug logging")
	flag.Parse()

	// 如果启用了调试模式，则设置日志显示文件名和行号
	if *debug {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
	} else {
		log.SetOutput(os.Stdout)
	}

	c, err := config.LoadConfig(*configFile)
	if err != nil {
		fmt.Printf("[ERR] load config err: %v\n", err)
		return
	}

	if err = qbittorrent.Init(c); err != nil {
		fmt.Printf("[ERR] login to qbittorrent err %v\n", err)
		return
	}

	offset := 0
	limit := 1000
	totalProcessed := 0
	for {
		params := map[string]string{
			"filter": "all",
			"sort":   "added_on",
			"limit":  strconv.Itoa(limit),
			"offset": strconv.Itoa(offset),
		}
		torrentList, err := qbittorrent.Api.GetTorrentList(params)
		if err != nil {
			log.Printf("[ERR] get torrent list err %v\n", err)
			return
		}

		log.Printf("[INFO] get torrent list count: %d\n", len(torrentList))
		if len(torrentList) == 0 {
			log.Printf("[INFO] No more torrents to process, exiting loop\n")
			break
		}
		
		for i, torrent := range torrentList {
			log.Printf("[DEBUG] Processing torrent #%d of batch: %s\n", i+1, torrent.Name)
			tool.AutoCategory(c, torrent)
			tool.DomainTag(c, torrent)
			tool.SeedingLimits(c, torrent)
			tool.StatusTag(c, torrent)
			totalProcessed++
		}

		if len(torrentList) < limit {
			log.Printf("[INFO] Processed last batch of torrents, total processed: %d\n", totalProcessed)
			break
		}
		offset += limit
		log.Printf("[INFO] Moving to next batch, offset: %d\n", offset)
	}
	
	log.Printf("[INFO] Finished processing all torrents, total: %d\n", totalProcessed)
}