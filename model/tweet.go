package model

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
