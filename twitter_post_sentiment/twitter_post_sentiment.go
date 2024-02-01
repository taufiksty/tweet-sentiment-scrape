package twitterpostsentiment

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	"github.com/jonreiter/govader"
)

type Tweet struct {
	Name              string  `json:"name"`
	Username          string  `json:"username"`
	Message           string  `json:"message"`
	Sentiment         string  `json:"sentiment"`
	CompoundSentiment float64 `json:"compound_sentiment"`
	PositiveSentiment float64 `json:"positive_sentiment"`
	NegativeSentiment float64 `json:"negative_sentiment"`
	NeutralSentiment  float64 `json:"neutral_sentiment"`
}

func ScrapeTwitterPostSentiment(url string) {
	var tweet Tweet

	ctx, cf := chromedp.NewContext(context.Background())
	defer cf()

	if err := chromedp.Run(ctx, chromedp.Navigate(url)); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Navigated to page", url)

	queryName := "div[data-testid='User-Name'] > div:nth-child(1) > div > a > div > div:nth-child(1) > span > span"
	queryUsername := "div[data-testid='User-Name'] > div:nth-child(2) > div > div > a > div > span"
	queryTweet := "div[data-testid='tweetText'] > span"

	fmt.Println("Wait to query User-Name visible...")
	if err := chromedp.Run(ctx, chromedp.WaitVisible("div[data-testid='User-Name']")); err != nil {
		log.Fatal(err)
	}
	fmt.Println("User-Name is visible")

	fmt.Println("Wait to query tweetText visible...")
	if err := chromedp.Run(ctx, chromedp.WaitVisible("div[data-testid='tweetText']")); err != nil {
		log.Fatal(err)
	}
	fmt.Println("tweetText is visible")

	fmt.Println("Getting the name")
	if err := chromedp.Run(ctx, chromedp.Text(
		queryName,
		&tweet.Name,
	)); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Getting the username")
	if err := chromedp.Run(ctx, chromedp.Text(
		queryUsername,
		&tweet.Username,
	)); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Getting the tweet")
	var tweetTextNodes []*cdp.Node
	if err := chromedp.Run(ctx, chromedp.Nodes(
		queryTweet,
		&tweetTextNodes,
		chromedp.ByQueryAll,
	)); err != nil {
		log.Fatal(err)
	}

	var tweetTexts []string
	for _, textNode := range tweetTextNodes {
		tweetTexts = append(tweetTexts, textNode.Children[0].NodeValue)
	}

	tweet.Message = strings.Join(tweetTexts, "")

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
}
