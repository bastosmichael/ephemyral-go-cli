package gpt4client

import (
    "bytes"
    "crypto/tls"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
    "os"
    "time"

    "github.com/joho/godotenv"
)

const apiURL = "https://api.openai.com/v1/chat/completions"

var debug bool = false

func SetDebug(enabled bool) {
    debug = enabled
}

func DebugLog(format string, v ...interface{}) {
    if debug {
        fmt.Printf(format+"\n", v...)
    }
}

func GetGPT4ResponseWithPrompt(prompt string) (string, error) {
    if err := godotenv.Load(); err != nil {
        DebugLog("Error loading .env file: %v", err)
    }

    apiKey := os.Getenv("OPENAI_API_KEY")
    if apiKey == "" {
        return "", fmt.Errorf("API key not set")
    }

    messages := []map[string]interface{}{
        {"role": "system", "content": "You are writing software code."},
        {"role": "user", "content": prompt},
    }

    payload := map[string]interface{}{
        "model":    "gpt-4-turbo-preview",
        "messages": messages,
    }
    payloadBytes, err := json.Marshal(payload)
    if err != nil {
        return "", fmt.Errorf("error preparing request payload: %w", err)
    }

    DebugLog("Request payload: %s", string(payloadBytes))

    // Custom TLS configuration
    tlsConfig := &tls.Config{
        MinVersion:               tls.VersionTLS12,
        CipherSuites:             []uint16{tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384, tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256},
        PreferServerCipherSuites: true,
        InsecureSkipVerify:       false,
    }

    client := &http.Client{
        Transport: &http.Transport{
            TLSClientConfig: tlsConfig,
            Proxy:           http.ProxyFromEnvironment,
        },
        Timeout: 30 * time.Second,
    }

    req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(payloadBytes))
    if err != nil {
        return "", fmt.Errorf("error creating request: %w", err)
    }

    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "Bearer "+apiKey)

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
