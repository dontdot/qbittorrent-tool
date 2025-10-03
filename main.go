package main

import (
	"fengqi/qbittorrent-tool/config"
	"fengqi/qbittorrent-tool/qbittorrent"
	"fengqi/qbittorrent-tool/tool"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"time"
)

func main() {
	configFile := flag.String("c", "./config.json", "配置文件路径")
	logFile := flag.String("log", "./qbittorrent-tool.log", "日志文件路径")
	flag.Parse()

	// 创建日志文件
	file, err := os.OpenFile(*logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("[ERR] 无法创建日志文件: %v\n", err)
		return
	}
	defer file.Close()

	// 设置日志输出格式
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	
	// 创建一个多写入器，同时写入文件和控制台
	multiWriter := &MultiWriter{[]io.Writer{file, os.Stdout}}
	log.SetOutput(multiWriter)

	c, err := config.LoadConfig(*configFile)
	if err != nil {
		fmt.Printf("[ERR] 加载配置错误: %v\n", err)
		return
	}

	if err = qbittorrent.Init(c); err != nil {
		fmt.Printf("[ERR] 登录qbittorrent错误 %v\n", err)
		return
	}

	log.Printf("[INFO] 程序开始执行 %s\n", time.Now().Format("2006-01-02 15:04:05"))

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
			log.Printf("[ERR] 获取种子列表错误: %v\n", err)
			return
		}

		log.Printf("[INFO] 获取种子列表数量: %d\n", len(torrentList))
		if len(torrentList) == 0 {
			log.Printf("[INFO] 没有更多种子需要处理，退出循环\n")
			break
		}
		
		for i, torrent := range torrentList {
			log.Printf("[DEBUG] 处理第 %d 个种子: %s\n", i+1, torrent.Name)
			tool.AutoCategory(c, torrent)
			tool.DomainTag(c, torrent)
			tool.SeedingLimits(c, torrent)
			tool.StatusTag(c, torrent)
			totalProcessed++
		}

		if len(torrentList) < limit {
			log.Printf("[INFO] 处理完最后一批种子，总共处理: %d\n", totalProcessed)
			break
		}
		offset += limit
		log.Printf("[INFO] 移动到下一批，偏移量: %d\n", offset)
	}
	
	log.Printf("[INFO] 完成所有种子处理，总计: %d\n", totalProcessed)
	log.Printf("[INFO] 程序执行结束 %s\n", time.Now().Format("2006-01-02 15:04:05"))
}

// MultiWriter 实现同时写入多个Writer
type MultiWriter struct {
	writers []io.Writer
}

func (t *MultiWriter) Write(p []byte) (n int, err error) {
	for _, w := range t.writers {
		n, err = w.Write(p)
		if err != nil {
			return
		}
		// 注意：这里应该返回最后一次写入的字节数，而不是循环中的值
	}
	return len(p), nil
}