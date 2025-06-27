package cmd

import (
	"testing"
)

const (
	testVersion = "1.0.0-test"
	testCommit  = "abc123"
	testDate    = "2025-01-01"
)

func TestRootCmd(t *testing.T) {
	t.Run("root command is configured", func(t *testing.T) {
		if rootCmd == nil {
			t.Fatal("rootCmd should not be nil")
		}

		if rootCmd.Use != "gos" {
			t.Errorf("expected Use to be 'gos', got %q", rootCmd.Use)
		}

		if rootCmd.Short == "" {
			t.Error("Short description should not be empty")
		}
	})
}

func TestSetVersionInfo(t *testing.T) {
	t.Run("version info can be set", func(t *testing.T) {
		SetVersionInfo(testVersion, testCommit, testDate)

		if version != testVersion {
			t.Errorf("expected version %q, got %q", testVersion, version)
		}

		if commit != testCommit {
			t.Errorf("expected commit %q, got %q", testCommit, commit)
		}

		if date != testDate {
			t.Errorf("expected date %q, got %q", testDate, date)
		}
	})
}

func TestGetVersionString(t *testing.T) {
	t.Run("version string is formatted correctly", func(t *testing.T) {
		SetVersionInfo(testVersion, testCommit, testDate)

		versionStr := getVersionString()
		if versionStr == "" {
			t.Error("version string should not be empty")
		}

		// Should contain all the version info
		if !contains(versionStr, testVersion) {
			t.Error("version string should contain version")
		}

		if !contains(versionStr, testCommit) {
			t.Error("version string should contain commit")
		}

		if !contains(versionStr, testDate) {
			t.Error("version string should contain date")
		}
	})
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
