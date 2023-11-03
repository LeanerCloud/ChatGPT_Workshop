package cmd

import (
	"errors"
	"testing"
	"time"

	"github.com/turnage/graw/reddit"
)

func TestAnnouncerPost(t *testing.T) {
	tests := []struct {
		name        string
		companyName string
		postTitle   string
		expected    bool
	}{
		{
			name:        "post contains company name",
			companyName: "testCompany",
			postTitle:   "This is a post about testCompany",
			expected:    true,
		},
		{
			name:        "post does not contain company name",
			companyName: "testCompany",
			postTitle:   "This is a post about anotherCompany",
			expected:    false,
		},
		{
			name:        "case insensitive match",
			companyName: "testCompany",
			postTitle:   "This is a post about TESTCOMPANY",
			expected:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &announcer{companyName: tt.companyName}
			post := &reddit.Post{
				Title:  tt.postTitle,
				Author: "testUser",
			}

			err := a.Post(post)
			if tt.expected && err != nil {
				t.Fatalf("Expected no error, got %v", err)
			} else if !tt.expected && err == nil {
				t.Fatalf("Expected an error but got none")
			}
		})
	}
}

func TestStreamReddit(t *testing.T) {
	companyName = "testCompany"

	mockBot := &MockRedditBot{
		PostChannel: make(chan *reddit.Post, 1),
	}

	// Simulate a post being streamed from Reddit
	go func() {
		post := &reddit.Post{
			Title:  "This is a post about testCompany",
			Author: "testUser",
		}
		mockBot.PostChannel <- post
		time.Sleep(1 * time.Second) // Give some time for processing
		close(mockBot.PostChannel)  // Close the channel to simulate end of stream
	}()

	err := streamRedditWithBot(mockBot)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Test error scenario
	mockErrorBot := &MockRedditBot{
		PostChannel: nil, // Simulate an error scenario
		Error:       errors.New("simulated error"),
	}

	err = streamRedditWithBot(mockErrorBot)
	if err == nil {
		t.Fatalf("Expected an error but got none")
	}
}

func TestStreamReddit_ErrorScenario(t *testing.T) {
	companyName = "testCompany"

	mockErrorBot := &MockRedditBot{
		PostChannel: nil, // Simulate an error scenario
		Error:       errors.New("simulated error"),
	}

	// Inject the mockErrorBot into streamReddit (you might need to modify the original code to support this)
	// For the sake of this example, let's assume streamReddit can accept a bot interface.

	err := streamRedditWithBot(mockErrorBot)
	if err == nil {
		t.Fatalf("Expected an error but got none")
	}
}
