package aflDataAccess

import (
	"math"
	"sort"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/brentonmcs/afl/aflShared"
)

//GetCurrentRound returns the current Round and Year
func GetCurrentRound() aflShared.ActiveRound {

	var years, rounds []int

	aflShared.Find("forecast", nil, func(q *mgo.Query) {
		aflShared.HandleError(q.Sort("-year").Distinct("year", &years))
	})

	sort.Sort(sort.Reverse(sort.IntSlice(years)))

	if len(years) == 0 {
		return aflShared.ActiveRound{Year: 0, Round: 0}
	}

	aflShared.Find("forecast", bson.M{"year": years[0], "resultmodel.winteam": bson.M{"$ne": ""}}, func(q *mgo.Query) {
		aflShared.HandleError(q.Sort("-round").Distinct("round", &rounds))
	})

	sort.Sort(sort.Reverse(sort.IntSlice(rounds)))

	return aflShared.ActiveRound{Year: years[0], Round: rounds[0] + 1}
}

//GetCurrentRoundDetails gets the forecasts for the current round
func GetCurrentRoundDetails() []aflShared.ForecastModel {

	activeRound := GetCurrentRound()

	var result []aflShared.ForecastModel
	aflShared.Find("forecast", bson.M{"year": activeRound.Year, "round": activeRound.Round}, func(q *mgo.Query) {
		aflShared.HandleError(q.Sort("order").All(&result))
	})
	return result
}

//GetPercentageForLine - Finds the line Price calculation
func GetPercentageForLine(pointHigh int, pointLow int, linePoints float32) aflShared.LinePointsAggregate {

	eqWon := bson.M{"$eq": []interface{}{"$resultmodel.won", true}}
	eqOverLine := bson.M{"$gte": []interface{}{"$resultmodel.winpoints", math.Abs(float64(linePoints))}}
	and := bson.M{"$and": []interface{}{eqWon, eqOverLine}}

	query := []interface{}{
		bson.M{"$match": bson.M{"resultmodel.winteam": bson.M{"$ne": ""},
			"winpoints": bson.M{"$lte": pointHigh, "$gte": pointLow}}},

		bson.M{"$group": bson.M{
			"_id":         "null",
			"betTotal":    bson.M{"$sum": 1},
			"wonOverLine": bson.M{"$sum": bson.M{"$cond": []interface{}{and, 1, 0}}}}},
		bson.M{"$project": bson.M{"betTotal": "$betTotal", "wonOverLine": "$wonOverLine",
			"wonOverLinePercentage": bson.M{"$divide": []interface{}{"$wonOverLine", "$betTotal"}}}}}

	var result aflShared.LinePointsAggregate
	aflShared.Pipe("forecast", query, func(q *mgo.Pipe) {
		aflShared.HandleError(q.One(&result))
	})
	return result

}
