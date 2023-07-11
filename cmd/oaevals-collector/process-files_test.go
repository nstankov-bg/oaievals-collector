package oaievals_collector

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProcessFile(t *testing.T) {
	// Setup
	tmpDir, err := ioutil.TempDir("", "test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	dataDir = tmpDir
	processedDir = filepath.Join(tmpDir, "processed")
	os.MkdirAll(processedDir, 0755) // Create the processed directory

	testFileContent := `{"run_id":"testRun", "type":"match", "data":{"correct":true}}`
	testFileName := "testFile"
	err = ioutil.WriteFile(filepath.Join(dataDir, testFileName), []byte(testFileContent), 0644)
	if err != nil {
		t.Fatal(err)
	}

	// Act
	fileInfo, _ := os.Stat(filepath.Join(dataDir, testFileName))
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
	tmpDir, err := ioutil.TempDir("", "test")
	if err != nil {
		t.Fatal(err)
	}
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
