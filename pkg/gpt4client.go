package gpt4client

import (
    "bytes"
    "crypto/tls"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
    // "os"
    "time"
    // "github.com/joho/godotenv"
    "github.com/fatih/color" // for colored output
    "sync" // for concurrency
)

const (
    apiURL   = "https://api.openai.com/v1/chat/completions"
    model    = "gpt-4-turbo-preview"
    roleSys  = "system"
    roleUser = "user"
    roleSysContent = "You are writing software code."
)

var (
    debug bool
    stopSpinner = make(chan bool)
    spinnerDone sync.WaitGroup
)

// SetDebug enables or disables debug output.
func SetDebug(enabled bool) {
    debug = enabled
}

// debugLog prints debug information if debug mode is enabled.
func debugLog(format string, v ...interface{}) {
    if debug {
        fmt.Printf(format+"\n", v...)
    }
}

// startSpinner starts a spinner in a separate goroutine.
func startSpinner() {
    spinnerDone.Add(1)
    go func() {
        defer spinnerDone.Done()
        spinnerChars := []string{"|", "/", "-", "\\"}
        color := color.New(color.FgCyan).SprintFunc() // Cyan spinner
        i := 0
        for {
            select {
            case <-stopSpinner:
                return
            default:
                fmt.Printf("\r%s", color(spinnerChars[i%len(spinnerChars)]))
                time.Sleep(100 * time.Millisecond)
                i++
            }
        }
    }()
}

// stopSpinner stops the spinner.
func stopSpinnerFunc() {
    stopSpinner <- true
    spinnerDone.Wait() // Wait for the spinner goroutine to finish
}

func getAPIKey() (string, error) {
    // if err := godotenv.Load(); err != nil {
    //     debugLog("Error loading .env file: %v", err)
    // }
    // apiKey := os.Getenv("OPENAI_API_KEY")
    // if apiKey == "" {
    //     return "", fmt.Errorf("API key not set")
    // }
    // return apiKey, nil
    return "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", nil
}

func preparePayload(prompt string) ([]byte, error) {
    messages := []map[string]interface{}{
        { "role": roleSys,  "content": roleSysContent },
        { "role": roleUser, "content": prompt },
    }

    payload := map[string]interface{}{
        "model":    model,
        "messages": messages,
    }

    return json.Marshal(payload)
}

func createHTTPClient() *http.Client {
    tlsConfig := &tls.Config{
        MinVersion:               tls.VersionTLS12,
        CipherSuites:             []uint16{tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384, tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256},
        PreferServerCipherSuites: true,
    }
    return &http.Client{
        Transport: &http.Transport{
            TLSClientConfig: tlsConfig,
        },
        Timeout: 30 * time.Second,
    }
}

func doPostRequest(client *http.Client, payloadBytes []byte, apiKey string) (*http.Response, error) {
    req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(payloadBytes))
    if err != nil {
        return nil, err
    }

    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "Bearer " + apiKey)

    return client.Do(req)
}

func GetGPT4ResponseWithPrompt(prompt string) (string, error) {
    apiKey, err := getAPIKey()
    if err != nil {
        return "", err
    }

    payloadBytes, err := preparePayload(prompt)
    if err != nil {
        return "", fmt.Errorf("error preparing request payload: %v", err)
    }

    debugLog("Request payload: %s", string(payloadBytes))

    client := createHTTPClient()

    startSpinner() // Start the spinner
    defer stopSpinnerFunc() // Ensure spinner stops after processing

    resp, err := doPostRequest(client, payloadBytes, apiKey)
    if err != nil {
        return "", fmt.Errorf("error sending request to the API: %v", err)
    }
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return "", fmt.Errorf("error reading response body: %v", err)
    }

    var responseMap map[string]interface{}
    if err := json.Unmarshal(body, &responseMap); err != nil {
        return "", fmt.Errorf("error parsing JSON response: %v", err)
    }

    return extractContentFromResponse(responseMap)
}

func extractContentFromResponse(responseMap map[string]interface{}) (string, error) {
    choices, ok := responseMap["choices"].([]interface{})
    if ok && len(choices) > 0 {
        firstChoice, ok := choices[0].(map[string]interface{})
        if ok {
            message, ok := firstChoice["message"].(map[string]interface{})
            if ok {
                content, ok := message["content"].(string)
                if ok {
                    return content, nil
                }
            }
        }
    }
    return "", fmt.Errorf("no valid content found")
}
