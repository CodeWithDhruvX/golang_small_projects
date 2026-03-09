package auth

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestHashPassword(t *testing.T) {
	password := "testpassword123"

	hashedPassword, err := HashPassword(password)
	require.NoError(t, err)
	assert.NotEmpty(t, hashedPassword)
	assert.NotEqual(t, password, hashedPassword)

	// Verify the hash works with bcrypt
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	assert.NoError(t, err)
}

func TestHashPasswordEmpty(t *testing.T) {
	password := ""

	hashedPassword, err := HashPassword(password)
	require.NoError(t, err)
	assert.NotEmpty(t, hashedPassword)
}

func TestCheckPassword(t *testing.T) {
	password := "testpassword123"

	// Hash the password
	hashedPassword, err := HashPassword(password)
	require.NoError(t, err)

	// Check correct password
	isValid := CheckPassword(password, hashedPassword)
	assert.True(t, isValid)

	// Check incorrect password
	isValid = CheckPassword("wrongpassword", hashedPassword)
	assert.False(t, isValid)
}

func TestGenerateJWT(t *testing.T) {
	userID := "user123"
	username := "testuser"
	secret := "test-secret"

	token, err := GenerateJWT(userID, username, secret)
	require.NoError(t, err)
	assert.NotEmpty(t, token)

	// Parse the token to verify its contents
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	require.NoError(t, err)
	require.True(t, parsedToken.Valid)

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	require.True(t, ok)

	assert.Equal(t, userID, claims["user_id"])
	assert.Equal(t, username, claims["username"])
	assert.Equal(t, "knowledge-base", claims["iss"])
	assert.Equal(t, "knowledge-base", claims["aud"])
}

func TestGenerateJWTEmptySecret(t *testing.T) {
	userID := "user123"
	username := "testuser"
	secret := ""

	_, err := GenerateJWT(userID, username, secret)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "secret key is required")
}

func TestValidateJWT(t *testing.T) {
	userID := "user123"
	username := "testuser"
	secret := "test-secret"

	// Generate token
	token, err := GenerateJWT(userID, username, secret)
	require.NoError(t, err)

	// Validate token
	claims, err := ValidateJWT(token, secret)
	require.NoError(t, err)
	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, username, claims.Username)
}

func TestValidateJWTInvalidToken(t *testing.T) {
	secret := "test-secret"

	// Test with invalid token
	invalidToken := "invalid.jwt.token"
	_, err := ValidateJWT(invalidToken, secret)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid token")
}

func TestValidateJWTWrongSecret(t *testing.T) {
	userID := "user123"
	username := "testuser"
	secret := "test-secret"
	wrongSecret := "wrong-secret"

	// Generate token with correct secret
	token, err := GenerateJWT(userID, username, secret)
	require.NoError(t, err)

	// Try to validate with wrong secret
	_, err = ValidateJWT(token, wrongSecret)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid token")
}

func TestValidateJWTExpiredToken(t *testing.T) {
	userID := "user123"
	username := "testuser"
	secret := "test-secret"

	// Create expired token manually
	expiredTime := time.Now().Add(-time.Hour).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  userID,
		"username": username,
		"iss":      "knowledge-base",
		"aud":      "knowledge-base",
		"exp":      expiredTime,
	})

	tokenString, err := token.SignedString([]byte(secret))
	require.NoError(t, err)

	// Try to validate expired token
	_, err = ValidateJWT(tokenString, secret)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "token is expired")
}

func TestRefreshToken(t *testing.T) {
	userID := "user123"
	username := "testuser"
	secret := "test-secret"

	// Generate original token
	originalToken, err := GenerateJWT(userID, username, secret)
	require.NoError(t, err)

	// Wait a bit to ensure different expiration times
	time.Sleep(10 * time.Millisecond)

	// Refresh token
	newToken, err := RefreshToken(originalToken, secret)
	require.NoError(t, err)
	assert.NotEmpty(t, newToken)
	assert.NotEqual(t, originalToken, newToken)

	// Validate new token
	claims, err := ValidateJWT(newToken, secret)
	require.NoError(t, err)
	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, username, claims.Username)

	// Ensure expiration times are different
	originalClaims, err := ValidateJWT(originalToken, secret)
	require.NoError(t, err)
	newClaims, err := ValidateJWT(newToken, secret)
	require.NoError(t, err)
	assert.True(t, newClaims.ExpiresAt.After(originalClaims.ExpiresAt))
}

func TestRefreshTokenInvalid(t *testing.T) {
	secret := "test-secret"

	// Try to refresh invalid token
	_, err := RefreshToken("invalid.token", secret)
	assert.Error(t, err)
}

func TestClaimsStruct(t *testing.T) {
	userID := "user123"
	username := "testuser"
	secret := "test-secret"

	token, err := GenerateJWT(userID, username, secret)
	require.NoError(t, err)

	claims, err := ValidateJWT(token, secret)
	require.NoError(t, err)

	// Test Claims struct methods
	assert.Equal(t, userID, claims.GetUserID())
	assert.Equal(t, username, claims.GetUsername())
	assert.Equal(t, "knowledge-base", claims.GetIssuer())
	assert.Equal(t, "knowledge-base", claims.GetAudience())
	assert.True(t, claims.IsExpired() == false)
}

func TestClaimsExpired(t *testing.T) {
	// Create expired claims
	claims := &Claims{
		UserID:   "user123",
		Username: "testuser",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Hour)),
		},
	}

	assert.True(t, claims.IsExpired())
}

func TestPasswordSecurity(t *testing.T) {
	passwords := []string{
		"simple",
		"complex123!@#",
		"verylongpasswordwithmultiplesymbolsandnumbers123456",
		"🔒passwordwithunicode",
	}

	for _, password := range passwords {
		hashed, err := HashPassword(password)
		require.NoError(t, err)

		// Verify password works
		isValid := CheckPassword(password, hashed)
		assert.True(t, isValid)

		// Verify different password doesn't work
		isValid = CheckPassword(password+"different", hashed)
		assert.False(t, isValid)

		// Verify hash is different each time (due to salt)
		hashed2, err := HashPassword(password)
		require.NoError(t, err)
		assert.NotEqual(t, hashed, hashed2)

		// But both work with the same password
		assert.True(t, CheckPassword(password, hashed2))
	}
}

func BenchmarkHashPassword(b *testing.B) {
	password := "testpassword123"
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		HashPassword(password)
	}
}

func BenchmarkCheckPassword(b *testing.B) {
	password := "testpassword123"
	hashedPassword, _ := HashPassword(password)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		CheckPassword(password, hashedPassword)
	}
}

func BenchmarkGenerateJWT(b *testing.B) {
	userID := "user123"
	username := "testuser"
	secret := "test-secret"
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GenerateJWT(userID, username, secret)
	}
}

func BenchmarkValidateJWT(b *testing.B) {
	userID := "user123"
	username := "testuser"
	secret := "test-secret"
	token, _ := GenerateJWT(userID, username, secret)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ValidateJWT(token, secret)
	}
}
