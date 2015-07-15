package aflDataAccess

import ("github.com/brentonmcs/afl/aflShared"
"log")

//FindStat - extracts Stat that matches the point range

func FindStat(stats []aflShared.StatsModel, winPoints int, linePoints float32) (aflShared.StatsModel, aflShared.LinePointsAggregate) {
	for _, s := range stats {
		if s.PointHigh >= winPoints && s.PointLow <= winPoints {

			return s, GetPercentageForLine(s.PointHigh, s.PointLow, linePoints)
		}
	}

	log.Println("No Stats found?", winPoints)
	return aflShared.StatsModel{}, aflShared.LinePointsAggregate{}
}