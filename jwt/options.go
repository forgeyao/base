package jwt

import "errors"

// Config defines configuration used by the jwt package.
type Config struct {
	Secret     string `yaml:"secret"`
	Issuer     string `yaml:"issuer,omitempty"`
	ExpireDays int    `yaml:"expire_days,omitempty"`
}

// Validate normalizes and validates the JWT config.
func (c *Config) Validate() error {
	if c.Secret == "" {
		return errors.New("conf.JWT.Secret must provided")
	}
	if c.ExpireDays <= 0 {
		c.ExpireDays = 7
	}
	return nil
}
