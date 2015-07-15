package aflScraper

import (
	"log"
	"strconv"
	"strings"

	"github.com/brentonmcs/afl/aflShared"
)

func getPercentage(percentageStr, winTeam string) float32 {
	percentageSplit := strings.SplitAfter(percentageStr, "%")

	var teamPercentageIndex int

	if strings.Contains(percentageSplit[0], winTeam) {
		teamPercentageIndex = 0
	} else {
		teamPercentageIndex = 1
	}

	extractedPercentage := strings.Split(percentageSplit[teamPercentageIndex], " ")
	percentage, err := strconv.ParseFloat(strings.Trim(extractedPercentage[len(extractedPercentage)-1], "%"), 32)
	aflShared.HandleError(err)
	return float32(percentage)
}

func strToPrice(strPrice string) float32 {

	strPrice = strings.Trim(strPrice, " ")
	strPrice = strings.Trim(strPrice, "\n")
	if strPrice == "" {
		return 0
	}
	price, err := strconv.ParseFloat(strPrice, 32)

	if err != nil {
		log.Print("Price was : " + strPrice)
		aflShared.HandleError(err)
	}
	return float32(price)
}
