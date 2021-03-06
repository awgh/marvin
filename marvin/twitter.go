package marvin

import (
	"regexp"
	"strings"

	"github.com/gocolly/colly"
)

var rTweet = regexp.MustCompile(`twitter.com/.*/status/\d+`)

// CheckForTweet - checks a text blob for links to tweets
func CheckForTweet(text string) (bool, string) {
	tweet := rTweet.FindString(text)
	return tweet != "", tweet
}

// GetTweetFromURL - fetches a tweet
func GetTweetFromURL(url string) string {
	c := colly.NewCollector(
		colly.AllowedDomains("mobile.twitter.com"),
		colly.UserAgent("Mozilla/5.0 (iPhone; CPU iPhone OS 6_0 like Mac OS X) AppleWebKit/536.26 (KHTML, like Gecko) Version/6.0 Mobile/10A5376e Safari/8536.25"),
	)

	idx := strings.LastIndexByte(url, '/')
	dataID := url[idx+1:]

	text := ""
	c.OnHTML("div.tweet-text", func(e *colly.HTMLElement) {
		if text != "" || e.Attr("data-id") != dataID {
			return
		}
		text = strings.Trim(string(e.Text), " \n\t\v")
		href := e.ChildAttr("a[href]", "href")
		if href != "" {
			text = text + " " + href
		}
	})
	c.Visit(url)

	return text
}

// GetTweetFromText - extracts tweet links from text and fetches the tweet
func GetTweetFromText(text string) string {
	if ok, s := CheckForTweet(text); ok {
		tweet := GetTweetFromURL("https://mobile." + s)
		return tweet
	}
	return ""
}
