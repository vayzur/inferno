package httputil

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	client *http.Client
}

func New(timeout time.Duration) *Client {
	return &Client{
		client: &http.Client{Timeout: timeout},
	}
}

func buildAuthHeader(token string) string {
	var b []byte
	ts := time.Now().Unix()
	b = fmt.Appendf(b, "%d:%s", ts, token)
	hash := sha256.Sum256(b)
	sig := hex.EncodeToString(hash[:])
	return fmt.Sprintf("rolling %d:%s", ts, sig)
}

func (c *Client) Do(method, url, token string, body any) (int, []byte, error) {
	var requestBody []byte
	var err error

	if body != nil {
		requestBody, err = json.Marshal(body)
		if err != nil {
			return 0, nil, fmt.Errorf("marshal error: %w", err)
		}
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(requestBody))
	if err != nil {
		return 0, nil, fmt.Errorf("request creation error: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", buildAuthHeader(token))

	resp, err := c.client.Do(req)
	if err != nil {
		return 0, nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, nil, fmt.Errorf("read error: %w", err)
	}

	return resp.StatusCode, data, nil
}
