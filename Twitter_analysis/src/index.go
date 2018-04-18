package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"regexp"
	"sort"
	"strings"
)

var ss []kv

type kv struct {
	Key   string
	Value int
}

func main() {

	//We load the csv file with the data
	file, err := os.Open("../data/tweets.csv")
	var validID = regexp.MustCompile(`@([0-9A-Za-z_]+)[== \t]`)
	var validInfo = regexp.MustCompile(`>([0-9A-Za-z_\s"<>:/.=]+)</a>`)
	var validRT = regexp.MustCompile(`RT[== \t]@([0-9A-Za-z_:]+)[== \t]`)

	if err != nil {
		fmt.Println("Error: ", err)

		return
	}

	defer file.Close()

	reader := csv.NewReader(file)

	//CSV fields will be sepparated by commas
	reader.Comma = ','

	line := ""
	lineRT := ""
	numRT := 0
	numTweets := 0
	numInteractedTweets := 0
	numAndroid := 0
	numWeb := 0
	numTweetDeck := 0
	numAnotherApps := 0

	for {
		record, err := reader.Read()

		numTweets++ //Every different line is a tweet

		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println("Error:", err)
			return
		}

		// We go through all the csv lines
		for i := 0; i < len(record); i++ {

			line += strings.ToLower(validID.FindString(record[i]))      //Interactions
			lineRT += strings.Trim(validRT.FindString(record[i]), "RT") //Retweets

			//Number of RT
			if strings.HasPrefix(record[i], "RT") {
				numRT++
			}
			//Number of tweets with interactions
			if strings.HasPrefix(record[i], "@") {
				numInteractedTweets++
			}

			if validInfo.MatchString(record[i]) {
				if strings.Contains(record[i], "Twitter Web Client") {
					numWeb++
				} else if strings.Contains(record[i], "Twitter for Android") {
					numAndroid++
				} else if strings.Contains(record[i], "TweetDeck") {
					numTweetDeck++
				} else {
					numAnotherApps++
				}
			}
		}
	}

	name := topUser(line)

	percentageRT := float64(numRT) / float64(numTweets) * 100
	percentageInter := float64(numInteractedTweets) / float64(numTweets) * 100
	ownTweets := numTweets - (numRT + numInteractedTweets)
	percentageTweets := float64(ownTweets) / float64(numTweets) * 100

	fmt.Println()

	fmt.Println("---------------------DATA EXTRACTED---------------------------------")

	fmt.Println()

	fmt.Println("Number of tweets: ", numTweets)
	fmt.Println("Most interacted user: ", name)
	fmt.Println()
	fmt.Println("Top 5 most directly interacted users: ")
	showTopFive()
	fmt.Println()
	name = topUser(lineRT)
	fmt.Println()
	fmt.Println("Most retweeted user: ", name)
	fmt.Println("Top 5 most Retweeted users: ")
	showTopFive()
	fmt.Println()
	fmt.Printf("Number of RT: %v (%.2f%% of tweets)\n", numRT, percentageRT)
	fmt.Printf("Number of tweets with interactions: %v (%.2f%% of tweets)\n", numInteractedTweets, percentageInter)
	fmt.Printf("Your own tweets: %v (%.2f%% of tweets)\n", ownTweets, percentageTweets)
	fmt.Printf("Tweets from Android application: %v\n", numAndroid)
	fmt.Printf("Tweets from Twitter web: %v\n", numWeb)
	fmt.Printf("Tweets from TweetDeck: %v\n", numTweetDeck)
	fmt.Printf("Tweets from other Apps: %v\n", numAnotherApps)

	fmt.Println()

	fmt.Println("---------------------------------------------------------------------")

	fmt.Println()

}

func topUser(line string) string {
	max := 0
	name := ""

	for index, element := range userCount(line) {
		if element > max {
			max = element
			name = index
		}
	}

	return name
}

func showTopFive() {

	m := 0

	for m = 0; m < 5; m++ {
		fmt.Println(ss[m])
	}
}

func userCount(str string) map[string]int {
	wordList := strings.Fields(str)
	counts := make(map[string]int)

	for _, word := range wordList {
		_, ok := counts[word]

		if ok {
			counts[word]++
		} else {
			counts[word] = 1
		}
	}

	ss = nil

	for k, v := range counts {
		ss = append(ss, kv{k, v})
	}

	sort.Slice(ss, func(i, j int) bool {
		return ss[i].Value > ss[j].Value
	})

	return counts
}
