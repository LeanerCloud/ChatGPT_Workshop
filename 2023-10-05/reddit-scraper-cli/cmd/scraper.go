package cmd

import (
	"fmt"
	"strings"

	"github.com/turnage/graw"
	"github.com/turnage/graw/reddit"
)

type announcer struct {
	companyName string
}

func (a *announcer) Post(post *reddit.Post) error {
	if strings.Contains(strings.ToLower(post.Title), strings.ToLower(a.companyName)) {
		fmt.Printf("Found mention of %s in post title: %s by %s\n", a.companyName, post.Title, post.Author)
	}
	return nil
}

func streamRedditWithBot(bot reddit.Bot) error {
	a := &announcer{companyName: companyName}

	cfg := graw.Config{Subreddits: []string{"all"}} // Listen to all subreddits

	_, wait, err := graw.Scan(a, bot, cfg)
	if err != nil {
		return fmt.Errorf("failed to start graw scan: %v", err)
	}

	// This time, let's block so the bot will announce (ideally) forever.
	if err := wait(); err != nil {
		return fmt.Errorf("graw run encountered an error: %v", err)
	}

	return nil
}
