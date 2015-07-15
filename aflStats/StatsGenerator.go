package aflStats

import (
	"time"

	"github.com/brentonmcs/afl/aflDataAccess"
	"github.com/brentonmcs/afl/aflShared"
)

var pointLow, betCount, wonCount, won40plus, wonUnder40, lose40plus, loseUnder40 int

//GenerateStats loads up the array for the chart on the site
func GenerateStats() []aflShared.StatsModel {
	i := 0
	var result []aflShared.StatsModel

	groupByPoint := groupMatchesByPoints()

	for _, p := range groupByPoint {
		if i == 0 {
			pointLow = p.Point
		}

		betCount += p.BetTotal
		wonCount += p.WonCount
		won40plus += p.WonOver40Count
		wonUnder40 += p.WonUnder40Count
		lose40plus += p.LoseOver40Count
		loseUnder40 += p.LoseUnder40Count

		i++
		if i > 4 {
			result = addToArray(result, p.Point)
			reset()
			i = 0
		}
	}

	// Add Leftover results
	if i > 0 {
		result = addToArray(result, groupByPoint[len(groupByPoint)-1].Point)
	}
	return result
}

// GenerateCurrentRoundStats - get the current round details with stakes and stats
func GenerateCurrentRoundStats() []aflShared.MatchStatPriceModel {
	currentRound := aflDataAccess.GetCurrentRound()
	currentRoundPrices := getCurrentRoundPrices(currentRound)
	stats := GenerateStats()

	var result []aflShared.MatchStatPriceModel
	for _, cur := range aflDataAccess.GetCurrentRoundDetails(currentRound) {

		matchPrices, matchDate := getRoundPrices(convertName(cur.WinTeam), currentRoundPrices)
		if (matchPrices == aflShared.FavPrices{}) {
			continue
		}
		stat, lineStat := aflDataAccess.FindStat(stats, cur.WinPoints, matchPrices.Favourite.LinePoints)

		result = append(result, getStakingInformation(matchPrices, stat, cur, lineStat, matchDate))
	}
	return result
}

func getStakingInformation(matchPrices aflShared.FavPrices, stats aflShared.StatsModel, forecast aflShared.ForecastModel, lineStats aflShared.LinePointsAggregate, matchDate time.Time) aflShared.MatchStatPriceModel {

	return aflShared.MatchStatPriceModel{
		FavPrices:       matchPrices,
		Stats:           stats,
		FavouriteStakes: getTeamStakes(matchPrices.Favourite, stats.WinStat, lineStats),
		OtherTeamStakes: getTeamStakes(matchPrices.OtherTeam, stats.LoseStat, lineStats),
		Forecast:        forecast,
		LineStats:       lineStats,
		MatchDate:       matchDate}
}

func getTeamStakes(prices aflShared.PriceModel, stat aflShared.Stat, line aflShared.LinePointsAggregate) aflShared.StakeModel {

	percentage := line.WonOverLinePercentage
	if prices.LinePoints > 0 {
		percentage = 1.00 - percentage
	}

	return aflShared.StakeModel{HeadToHead: aflShared.KellyCriterion(prices.HeadToHead, stat.WinPercentage),
		Under40: aflShared.KellyCriterion(prices.Under39, stat.WinUnder40Percentage),
		Over40:  aflShared.KellyCriterion(prices.Over40, stat.WinOver40Percentage),
		Line:    aflShared.KellyCriterion(prices.LinePrice, percentage)}
}

func reset() {
	betCount = 0
	wonCount = 0
	won40plus = 0
	wonUnder40 = 0
	lose40plus = 0
	loseUnder40 = 0
}

func addToArray(result []aflShared.StatsModel, pointHigh int) []aflShared.StatsModel {
	return append(result,
		aflShared.StatsModel{PointHigh: pointLow,
			PointLow: pointHigh,
			BetCount: betCount,

			WinStat: aflShared.Stat{
				WonCount:             wonCount,
				Under40:              wonUnder40,
				Plus40:               won40plus,
				WinPercentage:        float32(wonCount) / float32(betCount),
				WinOver40Percentage:  float32(won40plus) / float32(betCount),
				WinUnder40Percentage: float32(wonUnder40) / float32(betCount),
			},
			LoseStat: aflShared.Stat{
				WonCount:             betCount - wonCount,
				Under40:              loseUnder40,
				Plus40:               lose40plus,
				WinPercentage:        float32(betCount-wonCount) / float32(betCount),
				WinOver40Percentage:  float32(lose40plus) / float32(betCount),
				WinUnder40Percentage: float32(loseUnder40) / float32(betCount)}})
}
