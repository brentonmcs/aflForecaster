package aflPerfectBalance

import (
	"log"
	"strings"
	"time"

	"github.com/brentonmcs/afl/aflShared"
	"github.com/brentonmcs/afl/aflStats"

	"github.com/tealeg/xlsx"
)

func savePrices() {
	excelFileName := "/Users/brentonmcsweyn/afl.xlsx"
	xlFile, err := xlsx.OpenFile(excelFileName)
	aflShared.HandleError(err)

	sheet := xlFile.Sheets[0]

	if sheet.Name != "Data" {
		return
	}

	for i := 2; i < len(sheet.Rows); i++ {

		curRow := sheet.Rows[i]

		if strings.Contains(curRow.Cells[0].String(), "-15") {
			homeTeam := aflShared.PriceModel{Name: curRow.Cells[2].String(), HeadToHead: cellToFloat(curRow.Cells[18]), LinePoints: cellToFloat(curRow.Cells[26]), LinePrice: 1.92}
			awayTeam := aflShared.PriceModel{Name: curRow.Cells[3].String(), HeadToHead: cellToFloat(curRow.Cells[22]), LinePoints: cellToFloat(curRow.Cells[30]), LinePrice: 1.92}
			matchPrice := aflShared.MatchPrices{HomeTeam: homeTeam, AwayTeam: awayTeam, MatchDate: getDate(curRow.Cells[0])}
			addHistoricalPriceRecord(&matchPrice)
		}
	}
}

func getDate(cell *xlsx.Cell) time.Time {

	trimStr := strings.Trim(cell.FormattedValue(), "[$-c09]")
	trimStr = strings.Trim(trimStr, ";@")
	trimStr = strings.Replace(trimStr, "\\", "", 4)

	if len(trimStr) < 8 {
		joinStr := []string{"09", trimStr}
		trimStr = strings.Join(joinStr, "")
	}
	date, err := time.Parse("2-Jan-06", trimStr)
	aflShared.HandleError(err)
	return date
}

func cellToFloat(cell *xlsx.Cell) float32 {
	headToHead, err := cell.Float()

	aflShared.HandleError(err)

	return float32(headToHead)
}

func determineWinProfit() (float32, float32) {

	matches := get2015Matches()
	stats := aflStats.GenerateStats()
	prices := get2015Prices()
	rounds := getRounds()

	balance := float32(100.00)
	deposited := float32(100.00)

	for _, m := range matches {
		if balance < 10 {
			balance += 100
			deposited += 100
		}

		matchPrices := findPrice(prices, m.WinTeam, findRound(rounds, m.Round))
		matchStat, lineStats := findStat(stats, m.WinPoints, matchPrices.Favourite.LinePoints)

		winTeamName, winStake, winPrice, winPoints := getStake(balance, matchPrices, matchStat, lineStats, m.WinPercentage)

		if winStake <= 0 {
			log.Println("No H2H for", winTeamName)
			continue
		}

		if m.aflShared.ResultModel.WinTeam == winTeamName && (float32(m.aflShared.ResultModel.WinPoints)+winPoints) > 0 {
			balance += ((winPrice - 1) * winStake)
		} else {
			balance -= winStake
		}

		log.Printf("Round: %d Team : %s,  Stake: %f, Price : %f,  Bal: %f, Dep: %f", m.Round, winTeamName, winStake, winPrice, balance, deposited)
	}

	return balance, deposited
}

func getStake(balance float32, matchPrices aflShared.FavPrices, matchStat aflShared.StatsModel, lineStat aflShared.LinePointsAggregate, WinPercentage float32) (string, float32, float32, float32) {
	log.Println(matchStat.WinStat.WinPercentage)
	winStake := aflShared.KellyCriterion(matchPrices.Favourite.HeadToHead, 0.71)
	winTeamName := matchPrices.Favourite.Name
	winPrice := matchPrices.Favourite.HeadToHead

	return winTeamName, winStake * balance, winPrice, 0
	// if winStake < 0 {
	// 	winTeamName = matchPrices.OtherTeam.Name
	// 	winPrice = matchPrices.OtherTeam.HeadToHead
	// 	winStake = kellyCriterion(matchPrices.OtherTeam.HeadToHead, matchStat.LoseStat.WinPercentage)
	// }
	//
	// winLineStake := aflShared.KellyCriterion(aflShared.MatchPrices.Favourite.LinePrice, lineStat.WonOverLinePercentage)
	// winLineTeam := matchPrices.Favourite.Name
	// winLinePoints := matchPrices.Favourite.LinePoints
	//
	// if winLineStake < 0 {
	// 	winLineStake = aflShared.KellyCriterion(aflShared.MatchPrices.OtherTeam.LinePrice, 1.00-lineStat.WonOverLinePercentage)
	// 	winLineTeam = aflShared.MatchPrices.OtherTeam.Name
	// 	winLinePoints = aflShared.MatchPrices.OtherTeam.LinePoints
	// }
	//
	// if winLineStake > winStake {
	// 	return winLineTeam, winLineStake * balance, 1.92, winLinePoints
	// }
	//
	// return winTeamName, winStake * balance, winPrice, 0

}


func findPrice(prices []aflShared.MatchPrices, winTeam string, round aflShared.RoundInfo) aflShared.FavPrices {

	winTeam = convertToHistoricalPriceName(winTeam)

	for _, p := range prices {
		if p.MatchDate.Before(round.End.AddDate(0, 0, 1)) && p.MatchDate.After(round.Start.AddDate(0, 0, -1)) {
			if p.HomeTeam.Name == winTeam {
				return aflShared.FavPrices{Favourite: p.HomeTeam, OtherTeam: p.AwayTeam}
			}

			if p.AwayTeam.Name == winTeam {
				return aflShared.FavPrices{OtherTeam: p.HomeTeam, Favourite: p.AwayTeam}
			}
		}
	}

	log.Println("No Price for "+winTeam, round)

	return aflShared.FavPrices{}
}

func findRound(rounds []aflShared.RoundInfo, curRound int) aflShared.RoundInfo {
	for _, r := range rounds {
		if r.Round == curRound {
			return r
		}
	}
	return aflShared.RoundInfo{}
}
