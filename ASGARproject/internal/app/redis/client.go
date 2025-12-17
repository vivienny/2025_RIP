package redis

import (
	"context"
	"fmt"
	"time"

	"decode/internal/app/config"

	"github.com/go-redis/redis/v8"
)

type Client struct {
	client *redis.Client
	prefix string
}

func NewClient(ctx context.Context, cfg config.RedisConfig) (*Client, error) {
	if cfg.Host == "" {
		cfg.Host = "localhost"
	}
	if cfg.Port == 0 {
		cfg.Port = 6379
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password:     cfg.Password,
		DB:           cfg.DB,
		DialTimeout:  cfg.DialTimeout,
		ReadTimeout:  cfg.ReadTimeout,
		PoolSize:     10,
		PoolTimeout:  30 * time.Second,
		MinIdleConns: 5,
	})

	// Test connection
	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &Client{
		client: rdb,
		prefix: "asgar:",
	}, nil
}

func (c *Client) Close() error {
	if c.client != nil {
		return c.client.Close()
	}
	return nil
}

// JWT Blacklist methods
func (c *Client) AddToBlacklist(ctx context.Context, token string, expiration time.Duration) error {
	key := c.prefix + "blacklist:" + token
	return c.client.Set(ctx, key, "1", expiration).Err()
}

func (c *Client) IsInBlacklist(ctx context.Context, token string) (bool, error) {
	key := c.prefix + "blacklist:" + token
	exists, err := c.client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return exists > 0, nil
}

// Session storage methods
func (c *Client) StoreSession(ctx context.Context, sessionID string, userID uint, expiration time.Duration) error {
	key := c.prefix + "session:" + sessionID
	return c.client.Set(ctx, key, userID, expiration).Err()
}

func (c *Client) GetSession(ctx context.Context, sessionID string) (uint, error) {
	key := c.prefix + "session:" + sessionID
	val, err := c.client.Get(ctx, key).Uint64()
	if err != nil {
		if err == redis.Nil {
			return 0, nil
		}
		return 0, err
	}
	return uint(val), nil
}

func (c *Client) DeleteSession(ctx context.Context, sessionID string) error {
	key := c.prefix + "session:" + sessionID
	return c.client.Del(ctx, key).Err()
}

// Get all keys (for debugging)
func (c *Client) GetAllKeys(ctx context.Context, pattern string) ([]string, error) {
	return c.client.Keys(ctx, c.prefix+pattern).Result()
}
