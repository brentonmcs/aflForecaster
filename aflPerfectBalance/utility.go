package aflPerfectBalance

func convertToHistoricalPriceName(winTeam string) string {

	if winTeam == "GW Sydney" {
		return "GWS Giants"
	}
	if winTeam == "Wstn Bulldogs" {
		return "Western Bulldogs"
	}
	if winTeam == "Nth Melbourne" {
		return "North Melbourne"
	}

	return winTeam
}
