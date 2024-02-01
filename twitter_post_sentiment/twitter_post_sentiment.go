package twitterpostsentiment

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

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

const queryName = "div[data-testid='User-Name'] > div:nth-child(1) > div > a > div > div:nth-child(1) > span > span"
const queryUsername = "div[data-testid='User-Name'] > div:nth-child(2) > div > div > a > div > span"
const queryTweet = "body div[data-testid='tweetText']"

func waitVisible(ctx context.Context, testId string) {
	fmt.Printf("Wait to query %s visible...\n", testId)
	if err := chromedp.Run(ctx, chromedp.WaitVisible(fmt.Sprintf("div[data-testid='%s']", testId))); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s is visible", testId)
}

func retrieveTweet(ctx context.Context, tweet *Tweet) *Tweet {
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

	fmt.Println(len(tweetTextNodes))
	var tweetTexts []string
	for _, textNode := range tweetTextNodes {
		tweetTexts = append(tweetTexts, textNode.Children[0].NodeValue)
	}

	tweet.Message = strings.Join(tweetTexts, "")

	return tweet
}

func ScrapeTwitterPostSentiment(url string) {
	// var tweets []Tweet
	var tweet Tweet

	ctx, cf := chromedp.NewContext(context.Background())
	defer cf()

	ctx, cf = context.WithTimeout(ctx, 30*time.Second)
	defer cf()

	if err := chromedp.Run(ctx, chromedp.Navigate(url)); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Navigated to page", url)

	waitVisible(ctx, "User-Name")
	waitVisible(ctx, "tweetText")

	retrieveTweet(ctx, &tweet)

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
