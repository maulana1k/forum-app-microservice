package tests

import (
	"os"
	"testing"

	"github.com/maulana1k/forum-app/tests/helper"
)

// TestMain will set up and tear down the test database for all tests in this package
func TestMain(m *testing.M) {
	helper.Setup()
	code := m.Run()
	// m.Run()
	helper.Teardown()
	os.Exit(code)
}
