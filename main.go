package main

import (
	"flag"
	"fmt"
	"time"

	twitterpostsentiment "github.com/taufiksty/web-scraper/twitter_post_sentiment"
)

func main() {
	startTime := time.Now()

	url := flag.String("url", "", "url tweet post to scrape")
	if url == nil {
		fmt.Println("url is required")
		return
	}

	flag.Parse()

	twitterpostsentiment.ScrapeTwitterPostSentiment(*url)

	endTime := time.Now()
	duration := endTime.Sub(startTime)

	fmt.Printf("Duration of execution : %.2f s \n", duration.Seconds())
}
