package factory

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSessionMem(t *testing.T) {
	sess := Session(MemSessionNew())
	sid := sess.Create("hello world")
	v, err := sess.Get(sid)
	assert.Nil(t, err)
	assert.Equal(t, "hello world", v)
}

func TestSesstionRedis(t *testing.T) {
	sess := Session(RedisSessionNew())
	sid := sess.Create("hello world")
	v, err := sess.Get(sid)
	assert.Nil(t, err)
	assert.Equal(t, "hello world", v)
}

func TestSessionStore(t *testing.T) {
	err := Register("mem", MemSessionNew())
	assert.Nil(t, err)
	err = Register("redis", RedisSessionNew())
	assert.Nil(t, err)
	_, err = GetSessionStore("mongo")
	assert.NotNil(t, err)
	sess, err := GetSessionStore("mem")
	assert.Nil(t, err)

	sid := sess.Create("hello world")
	v, err := sess.Get(sid)
	assert.Nil(t, err)
	assert.Equal(t, "hello world", v)
}
