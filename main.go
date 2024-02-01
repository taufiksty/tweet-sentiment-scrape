package main

import (
	"fmt"
	"time"

	twitterpostsentiment "github.com/taufiksty/web-scraper/twitter_post_sentiment"
)

func main() {
	startTime := time.Now()

	twitterpostsentiment.ScrapeTwitterPostSentiment("https://twitter.com/taylorswift13/status/1734927366378439057")

	endTime := time.Now()
	duration := endTime.Sub(startTime)

	fmt.Printf("Duration of execution : %.2f s \n", duration.Seconds())
}
