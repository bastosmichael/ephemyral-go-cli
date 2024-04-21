package gpt4client

import (
    "testing"
)

// TestSetDebug tests the SetDebug function.
func TestSetDebug(t *testing.T) {
    SetDebug(true)
    if !debug {
        t.Error("SetDebug(true) failed, expected debug to be true")
    }
    SetDebug(false)
    if debug {
        t.Error("SetDebug(false) failed, expected debug to be false")
    }
}

// TestGetGPT4ResponseWithPrompt tests the GetGPT4ResponseWithPrompt function.
func TestGetGPT4ResponseWithPrompt(t *testing.T) {
    // Here you would normally mock the http.Client or use an interface to control the output
    // For now, let's assume the environment is not set up, and an error should occur because of missing API key
    response, err := GetGPT4ResponseWithPrompt("Hello")
    if err == nil {
        t.Error("Expected error due to missing API key, got nil")
    }
    if response != "" {
        t.Errorf("Expected no response due to error, got: %s", response)
    }
}

// You might use a mocking library like gomock or httptest to simulate API responses for more comprehensive testing.
