// gpt4client.go
package gpt4client

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
    "os"

    "github.com/joho/godotenv"
)

// OpenAI Chat Completions URL for GPT-4
const apiURL = "https://api.openai.com/v1/chat/completions"

var debug bool = false

// SetDebug enables or disables debug mode
func SetDebug(enabled bool) {
    debug = enabled
}

// DebugLog prints debug information if debug mode is enabled
func DebugLog(format string, v ...interface{}) {
    if debug {
        fmt.Printf(format+"\n", v...)
    }
}

// GetGPT4ResponseWithPrompt sends a single prompt to GPT-4 and returns the response.
func GetGPT4ResponseWithPrompt(prompt string) (string, error) {
    // Load the .env file
    if err := godotenv.Load(); err != nil {
        DebugLog("Error loading .env file: %v", err)
    }

    apiKey := os.Getenv("OPENAI_API_KEY")
    if apiKey == "" {
        return "", fmt.Errorf("API key not set")
    }

    // Constructing messages from the provided prompt
    messages := []map[string]interface{}{
        {
            "role":    "system",
            "content": "You are a helpful assistant.",
        },
        {
            "role":    "user",
            "content": prompt,
        },
    }

    payload := map[string]interface{}{
        "model":    "gpt-4-turbo-preview", // Use gpt-4-turbo-preview or other appropriate GPT-4 model
        "messages": messages,
    }
    payloadBytes, err := json.Marshal(payload)
    if err != nil {
        return "", fmt.Errorf("error preparing request payload: %w", err)
    }

    DebugLog("Request payload: %s", string(payloadBytes))

    req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(payloadBytes))
    if err != nil {
        return "", fmt.Errorf("error creating request: %w", err)
    }

    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "Bearer "+apiKey)

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return "", fmt.Errorf("error sending request to the API: %w", err)
    }
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return "", fmt.Errorf("error reading response body: %w", err)
    }

    var responseMap map[string]interface{}
    if err := json.Unmarshal(body, &responseMap); err != nil {
        return "", fmt.Errorf("error parsing JSON response: %w", err)
    }

    if choices, ok := responseMap["choices"].([]interface{}); ok && len(choices) > 0 {
        if firstChoice, ok := choices[0].(map[string]interface{}); ok {
            if message, ok := firstChoice["message"].(map[string]interface{}); ok {
                if content, ok := message["content"].(string); ok {
                    return content, nil
                }
            }
        }
    }

    return "", fmt.Errorf("no suggestion found")
}
