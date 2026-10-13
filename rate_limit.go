package ozonCore

import (
	"context"
	"errors"
	"strings"
	"sync"
	"time"
)

// RateRule 路径前缀的限流规则。
type RateRule struct {
	// MinInterval 同一账号（Client-Id）对匹配路径相邻两次请求的最小间隔；<= 0 表示不限频。
	MinInterval time.Duration

	// DailyLimit 滚动 24 小时窗口内的最大请求次数；<= 0 表示不限次数。
	DailyLimit int
}

// ErrDailyRateLimitExceeded 当天配额已用完（隔天自动恢复）。
var ErrDailyRateLimitExceeded = errors.New("ozon daily rate limit exceeded")

// FallbackMinInterval 未配置具体规则的接口的兜底限频间隔：
// 统一 ≤5 次/秒 → 相邻两次请求间隔 ≥ 210 毫秒（≈4.76 次/秒，留余量严格小于 5 次/秒）。
// 兜底按 Client-Id + 具体路径 隔离配额。
const FallbackMinInterval = 210 * time.Millisecond

// defaultRateRules 内置默认限流规则（来源：docs.ozon.ru Seller API 中文文档各方法的
// 「描述和范例」及「接口使用规范」模块）。
//
// 有明确数字限制的接口：
//
//   - /v1/analytics/data           ≤1 次/分钟（文档）；无订阅卖家另 ≤50 次/天
//     （2026-09-16 dev.ozon.ru news 797）→ 间隔 ≥61 秒 + 每天 ≤50 次
//   - /v1/report/discounted/create ≤1 次/分钟（文档）→ 间隔 ≥61 秒
//   - /v1/barcode/add              ≤20 次/分钟（文档）→ 间隔 ≥3.1 秒
//   - /v1/barcode/generate         ≤20 次/分钟（文档）→ 间隔 ≥3.1 秒
//   - /v2/products/stocks          ≤80 次/分钟（文档）→ 间隔 ≥1 秒；
//     另：同一「商品-仓库」30 秒内只能更新 1 次库存，需调用方按货位合并请求保证
//
// 不做固定限速的接口（文档未给出固定数字，避免一刀切）：
//
//   - /v3/product/import、/v1/product/import-by-sku、/v1/product/attributes/update 等
//     商品创建/更新类方法：每分钟/每天的限额因账号而异（用 /v4/product/info/limit 查询），
//     超限时 429 响应带 Item-Retry-After 头，SDK 的 429 退避重试会按其等待；
//   - /v1/product/unarchive：每天 ≤100 件「自动归档」商品，但手动归档无限制，
//     且限的是商品件数而非请求次数，传输层无法区分，交由 Ozon 429 处理；
//   - 其余接口（如 /v4/posting/fbs/list 等）：文档只要求全局 ≤50 次/秒/Client ID
//     并提示大流量会触发 "Circle is open"。这些接口没有具体规则，统一走
//     FallbackMinInterval 兜底限频（≤5 次/秒）。
//
// 兜底规则：所有未在上表配置的接口，按 ≤5 次/秒（间隔 ≥210 毫秒）限频，
// 按 Client-Id + 路径 隔离配额。
var defaultRateRules = map[string]RateRule{
	"/v1/analytics/data": {
		MinInterval: 61 * time.Second,
		DailyLimit:  50,
	},
	"/v1/report/discounted/create": {MinInterval: 61 * time.Second},
	"/v1/barcode/add":              {MinInterval: 3100 * time.Millisecond},
	"/v1/barcode/generate":         {MinInterval: 3100 * time.Millisecond},
	"/v2/products/stocks":          {MinInterval: 1 * time.Second},
}

type rateLimitEntry struct {
	lastCall time.Time
	dayStart time.Time
	dayCount int
}

// rateLimiter 进程级限流器：所有 Client 实例共享，按 Client-Id + 路径前缀 隔离配额。
type rateLimiter struct {
	mu    sync.Mutex
	rules map[string]RateRule
	state map[string]*rateLimitEntry
}

var globalRateLimiter = func() *rateLimiter {
	l := &rateLimiter{
		rules: map[string]RateRule{},
		state: map[string]*rateLimitEntry{},
	}
	for path, rule := range defaultRateRules {
		l.rules[path] = rule
	}
	return l
}()

// RegisterRateRule 注册/覆盖某个路径前缀的限流规则。
// 例如 Premium 订阅卖家可覆盖 "/v1/analytics/data" 解除每日 50 次上限：
//
//	ozonCore.RegisterRateRule("/v1/analytics/data", ozonCore.RateRule{MinInterval: 61 * time.Second})
func RegisterRateRule(path string, rule RateRule) {
	globalRateLimiter.mu.Lock()
	defer globalRateLimiter.mu.Unlock()
	globalRateLimiter.rules[path] = rule
}

// ResetRateRules 恢复为内置默认规则并清空配额状态（测试用）。
func ResetRateRules() {
	globalRateLimiter.mu.Lock()
	defer globalRateLimiter.mu.Unlock()

	globalRateLimiter.rules = map[string]RateRule{}
	for path, rule := range defaultRateRules {
		globalRateLimiter.rules[path] = rule
	}
	globalRateLimiter.state = map[string]*rateLimitEntry{}
}

// matchRule 按最长前缀匹配规则，返回匹配到的规则路径。
func (l *rateLimiter) matchRule(path string) (string, RateRule, bool) {
	l.mu.Lock()
	defer l.mu.Unlock()

	bestPath := ""
	bestRule := RateRule{}
	found := false
	for rulePath, rule := range l.rules {
		if strings.HasPrefix(path, rulePath) && len(rulePath) > len(bestPath) {
			bestPath = rulePath
			bestRule = rule
			found = true
		}
	}
	return bestPath, bestRule, found
}

// ruleFor 返回该路径适用的限流规则：
// 优先匹配已注册的具体规则；未配置的接口使用兜底规则（≤5 次/秒）。
// 返回的 keyPath 用于隔离配额：具体规则用规则路径，兜底规则用请求的完整路径。
func (l *rateLimiter) ruleFor(path string) (string, RateRule) {
	if rulePath, rule, ok := l.matchRule(path); ok {
		return rulePath, rule
	}
	return path, RateRule{MinInterval: FallbackMinInterval}
}

// acquire 获取一次请求配额：
//   - 距上次请求不足 MinInterval 时等待（可被 ctx 取消）；
//   - 超过 DailyLimit 时立即返回 ErrDailyRateLimitExceeded。
func (l *rateLimiter) acquire(ctx context.Context, key string, rule RateRule) error {
	for {
		now := time.Now()

		l.mu.Lock()
		entry, ok := l.state[key]
		if !ok {
			entry = &rateLimitEntry{}
			l.state[key] = entry
		}

		if rule.DailyLimit > 0 {
			if now.Sub(entry.dayStart) >= 24*time.Hour {
				entry.dayStart = now
				entry.dayCount = 0
			}
			if entry.dayCount >= rule.DailyLimit {
				l.mu.Unlock()
				return ErrDailyRateLimitExceeded
			}
		}

		wait := time.Duration(0)
		if rule.MinInterval > 0 && !entry.lastCall.IsZero() {
			if d := entry.lastCall.Add(rule.MinInterval).Sub(now); d > 0 {
				wait = d
			}
		}

		if wait <= 0 {
			if rule.DailyLimit > 0 {
				entry.dayCount++
			}
			entry.lastCall = now
			l.mu.Unlock()
			return nil
		}
		l.mu.Unlock()

		select {
		case <-time.After(wait):
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
