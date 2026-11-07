// Package shorturl 是一个进程内短链表。
// 短码是递增计数的 62 进制，方便测试看懂。线上通常改成随机或哈希，并在冲突时重试。
package shorturl

import (
	"errors"
	"sync"
	"time"
)

// ErrNotFound 表示短码不存在或已经过期。
var ErrNotFound = errors.New("shorturl: not found")

type record struct {
	url     string
	expires time.Time
}

// Store 保存短码。ttl <= 0 表示不过期。now 由调用方注入。
type Store struct {
	mu     sync.Mutex
	byCode map[string]record
	byURL  map[string]string
	ttl    time.Duration
	now    func() time.Time
	seq    uint64
}

// New 返回空表。
func New(ttl time.Duration, now func() time.Time) *Store {
	if now == nil {
		now = time.Now
	}
	return &Store{
		byCode: map[string]record{},
		byURL:  map[string]string{},
		ttl:    ttl,
		now:    now,
	}
}

// Put 为 raw 分配短码。未过期的同一个 URL 返回原来的短码。
func (s *Store) Put(raw string) (string, error) {
	if raw == "" {
		return "", errors.New("shorturl: empty url")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.dropExpiredLocked()
	if code, ok := s.byURL[raw]; ok {
		return code, nil
	}
	s.seq++
	code := encode(s.seq)
	for _, exists := s.byCode[code]; exists; _, exists = s.byCode[code] {
		s.seq++
		code = encode(s.seq)
	}
	rec := record{url: raw}
	if s.ttl > 0 {
		rec.expires = s.now().Add(s.ttl)
	}
	s.byCode[code] = rec
	s.byURL[raw] = code
	return code, nil
}

// Resolve 用短码换回原始 URL。
func (s *Store) Resolve(code string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.dropExpiredLocked()
	rec, ok := s.byCode[code]
	if !ok {
		return "", ErrNotFound
	}
	return rec.url, nil
}

func (s *Store) dropExpiredLocked() {
	if s.ttl <= 0 {
		return
	}
	now := s.now()
	for code, rec := range s.byCode {
		if !rec.expires.IsZero() && !now.Before(rec.expires) {
			delete(s.byURL, rec.url)
			delete(s.byCode, code)
		}
	}
}

const alphabet = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func encode(n uint64) string {
	if n == 0 {
		return "0"
	}
	var buf [16]byte
	i := len(buf)
	for n > 0 {
		i--
		buf[i] = alphabet[n%62]
		n /= 62
	}
	return string(buf[i:])
}
