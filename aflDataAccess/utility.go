package aflDataAccess

import (
	"log"

	"github.com/brentonmcs/afl/aflShared"
)

//FindStat - finds the stat that fits the range of the winPoints
func FindStat(stats []aflShared.StatsModel, winPoints int, linePoints float32) (aflShared.StatsModel, aflShared.LinePointsAggregate) {
	for _, s := range stats {
		if s.PointHigh >= winPoints && s.PointLow <= winPoints {

			return s, GetPercentageForLine(s.PointHigh, s.PointLow, linePoints)
		}
	}

	log.Println("No Stats found?", winPoints)
	return aflShared.StatsModel{}, aflShared.LinePointsAggregate{}
}
