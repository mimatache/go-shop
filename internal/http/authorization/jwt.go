package authorization

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

const (
	jwtKey    = "Zq4t7w!z%C*F-J@NcRfUjXn2r5u8x/A?D(G+KbPdSgVkYp3s6v9y$B&E)H@McQfThWmZq4t7w!z%C*F-JaNdRgUkXn2r5u8x/A?D(G+KbPeShVmYq3s6v9y$B&E)H@Mc"
	userIDKey = "user-id"
	authToken = "auth-token"
)

var blackListedTokens = map[string]struct{}{}


func init() {
	go func() {
		timer := time.NewTimer(time.Minute * 5)
		for range timer.C{
				clearExpiredFromBlacklist()
		}
	}()
}

// Claim uses the standard JWT Claim to create a custom claim
type Claim struct {
	Username string `json:"username"`
	ID       uint   `json:"ID"`
	jwt.StandardClaims
}

// GenerateJWTToken generates a JWT token for the provided username
func GenerateJWTToken(username string, ID uint) (string, time.Time, error) {
	expirationTime := time.Now().Add(5 * time.Minute)
	claims := &Claim{
		Username: username,
		ID:       ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(jwtKey))
	return tokenString, expirationTime, err
}

// ValidateToken validates whether the value of the receoved key is a valid token
func ValidateToken(tknStr string) (bool, *Claim, error) {
	claim := &Claim{}
	tkn, err := jwt.ParseWithClaims(tknStr, claim, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtKey), nil
	})
	if err != nil {
		return false, nil, err
	}
	return tkn.Valid, claim, nil
}

func AddUserIdHeader(r *http.Request, claim *Claim) {
	r.Header.Add(userIDKey, fmt.Sprint(claim.ID))
}

func RemoveUserIDHeader(r *http.Request) {
	r.Header.Del(userIDKey)

}

func GetUserIdFromRequest(r *http.Request) (uint, error) {
	id := r.Header.Get(userIDKey)
	if id == "" {
		return 0, fmt.Errorf("request is missing user ID")
	}
	uid, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(uid), nil
}

func SetAuthCookie(w http.ResponseWriter, tokenString string, expires time.Time) {
	http.SetCookie(w, &http.Cookie{
		Name:    authToken,
		Value:   tokenString,
		Expires: expires,
	})
}

func GetAuthToken(r *http.Request) (string, error) {
	c, err := r.Cookie(authToken)
	if err != nil {
		return "", err
	}
	return c.Value, nil
}

func BlacklistToken(token string) {
	blackListedTokens[token] = struct{}{}
}

func IsBlacklisted(token string) bool {
	_, ok := blackListedTokens[token]
	return ok
} 

func clearExpiredFromBlacklist() {
	validTokens := map[string]struct{}{}
	for token := range blackListedTokens {
		valid, _, err := ValidateToken(token)
		if valid && err == nil {
			validTokens[token] = struct{}{}
		}
	}
	blackListedTokens = validTokens
}
