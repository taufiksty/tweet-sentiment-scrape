package instagrampostsentiment

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"
	"github.com/jonreiter/govader"
)

type Post struct {
	Username      string  `json:"username"`
	Message       string  `json:"message"`
	CompoundScore float64 `json:"compund_score"`
	PositiveScore float64 `json:"positive_score"`
	NeutralScore  float64 `json:"neutral_score"`
	NegativeScore float64 `json:"negative_score"`
}

func ScrapeInstagramSentiment(link string) {
	c := colly.NewCollector(
		colly.AllowedDomains("instagram.com", "www.instagram.com"),
		colly.CacheDir("./instagram_post_sentiment/instagram_cache"),
	)

	// avoid bot detection
	c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36"

	c.OnRequest(func(r *colly.Request) {
		log.Println("visiting", r.URL.String())
	})

	queryStringUsername := "_ap3a _aaco _aacw _aacx _aad7 _aade"
	queryStringMessage := "x193iq5w xeuugli x1fj9vlw x13faqbe x1vvkbs xt0psk2 x1i0vuye xvs91rp xo1l8bm x5n08af x10wh9bi x1wdrske x8viiok x18hxmgj"

	queryUsername := fmt.Sprintf("span.%s:nth-child(1)", strings.Join(strings.Split(queryStringUsername, " "), "."))
	queryMessage := fmt.Sprintf("span.%s", strings.Join(strings.Split(queryStringMessage, " "), "."))

	fmt.Println(queryUsername)
	fmt.Println(queryMessage)

	var post Post

	c.OnHTML(queryUsername, func(h *colly.HTMLElement) {
		log.Println("User found on", h.Request.URL)

		post.Username = h.Text
	})

	c.OnHTML("body", func(h *colly.HTMLElement) {
		log.Println("Message found on", h.Request.URL)

		var msgElement []string
		h.ForEach(queryMessage, func(i int, t *colly.HTMLElement) {
			log.Println("Message found", strconv.Itoa(i))
			msgElement = append(msgElement, t.Text)
			if i >= 49 {
				return
			}
		})

		messagesText := strings.Join(msgElement, " ")

		analyzer := govader.NewSentimentIntensityAnalyzer()
		sentiment := analyzer.PolarityScores(messagesText)

		fmt.Println("Posting by", post.Username)
		fmt.Println("Compound score:", sentiment.Compound)
		fmt.Println("Positive score:", sentiment.Positive)
		fmt.Println("Neutral score:", sentiment.Neutral)
		fmt.Println("Negative score:", sentiment.Negative)

		post.CompoundScore = sentiment.Compound
		post.PositiveScore = sentiment.Positive
		post.NeutralScore = sentiment.Neutral
		post.NegativeScore = sentiment.Negative
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Println("Error: ", err)
	})

	err := c.Visit(link)
	if err != nil {
		log.Fatal(err)
	}

	saveToJSON(post)
}

func saveToJSON(tweet Post) {
	fileName := fmt.Sprintf("./instagram_post_sentiment/%s.json", tweet.Username)
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

	fmt.Printf("Post's sentiment saved to %s\n", fileName)
}
