package logger

import "testing"

func TestConfigValidateDefaultsSingleFile(t *testing.T) {
	cfg := Config{Filename: "app.log"}
	if err := cfg.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
	if cfg.MaxSize != 100 {
		t.Fatalf("Validate() MaxSize = %d, want 100", cfg.MaxSize)
	}
	if cfg.MaxAge != 7 {
		t.Fatalf("Validate() MaxAge = %d, want 7", cfg.MaxAge)
	}
	if cfg.MaxBackups != 7 {
		t.Fatalf("Validate() MaxBackups = %d, want 7", cfg.MaxBackups)
	}
}

func TestConfigValidateRequiresFilenameInSingleFileMode(t *testing.T) {
	cfg := Config{}
	if err := cfg.Validate(); err == nil {
		t.Fatal("Validate() error = nil, want error")
	}
}

func TestConfigValidateRequiresNamedLoggers(t *testing.T) {
	cfg := Config{
		Loggers: []Entry{
			{Filename: "app.log"},
		},
	}
	if err := cfg.Validate(); err == nil {
		t.Fatal("Validate() error = nil, want error")
	}
}

func TestConfigValidateRequiresLoggerFilename(t *testing.T) {
	cfg := Config{
		Loggers: []Entry{
			{Name: "app"},
		},
	}
	if err := cfg.Validate(); err == nil {
		t.Fatal("Validate() error = nil, want error")
	}
}
