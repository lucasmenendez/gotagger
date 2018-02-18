package gotagger

import (
	"testing"
	"os"
	"path/filepath"
)

func TestLoadLanguage(t *testing.T) {
	if _, err := loadLanguage("false-language-code"); err == nil {
		t.Error("Expected error, got nil")
	}

	if _, err := loadLanguage("en"); err != nil {
		t.Errorf("Expected nil, got %s", err.Error())
	}

	var dir string = filepath.FromSlash("/nowhere")
	os.Setenv("STOPWORDS", dir)
	if _, err := loadLanguage("no-matter"); err == nil {
		t.Error("Expected error opening '/nowhere' folder, got nil")
	}
	os.Unsetenv("STOPWORDS")
}
