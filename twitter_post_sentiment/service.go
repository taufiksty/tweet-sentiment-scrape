package twitterpostsentiment

import (
	"context"
	"fmt"
	"log"

	"github.com/chromedp/chromedp"
	"github.com/jonreiter/govader"
	"github.com/taufiksty/web-scraper/model"
)

const queryName = "div[data-testid='User-Name'] > div:nth-child(1) > div > a > div > div:nth-child(1) > span > span"
const queryUsername = "div[data-testid='User-Name'] > div:nth-child(2) > div > div > a > div > span"
const queryTweet = "div[data-testid='tweetText'] > span"

func waitVisible(ctx context.Context, testId string) {
	fmt.Printf("Wait to query %s visible...\n", testId)
	if err := chromedp.Run(ctx, chromedp.WaitVisible(fmt.Sprintf("div[data-testid='%s']", testId))); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s is visible", testId)
}

func ScrapeTwitterPostSentiment(url string) {
	var tweet model.Tweet

	ctx, cf := chromedp.NewContext(context.Background())
	defer cf()

	if err := chromedp.Run(ctx, chromedp.Navigate(url)); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Navigated to page", url)

	waitVisible(ctx, "User-Name")
	waitVisible(ctx, "tweetText")

	RetrieveTweet(ctx, &tweet)

	analyzer := govader.NewSentimentIntensityAnalyzer()
	sentiment := analyzer.PolarityScores(tweet.Message)

	tweet.CompoundSentiment = sentiment.Compound
	tweet.PositiveSentiment = sentiment.Positive
	tweet.NegativeSentiment = sentiment.Negative
	tweet.NeutralSentiment = sentiment.Neutral

	if tweet.CompoundSentiment > 6 {
		tweet.Sentiment = "Very Positive"
	} else if tweet.CompoundSentiment > 2 {
		tweet.Sentiment = "Positive"
	} else if tweet.CompoundSentiment > -2 {
		tweet.Sentiment = "Neutral"
	} else if tweet.CompoundSentiment > -6 {
		tweet.Sentiment = "Negative"
	} else {
		tweet.Sentiment = "Very Negative"
	}

	fmt.Printf("Found name : %s\n", tweet.Name)
	fmt.Printf("Found username : %s\n", tweet.Username)
	fmt.Printf("Found message : \n%s\n", tweet.Message)
	fmt.Printf("Sentiment : %s\n", tweet.Sentiment)
	fmt.Printf("Sentiment compound : %f\n", tweet.CompoundSentiment)

	SaveToJSON(&tweet)
}
