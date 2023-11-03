package cmd

import (
	"github.com/turnage/graw/reddit"
)

type MockRedditBot struct {
	PostChannel chan *reddit.Post
	Error       error
}

func (m *MockRedditBot) Stream() (postChan <-chan *reddit.Post, commentChan <-chan *reddit.Comment, messageChan <-chan *reddit.Message, err error) {
	if m.Error != nil {
		return nil, nil, nil, m.Error
	}
	return m.PostChannel, nil, nil, nil
}
