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

const (
	baseURL             string = "https://api.openai.com/v1/"
	defaultInstructions string = "Format your response using Markdown syntax."
)

type apiEndpoints struct {
	responses string
}

var api = apiEndpoints{
	responses: baseURL + "responses",
}

type streamingRequest struct {
	PreviousResponseId string `json:"previous_response_id,omitempty"`
	Model              string `json:"model"`
	Input              string `json:"input"`
	Stream             bool   `json:"stream"`
	Instructions       string `json:"instructions"`
}

type Client struct {
	previousResponseId string
	res                chan string
	err                chan error
}

func NewClient() *Client {
	return &Client{
		res: make(chan string),
		err: make(chan error),
	}
}

func (c *Client) Close() {
	defer func() {
		if r := recover(); r != nil {
			// Silently recover from "close of closed channel" panic
		}
	}()

	// Close both channels
	close(c.res)
	close(c.err)
}

func (c *Client) Subscribe() (res <-chan string, err <-chan error) {
	return c.res, c.err
}

func (c *Client) Ask(prompt string) {
	gptConfig := config.Get().ChatGPT

	go func() {

		// Prepare the request payload
		requestBody := streamingRequest{
			Model:        gptConfig.Model,
			Input:        prompt,
			Stream:       true,
			Instructions: defaultInstructions,
		}

		if c.previousResponseId != "" {
			requestBody.PreviousResponseId = c.previousResponseId
		}

		jsonData, err := json.Marshal(requestBody)
		if err != nil {
			c.err <- fmt.Errorf("failed to marshal request: %w", err)
			return
		}

		// Create HTTP request
		req, err := http.NewRequest("POST", api.responses, bytes.NewBuffer(jsonData))
		if err != nil {
			c.err <- fmt.Errorf("failed to create request: %w", err)
			return
		}

		// Set headers
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+gptConfig.ApiKey)

		// Make the request
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			c.err <- fmt.Errorf("failed to make request: %w", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			c.err <- fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
			return
		}

		nextEvent := ""

		// Read streaming response
		scanner := bufio.NewScanner(resp.Body)
		for scanner.Scan() {
			line := scanner.Text()

			if line == "" {
				continue
			}

			line = strings.TrimPrefix(line, "data: ")

			if after, ok := strings.CutPrefix(line, "event: "); ok {
				nextEvent = after
				continue
			}

			if nextEvent != "" {

				switch nextEvent {
				case "response.created":
					var res struct {
						Response struct {
							ID string `json:"id"`
						} `json:"response"`
					}
					json.Unmarshal([]byte(line), &res)
					c.previousResponseId = res.Response.ID
				case "response.output_text.delta":
					var res struct {
						Delta string `json:"delta"`
					}
					json.Unmarshal([]byte(line), &res)
					c.res <- res.Delta
				}
			}

		}

		if err := scanner.Err(); err != nil {
			c.err <- fmt.Errorf("error reading response: %w", err)
		}
	}()
}
