package managers

import (
	"sort"
	"strconv"

	"github.com/CalebRose/SimFBA/dbprovider"
	"github.com/CalebRose/SimFBA/structs"
)

func GetAllCollegePollsByWeekIDAndSeasonID(weekID, seasonID string) []structs.CollegePollSubmission {
	db := dbprovider.GetInstance().GetDB()

	submissions := []structs.CollegePollSubmission{}

	err := db.Where("week_id = ? AND season_id = ?", weekID, seasonID).Find(&submissions).Error
	if err != nil {
		return []structs.CollegePollSubmission{}
	}

	return submissions
}

func GetPollSubmissionBySubmissionID(id string) structs.CollegePollSubmission {
	db := dbprovider.GetInstance().GetDB()

	submission := structs.CollegePollSubmission{}

	err := db.Where("id = ?", id).Find(&submission).Error
	if err != nil {
		return structs.CollegePollSubmission{}
	}

	return submission
}

func GetPollSubmissionByUsernameWeekAndSeason(username string) structs.CollegePollSubmission {
	db := dbprovider.GetInstance().GetDB()
	ts := GetTimestamp()
	weekID := strconv.Itoa(int(ts.CollegeWeekID + 1))
	seasonID := strconv.Itoa(int(ts.CollegeSeasonID))

	submission := structs.CollegePollSubmission{}

	err := db.Where("username = ? AND week_id = ? AND season_id = ?", username, weekID, seasonID).Find(&submission).Error
	if err != nil {
		return structs.CollegePollSubmission{}
	}

	return submission
}

func SyncCollegePollSubmissionForCurrentWeek(week, weekID, seasonID uint) {
	db := dbprovider.GetInstance().GetDB()

	weekIDStr := strconv.Itoa(int(weekID))
	seasonIDStr := strconv.Itoa(int(seasonID))
	standingsMap := GetCollegeStandingsMap(seasonIDStr)

	submissions := GetAllCollegePollsByWeekIDAndSeasonID(weekIDStr, seasonIDStr)

	allCollegeTeams := GetAllCollegeTeams()

	voteMap := make(map[uint]*structs.TeamVote)

	for _, t := range allCollegeTeams {
		voteMap[t.ID] = &structs.TeamVote{TeamID: t.ID, Team: t.TeamAbbr}
	}

	for _, s := range submissions {
		if week > 3 {
			// Invalid check
			s1standings := standingsMap[s.Rank1ID]
			if s1standings.TotalWins == 0 || !s1standings.IsFBS {
				continue
			}
			s2standings := standingsMap[s.Rank2ID]
			if s2standings.TotalWins == 0 || !s2standings.IsFBS {
				continue
			}
			s3standings := standingsMap[s.Rank3ID]
			if s3standings.TotalWins == 0 || !s3standings.IsFBS {
				continue
			}
			s4standings := standingsMap[s.Rank4ID]
			if s4standings.TotalWins == 0 || !s4standings.IsFBS {
				continue
			}
			s5standings := standingsMap[s.Rank5ID]
			if s5standings.TotalWins == 0 || !s5standings.IsFBS {
				continue
			}
			s6standings := standingsMap[s.Rank6ID]
			if s6standings.TotalWins == 0 || !s6standings.IsFBS {
				continue
			}
			s7standings := standingsMap[s.Rank7ID]
			if s7standings.TotalWins == 0 || !s7standings.IsFBS {
				continue
			}
			s8standings := standingsMap[s.Rank8ID]
			if s8standings.TotalWins == 0 || !s8standings.IsFBS {
				continue
			}
			s9standings := standingsMap[s.Rank9ID]
			if s9standings.TotalWins == 0 || !s9standings.IsFBS {
				continue
			}
			s10standings := standingsMap[s.Rank10ID]
			if s10standings.TotalWins == 0 || !s10standings.IsFBS {
				continue
			}
			s11standings := standingsMap[s.Rank11ID]
			if s11standings.TotalWins == 0 || !s11standings.IsFBS {
				continue
			}
			s12standings := standingsMap[s.Rank12ID]
			if s12standings.TotalWins == 0 || !s12standings.IsFBS {
				continue
			}
			s13standings := standingsMap[s.Rank13ID]
			if s13standings.TotalWins == 0 || !s13standings.IsFBS {
				continue
			}
			s14standings := standingsMap[s.Rank14ID]
			if s14standings.TotalWins == 0 || !s14standings.IsFBS {
				continue
			}
			s15standings := standingsMap[s.Rank15ID]
			if s15standings.TotalWins == 0 || !s15standings.IsFBS {
				continue
			}
			s16standings := standingsMap[s.Rank16ID]
			if s16standings.TotalWins == 0 || !s16standings.IsFBS {
				continue
			}
			s17standings := standingsMap[s.Rank17ID]
			if s17standings.TotalWins == 0 || !s17standings.IsFBS {
				continue
			}
			s18standings := standingsMap[s.Rank18ID]
			if s18standings.TotalWins == 0 || !s18standings.IsFBS {
				continue
			}
			s19standings := standingsMap[s.Rank19ID]
			if s19standings.TotalWins == 0 || !s19standings.IsFBS {
				continue
			}
			s20standings := standingsMap[s.Rank20ID]
			if s20standings.TotalWins == 0 || !s20standings.IsFBS {
				continue
			}
			s21standings := standingsMap[s.Rank21ID]
			if s21standings.TotalWins == 0 || !s21standings.IsFBS {
				continue
			}
			s22standings := standingsMap[s.Rank22ID]
			if s22standings.TotalWins == 0 || !s22standings.IsFBS {
				continue
			}
			s23standings := standingsMap[s.Rank23ID]
			if s23standings.TotalWins == 0 || !s23standings.IsFBS {
				continue
			}
			s24standings := standingsMap[s.Rank24ID]
			if s24standings.TotalWins == 0 || !s24standings.IsFBS {
				continue
			}
			s25standings := standingsMap[s.Rank25ID]
			if s25standings.TotalWins == 0 || !s25standings.IsFBS {
				continue
			}
		}
		voteMap[s.Rank1ID].AddVotes(1)
		voteMap[s.Rank2ID].AddVotes(2)
		voteMap[s.Rank3ID].AddVotes(3)
		voteMap[s.Rank4ID].AddVotes(4)
		voteMap[s.Rank5ID].AddVotes(5)
		voteMap[s.Rank6ID].AddVotes(6)
		voteMap[s.Rank7ID].AddVotes(7)
		voteMap[s.Rank8ID].AddVotes(8)
		voteMap[s.Rank9ID].AddVotes(9)
		voteMap[s.Rank10ID].AddVotes(10)
		voteMap[s.Rank11ID].AddVotes(11)
		voteMap[s.Rank12ID].AddVotes(12)
		voteMap[s.Rank13ID].AddVotes(13)
		voteMap[s.Rank14ID].AddVotes(14)
		voteMap[s.Rank15ID].AddVotes(15)
		voteMap[s.Rank16ID].AddVotes(16)
		voteMap[s.Rank17ID].AddVotes(17)
		voteMap[s.Rank18ID].AddVotes(18)
		voteMap[s.Rank19ID].AddVotes(19)
		voteMap[s.Rank20ID].AddVotes(20)
		voteMap[s.Rank21ID].AddVotes(21)
		voteMap[s.Rank22ID].AddVotes(22)
		voteMap[s.Rank23ID].AddVotes(23)
		voteMap[s.Rank24ID].AddVotes(24)
		voteMap[s.Rank25ID].AddVotes(25)
	}

	allVotes := []structs.TeamVote{}

	for _, t := range allCollegeTeams {
		v := voteMap[t.ID]
		if v.TotalVotes == 0 {
			continue
		}
		newVoteObj := structs.TeamVote{TeamID: v.TeamID, Team: v.Team, TotalVotes: v.TotalVotes, Number1Votes: v.Number1Votes}

		allVotes = append(allVotes, newVoteObj)
	}

	sort.Slice(allVotes, func(i, j int) bool {
		return allVotes[i].TotalVotes > allVotes[j].TotalVotes
	})

	officialPoll := structs.CollegePollOfficial{
		WeekID:   weekID,
		Week:     week,
		SeasonID: seasonID,
	}
	count := 0
	for idx, v := range allVotes {
		if count > 24 {
			break
		}

		count += 1
		officialPoll.AssignRank(idx, v)
		// Get Standings
		teamID := strconv.Itoa(int(v.TeamID))
		teamStandings := GetCollegeStandingsRecordByTeamID(teamID, seasonIDStr)
		rank := idx + 1
		teamStandings.AssignRank(rank)
		db.Save(&teamStandings)

		matches := GetCollegeGamesByTeamIdAndSeasonId(teamID, seasonIDStr, false)

		for _, m := range matches {
			if m.Week < int(week) {
				continue
			}
			if m.Week > int(week) {
				break
			}
			m.AssignRank(v.TeamID, uint(rank))
			db.Save(&m)
		}
	}

	db.Create(&officialPoll)
}

func CreatePoll(dto structs.CollegePollSubmission) structs.CollegePollSubmission {
	db := dbprovider.GetInstance().GetDB()
	existingPoll := GetPollSubmissionBySubmissionID(strconv.Itoa(int(dto.ID)))

	if existingPoll.ID > 0 {
		dto.AssignID(existingPoll.ID)
		db.Save(&dto)
	} else {
		db.Create(&dto)
	}

	return dto
}

func GetAllOfficialPolls() []structs.CollegePollOfficial {
	db := dbprovider.GetInstance().GetDB()
	officialPoll := []structs.CollegePollOfficial{}

	err := db.Find(&officialPoll).Error
	if err != nil {
		return []structs.CollegePollOfficial{}
	}

	return officialPoll
}

func GetOfficialPollBySeasonID(seasonID string) []structs.CollegePollOfficial {
	db := dbprovider.GetInstance().GetDB()
	officialPoll := []structs.CollegePollOfficial{}

	err := db.Where("season_id = ?", seasonID).Find(&officialPoll).Error
	if err != nil {
		return []structs.CollegePollOfficial{}
	}

	return officialPoll
}
