package storage

import "time"

type data struct {
	Expiry time.Time
	Data   string
}

type Storage struct {
	m map[string]*data
}

func NewStorage() *Storage {
	s := &Storage{
		m: make(map[string]*data),
	}
	return s
}

func (s *Storage) Set(key string, value string) {
	s.m[key] = &data{
		Data:   value,
		Expiry: time.Unix(1<<63-1, 0),
	}
}

func (s *Storage) SetWithExpiry(key string, value string, expiry time.Duration) {
	s.m[key] = &data{
		Data:   value,
		Expiry: time.Now().Add(expiry),
	}
}

func (s *Storage) Get(key string) string {
	dat := s.m[key]
	if dat == nil {
		return ""
	}

	if !dat.Expiry.After(time.Now()) {
		delete(s.m, key)
		return ""
	}

	return s.m[key].Data
}
