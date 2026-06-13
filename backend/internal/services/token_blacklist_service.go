package services

import (
	"context"
	"sync"
	"time"
)

type TokenBlacklist interface {
	Revoke(tokenID string, expiresAt time.Time)
	IsRevoked(tokenID string) bool
	StartCleanup(ctx context.Context, interval time.Duration)
}

type inMemoryTokenBlacklist struct {
	mu            sync.RWMutex
	revokedTokens map[string]time.Time
}

func NewInMemoryTokenBlacklist() TokenBlacklist {
	return &inMemoryTokenBlacklist{
		revokedTokens: make(map[string]time.Time),
	}
}

func (b *inMemoryTokenBlacklist) Revoke(tokenID string, expiresAt time.Time) {
	if tokenID == "" {
		return
	}

	b.mu.Lock()
	defer b.mu.Unlock()

	b.revokedTokens[tokenID] = expiresAt
}

func (b *inMemoryTokenBlacklist) IsRevoked(tokenID string) bool {
	if tokenID == "" {
		return false
	}

	b.mu.RLock()
	expiresAt, exists := b.revokedTokens[tokenID]
	b.mu.RUnlock()

	if !exists {
		return false
	}

	if time.Now().After(expiresAt) {
		b.mu.Lock()
		delete(b.revokedTokens, tokenID)
		b.mu.Unlock()

		return false
	}

	return true
}

func (b *inMemoryTokenBlacklist) StartCleanup(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)

	go func() {
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return

			case <-ticker.C:
				b.removeExpiredTokens()
			}
		}
	}()
}

func (b *inMemoryTokenBlacklist) removeExpiredTokens() {
	now := time.Now()

	b.mu.Lock()
	defer b.mu.Unlock()

	for tokenID, expiresAt := range b.revokedTokens {
		if now.After(expiresAt) {
			delete(b.revokedTokens, tokenID)
		}
	}
}
