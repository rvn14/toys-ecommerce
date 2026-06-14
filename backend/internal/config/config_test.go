package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadDotEnvAllowsMissingFile(t *testing.T) {
	withWorkingDirectory(t, t.TempDir())

	if err := loadDotEnv(); err != nil {
		t.Fatalf("loadDotEnv() returned an error for a missing file: %v", err)
	}
}

func TestLoadDotEnvReturnsParseErrors(t *testing.T) {
	dir := t.TempDir()
	withWorkingDirectory(t, dir)

	if err := os.WriteFile(filepath.Join(dir, ".env"), []byte("INVALID LINE\n"), 0o600); err != nil {
		t.Fatalf("write .env: %v", err)
	}

	if err := loadDotEnv(); err == nil {
		t.Fatal("loadDotEnv() returned nil for an invalid .env file")
	}
}

func withWorkingDirectory(t *testing.T, dir string) {
	t.Helper()

	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("get working directory: %v", err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("change working directory: %v", err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(originalDir); err != nil {
			t.Errorf("restore working directory: %v", err)
		}
	})
}
