package repository

import (
	"github.com/CalebRose/SimFBA/dbprovider"
	"github.com/CalebRose/SimFBA/structs"
)

type TradeClauses struct {
	IsAccepted          bool
	IsRejected          bool
	PreloadTradeOptions bool
}

func FindAllTradeProposalsRecords(clauses TradeClauses) []structs.NFLTradeProposal {
	db := dbprovider.GetInstance().GetDB()
	proposal := []structs.NFLTradeProposal{}

	query := db.Model(&proposal)

	if clauses.PreloadTradeOptions {
		query = query.Preload("NFLTeamTradeOptions").Preload("RecepientTeamTradeOptions")
	}

	if clauses.IsAccepted {
		query = query.Where("is_trade_accepted = ?", true)
	}

	if clauses.IsRejected {
		query = query.Where("is_trade_rejected = ?", true)
	}

	if err := query.Find(&proposal).Error; err != nil {
		return []structs.NFLTradeProposal{}
	}

	return proposal
}
