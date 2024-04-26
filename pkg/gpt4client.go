package gpt4client

import (
    "bytes"
    "crypto/tls"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
    "time"
    "github.com/fatih/color"
    "sync"
)

const (
    apiURL   = "https://api.openai.com/v1/chat/completions"
    model    = "gpt-4-turbo"
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
        color := color.New(color.FgCyan).SprintFunc()
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
    spinnerDone.Wait()
}

func getAPIKey() (string, error) {
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
        return "", err
    }

    debugLog("Request payload: %s", string(payloadBytes))

    client := createHTTPClient()

    startSpinner() // Start the spinner
    defer stopSpinnerFunc() // Ensure spinner stops after processing

    resp, err := doPostRequest(client, payloadBytes, apiKey)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return "", err
    }

    var responseMap map[string]interface{}
    if err := json.Unmarshal(body, &responseMap); err != nil {
        return "", err
    }

    return extractContentFromResponse(responseMap)
}

func extractContentFromResponse(responseMap map[string]interface{}) (string, error) {
    if content, ok := responseMap["error"].(string); ok {
        return "", fmt.Errorf(content)
    }

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

    return "", fmt.Errorf("Your current trial usage of Ephemyral may have either expired or been revoked, please reach out for an updated version and or contact us.")
}
