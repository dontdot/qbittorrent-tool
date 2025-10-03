package tool

import (
	"fengqi/qbittorrent-tool/config"
	"fengqi/qbittorrent-tool/qbittorrent"
	"fengqi/qbittorrent-tool/util"
	"fmt"
	"log"
	"strings"
	"time"
)

// SeedingLimits 做种限制加强版
// 相比较于qb自带，增加根据标签、分类、关键字精确限制
func SeedingLimits(c *config.Config, torrent *qbittorrent.Torrent) {
	if !c.SeedingLimits.Enable || len(c.SeedingLimits.Rules) == 0 {
		log.Printf("[DEBUG] 跳过做种限制 %s: 启用=%v, 规则数量=%d\n", torrent.Name, c.SeedingLimits.Enable, len(c.SeedingLimits.Rules))
		return
	}

	action, limits := matchRule(torrent, c.SeedingLimits.Rules)
	log.Printf("[DEBUG] 匹配规则 %s: 动作=%d\n", torrent.Name, action)
	if action == 0 {
		if !strings.Contains(torrent.State, "paused") || !c.SeedingLimits.Resume {
			log.Printf("[DEBUG] 跳过恢复动作 %s: 已暂停=%v, 恢复设置=%v\n", torrent.Name, strings.Contains(torrent.State, "paused"), c.SeedingLimits.Resume)
			return
		}
	}

	if action == 1 && strings.Contains(torrent.State, "paused") {
		log.Printf("[DEBUG] 种子 %s 已暂停，跳过暂停动作\n", torrent.Name)
		return
	}

	fmt.Printf("动作:%d %s\n", action, torrent.Name)
	executeAction(torrent, action, limits)
}

// 规则至少有一个生效，且生效的全部命中，action才有效，后面的规则会覆盖前面的
func matchRule(torrent *qbittorrent.Torrent, rules []config.SeedingLimitsRule) (int, *config.Limits) {
	action := 0
	var limits *config.Limits
	loc, _ := time.LoadLocation("Asia/Shanghai")

	for i, rule := range rules {
		log.Printf("[DEBUG] 检查规则 #%d %s\n", i, torrent.Name)
		score := 0

		// 分享率
		if rule.Ratio > 0 {
			if torrent.Ratio < rule.Ratio {
				log.Printf("[DEBUG] 规则 #%d 分享率检查失败 %s: 种子分享率=%.2f, 要求=%.2f\n", i, torrent.Name, torrent.Ratio, rule.Ratio)
				continue
			}
			log.Printf("[DEBUG] 规则 #%d 分享率检查通过 %s: 种子分享率=%.2f, 要求=%.2f\n", i, torrent.Name, torrent.Ratio, rule.Ratio)
			score += 1
		}

		// 做种时间，从下载完成算起
		if rule.SeedingTime > 0 {
			if torrent.CompletionOn <= 0 {
				log.Printf("[DEBUG] 规则 #%d 做种时间检查失败 %s: 完成时间=%d\n", i, torrent.Name, torrent.CompletionOn)
				continue
			}
			completionOn := time.Unix(int64(torrent.CompletionOn), 0).In(loc)
			deadOn := completionOn.Add(time.Minute * time.Duration(rule.SeedingTime))
			if time.Now().In(loc).Before(deadOn) {
				log.Printf("[DEBUG] 规则 #%d 做种时间检查失败 %s: 当前时间=%v, 截止时间=%v\n", i, torrent.Name, time.Now().In(loc), deadOn)
				continue
			}
			log.Printf("[DEBUG] 规则 #%d 做种时间检查通过 %s: 完成时间=%v, 要求分钟数=%d\n", i, torrent.Name, completionOn, rule.SeedingTime)
			score += 1
		}

		// 最后活动时间，上传下载等都算
		if rule.ActivityTime > 0 {
			activityOn := time.Unix(int64(torrent.LastActivity), 0).In(loc)
			deadOn := activityOn.Add(time.Minute * time.Duration(rule.ActivityTime))
			if time.Now().In(loc).Before(deadOn) {
				log.Printf("[DEBUG] 规则 #%d 活动时间检查失败 %s: 当前时间=%v, 截止时间=%v\n", i, torrent.Name, time.Now().In(loc), deadOn)
				continue
			}
			log.Printf("[DEBUG] 规则 #%d 活动时间检查通过 %s: 最后活动时间=%v, 要求分钟数=%d\n", i, torrent.Name, activityOn, rule.ActivityTime)
			score += 1
		}

		// 标签
		if len(rule.Tag) != 0 && torrent.Tags != "" {
			tags := strings.Split(torrent.Tags, ",")
			hit := false
		jump:
			for _, item := range rule.Tag {
				for _, item2 := range tags {
					if item == item2 {
						hit = true
						break jump
					}
				}
			}
			if !hit {
				log.Printf("[DEBUG] 规则 #%d 标签检查失败 %s: 种子标签=%v, 要求标签=%v\n", i, torrent.Name, tags, rule.Tag)
				continue
			}
			log.Printf("[DEBUG] 规则 #%d 标签检查通过 %s: 种子标签=%v, 要求标签=%v\n", i, torrent.Name, tags, rule.Tag)
			score += 1
		}

		// 分类
		if len(rule.Category) != 0 && torrent.Category != "" {
			if !util.InArray(torrent.Category, rule.Category) {
				log.Printf("[DEBUG] 规则 #%d 分类检查失败 %s: 种子分类=%s, 要求分类=%v\n", i, torrent.Name, torrent.Category, rule.Category)
				continue
			}
			log.Printf("[DEBUG] 规则 #%d 分类检查通过 %s: 种子分类=%s, 要求分类=%v\n", i, torrent.Name, torrent.Category, rule.Category)
			score += 1
		}

		// tracker  TODO 可能有多个tracker的情况要处理
		tracker, _ := torrent.GetTrackerHost()
		if len(rule.Tracker) != 0 && tracker != "" {
			if !util.InArray(tracker, rule.Tracker) {
				log.Printf("[DEBUG] 规则 #%d Tracker检查失败 %s: 种子Tracker=%s, 要求Trackers=%v\n", i, torrent.Name, tracker, rule.Tracker)
				continue
			}
			log.Printf("[DEBUG] 规则 #%d Tracker检查通过 %s: 种子Tracker=%s, 要求Trackers=%v\n", i, torrent.Name, tracker, rule.Tracker)
			score += 1
		}

		// 做种数大于
		if rule.SeedsGt > 0 {
			if torrent.NumComplete < rule.SeedsGt {
				log.Printf("[DEBUG] 规则 #%d 做种数大于检查失败 %s: 做种数=%d, 要求=%d\n", i, torrent.Name, torrent.NumComplete, rule.SeedsGt)
				continue
			}
			log.Printf("[DEBUG] 规则 #%d 做种数大于检查通过 %s: 做种数=%d, 要求=%d\n", i, torrent.Name, torrent.NumComplete, rule.SeedsGt)
			score += 1
		}

		// 做种数小于
		if rule.SeedsLt > 0 {
			if torrent.NumComplete > rule.SeedsLt {
				log.Printf("[DEBUG] 规则 #%d 做种数小于检查失败 %s: 做种数=%d, 要求=%d\n", i, torrent.Name, torrent.NumComplete, rule.SeedsLt)
				continue
			}
			log.Printf("[DEBUG] 规则 #%d 做种数小于检查通过 %s: 做种数=%d, 要求=%d\n", i, torrent.Name, torrent.NumComplete, rule.SeedsLt)
			score += 1
		}

		// 关键字
		if len(rule.Keyword) != 0 {
			if !util.ContainsArray(torrent.Name, rule.Keyword) {
				log.Printf("[DEBUG] 规则 #%d 关键字检查失败 %s: 种子名称=%s, 要求关键字=%v\n", i, torrent.Name, torrent.Name, rule.Keyword)
				continue
			}
			log.Printf("[DEBUG] 规则 #%d 关键字检查通过 %s: 种子名称=%s, 要求关键字=%v\n", i, torrent.Name, torrent.Name, rule.Keyword)
			score += 1
		}

		if score > 0 {
			action = rule.Action
			log.Printf("[DEBUG] 规则 #%d 匹配成功，得分 %d %s, 设置动作为 %d\n", i, score, torrent.Name, rule.Action)
		}
		if action == 0 && limits == nil {
			limits = rule.Limits
			log.Printf("[DEBUG] 从规则 #%d 设置限制 %s\n", i, torrent.Name)
		}
	}

	log.Printf("[DEBUG] 最终动作 %s: %d\n", torrent.Name, action)
	return action, limits
}

func executeAction(torrent *qbittorrent.Torrent, action int, limits *config.Limits) {
	log.Printf("[DEBUG] 执行动作 %d %s\n", action, torrent.Name)
	switch action {
	case 0:
		_ = qbittorrent.Api.ResumeTorrents(torrent.Hash)
		if limits == nil {
			log.Printf("[DEBUG] 恢复种子 %s 无限制\n", torrent.Name)
			break
		}
		if limits.Download != torrent.DlLimit {
			_ = qbittorrent.Api.SetDownloadLimit(torrent.Hash, limits.Download)
			log.Printf("[DEBUG] 设置下载限制 %s: %d\n", torrent.Name, limits.Download)
		}
		if limits.Upload != torrent.UpLimit {
			_ = qbittorrent.Api.SetUploadLimit(torrent.Hash, limits.Download)
			log.Printf("[DEBUG] 设置上传限制 %s: %d\n", torrent.Name, limits.Upload)
		}

		flag := false
		radio := torrent.RatioLimit
		if limits.Ratio != torrent.RatioLimit {
			flag = true
			radio = limits.Ratio
			log.Printf("[DEBUG] 分享率限制已更改 %s: %f -> %f\n", torrent.Name, torrent.RatioLimit, limits.Ratio)
		}
		seedingTimeLimit := torrent.SeedingTimeLimit
		if limits.SeedingTime != torrent.SeedingTimeLimit {
			flag = true
			seedingTimeLimit = limits.SeedingTime
			log.Printf("[DEBUG] 做种时间限制已更改 %s: %d -> %d\n", torrent.Name, torrent.SeedingTimeLimit, limits.SeedingTime)
		}
		inactiveSeedingTimeLimit := torrent.InactiveSeedingTimeLimit
		if limits.InactiveSeedingTime != torrent.InactiveSeedingTimeLimit {
			flag = true
			inactiveSeedingTimeLimit = limits.InactiveSeedingTime
			log.Printf("[DEBUG] 不活跃做种时间限制已更改 %s: %d -> %d\n", torrent.Name, torrent.InactiveSeedingTimeLimit, limits.InactiveSeedingTime)
		}
		if flag {
			_ = qbittorrent.Api.SetShareLimit(torrent.Hash, radio, seedingTimeLimit, inactiveSeedingTimeLimit)
			log.Printf("[DEBUG] 更新分享限制 %s\n", torrent.Name)
		}

		break

	case 1:
		_ = qbittorrent.Api.PauseTorrents(torrent.Hash)
		log.Printf("[DEBUG] 暂停种子 %s\n", torrent.Name)
		break

	case 2:
		_ = qbittorrent.Api.DeleteTorrents(torrent.Hash, false)
		log.Printf("[DEBUG] 删除种子 %s (不含文件)\n", torrent.Name)
		break

	case 3:
		_ = qbittorrent.Api.DeleteTorrents(torrent.Hash, true)
		log.Printf("[DEBUG] 删除种子 %s (含文件)\n", torrent.Name)
		break

	case 4:
		_ = qbittorrent.Api.SetSuperSeeding(torrent.Hash, true)
		log.Printf("[DEBUG] 启用超级做种 %s\n", torrent.Name)
		break
	}
}
