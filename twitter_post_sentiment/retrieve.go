package twitterpostsentiment

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	"github.com/taufiksty/web-scraper/model"
)

func RetrieveTweet(ctx context.Context, tweet *model.Tweet) *model.Tweet {
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

	return tweet
}
