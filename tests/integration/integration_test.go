package integration

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// Wait for all services to be ready before running tests
	if err := WaitForServices(); err != nil {
		os.Exit(1)
	}

	// Run tests
	code := m.Run()

	// Exit with test result code
	os.Exit(code)
}
