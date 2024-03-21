// gpt3client.go
package gpt3client

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
    "os"
)

// OpenAI URL
const apiURL = "https://api.openai.com/v1/completions"

var debug bool = false // Debug flag default to false

// SetDebug enables or disables debug mode
func SetDebug(enabled bool) {
    debug = enabled
}

// DebugLog prints debug information if debug mode is enabled
func DebugLog(format string, v ...interface{}) {
    if debug {
        fmt.Printf(format, v...)
    }
}

// GetLLMSuggestion sends a prompt to GPT-4 and returns the suggestion.
func GetLLMSuggestion(prompt string) (string, error) {
    apiKey := os.Getenv("OPENAI_API_KEY")
    if apiKey == "" {
        if debug {
            fmt.Println("API key not set. Please set the OPENAI_API_KEY environment variable.")
        }
        return "", fmt.Errorf("API key not set")
    }
    
    payload := map[string]interface{}{
        "model":       "gpt-3.5-turbo-instruct", // Replace with actual GPT-4 model name when available
        "prompt":      prompt,
        "temperature": 0.7,
        "max_tokens":  2048,
        "top_p":       1,
        "frequency_penalty": 0,
        "presence_penalty": 0,
    }
    payloadBytes, err := json.Marshal(payload)
    if err != nil {
        DebugLog("Error preparing request payload: %v\n", err)
        return "", fmt.Errorf("error preparing request payload: %w", err)
    }

    DebugLog("Request payload: %s\n", string(payloadBytes))

    req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(payloadBytes))
    if err != nil {
        DebugLog("Error creating request: %v\n", err)
        return "", fmt.Errorf("error creating request: %w", err)
    }

    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "Bearer "+apiKey)

    DebugLog("Sending request to API URL: %s\n", apiURL)

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        DebugLog("Error sending request to the API: %v\n", err)
        return "", fmt.Errorf("error sending request to the API: %w", err)
    }
    defer resp.Body.Close()

    DebugLog("Request sent, reading response body...\n")

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        DebugLog("Error reading response body: %v\n", err)
        return "", fmt.Errorf("error reading response body: %w", err)
    }

    DebugLog("Response body: %s\n", string(body))

    var responseMap map[string]interface{}
    if err := json.Unmarshal(body, &responseMap); err != nil {
        DebugLog("Error parsing JSON response: %v\n", err)
        return "", fmt.Errorf("error parsing JSON response: %w", err)
    }

    if choices, ok := responseMap["choices"].([]interface{}); ok && len(choices) > 0 {
        if firstChoice, ok := choices[0].(map[string]interface{}); ok {
            if text, ok := firstChoice["text"].(string); ok {
                return text, nil
            }
        }
    }

    DebugLog("No suggestion found in the response.\n")
    return "", fmt.Errorf("no suggestion found")
}
