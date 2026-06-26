package jwt

import (
	"testing"
	"time"
)

func TestGenerateAndParseToken(t *testing.T) {
	secret := "test-secret"
	userID := uint64(12345)

	token, err := GenerateToken(userID, secret, 3600)
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}
	if token == "" {
		t.Fatal("token is empty")
	}

	claims, err := ParseToken(token, secret)
	if err != nil {
		t.Fatalf("ParseToken failed: %v", err)
	}

	if claims.UserID != userID {
		t.Errorf("userID = %d, want %d", claims.UserID, userID)
	}
}

func TestParseToken_InvalidSecret(t *testing.T) {
	token, _ := GenerateToken(1, "secret-a", 3600)

	_, err := ParseToken(token, "secret-b")
	if err == nil {
		t.Error("should fail with wrong secret")
	}
}

func TestParseToken_Expired(t *testing.T) {
	token, _ := GenerateToken(1, "secret", -1) // already expired

	_, err := ParseToken(token, "secret")
	if err != ErrTokenExpired {
		t.Errorf("expected ErrTokenExpired, got %v", err)
	}
}

func TestParseToken_InvalidFormat(t *testing.T) {
	_, err := ParseToken("not.a.token", "secret")
	if err != ErrTokenInvalid {
		t.Errorf("expected ErrTokenInvalid, got %v", err)
	}

	_, err = ParseToken("", "secret")
	if err != ErrTokenInvalid {
		t.Errorf("expected ErrTokenInvalid for empty token, got %v", err)
	}
}

func TestTokenExpiry(t *testing.T) {
	token, err := GenerateToken(42, "secret", 2)
	if err != nil {
		t.Fatal(err)
	}

	claims, err := ParseToken(token, "secret")
	if err != nil {
		t.Fatal(err)
	}

	if claims.ExpiresAt == nil {
		t.Fatal("ExpiresAt is nil")
	}

	expectedExp := time.Now().Add(2 * time.Second)
	diff := claims.ExpiresAt.Time.Sub(expectedExp).Abs()
	if diff > time.Second {
		t.Errorf("expiry time mismatch: got %v, expected ~%v", claims.ExpiresAt.Time, expectedExp)
	}
}
