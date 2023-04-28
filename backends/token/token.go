package token

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

type TokenClient struct {
	redisClient *redis.Client
}

func NewClient() (*TokenClient, error) {
	redisURL := os.Getenv("REDIS_URL")
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, err
	}
	return &TokenClient{redis.NewClient(opt)}, nil
}

func (t *TokenClient) CreateToken(ctx context.Context, bucket, objectName string) (string, error) {
	token := generateSecureToken(32)
	key := key(bucket, objectName)
	if err := t.redisClient.Set(ctx, key, token, 120*time.Second).Err(); err != nil {
		return "", err
	}
	return fmt.Sprintf(`{"token":"%s", "path":"%s"}`, token, key), nil
}

func (t *TokenClient) GetToken(ctx context.Context, bucket, objectName string) (string, error) {
	key := key(bucket, objectName)
	tk, err := t.redisClient.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", fmt.Errorf("key does not exist")
	} else if err != nil {
		return "", err
	}
	return tk, nil
}

func key(bucket, objectName string) string {
	return fmt.Sprintf("%s/%s", bucket, objectName)
}

func generateSecureToken(length int) string {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}
