package auth

import (
	"errors"
	"time"

	"github.com/future-mcp/future-mcp-server/internal/types"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// Service 认证服务
type Service struct {
	jwtSecret []byte
}

// NewService 创建认证服务
func NewService(jwtSecret string) *Service {
	return &Service{
		jwtSecret: []byte(jwtSecret),
	}
}

// JWTClaims JWT声明
type JWTClaims struct {
	UserID   uuid.UUID `json:"user_id"`
	Username string    `json:"username"`
	Role     string    `json:"role"`
	jwt.RegisteredClaims
}

// GenerateToken 生成JWT令牌
func (s *Service) GenerateToken(userID uuid.UUID, username, role string, expire time.Duration) (string, error) {
	claims := JWTClaims{
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expire)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "talink-mcp-server",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}

// ValidateToken 验证JWT令牌
func (s *Service) ValidateToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return s.jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// ValidateAPIKey 验证API密钥（简化实现）
func (s *Service) ValidateAPIKey(apiKey string) (*types.User, error) {
	// 简化实现，实际应该从数据库验证
	if apiKey == "" {
		return nil, errors.New("invalid API key")
	}

	// 返回模拟用户
	return &types.User{
		ID:       uuid.New(),
		Username: "api-user",
		Role:     types.UserRoleDeveloper,
	}, nil
}
