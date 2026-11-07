// Package shorturl 是进程内短链。短码是递增的 62 进制，方便测试看懂下一个码。
// 公开服务应改成随机或哈希，冲突时重试。递增计数会暴露生成顺序。
package shorturl

import (
	"errors"
	"sync"
	"time"
)

var ErrNotFound = errors.New("shorturl: not found")

type record struct {
	url     string
	expires time.Time
}

type Store struct {
	mu     sync.Mutex
	byCode map[string]record
	byURL  map[string]string
	ttl    time.Duration
	now    func() time.Time
	seq    uint64
}

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

func (s *Store) Put(raw string) (string, error) {
	if raw == "" {
		return "", errors.New("shorturl: empty url")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.dropExpired()
	if code, ok := s.byURL[raw]; ok {
		return code, nil
	}
	s.seq++
	code := encode(s.seq)
	rec := record{url: raw}
	if s.ttl > 0 {
		rec.expires = s.now().Add(s.ttl)
	}
	s.byCode[code] = rec
	s.byURL[raw] = code
	return code, nil
}

func (s *Store) Resolve(code string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.dropExpired()
	rec, ok := s.byCode[code]
	if !ok {
		return "", ErrNotFound
	}
	return rec.url, nil
}

func (s *Store) dropExpired() {
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
