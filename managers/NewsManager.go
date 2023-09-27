package managers

import (
	"sort"
	"strconv"

	"github.com/CalebRose/SimFBA/dbprovider"
	"github.com/CalebRose/SimFBA/structs"
)

func GetNFLRelatedNews(TeamID string) []structs.NewsLog {
	ts := GetTimestamp()

	newsLogs := GetAllNFLNewsLogs()

	sort.Slice(newsLogs, func(i, j int) bool {
		return newsLogs[i].CreatedAt.Unix() > newsLogs[j].CreatedAt.Unix()
	})

	newsFeed := []structs.NewsLog{}

	recentEventsCount := 0
	personalizedNewsCount := 0
	for _, news := range newsLogs {
		if recentEventsCount == 5 && personalizedNewsCount == 5 {
			break
		}
		if news.SeasonID != ts.NFLSeasonID && news.League != "NFL" {
			continue
		}
		if recentEventsCount < 5 {
			newsFeed = append(newsFeed, news)
			recentEventsCount += 1
		} else if news.TeamID > 0 && strconv.Itoa(news.TeamID) == TeamID && personalizedNewsCount < 5 {
			newsFeed = append(newsFeed, news)
			personalizedNewsCount += 1
		}
	}

	return newsFeed
}

func GetCFBRelatedNews(TeamID string) []structs.NewsLog {
	ts := GetTimestamp()

	newsLogs := GetAllNewsLogs()

	sort.Slice(newsLogs, func(i, j int) bool {
		return newsLogs[i].CreatedAt.Unix() > newsLogs[j].CreatedAt.Unix()
	})

	newsFeed := []structs.NewsLog{}

	recentEventsCount := 0
	personalizedNewsCount := 0
	for _, news := range newsLogs {
		if recentEventsCount == 5 && personalizedNewsCount == 5 {
			break
		}
		if news.SeasonID != ts.CollegeSeasonID && news.League != "CFB" {
			continue
		}
		if news.TeamID == 0 && recentEventsCount < 5 {
			newsFeed = append(newsFeed, news)
			recentEventsCount += 1
		} else if news.TeamID > 0 && strconv.Itoa(news.TeamID) == TeamID && personalizedNewsCount < 5 {
			newsFeed = append(newsFeed, news)
			personalizedNewsCount += 1
		}
	}

	return newsFeed
}

func CreateNewsLog(league, message, messageType string, teamID int, ts structs.Timestamp) {
	db := dbprovider.GetInstance().GetDB()

	seasonID := 0
	weekID := 0
	week := 0
	if league == "CFB" {
		seasonID = ts.CollegeSeasonID
		weekID = ts.CollegeWeekID
		week = ts.CollegeWeek
	} else {
		seasonID = ts.NFLSeasonID
		weekID = ts.NFLWeekID
		week = ts.NFLWeek
	}

	news := structs.NewsLog{
		League:      league,
		Message:     message,
		MessageType: messageType,
		SeasonID:    seasonID,
		WeekID:      weekID,
		Week:        week,
		TeamID:      teamID,
	}

	db.Create(&news)
}
