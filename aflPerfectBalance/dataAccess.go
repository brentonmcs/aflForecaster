package aflPerfectBalance

import (
	"../aflShared"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func addHistoricalPriceRecord(matchRecord *aflShared.MatchPrices) {

	aflShared.NewConnect(func(db *mgo.Database) interface{} {
		c := db.C("historicalprices")
		aflShared.AddIndex(c, []string{"awayteam.name", "hometeam.name", "matchdate"})

		var query = bson.M{"awayteam.name": matchRecord.AwayTeam.Name,
			"matchdate":     matchRecord.MatchDate,
			"hometeam.name": matchRecord.HomeTeam.Name}
		aflShared.UpdateRecord(c, query, &matchRecord)
		return 0
	})


}
