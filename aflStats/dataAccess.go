package aflStats

import (
	"github.com/brentonmcs/afl/aflShared"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func getCurrentRoundPrices(currentRound aflShared.ActiveRound) []aflShared.MatchPrices {

	var roundInfo aflShared.RoundInfo
	aflShared.Find("round", bson.M{"round": currentRound.Round}, func(q *mgo.Query) {
		aflShared.HandleError(q.One(&roundInfo))
	})

	var result []aflShared.MatchPrices
	aflShared.Find("prices", bson.M{"matchdate": bson.M{"$lt": roundInfo.End, "$gte": roundInfo.Start}}, func(q *mgo.Query) {
		q.All(&result)
	})
	return result
}

func groupMatchesByPoints() []aflShared.AggregatePoints {

	eqWon := bson.M{"$eq": []interface{}{"$resultmodel.won", true}}
	eqLose := bson.M{"$eq": []interface{}{"$resultmodel.won", false}}
	eqOver40 := bson.M{"$gte": []interface{}{"$resultmodel.winpoints", 40}}
	eqUnder40 := bson.M{"$lt": []interface{}{"$resultmodel.winpoints", 40}}

	and := bson.M{"$and": []interface{}{eqWon, eqOver40}}
	andUnder := bson.M{"$and": []interface{}{eqWon, eqUnder40}}
	andLose := bson.M{"$and": []interface{}{eqLose, eqOver40}}
	andUnderLose := bson.M{"$and": []interface{}{eqLose, eqUnder40}}

	query := []interface{}{
		bson.M{"$match": bson.M{"resultmodel.winteam": bson.M{"$ne": ""}, "round": bson.M{"$gte": 7}}},
		bson.M{"$group": bson.M{
			"_id":              "$winpoints",
			"betTotal":         bson.M{"$sum": 1},
			"wonCount":         bson.M{"$sum": bson.M{"$cond": []interface{}{eqWon, 1, 0}}},
			"wonOver40Count":   bson.M{"$sum": bson.M{"$cond": []interface{}{and, 1, 0}}},
			"wonUnder40Count":  bson.M{"$sum": bson.M{"$cond": []interface{}{andUnder, 1, 0}}},
			"loseOver40Count":  bson.M{"$sum": bson.M{"$cond": []interface{}{andLose, 1, 0}}},
			"loseUnder40Count": bson.M{"$sum": bson.M{"$cond": []interface{}{andUnderLose, 1, 0}}}}},
		bson.M{"$sort": bson.M{"_id": -1}}}

	var results []aflShared.AggregatePoints
	aflShared.Pipe("forecast", query, func(q *mgo.Pipe) {
		aflShared.HandleError(q.All(&results))
	})

	return results
}
