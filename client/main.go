package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
)

func main() {
	var offset int
	for {
		slog.Info("read next message", "offset", offset)
		offset = handler(offset)
	}
}

func handler(offset int) int {
	url := fmt.Sprintf("http://127.0.0.1:9000/api/messages?offset=%d", offset)
	resp, err := http.Get(url)
	if err != nil {
		slog.Error("do request", "err", err)
		return offset
	}
	defer resp.Body.Close()

	var respBody Response
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		slog.Error("parse response body", "err", err)
		return offset
	}

	for id, msg := range respBody.Messages {
		slog.Info("received message",
			"id", id,
			"message", msg,
		)
	}

	return respBody.Offset
}

type Response struct {
	Offset   int      `json:"offset"`
	Messages []string `json:"messages"`
}
