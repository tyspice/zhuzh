package chatgpt

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/tyspice/zhuzh/internal/config"
)

const baseURL string = "https://api.openai.com/v1/"

type apiEndpoints struct {
	responses string
}

var api = apiEndpoints{
	responses: baseURL + "responses",
}

type streamingRequest struct {
	Model  string `json:"model"`
	Input  string `json:"input"`
	Stream bool   `json:"stream"`
}

type Client struct {
	Res chan string
	Err chan error
}

func NewClient() *Client {
	return &Client{
		Res: make(chan string),
		Err: make(chan error),
	}
}

func (c *Client) Close() {
	defer func() {
		if r := recover(); r != nil {
			// Silently recover from "close of closed channel" panic
		}
	}()

	// Close both channels
	close(c.Res)
	close(c.Err)
}

func (c *Client) Ask(prompt string) {
	gptConfig := config.Get().ChatGPT

	go func() {

		// Prepare the request payload
		requestBody := streamingRequest{
			Model:  gptConfig.Model,
			Input:  prompt,
			Stream: true,
		}

		jsonData, err := json.Marshal(requestBody)
		if err != nil {
			c.Err <- fmt.Errorf("failed to marshal request: %w", err)
			return
		}

		// Create HTTP request
		req, err := http.NewRequest("POST", api.responses, bytes.NewBuffer(jsonData))
		if err != nil {
			c.Err <- fmt.Errorf("failed to create request: %w", err)
			return
		}

		// Set headers
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+gptConfig.ApiKey)

		// Make the request
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			c.Err <- fmt.Errorf("failed to make request: %w", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			c.Err <- fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
			return
		}

		nextEvent := ""

		// Read streaming response
		scanner := bufio.NewScanner(resp.Body)
		for scanner.Scan() {
			line := scanner.Text()

			// Skip empty lines and check for end of stream
			if line == "" {
				continue
			}

			// Remove "data: " prefix
			if strings.HasPrefix(line, "data: ") {
				line = strings.TrimPrefix(line, "data: ")
			}

			if strings.HasPrefix(line, "event: ") {
				line = strings.TrimPrefix(line, "event: ")
				nextEvent = line
				continue
			}

			if nextEvent != "" {
				if nextEvent == "response.output_text.delta" {
					var res struct {
						Delta string `json:"delta"`
					}
					json.Unmarshal([]byte(line), &res)
					c.Res <- res.Delta
				}
			}

		}

		if err := scanner.Err(); err != nil {
			c.Err <- fmt.Errorf("error reading response: %w", err)
		}
	}()
}
