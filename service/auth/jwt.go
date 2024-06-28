package auth

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/davidado/go-api-reference/config"
	"github.com/davidado/go-api-reference/netjson"
	"github.com/davidado/go-api-reference/types"
	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

// UserKey is the key for the user ID in the context
const UserKey contextKey = "userID"

// CreateJWT creates a JWT token
func CreateJWT(secret []byte, userID int) (string, error) {
	expiration := time.Second * time.Duration(config.Envs.JWTExpirationInSeconds)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID":    strconv.Itoa(userID),
		"expiredAt": time.Now().Add(expiration).Unix(),
	})

	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// WithJWTAuth adds JWT authentication to a handler
func WithJWTAuth(handlerFunc http.HandlerFunc, store types.UserStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the token from the user request.
		tokenString := getTokenFromRequest(r)
		if tokenString == "" {
			netjson.WriteError(w, http.StatusUnauthorized, fmt.Errorf("no token provided"))
			return
		}

		// Validate the JWT.
		token, err := validateToken(tokenString)
		if err != nil {
			log.Printf("failed to validate token: %v", err)
			permissionDenied(w)
			return
		}

		if !token.Valid {
			log.Println("invalid token")
			permissionDenied(w)
			return
		}

		// Fetch the userID from the db using the ID from the token.
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			permissionDenied(w)
			return
		}

		userID, err := strconv.Atoi(claims["userID"].(string))
		if err != nil {
			permissionDenied(w)
			return
		}

		u, err := store.GetUserByID(userID)
		if err != nil {
			log.Printf("failed to get user by ID: %v", err)
			permissionDenied(w)
			return
		}

		// Set context "userID" to the user ID.
		ctx := context.WithValue(r.Context(), UserKey, u.ID)
		r = r.WithContext(ctx)

		handlerFunc(w, r)
	}
}

func getTokenFromRequest(r *http.Request) string {
	token := r.Header.Get("Authorization")
	if token == "" {
		return ""
	}

	return token
}

func validateToken(tokenString string) (*jwt.Token, error) {
	secret := []byte(config.Envs.JWTSecret)
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secret, nil
	})
}

func permissionDenied(w http.ResponseWriter) {
	netjson.WriteError(w, http.StatusUnauthorized, fmt.Errorf("permission denied"))
}

// GetUserIDFromContext gets the user ID from the context
func GetUserIDFromContext(ctx context.Context) int {
	userID := ctx.Value(UserKey)
	if userID == nil {
		return 0
	}

	return userID.(int)
}
