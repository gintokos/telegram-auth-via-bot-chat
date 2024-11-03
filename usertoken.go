package telegramauth

import (
	"sync"
	"time"
)

type userTokenStore struct {
	userTokens map[string]*userToken
	mu         sync.Mutex
	ttl        time.Duration
}

type userToken struct {
	expiredAt time.Time
}

func newTokenStore(ttl time.Duration) *userTokenStore {
	st := &userTokenStore{
		userTokens: make(map[string]*userToken, 10),
		ttl:        ttl,
	}
	st.startCleanupRoutine()

	return st
}

func (st *userTokenStore) SaveUserToken(utoken string) {
	st.mu.Lock()
	defer st.mu.Unlock()

	st.userTokens[utoken] = &userToken{
		expiredAt: time.Now().Add(st.ttl),
	}
}

func (st *userTokenStore) ValidateUserToken(utoken string) bool {
	st.mu.Lock()
	defer st.mu.Unlock()

	token, exists := st.userTokens[utoken]
	if !exists || time.Now().After(token.expiredAt) {
		return false
	}

	delete(st.userTokens, utoken)

	return true
}

func (st *userTokenStore) cleanupExpiredTokens() {
	st.mu.Lock()
	defer st.mu.Unlock()

	now := time.Now()
	for token, usertoken := range st.userTokens {
		if now.After(usertoken.expiredAt) {
			delete(st.userTokens, token)
		}
	}
}

func (st *userTokenStore) startCleanupRoutine() {
	ticker := time.NewTicker(st.ttl)
	go func() {
		for {
			<-ticker.C
			st.cleanupExpiredTokens()
		}
	}()
}
