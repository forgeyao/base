package jwt

import "testing"

func TestConfigValidate(t *testing.T) {
	cfg := Config{Secret: "secret"}
	if err := cfg.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
	if cfg.ExpireDays != 7 {
		t.Fatalf("Validate() ExpireDays = %d, want 7", cfg.ExpireDays)
	}
}

func TestConfigValidateRequiresSecret(t *testing.T) {
	cfg := Config{}
	if err := cfg.Validate(); err == nil {
		t.Fatal("Validate() error = nil, want error")
	}
}
