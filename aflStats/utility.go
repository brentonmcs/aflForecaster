package aflStats

import (
	"time"

	"github.com/brentonmcs/afl/aflShared"
)

func getRoundPrices(winTeam string, prices []aflShared.MatchPrices) (aflShared.FavPrices, time.Time) {

	for _, r := range prices {

		if r.HomeTeam.Name == winTeam {
			return aflShared.FavPrices{Favourite: r.HomeTeam, OtherTeam: r.AwayTeam}, r.MatchDate
		}
		if r.AwayTeam.Name == winTeam {
			return aflShared.FavPrices{Favourite: r.AwayTeam, OtherTeam: r.HomeTeam}, r.MatchDate
		}
	}

	return aflShared.FavPrices{}, time.Time{}
}

func convertName(winTeam string) string {
	if winTeam == "GW Sydney" {
		return "Greater Western Sydney"
	}
	if winTeam == "Wstn Bulldogs" {
		return "Western Bulldogs"
	}
	if winTeam == "Nth Melbourne" {
		return "North Melbourne"
	}
	return winTeam
}
