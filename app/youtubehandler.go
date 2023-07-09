package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

type YoutubeStats struct {
	Subscribers int    `json:"subscribers"`
	ChannelName string `json:"channelName"`
	// MinutesWatched int    `json:"minutesWatched"`
	Views int `json:"views"`
}

func getChannelStats(k string) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		// w.Write([]byte("Response..."))

		ctx := context.Background()
		youtubeService, err := youtube.NewService(ctx, option.WithAPIKey(k))
		if err != nil {
			fmt.Println("Failed to create service")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		listCall := youtubeService.Channels.List([]string{"snippet", "contentDetails", "statistics"})
		response, err := listCall.Id("UCPrJmZJ4uwx6XaEZ1TF-Ojw").Do()
		if err != nil {
			fmt.Println("Failed to request channel data")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		yt := YoutubeStats{
			Subscribers: int(response.Items[0].Statistics.SubscriberCount),
			ChannelName: response.Items[0].Snippet.Title,
			Views:       int(response.Items[0].Statistics.ViewCount),
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(w).Encode(yt); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
}
