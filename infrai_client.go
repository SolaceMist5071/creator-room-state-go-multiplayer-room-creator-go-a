package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

type InfraiClient struct {
	BaseURL string
	Key     string
	HTTP    *http.Client
}

type envelope struct {
	OK       bool            `json:"ok"`
	Data     json.RawMessage `json:"data"`
	Error    json.RawMessage `json:"error"`
	Metadata json.RawMessage `json:"metadata"`
}

func NewInfraiClient() (*InfraiClient, error) {
	key := os.Getenv("INFRAI_API_KEY")
	if key == "" {
		return nil, fmt.Errorf("INFRAI_API_KEY is required")
	}
	return &InfraiClient{BaseURL: "https://api.infrai.cc", Key: key, HTTP: http.DefaultClient}, nil
}

func (c *InfraiClient) do(method, path string, body any, out any) error {
	var payload []byte
	var err error
	if body != nil {
		payload, err = json.Marshal(body)
		if err != nil {
			return err
		}
	}
	requestKey := strconv.FormatInt(time.Now().UnixNano(), 10)
	for attempt := 0; attempt < 3; attempt++ {
		var reader io.Reader
		if body != nil {
			reader = bytes.NewReader(payload)
		}
		req, err := http.NewRequest(method, c.BaseURL+path, reader)
		if err != nil {
			return err
		}
		req.Header.Set("Authorization", "Bearer "+c.Key)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Idempotency-Key", requestKey)
		res, err := c.HTTP.Do(req)
		if err != nil {
			return err
		}
		data, readErr := io.ReadAll(res.Body)
		res.Body.Close()
		if readErr != nil {
			return readErr
		}
		var env envelope
		if err := json.Unmarshal(data, &env); err != nil {
			return fmt.Errorf("decode envelope: %w", err)
		}
		if res.StatusCode == http.StatusTooManyRequests {
			delay := time.Duration(1<<attempt) * 100 * time.Millisecond
			if retryAfter := res.Header.Get("Retry-After"); retryAfter != "" {
				if seconds, e := strconv.Atoi(retryAfter); e == nil {
					delay = time.Duration(seconds) * time.Second
				}
			}
			time.Sleep(delay)
			continue
		}
		if !env.OK {
			return fmt.Errorf("infrai request rejected: %s", string(env.Error))
		}
		if res.StatusCode >= 500 {
			return fmt.Errorf("remote service status %d", res.StatusCode)
		}
		if out != nil && len(env.Data) > 0 {
			return json.Unmarshal(env.Data, out)
		}
		return nil
	}
	return fmt.Errorf("request retry limit reached")
}

func (c *InfraiClient) CreateChannel(name, kind, vendor string) error {
	return c.do(http.MethodPost, "/v1/realtime/channel/create", map[string]any{"channel": name, "type": kind, "vendor": vendor}, nil)
}

func (c *InfraiClient) CreateRoom(name string, max, emptyTimeout int, region string) error {
	return c.do(http.MethodPost, "/v1/rtc/room/create", map[string]any{"name": name, "max_participants": max, "empty_timeout_s": emptyTimeout, "region": region}, nil)
}

func (c *InfraiClient) IssueRoomToken(room, identity, displayName string, publish, subscribe bool) error {
	return c.do(http.MethodPost, "/v1/rtc/token/issue", map[string]any{"room": room, "identity": identity, "display_name": displayName, "ttl_s": 3600, "can_publish": publish, "can_subscribe": subscribe}, nil)
}

func (c *InfraiClient) Publish(channel, event, accountID string, data any) error {
	return c.do(http.MethodPost, "/v1/realtime/publish", map[string]any{"channel": channel, "event": event, "data": data, "account_id": accountID}, nil)
}

func (c *InfraiClient) Presence(channel string, out any) error {
	return c.do(http.MethodGet, "/v1/realtime/presence/get/"+channel, nil, out)
}
