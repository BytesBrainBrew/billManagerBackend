package utils

import (
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte(os.Getenv("JWT_SECRET"))
var jwtTokenFromEmail = make(map[string]*jwt.Token)

// Claims struct used to encode and decode JWT claims
type Claims struct {
    Email string `json:"email"`
    jwt.StandardClaims
}

// GenerateJWT generates a new JWT token for a user
func GenerateJWT(email string) (string, error) {
    expirationTime := time.Now().Add(15 * time.Minute)
    claims := &Claims{
        Email: email,
        StandardClaims: jwt.StandardClaims{
            ExpiresAt: expirationTime.Unix(),
        },
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    jwtTokenFromEmail[email] = token
    return token.SignedString(jwtKey)
}

// ValidateToken validates a given JWT token
func ValidateToken(tokenStr string) (*Claims, error) {
    claims := &Claims{}
    _, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
        return jwtKey, nil
    })

    return claims, err
}

func ClearExpiredTokens(tokenStr string) {
    claims, _ := ValidateToken(tokenStr)
    delete(jwtTokenFromEmail, claims.Email)
}

func IsTokenExpired(tokenStr string) (string,bool) {
    claims, _ := ValidateToken(tokenStr)
    var ok bool = true
    //check if email mapped with any token in jwtTokenFromEmail
    if _, ok = jwtTokenFromEmail[claims.Email]; ok {
        return claims.Email,true
    }
    return claims.Email,time.Unix(claims.ExpiresAt, 0).Before(time.Now()) && ok
}