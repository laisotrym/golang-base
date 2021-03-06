package auth

import (
    "fmt"
    "time"
    
    "github.com/dgrijalva/jwt-go"
    
    "safeweb.app/model"
)

type IJWTManager interface {
    Generate(user *model.User, role string) (string, error)
    Verify(accessToken string) (*UserClaims, error)
}

// JWTManager is a JSON web token manager
type JWTManager struct {
    secretKey     string
    tokenDuration time.Duration
}

// UserClaims is a custom JWT claims that contains some user's information
type UserClaims struct {
    jwt.StandardClaims
    Username string `json:"username"`
    Role     string `json:"role"`
}

// NewJWTManager returns a new JWT manager
func NewJWTManager(secretKey string, tokenDuration time.Duration) *JWTManager {
    return &JWTManager{secretKey, tokenDuration}
}

// Generate generates and signs a new token for a user
func (jm *JWTManager) Generate(user *model.User, role string) (string, error) {
    claims := UserClaims{
        StandardClaims: jwt.StandardClaims{
            ExpiresAt: time.Now().Add(jm.tokenDuration).Unix(),
        },
        Username: user.UserName,
        // Role:     user.Role,
        Role: role,
    }
    
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(jm.secretKey))
}

// Verify verifies the access token string and return a user claim if the token is valid
func (jm *JWTManager) Verify(accessToken string) (*UserClaims, error) {
    token, err := jwt.ParseWithClaims(
        accessToken,
        &UserClaims{},
        func(token *jwt.Token) (interface{}, error) {
            _, ok := token.Method.(*jwt.SigningMethodHMAC)
            if !ok {
                return nil, fmt.Errorf("unexpected token signing method")
            }
            
            return []byte(jm.secretKey), nil
        },
    )
    
    if err != nil {
        return nil, fmt.Errorf("invalid token: %w", err)
    }
    
    claims, ok := token.Claims.(*UserClaims)
    if !ok {
        return nil, fmt.Errorf("invalid token claims")
    }
    
    return claims, nil
}
