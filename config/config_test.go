package config

import "testing"

func TestJWTValidate(t *testing.T) {
	cfg := JWT{Secret: "secret"}
	if err := cfg.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
	if cfg.ExpireDays != 7 {
		t.Fatalf("Validate() ExpireDays = %d, want 7", cfg.ExpireDays)
	}
}

func TestJWTValidateRequiresSecret(t *testing.T) {
	cfg := JWT{}
	if err := cfg.Validate(); err == nil {
		t.Fatal("Validate() error = nil, want error")
	}
}
