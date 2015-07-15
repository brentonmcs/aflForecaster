package aflShared

import "time"

//AggregatePoints is the Points Stats
type AggregatePoints struct {
	Point            int `bson:"_id" json:"_id"`
	WonCount         int `bson:"wonCount" json:"wonCount"`
	BetTotal         int `bson:"betTotal" json:"betTotal"`
	WonOver40Count   int `bson:"wonOver40Count" json:"wonOver40Count"`
	WonUnder40Count  int `bson:"wonUnder40Count" json:"wonUnder40Count"`
	LoseOver40Count  int `bson:"loseOver40Count" json:"loseOver40Count"`
	LoseUnder40Count int `bson:"loseUnder40Count" json:"loseUnder40Count"`
}

//LinePointsAggregate is the Points Stats
type LinePointsAggregate struct {
	Point                 int     `bson:"_id" json:"_id"`
	BetTotal              int     `bson:"betTotal" json:"betTotal"`
	WonOverLine           int     `bson:"wonOverLine" json:"wonOverLine"`
	WonOverLinePercentage float32 `bson:"wonOverLinePercentage" json:"wonOverLinePercentage"`
}

// ActiveRound holds the round and year that is being access
type ActiveRound struct {
	Round int
	Year  int
}

// ForecastModel is the database Structure for reading off the website
type ForecastModel struct {
	WinTeam       string
	WinPercentage float32
	WinPoints     int
	Round         int
	ResultModel   ResultModel
	Year          int
	Updated       time.Time
	Order         int
}

// ResultModel is the results from the website
type ResultModel struct {
	Won       bool
	WinTeam   string
	WinPoints int
}

// StatsModel is the generated model for the Graph
type StatsModel struct {
	PointLow  int
	PointHigh int
	BetCount  int
	WinStat   Stat
	LoseStat  Stat
}

//Stat individual Stat
type Stat struct {
	WonCount             int
	Under40              int
	Plus40               int
	WinPercentage        float32
	WinOver40Percentage  float32
	WinUnder40Percentage float32
}

// PriceModel is the prices for each of the bets
type PriceModel struct {
	Name       string
	HeadToHead float32
	Under39    float32
	Over40     float32
	LinePrice  float32
	LinePoints float32
}

// FavPrices is the prices for each of the bets for the favourite and the Other team
type FavPrices struct {
	Favourite PriceModel
	OtherTeam PriceModel
}

// MatchPrices is the prices for each of the bets for both teams in the match
type MatchPrices struct {
	HomeTeam  PriceModel
	AwayTeam  PriceModel
	MatchDate time.Time
}

// StakeModel is the amount to bet on each market
type StakeModel struct {
	HeadToHead float32
	Under40    float32
	Over40     float32
	Line       float32
}

// MatchStatPriceModel is the view model for the stats, prices and stakes for a match
type MatchStatPriceModel struct {
	Stats           StatsModel
	FavPrices       FavPrices
	FavouriteStakes StakeModel
	OtherTeamStakes StakeModel
	Forecast        ForecastModel
	LineStats       LinePointsAggregate
	MatchDate       time.Time
}

//ByDate is the sorting interface for Dates
type ByDate []MatchStatPriceModel

func (s ByDate) Len() int {
	return len(s)
}
func (s ByDate) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s ByDate) Less(i, j int) bool {
	return s[i].MatchDate.Before(s[j].MatchDate)
}

//RoundInfo is the data structure for the Round
type RoundInfo struct {
	Round int       `bson:"round" json:"round"`
	Start time.Time `bson:"start" json:"start"`
	End   time.Time `bson:"end" json:"end"`
}
