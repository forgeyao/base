package jwt

import (
	"errors"
	"testing"
	"time"

	conf "util/config"

	gjwt "github.com/golang-jwt/jwt/v5"
)

func newTestService(t *testing.T) *Service {
	t.Helper()

	svc, err := New(conf.JWT{
		Secret:     "secret",
		Issuer:     "util-test",
		ExpireDays: 1,
	})
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	return svc
}

func TestGenerateAndParseToken(t *testing.T) {
	svc := newTestService(t)

	token, jti, expireAt, err := svc.GenerateToken(42)
	if err != nil {
		t.Fatalf("GenerateToken() error = %v", err)
	}
	if jti == "" {
		t.Fatal("GenerateToken() jti = empty")
	}
	if expireAt.IsZero() {
		t.Fatal("GenerateToken() expireAt = zero")
	}

	claims, err := svc.ParseToken(token)
	if err != nil {
		t.Fatalf("ParseToken() error = %v", err)
	}
	if claims.UID != 42 {
		t.Fatalf("ParseToken() uid = %d, want 42", claims.UID)
	}
	if claims.JTI != jti {
		t.Fatalf("ParseToken() jti = %q, want %q", claims.JTI, jti)
	}
	if claims.Issuer != "util-test" {
		t.Fatalf("ParseToken() issuer = %q, want util-test", claims.Issuer)
	}
}

func TestRefreshToken(t *testing.T) {
	svc := newTestService(t)

	oldToken, oldJTI, _, err := svc.GenerateToken(7)
	if err != nil {
		t.Fatalf("GenerateToken() error = %v", err)
	}

	newToken, newJTI, _, err := svc.RefreshToken(oldToken)
	if err != nil {
		t.Fatalf("RefreshToken() error = %v", err)
	}
	if newToken == oldToken {
		t.Fatal("RefreshToken() token was not regenerated")
	}
	if newJTI == oldJTI {
		t.Fatal("RefreshToken() jti was not regenerated")
	}
}

func TestShouldRefreshTokenByExpireTime(t *testing.T) {
	svc := newTestService(t)

	if !svc.ShouldRefreshTokenByExpireTime(time.Now().Add(90 * time.Minute)) {
		t.Fatal("ShouldRefreshTokenByExpireTime() = false, want true")
	}
	if svc.ShouldRefreshTokenByExpireTime(time.Now().Add(3 * time.Hour)) {
		t.Fatal("ShouldRefreshTokenByExpireTime() = true, want false")
	}
}

func TestParseTokenExpired(t *testing.T) {
	svc := newTestService(t)
	token := gjwt.NewWithClaims(gjwt.SigningMethodHS256, Claims{
		UID: 1,
		JTI: "expired-jti",
		RegisteredClaims: gjwt.RegisteredClaims{
			ExpiresAt: gjwt.NewNumericDate(time.Now().Add(-time.Minute)),
			IssuedAt:  gjwt.NewNumericDate(time.Now().Add(-2 * time.Minute)),
			NotBefore: gjwt.NewNumericDate(time.Now().Add(-2 * time.Minute)),
			Issuer:    "util-test",
		},
	})
	tokenString, err := token.SignedString([]byte("secret"))
	if err != nil {
		t.Fatalf("SignedString() error = %v", err)
	}

	if _, err := svc.ParseToken(tokenString); !errors.Is(err, gjwt.ErrTokenExpired) {
		t.Fatalf("ParseToken() error = %v, want ErrTokenExpired", err)
	}
}
