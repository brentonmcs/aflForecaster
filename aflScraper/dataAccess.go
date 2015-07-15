package aflScraper

import (
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/brentonmcs/afl/aflShared"
)

func addPriceRecord(matchRecord *aflShared.MatchPrices) {
	query := bson.M{"awayteam.name": matchRecord.AwayTeam.Name, "matchdate": matchRecord.MatchDate, "hometeam.name": matchRecord.HomeTeam.Name}
	aflShared.Update("prices", []string{"awayteam.name", "hometeam.name", "matchdate"}, query, &matchRecord)
}

func getactiveRounds() []aflShared.ActiveRound {

	var results []aflShared.ForecastModel
	var round []aflShared.ActiveRound

	query := bson.M{"resultmodel.winteam": "", "round": bson.M{"$ne": 23}}
	aflShared.Find("forecast", query, func(q *mgo.Query) {
		aflShared.HandleError(q.Sort("year", "round").All(&results))
	})

	curYear := 0
	curRound := 0
	for _, e := range results {
		if curYear != e.Year || curRound != e.Round {
			curYear = e.Year
			curRound = e.Round

			round = append(round, aflShared.ActiveRound{Round: curRound, Year: curYear})
		}
	}
	return round
}

func addForecast(forecast *aflShared.ForecastModel) {
	aflShared.Update("forecast",
		[]string{"round", "winteam", "year"},
		bson.M{"round": forecast.Round, "winteam": forecast.WinTeam, "year": forecast.Year}, &forecast)
}

func removeOldForecasts(updatedTime time.Time, round, year int) {
	aflShared.RemoveAll("forecast", bson.M{"round": round, "year": year, "updated": bson.M{"$exists": false}})
	aflShared.RemoveAll("forecast", bson.M{"round": round, "year": year, "updated": bson.M{"$lt": updatedTime}})
}
