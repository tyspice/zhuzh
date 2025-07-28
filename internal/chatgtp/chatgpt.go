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

type StreamingRequest struct {
	Model  string `json:"model"`
	Input  string `json:"input"`
	Stream bool   `json:"stream"`
}

func StreamResponse(prompt string) (<-chan string, <-chan error) {
	gptConfig := config.Get().ChatGPT

	responseChan := make(chan string)
	errorChan := make(chan error, 1)

	go func() {
		defer close(responseChan)
		defer close(errorChan)

		// Prepare the request payload
		requestBody := StreamingRequest{
			Model:  gptConfig.Model,
			Input:  prompt,
			Stream: true,
		}

		jsonData, err := json.Marshal(requestBody)
		if err != nil {
			errorChan <- fmt.Errorf("failed to marshal request: %w", err)
			return
		}

		// Create HTTP request
		req, err := http.NewRequest("POST", api.responses, bytes.NewBuffer(jsonData))
		if err != nil {
			errorChan <- fmt.Errorf("failed to create request: %w", err)
			return
		}

		// Set headers
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+gptConfig.ApiKey)

		// Make the request
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			errorChan <- fmt.Errorf("failed to make request: %w", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			errorChan <- fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
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

			// Check for end of stream
			if line == "[DONE]" {
				break
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
					responseChan <- res.Delta
				}
			}

		}

		if err := scanner.Err(); err != nil {
			errorChan <- fmt.Errorf("error reading response: %w", err)
		}
	}()

	return responseChan, errorChan
}
