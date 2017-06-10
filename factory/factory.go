package factory

import (
	"errors"
	redis "gopkg.in/redis.v3"
	"math/rand"
	"time"
)

var (
	SessionFoundError = errors.New("Session not found")
	letterRunes       = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
)

func RandStringRunes(n int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

type Session interface {
	Create(value string) string
	Get(id string) (string, error)
	Delete(id string) error
}

type MemSession struct {
	DB map[string]string
}

func MemSessionNew() *MemSession {
	sess := &MemSession{DB: make(map[string]string)}
	return sess
}
func (m *MemSession) Create(v string) string {
	randStr := RandStringRunes(32)
	m.DB[randStr] = v
	return randStr
}
func (m *MemSession) Get(id string) (string, error) {
	if v, ok := m.DB[id]; ok {
		return v, nil
	}
	return "", SessionFoundError
}

func (m *MemSession) Delete(id string) error {
	if _, ok := m.DB[id]; ok {
		delete(m.DB, id)
		return nil
	}
	return SessionFoundError
}

func RedisSessionNew() *RedisSession {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	_, err := client.Ping().Result()
	if err != nil {
		panic(err)
	}
	return &RedisSession{cli: client}
}

type RedisSession struct {
	cli *redis.Client
}

func (r *RedisSession) Create(v string) string {
	randStr := RandStringRunes(32)
	r.cli.Set(randStr, v, 0)
	return randStr
}

func (r *RedisSession) Get(id string) (string, error) {
	v, err := r.cli.Get(id).Result()
	if err != nil {
		return "", err
	}
	return v, nil
}

func (r *RedisSession) Delete(id string) error {
	return r.cli.Del(id).Err()
}

var SessionStore = make(map[string]Session)

func Register(name string, sess Session) error {
	if sess == nil {
		return errors.New("Session not found")
	}
	_, ok := SessionStore[name]
	if ok {
		return errors.New("Session exists")
	}
	SessionStore[name] = sess
	return nil
}

func GetSessionStore(name string) (Session, error) {
	sess, ok := SessionStore[name]
	if !ok {
		return nil, errors.New("Session store not found :" + name)
	}
	return sess, nil
}
