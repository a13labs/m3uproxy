package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func createJWT(userID, role string) (string, error) {
	// Define token expiration time
	expirationTime := time.Now().Add(time.Hour * time.Duration(authConfig.ExpirationTime)) // 1 hour expiry

	// Define the claims
	claims := jwt.MapClaims{
		"sub":  userID,                // Subject or user ID
		"exp":  expirationTime.Unix(), // Expiration time
		"iat":  time.Now().Unix(),     // Issued at time
		"role": role,                  // Custom claim (e.g., user role)
	}

	// Create a new token using the HS256 signing method
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with your secret key
	tokenString, err := token.SignedString([]byte(authConfig.SecretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func verifyJWT(tokenString string) (jwt.MapClaims, error) {
	// Parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Ensure the signing method is HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(authConfig.SecretKey), nil
	})

	// Handle errors
	if err != nil {
		return nil, err
	}

	// Extract claims if the token is valid
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, fmt.Errorf("invalid token")
	}
}

func CreateToken(userId, password string) (string, error) {

	if CheckCredentials(userId, password) {
		role, err := GetRole(userId)
		if err != nil || role == "" {
			role = "viewer"
		}
		return createJWT(userId, role)
	}
	return "", fmt.Errorf("invalid credentials")
}

func VerifyUserToken(userId, token string) bool {
	claims, err := verifyJWT(token)
	if err != nil {
		return false
	}
	if sub, ok := claims["sub"].(string); ok {
		if sub == userId {
			return true
		}
	}
	return false
}

func VerifyToken(token string) bool {
	_, err := verifyJWT(token)
	return err == nil
}

func GetRoleFromToken(token string) (string, error) {
	claims, err := verifyJWT(token)
	if err != nil {
		return "", err
	}
	if role, ok := claims["role"].(string); ok {
		return role, nil
	}
	return "", fmt.Errorf("role not found")
}

func GetUserFromToken(token string) (string, error) {
	claims, err := verifyJWT(token)
	if err != nil {
		return "", err
	}
	if sub, ok := claims["sub"].(string); ok {
		return sub, nil
	}
	return "", fmt.Errorf("user id not found")
}
