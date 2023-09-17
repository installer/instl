package metrics

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2/log"
	"io"
	"net/http"
	"os"
)

type Metric struct {
	UserAgent    string            `json:"-"`
	ForwardedFor string            `json:"-"`
	EventName    string            `json:"name"`
	URL          string            `json:"url"`
	Props        map[string]string `json:"props"`
	Domain       string            `json:"domain"`
}

// Send sends a metric to a plausible instance
func Send(m Metric) error {
	endpoint := os.Getenv("PLAUSIBLE_ENDPOINT")
	domain := os.Getenv("PLAUSIBLE_DOMAIN")
	if endpoint == "" {
		return nil
	}

	if domain == "" {
		return nil
	} else {
		m.Domain = domain
	}

	endpoint += "/api/event"

	header := map[string]string{
		"Content-Type":    "application/json",
		"User-Agent":      m.UserAgent,
		"X-Forwarded-For": m.ForwardedFor,
	}

	if header["X-Forwarded-For"] == "" {
		header["X-Forwarded-For"] = "127.0.0.1"
	}

	// Post to endpoint
	client := &http.Client{}
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(marshal(m)))
	if err != nil {
		return err
	}

	for key, value := range header {
		req.Header.Set(key, value)
	}

	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return err
	}

	response, _ := io.ReadAll(resp.Body)
	log.Debug(fmt.Sprintf("Sent metric to %s. Status: %s, Response: %s", endpoint, resp.Status, string(response)))

	return nil
}

func marshal(m Metric) []byte {
	j, _ := json.Marshal(m)
	return j
}
