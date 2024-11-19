package auth

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"scratch/config"
	"scratch/types"
	"scratch/utils"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AuthContext string

const (
	USER_ID AuthContext = "userID"
)

var (
	ErrInvalidToken     = errors.New("invalid token")
	ErrPermissionDenied = errors.New("permission denied")
)

func CreateJWT(secret []byte, userId int) (string, error) {
	expiration := time.Second * time.Duration(config.Envs.JWTExpirationInSeconds)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID":    strconv.Itoa(userId),
		"expiredAt": time.Now().Add(expiration).Unix(),
	})

	tokenString, err := token.SignedString(secret)

	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// JWT auth middleware
func WithJWTAuth(handlerFunc http.HandlerFunc, store types.UserStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// get the token from the user request
		tokenString := getTokenFromRequest(r)
		// validate the token
		token, err := validateToken(tokenString)
		if err != nil || !token.Valid {
			permissionDenied(w)
			return
		}

		// extract userId and set
		claims := token.Claims.(jwt.MapClaims)
		str := claims["userID"].(string)
		userID, err := strconv.Atoi(str)
		if err != nil || !token.Valid {
			permissionDenied(w)
			return
		}
		user, err := store.GetUserByID(userID)
		if err != nil || !token.Valid {
			permissionDenied(w)
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, USER_ID, user.ID)
		r = r.WithContext(ctx)

		handlerFunc(w, r)
	}
}

func getTokenFromRequest(request *http.Request) string {
	authToken := request.Header.Get("Authorization")

	return authToken
}

func validateToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("%w: unexpected signing method", ErrInvalidToken)
		}

		return []byte(config.Envs.JWTSecret), nil
	})
}

func permissionDenied(w http.ResponseWriter) {
	utils.WriteError(w, http.StatusForbidden, ErrPermissionDenied)
}
