package oaievals_collector

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProcessFile(t *testing.T) {
	// Setup
	tmpDir := os.TempDir()

	dataDir = tmpDir
	processedDir = filepath.Join(tmpDir, "processed")
	err := os.MkdirAll(processedDir, 0o755) // Create the processed directory
	if err != nil {
		t.Fatal(err)
	}

	testFileContent := `{"run_id":"testRun", "type":"match", "data":{"correct":true}}`
	testFileName := "testFile"
	err = os.WriteFile(filepath.Join(dataDir, testFileName), []byte(testFileContent), 0o644)
	if err != nil {
		t.Fatal(err)
	}

	// Start a local HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Test request parameters
		assert.Equal(t, req.URL.String(), "/events")
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, string(body), "["+testFileContent+"]")

		// Send response to be tested
		rw.Write([]byte("OK"))
	}))

	// Close the server when test finishes
	defer server.Close()

	// Act
	fileInfo, err := os.Stat(filepath.Join(dataDir, testFileName))
	if err != nil {
		t.Fatal(err)
	}
	processFile(fileInfo)

	// Assert
	_, err = os.Stat(filepath.Join(dataDir, testFileName))
	assert.True(t, os.IsNotExist(err), "original file should be moved")

	_, err = os.Stat(filepath.Join(processedDir, testFileName))
	assert.False(t, os.IsNotExist(err), "processed file should exist")

	// More assertions can be added to check if the correct event processing happened.
}

func TestProcessFile_ErrorOpeningFile(t *testing.T) {
	// Setup
	tmpDir := os.TempDir()
	defer os.RemoveAll(tmpDir)

	dataDir = tmpDir
	processedDir = filepath.Join(tmpDir, "processed")

	// Act
	fileInfo, err := os.Stat(filepath.Join(dataDir, "nonexistentFile"))
	if err != nil && os.IsNotExist(err) {
		// This is the expected error
		err = nil // reset the error
	} else {
		t.Fatal("Expected error due to nonexistent file, got none")
	}

	if fileInfo != nil {
		processFile(fileInfo) // This should not panic
	}
	// No asserts here, because we only want to test that processFile does not panic
}
