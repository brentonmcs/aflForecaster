package aflServer

import (
	"os"

	"github.com/brentonmcs/afl/aflScraper"
	"github.com/brentonmcs/afl/aflStats"
	"github.com/gin-gonic/gin"
)

func determineListenAddress() string {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8000"
	}

	return ":" + port
}

func enableCORS(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type,Token")
	c.Next()
}

//StartHTTPServer adds all the Routes for the Server
func StartHTTPServer() {

	router := gin.Default()

	router.Use(enableCORS)

	router.OPTIONS("/*cors", func(c *gin.Context) {
		// Empty 200 response
	})
	router.GET("/stats", func(c *gin.Context) {
		c.JSON(200, aflStats.GenerateStats())
	})
	router.GET("/currentRoundStats", func(c *gin.Context) {
		c.JSON(200, aflStats.GenerateCurrentRoundStats())
	})

	router.GET("/updateStats", func(c *gin.Context) {
		c.JSON(200, aflScraper.ScrapePages())
	})
	router.GET("/seedPages", func(c *gin.Context) {
		c.JSON(200, aflScraper.SeedPages())
	})

	router.Run(determineListenAddress())
}
