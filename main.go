package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/labstack/echo/v5"
)

func main() {
	e := echo.New()
	cond := sync.NewCond(new(sync.Mutex))

	// read message
	e.GET("/api/messages", func(c *echo.Context) error {
		offset, err := strconv.Atoi(c.QueryParam("offset"))
		if err != nil || offset > len(messages) {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid offset")
		}

		ctx, cancel := context.WithTimeout(c.Request().Context(), 10*time.Second)
		defer cancel()
		msg := Poll(ctx, cond, offset)

		return c.JSON(http.StatusOK, PollingResponse{
			Offset:   len(messages),
			Messages: msg,
		})
	})

	// send message
	e.POST("/api/messages", func(c *echo.Context) error {
		var reqBody MessageRequest
		if err := c.Bind(&reqBody); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
		}

		cond.L.Lock()
		messages = append(messages, reqBody.Message)
		cond.L.Unlock()
		cond.Broadcast()

		return c.NoContent(http.StatusCreated)
	})

	if err := e.Start("127.0.0.1:9000"); err != nil {
		slog.Error("listen server", "err", err)
		os.Exit(1)
	}
}

func Poll(ctx context.Context, c *sync.Cond, offset int) []string {
	c.L.Lock()
	defer c.L.Unlock()

	done := make(chan struct{})
	defer close(done)

	go func() {
		select {
		case <-ctx.Done():
			c.Broadcast()
		case <-done:
		}
	}()

	for len(messages[offset:]) == 0 {
		select {
		case <-ctx.Done():
			return nil
		default:
		}
		c.Wait()
	}

	return messages[offset:]
}

type H map[string]any

type PollingResponse struct {
	Offset   int      `json:"offset"`
	Messages []string `json:"messages"`
}

type MessageRequest struct {
	Message string `json:"message"`
}

// message buffer
var messages []string
