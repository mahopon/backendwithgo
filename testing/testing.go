package testing

import (
	"os"
	"testing"
)

func CheckEnvironment(t *testing.T, env string) {
	switch env {
	case "UNIT":
		{
			if os.Getenv("UNIT") == "" {
				t.Skip("Skipping unit tests")
			}
		}
	case "INTEGRATION":
		{
			if os.Getenv("INTEGRATION") == "" {
				t.Skip("Skipping integration tests")
			}
		}
	}
}
