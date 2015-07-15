package aflScraper

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/brentonmcs/afl/aflShared"
)

func scapeSportsBet() {

	doc := goToPage("http://www.sportsbet.com.au/betting/australian-rules?QuickLinks")

	headerAndAccord := doc.Find(".accordion-main").Find(".bettypes-header, .accordion-body")
	var err error
	var currentDate time.Time

	headerAndAccord.Each(func(i int, s *goquery.Selection) {
		if s.HasClass("bettypes-header") {
			date := s.Find(".date").Text()
			currentDate, err = time.Parse("02/01/2006", strings.Split(date, " ")[1])
			if err != nil {
				log.Fatal(err)
			}
		} else {
			buttons := s.Find(".market-buttons")

			homeTeam := parsePrices(buttons.First())
			awayTeam := parsePrices(buttons.Last())

			matchPrice := aflShared.MatchPrices{HomeTeam: homeTeam, AwayTeam: awayTeam, MatchDate: currentDate}
			fmt.Println(matchPrice)
			addPriceRecord(&matchPrice)
		}
	})
}

func parsePrices(prices *goquery.Selection) aflShared.PriceModel {

	priceBox := prices.Find(".price-link")

	headToHeadBox := priceBox.First()
	headToHeadPrice := headToHeadBox.Find(".odd-val").Text()

	if headToHeadPrice == "" {
		return aflShared.PriceModel{}
	}

	priceBox = priceBox.Next()
	under40Price := priceBox.First().Find(".odd-val").Text()

	priceBox = priceBox.Next()
	over40Price := priceBox.First().Find(".odd-val").Text()

	priceBox = priceBox.Next()
	lineBox := priceBox.First()

	linePrice := lineBox.Find(".odd-val").Text()
	linePoints := lineBox.Find(".team-name").Text()

	log.Print(linePoints)
	linePoints = strings.Trim(strings.Trim(strings.Trim(strings.Trim(linePoints, "Line"), " "), "("), ")")

	return aflShared.PriceModel{Name: headToHeadBox.Find(".team-name").Text(),
		HeadToHead: strToPrice(headToHeadPrice),
		Under39:    strToPrice(under40Price),
		Over40:     strToPrice(over40Price),
		LinePrice:  strToPrice(linePrice),
		LinePoints: strToPrice(linePoints)}
}
