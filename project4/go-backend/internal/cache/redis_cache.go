package cache

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
)

// RedisCache handles caching operations using Redis
type RedisCache struct {
	client *redis.Client
}

// NewRedisCache creates a new Redis cache instance
func NewRedisCache(addr, password string, db int) *RedisCache {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	return &RedisCache{
		client: rdb,
	}
}

// CacheEntry represents a cached entry
type CacheEntry struct {
	Data      interface{} `json:"data"`
	CreatedAt time.Time   `json:"created_at"`
	TTL       int         `json:"ttl"`
}

// CacheAIResponse caches an AI response
func (rc *RedisCache) CacheAIResponse(ctx context.Context, prompt, model, response string, ttl int) error {
	key := rc.generateCacheKey("ai_response", prompt, model)
	
	entry := CacheEntry{
		Data:      response,
		CreatedAt: time.Now(),
		TTL:       ttl,
	}

	data, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("failed to marshal cache entry: %w", err)
	}

	err = rc.client.Set(ctx, key, data, time.Duration(ttl)*time.Second).Err()
	if err != nil {
		return fmt.Errorf("failed to cache AI response: %w", err)
	}

	logrus.Debugf("Cached AI response with key: %s, TTL: %d seconds", key, ttl)
	return nil
}

// GetCachedAIResponse retrieves a cached AI response
func (rc *RedisCache) GetCachedAIResponse(ctx context.Context, prompt, model string) (string, error) {
	key := rc.generateCacheKey("ai_response", prompt, model)
	
	data, err := rc.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return "", nil // Cache miss
		}
		return "", fmt.Errorf("failed to get cached response: %w", err)
	}

	var entry CacheEntry
	if err := json.Unmarshal([]byte(data), &entry); err != nil {
		return "", fmt.Errorf("failed to unmarshal cache entry: %w", err)
	}

	response, ok := entry.Data.(string)
	if !ok {
		return "", fmt.Errorf("invalid cache data type")
	}

	logrus.Debugf("Cache hit for AI response with key: %s", key)
	return response, nil
}

// CacheEmbedding caches text embeddings
func (rc *RedisCache) CacheEmbedding(ctx context.Context, text string, embedding []float64, ttl int) error {
	key := rc.generateCacheKey("embedding", text)
	
	entry := CacheEntry{
		Data:      embedding,
		CreatedAt: time.Now(),
		TTL:       ttl,
	}

	data, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("failed to marshal embedding cache entry: %w", err)
	}

	err = rc.client.Set(ctx, key, data, time.Duration(ttl)*time.Second).Err()
	if err != nil {
		return fmt.Errorf("failed to cache embedding: %w", err)
	}

	logrus.Debugf("Cached embedding with key: %s, dimensions: %d", key, len(embedding))
	return nil
}

// GetCachedEmbedding retrieves cached text embeddings
func (rc *RedisCache) GetCachedEmbedding(ctx context.Context, text string) ([]float64, error) {
	key := rc.generateCacheKey("embedding", text)
	
	data, err := rc.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // Cache miss
		}
		return nil, fmt.Errorf("failed to get cached embedding: %w", err)
	}

	// Parse the data directly as []float64
	var embedding []float64
	if err := json.Unmarshal([]byte(data), &embedding); err != nil {
		return nil, fmt.Errorf("failed to convert cached embedding: %w", err)
	}

	logrus.Debugf("Cache hit for embedding with key: %s", key)
	return embedding, nil
}

// CacheEmailClassification caches email classification results
func (rc *RedisCache) CacheEmailClassification(ctx context.Context, emailBody string, isRecruiter bool, confidence float64, ttl int) error {
	key := rc.generateCacheKey("email_classification", emailBody)
	
	classification := map[string]interface{}{
		"is_recruiter": isRecruiter,
		"confidence":   confidence,
	}
	
	entry := CacheEntry{
		Data:      classification,
		CreatedAt: time.Now(),
		TTL:       ttl,
	}

	data, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("failed to marshal classification cache entry: %w", err)
	}

	err = rc.client.Set(ctx, key, data, time.Duration(ttl)*time.Second).Err()
	if err != nil {
		return fmt.Errorf("failed to cache email classification: %w", err)
	}

	logrus.Debugf("Cached email classification with key: %s", key)
	return nil
}

// GetCachedEmailClassification retrieves cached email classification
func (rc *RedisCache) GetCachedEmailClassification(ctx context.Context, emailBody string) (bool, float64, error) {
	key := rc.generateCacheKey("email_classification", emailBody)
	
	data, err := rc.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return false, 0, nil // Cache miss
		}
		return false, 0, fmt.Errorf("failed to get cached classification: %w", err)
	}

	var entry CacheEntry
	if err := json.Unmarshal([]byte(data), &entry); err != nil {
		return false, 0, fmt.Errorf("failed to unmarshal classification cache entry: %w", err)
	}

	classification, ok := entry.Data.(map[string]interface{})
	if !ok {
		return false, 0, fmt.Errorf("invalid classification cache data type")
	}

	isRecruiter, ok1 := classification["is_recruiter"].(bool)
	confidence, ok2 := classification["confidence"].(float64)
	if !ok1 || !ok2 {
		return false, 0, fmt.Errorf("invalid classification cache data format")
	}

	logrus.Debugf("Cache hit for email classification with key: %s", key)
	return isRecruiter, confidence, nil
}

// CacheRequirementExtraction caches requirement extraction results
func (rc *RedisCache) CacheRequirementExtraction(ctx context.Context, emailBody string, requirements map[string]bool, ttl int) error {
	key := rc.generateCacheKey("requirements", emailBody)
	
	entry := CacheEntry{
		Data:      requirements,
		CreatedAt: time.Now(),
		TTL:       ttl,
	}

	data, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("failed to marshal requirements cache entry: %w", err)
	}

	err = rc.client.Set(ctx, key, data, time.Duration(ttl)*time.Second).Err()
	if err != nil {
		return fmt.Errorf("failed to cache requirements: %w", err)
	}

	logrus.Debugf("Cached requirement extraction with key: %s", key)
	return nil
}

// GetCachedRequirementExtraction retrieves cached requirement extraction
func (rc *RedisCache) GetCachedRequirementExtraction(ctx context.Context, emailBody string) (map[string]bool, error) {
	key := rc.generateCacheKey("requirements", emailBody)
	
	data, err := rc.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // Cache miss
		}
		return nil, fmt.Errorf("failed to get cached requirements: %w", err)
	}

	var entry CacheEntry
	if err := json.Unmarshal([]byte(data), &entry); err != nil {
		return nil, fmt.Errorf("failed to unmarshal requirements cache entry: %w", err)
	}

	requirements, ok := entry.Data.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid requirements cache data type")
	}

	// Convert interface{} map to bool map
	result := make(map[string]bool)
	for key, value := range requirements {
		if boolVal, ok := value.(bool); ok {
			result[key] = boolVal
		}
	}

	logrus.Debugf("Cache hit for requirement extraction with key: %s", key)
	return result, nil
}

// generateCacheKey generates a consistent cache key
func (rc *RedisCache) generateCacheKey(prefix string, parts ...string) string {
	combined := prefix + ":" + strings.Join(parts, ":")
	hash := sha256.Sum256([]byte(combined))
	return hex.EncodeToString(hash[:]) // Use first 32 characters of SHA-256
}

// InvalidatePattern invalidates cache entries matching a pattern
func (rc *RedisCache) InvalidatePattern(ctx context.Context, pattern string) error {
	keys, err := rc.client.Keys(ctx, pattern).Result()
	if err != nil {
		return fmt.Errorf("failed to get keys for pattern: %w", err)
	}

	if len(keys) == 0 {
		return nil
	}

	err = rc.client.Del(ctx, keys...).Err()
	if err != nil {
		return fmt.Errorf("failed to delete keys: %w", err)
	}

	logrus.Infof("Invalidated %d cache entries matching pattern: %s", len(keys), pattern)
	return nil
}

// GetCacheStats returns cache statistics
func (rc *RedisCache) GetCacheStats(ctx context.Context) (map[string]interface{}, error) {
	info, err := rc.client.Info(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get Redis info: %w", err)
	}

	// Parse Redis info for basic stats
	stats := map[string]interface{}{
		"redis_info": info,
		"connected":  true,
	}

	return stats, nil
}

// Close closes the Redis connection
func (rc *RedisCache) Close() error {
	return rc.client.Close()
}
