package jwt

import (
	"errors"
	"time"

	conf "util/config"

	gjwt "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const defaultRefreshBefore = 2 * time.Hour

var errInvalidSigningMethod = errors.New("invalid signing method")

// Claims contains the application-specific JWT payload.
type Claims struct {
	UID int64  `json:"uid"`
	JTI string `json:"jti"`
	gjwt.RegisteredClaims
}

// Service signs and validates JWT tokens using repository-local config.
type Service struct {
	secret        []byte
	issuer        string
	expire        time.Duration
	refreshBefore time.Duration
	now           func() time.Time
}

// New returns a JWT service from config.
func New(cfg conf.JWT) (*Service, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return &Service{
		secret:        []byte(cfg.Secret),
		issuer:        cfg.Issuer,
		expire:        time.Duration(cfg.ExpireDays) * 24 * time.Hour,
		refreshBefore: defaultRefreshBefore,
		now:           time.Now,
	}, nil
}

// GenerateToken creates a signed token for the given uid.
func (s *Service) GenerateToken(uid int64) (string, string, time.Time, error) {
	jti := uuid.NewString()
	now := s.now()
	expireTime := now.Add(s.expire)

	claims := Claims{
		UID: uid,
		JTI: jti,
		RegisteredClaims: gjwt.RegisteredClaims{
			ExpiresAt: gjwt.NewNumericDate(expireTime),
			IssuedAt:  gjwt.NewNumericDate(now),
			NotBefore: gjwt.NewNumericDate(now),
			Issuer:    s.issuer,
		},
	}

	token := gjwt.NewWithClaims(gjwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(s.secret)
	if err != nil {
		return "", "", time.Time{}, err
	}

	return tokenString, jti, expireTime, nil
}

// ParseToken parses and validates a token string.
func (s *Service) ParseToken(tokenString string) (*Claims, error) {
	token, err := gjwt.ParseWithClaims(tokenString, &Claims{}, func(token *gjwt.Token) (interface{}, error) {
		if token.Method != gjwt.SigningMethodHS256 {
			return nil, errInvalidSigningMethod
		}
		return s.secret, nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("invalid token")
	}
	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	return claims, nil
}

// ValidateToken validates a token and returns its claims.
func (s *Service) ValidateToken(tokenString string) (*Claims, error) {
	return s.ParseToken(tokenString)
}

// RefreshToken reissues a token for the same uid.
func (s *Service) RefreshToken(oldTokenString string) (string, string, time.Time, error) {
	claims, err := s.ParseToken(oldTokenString)
	if err != nil {
		return "", "", time.Time{}, err
	}

	return s.GenerateToken(claims.UID)
}

// ShouldRefreshToken reports whether a token is close enough to expiry to refresh.
func (s *Service) ShouldRefreshToken(tokenString string) (bool, error) {
	claims, err := s.ParseToken(tokenString)
	if err != nil {
		return false, err
	}

	return s.ShouldRefreshTokenByExpireTime(claims.ExpiresAt.Time), nil
}

// ShouldRefreshTokenByExpireTime reports whether the expiry time is within the refresh window.
func (s *Service) ShouldRefreshTokenByExpireTime(expireTime time.Time) bool {
	return time.Until(expireTime) < s.refreshBefore
}

// SetRefreshBefore overrides the refresh window.
func (s *Service) SetRefreshBefore(d time.Duration) {
	if d > 0 {
		s.refreshBefore = d
	}
}
