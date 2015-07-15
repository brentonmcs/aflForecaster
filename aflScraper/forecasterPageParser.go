package aflScraper

import (
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/brentonmcs/afl/aflShared"
)

func parseForecast(forecast, percentageStr, resultStr string, i, year int, order int, updatedTime time.Time) aflShared.ForecastModel {
	forecastSplit := strings.Split(forecast, "by")
	winTeam := strings.Trim(forecastSplit[0], " ")

	var resultModel aflShared.ResultModel
	if len(resultStr) > 10 {
		resultModel = getResult(resultStr, winTeam)
	}
	return aflShared.ForecastModel{
		WinTeam:       winTeam,
		WinPercentage: getPercentage(percentageStr, winTeam),
		WinPoints:     int(getPoints(forecastSplit[1])),
		Round:         i,
		Year:          year,
		Order:         order,
		Updated:       updatedTime,
		ResultModel:   resultModel}
}

func getPoints(pointsStr string) int {
	pointsSplit := strings.Split(strings.Trim(pointsStr, " "), " ")

	points, err := strconv.ParseInt(pointsSplit[0], 10, 0)
	if err != nil {
		log.Fatal(err)
	}
	return int(points)
}

func getResult(resultStr string, winTeam string) aflShared.ResultModel {
	bySplit := strings.Split(resultStr, "by")
	won := strings.Trim(bySplit[0], " ")
	return aflShared.ResultModel{Won: won == winTeam, WinPoints: int(getPoints(bySplit[1])), WinTeam: won}
}
