package twitterpostsentimentcolly

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gocolly/colly/v2"
	"github.com/jonreiter/govader"
)

type Tweet struct {
	User          string  `json:"user"`
	Username      string  `json:"username"`
	Tweet         string  `json:"tweet"`
	CompoundScore float64 `json:"compund_score"`
	PositiveScore float64 `json:"positive_score"`
	NeutralScore  float64 `json:"neutral_score"`
	NegativeScore float64 `json:"negative_score"`
}

func ScrapeTweetSentiment(link string) {
	c := colly.NewCollector(
		colly.AllowedDomains("twitter.com", "www.twitter.com"),
		colly.CacheDir("./twitter_post_sentiment/twitter_cache"),
	)

	// avoid bot detection
	c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36"

	c.OnRequest(func(r *colly.Request) {
		log.Println("visiting", r.URL.String())
	})

	var tweet Tweet

	c.OnHTML("div[data-testid='User-Name']", func(h *colly.HTMLElement) {
		log.Println("User found", h.Request.URL)

		tweet.User = h.ChildText("div:nth-child(1) > div > a > div > div:nth-child(1) > span > span")
		tweet.Username = h.ChildText("div:nth-child(2) > div > div > a > div > span")
	})

	c.OnHTML("div[data-testid='tweetText']", func(h *colly.HTMLElement) {
		log.Println("Tweet found", h.Request.URL)

		textElement := []string{}
		h.ForEach("span", func(_ int, t *colly.HTMLElement) {
			textElement = append(textElement, t.Text)
		})

		tweet.Tweet = strings.Join(textElement, "")

		analyzer := govader.NewSentimentIntensityAnalyzer()
		sentiment := analyzer.PolarityScores(tweet.Tweet)

		fmt.Println("Compound score:", sentiment.Compound)
		fmt.Println("Positive score:", sentiment.Positive)
		fmt.Println("Neutral score:", sentiment.Neutral)
		fmt.Println("Negative score:", sentiment.Negative)

		tweet.CompoundScore = sentiment.Compound
		tweet.PositiveScore = sentiment.Positive
		tweet.NeutralScore = sentiment.Neutral
		tweet.NegativeScore = sentiment.Negative
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Println("Error: ", err)
	})

	err := c.Visit(link)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(tweet)

	saveToJSON(tweet)
}

func saveToJSON(tweet Tweet) {
	fileName := fmt.Sprintln("./twitter_post_sentiment/result.json")
	file, err := os.Create(fileName)
	if err != nil {
		log.Fatal("Error creating file:", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	encoder.Encode(tweet)
	if err != nil {
		log.Fatal("Error encoding JSON:", err)
	}

	fmt.Printf("Tweet's sentiment saved to %s\n", fileName)
}
