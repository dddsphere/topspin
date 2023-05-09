package topspin_test

import (
	"os"
	"testing"
	"time"

	"github.com/fsnotify/fsnotify"

	"github.com/dddsphere/topspin"
)

func TestLoad(t *testing.T) {
	// Create a temporary directory and file for testing
	tmpFile, err := os.CreateTemp("", "test_load_*.yml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	tmpFile, err = os.OpenFile(tmpFile.Name(), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		t.Fatal(err)
	}

	// Update the configuration file
	_, err = tmpFile.WriteString("key1: 123\n")
	if err != nil {
		t.Fatal(err)
	}

	err = tmpFile.Close()
	if err != nil {
		t.Fatal(err)
	}

	// Load the configuration
	cfg := topspin.NewConfig("test")
	_, err = cfg.Load(tmpFile.Name())
	if err != nil {
		t.Fatal(err)
	}

	// Check that the configuration was loaded correctly
	val := cfg.GetString("key1")
	if val != "123" {
		t.Errorf("Expected test value to be 123, got %s", val)
	}
}

func TestOnConfigChange(t *testing.T) {
	cfg := topspin.NewConfig("test")
	tmpFile, err := os.CreateTemp("", "test_on_config_change_*.yml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	tmpFile, err = os.OpenFile(tmpFile.Name(), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		t.Fatal(err)
	}

	// Update the configuration file
	_, err = tmpFile.WriteString("key1: 123\n")
	if err != nil {
		t.Fatal(err)
	}

	err = tmpFile.Close()
	if err != nil {
		t.Fatal(err)
	}

	// Initial load
	_, err = cfg.Load(tmpFile.Name())
	if err != nil {
		t.Fatal(err)
	}

	// Watch for changes
	onChangeCalled := false
	cfg.SetOnConfigChange(func(e fsnotify.Event) {
		onChangeCalled = true
	})

	tmpFile, err = os.OpenFile(tmpFile.Name(), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		t.Fatal(err)
	}

	// Update the configuration file
	_, err = tmpFile.WriteString("key1: 321\n")
	if err != nil {
		t.Fatal(err)
	}

	err = tmpFile.Close()
	if err != nil {
		t.Fatal(err)
	}

	// It takes a bit to detect the changes on config
	time.Sleep(2 * time.Second)

	// Check if the onChange function was called
	if !onChangeCalled {
		t.Error("onChange function was not called")
	}

	// Check if the updated value was loaded
	val := cfg.GetInt("key1")
	if val != 321 {
		t.Errorf("unexpected value for key1: %d", val)
	}
}
