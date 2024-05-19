package storage

import (
	"strings"
	"time"

	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

type data struct {
	Expiry time.Time
	Data   string
}

type Storage struct {
	m map[string]*data
}

var (
	rw = resp.NewWriter()
)

func NewStorage() *Storage {
	s := &Storage{
		m: make(map[string]*data),
	}
	return s
}

func (s *Storage) Set(k resp.Value, val resp.Value) {
	key := rw.Encode(k)
	value := rw.Encode(val)
	s.m[key] = &data{
		Data:   value,
		Expiry: time.Unix(1<<61-1, 0),
	}
}

func (s *Storage) SetWithExpiry(k resp.Value, val resp.Value, expiry time.Duration) {
	key := rw.Encode(k)
	value := rw.Encode(val)
	s.m[key] = &data{
		Data:   value,
		Expiry: time.Now().Add(expiry),
	}
}

func (s *Storage) Get(k resp.Value) *resp.Value {
	key := rw.Encode(k)
	dat := s.m[key]
	if dat == nil {
		return nil
	}

	if !dat.Expiry.After(time.Now()) {
		delete(s.m, key)
		return nil
	}
	output, _ := resp.NewReader(strings.NewReader(s.m[key].Data)).Read()
	return &output
}
