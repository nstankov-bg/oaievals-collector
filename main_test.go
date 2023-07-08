package main

import (
	"testing"
)

func TestMainFunction(t *testing.T) {
    go main() // Run the main function in a goroutine so it doesn't block
    // Insert other testing logic here
}
