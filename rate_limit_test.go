package ozonCore

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRateLimiterMatchRule(t *testing.T) {
	ResetRateRules()
	defer ResetRateRules()

	// /v1/analytics/data：1 次/分钟 + 无订阅 50 次/天
	rulePath, rule, ok := globalRateLimiter.matchRule("/v1/analytics/data")
	assert.True(t, ok)
	assert.Equal(t, "/v1/analytics/data", rulePath)
	assert.Equal(t, 50, rule.DailyLimit)
	assert.Equal(t, 61*time.Second, rule.MinInterval)

	// /v1/report/discounted/create：1 次/分钟
	rulePath, rule, ok = globalRateLimiter.matchRule("/v1/report/discounted/create")
	assert.True(t, ok)
	assert.Equal(t, "/v1/report/discounted/create", rulePath)
	assert.Equal(t, 61*time.Second, rule.MinInterval)
	assert.Equal(t, 0, rule.DailyLimit)

	// /v1/barcode/generate：20 次/分钟 → 间隔 ≥3.1 秒
	rulePath, rule, ok = globalRateLimiter.matchRule("/v1/barcode/generate")
	assert.True(t, ok)
	assert.Equal(t, "/v1/barcode/generate", rulePath)
	assert.Equal(t, 3100*time.Millisecond, rule.MinInterval)

	// /v1/barcode/add：20 次/分钟
	rulePath, rule, ok = globalRateLimiter.matchRule("/v1/barcode/add")
	assert.True(t, ok)
	assert.Equal(t, "/v1/barcode/add", rulePath)
	assert.Equal(t, 3100*time.Millisecond, rule.MinInterval)

	// /v2/products/stocks：80 次/分钟 → 间隔 ≥1 秒
	rulePath, rule, ok = globalRateLimiter.matchRule("/v2/products/stocks")
	assert.True(t, ok)
	assert.Equal(t, "/v2/products/stocks", rulePath)
	assert.Equal(t, 1*time.Second, rule.MinInterval)

	// 文档未注明固定限速的接口：不匹配任何规则
	_, _, ok = globalRateLimiter.matchRule("/v1/analytics/product-queries")
	assert.False(t, ok)
	_, _, ok = globalRateLimiter.matchRule("/v2/analytics/stock_on_warehouses")
	assert.False(t, ok)
	_, _, ok = globalRateLimiter.matchRule("/v4/posting/fbs/list")
	assert.False(t, ok)
	_, _, ok = globalRateLimiter.matchRule("/v3/product/import")
	assert.False(t, ok)
}

func TestRateLimiterFallbackRule(t *testing.T) {
	ResetRateRules()
	defer ResetRateRules()

	// 未配置的接口：兜底 ≤5 次/秒（间隔 ≥210 毫秒），按完整路径隔离
	keyPath, rule := globalRateLimiter.ruleFor("/v4/posting/fbs/list")
	assert.Equal(t, "/v4/posting/fbs/list", keyPath)
	assert.Equal(t, FallbackMinInterval, rule.MinInterval)
	assert.Equal(t, 0, rule.DailyLimit)

	// 已配置的接口：具体规则优先于兜底
	keyPath, rule = globalRateLimiter.ruleFor("/v1/analytics/data")
	assert.Equal(t, "/v1/analytics/data", keyPath)
	assert.Equal(t, 61*time.Second, rule.MinInterval)

	keyPath, rule = globalRateLimiter.ruleFor("/v1/barcode/generate")
	assert.Equal(t, "/v1/barcode/generate", keyPath)
	assert.Equal(t, 3100*time.Millisecond, rule.MinInterval)
}

func TestRateLimiterFallbackSpacing(t *testing.T) {
	ResetRateRules()
	defer ResetRateRules()

	ctx := context.Background()
	start := time.Now()
	assert.Nil(t, globalRateLimiter.acquire(ctx, "cid|/v4/posting/fbs/list", RateRule{MinInterval: FallbackMinInterval}))
	assert.Nil(t, globalRateLimiter.acquire(ctx, "cid|/v4/posting/fbs/list", RateRule{MinInterval: FallbackMinInterval}))
	assert.True(t, time.Since(start) >= FallbackMinInterval)
}

func TestRateLimiterDailyLimit(t *testing.T) {
	ResetRateRules()
	defer ResetRateRules()

	RegisterRateRule("/test/daily", RateRule{MinInterval: 0, DailyLimit: 2})

	ctx := context.Background()
	assert.Nil(t, globalRateLimiter.acquire(ctx, "cid|/test/daily", RateRule{MinInterval: 0, DailyLimit: 2}))
	assert.Nil(t, globalRateLimiter.acquire(ctx, "cid|/test/daily", RateRule{MinInterval: 0, DailyLimit: 2}))
	err := globalRateLimiter.acquire(ctx, "cid|/test/daily", RateRule{MinInterval: 0, DailyLimit: 2})
	assert.True(t, errors.Is(err, ErrDailyRateLimitExceeded))

	// 不同账号互不影响
	assert.Nil(t, globalRateLimiter.acquire(ctx, "other|/test/daily", RateRule{MinInterval: 0, DailyLimit: 2}))
}

func TestRateLimiterMinInterval(t *testing.T) {
	ResetRateRules()
	defer ResetRateRules()

	RegisterRateRule("/test/spacing", RateRule{MinInterval: 100 * time.Millisecond, DailyLimit: 0})

	ctx := context.Background()
	rule := RateRule{MinInterval: 100 * time.Millisecond, DailyLimit: 0}

	start := time.Now()
	assert.Nil(t, globalRateLimiter.acquire(ctx, "cid|/test/spacing", rule))
	assert.Nil(t, globalRateLimiter.acquire(ctx, "cid|/test/spacing", rule))
	assert.True(t, time.Since(start) >= 100*time.Millisecond)
}

func TestRateLimiterAcquireCtxCancel(t *testing.T) {
	ResetRateRules()
	defer ResetRateRules()

	RegisterRateRule("/test/cancel", RateRule{MinInterval: 10 * time.Second, DailyLimit: 0})

	ctx := context.Background()
	assert.Nil(t, globalRateLimiter.acquire(ctx, "cid|/test/cancel", RateRule{MinInterval: 10 * time.Second}))

	ctxCancel, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	err := globalRateLimiter.acquire(ctxCancel, "cid|/test/cancel", RateRule{MinInterval: 10 * time.Second})
	assert.True(t, errors.Is(err, context.DeadlineExceeded))
}
