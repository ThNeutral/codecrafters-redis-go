package storage

type Storage struct {
	m map[string]string
}

func NewStorage() *Storage {
	s := &Storage{
		m: make(map[string]string),
	}
	return s
}

func (s *Storage) Set(key string, value string) {
	s.m[key] = value
}

func (s *Storage) Get(key string) string {
	return s.m[key]
}
