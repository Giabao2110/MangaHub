package utils

import (

    "errors"
    "time"
    "github.com/golang-jwt/jwt/v4"

)

// KHAI BÁO BIẾN NÀY Ở ĐÂY ĐỂ HẾT LỖI UNDEFINED
var SecretKey = []byte("mangahub_secret_2025") 

type Claims struct {
    UserID   string `json:"id"`
    Username string `json:"user"`
    Role     string `json:"role"`
    jwt.RegisteredClaims
}

// GenerateToken tạo token mới có hiệu lực 24h
func GenerateToken(userID, username string) (string, error) {
    expirationTime := time.Now().Add(24 * time.Hour)
    claims := &Claims{
        UserID: userID,
        Username: username,
        Role: "user", // Mặc định là user
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(expirationTime),
        },
    }   
    
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(SecretKey)
}

func ValidateToken(tokenString string) (*Claims, error) {
    claims := &Claims{}

    token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, errors.New("unexpected signing method")
        }
        return SecretKey, nil
    })

    if err != nil {
        return nil, err
    }

    if !token.Valid {
        return nil, errors.New("invalid token")
    }

    return claims, nil 
}