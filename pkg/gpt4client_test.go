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
