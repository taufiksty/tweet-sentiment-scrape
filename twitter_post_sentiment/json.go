package twitterpostsentiment

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/taufiksty/web-scraper/model"
)

func SaveToJSON(tweet *model.Tweet) {
	jsonData, err := json.Marshal(tweet)
	if err != nil {
		fmt.Println("Error marshalling JSON", err)
		return
	}

	filename := fmt.Sprintln("./result/tweet.json")
	file, err := os.Create(filename)
	if err != nil {
		fmt.Println("Error creating file", err)
		return
	}
	defer file.Close()

	_, err = file.Write(jsonData)
	if err != nil {
		fmt.Println("Error writing JSON to file", err)
		return
	}

	fmt.Println("JSON data has been written successfully")
}
