package filter

import (
	"github.com/spaolacci/murmur3"
	"github.com/gomodule/redigo/redis"
	"time"
	"github.com/Soul-Mate/gospider/util"
)

type RedisFilter struct {
	Pool    *redis.Pool
	bitKey  string
	hashNum uint
	offset  uint
}

const MAX_SIZE uint = 4294967296 - 1

func NewRedisFilter() *RedisFilter {
	rf := new(RedisFilter)
	rf.Pool = util.GetSharedPool()
	rf.hashNum = 5
	rf.bitKey = "redis_bloom_filter"
	rf.offset = MAX_SIZE
	if !rf.existBitSet() {
		if err := rf.newBitSet(); err != nil {
			panic(err)
		}
	}
	return rf
}

func (rf *RedisFilter) Add(data string) {
	rf.AddByte([]byte(data))
}

func (rf *RedisFilter) AddByte(data []byte) {
	h := hashes(data)
	for i := uint(0); i < rf.hashNum; i++ {
		rf.setbit(location(h, i, rf.offset))
	}
}

func (rf *RedisFilter) Contains(data string) bool {
	return rf.ContainsByte([]byte(data))
}

func (rf *RedisFilter) ContainsByte(data []byte) bool {
	h := hashes(data)
	for i := uint(0); i < rf.hashNum; i++ {
		if !rf.getbit(location(h, i, rf.offset)) {
			return false
		}
	}
	return true
}

func (rf *RedisFilter) newBitSet() error {
	var (
		cerr chan error
		t    <-chan time.Time
	)
	cerr = make(chan error)
	t = time.Tick(time.Second * 5)
	go func() {
		conn := rf.Pool.Get()
		defer conn.Close()
		_, err := conn.Do("SETBIT", rf.bitKey, rf.offset, 0)
		cerr <- err
	}()
	select {
	case err := <-cerr:
		return err
	case <-t:
		return nil
	}
}

func (rf *RedisFilter) existBitSet() bool {
	conn := rf.Pool.Get()
	defer conn.Close()

	reply, err := conn.Do("EXISTS", rf.bitKey)
	if err != nil {
		return false
	}
	if v, ok := reply.(int64); ok {
		if v != 1 {
			return false
		}
	}
	return true
}

func (rf *RedisFilter) setbit(location uint64) error {
	conn := rf.Pool.Get()
	defer conn.Close()
	_, err := conn.Do("SETBIT", rf.bitKey, location, 1)
	return err
}

func (rf *RedisFilter) getbit(location uint64) bool {
	conn := rf.Pool.Get()
	defer conn.Close()
	reply, err := conn.Do("GETBIT", rf.bitKey, location)
	if err != nil {
		return false
	}
	if v, ok := reply.(int64); ok {
		if v == 1 {
			return true
		} else {
			return false
		}
	}
	return false
}

func hashes(data []byte) [4]uint64 {
	a1 := []byte{1} // to grab another bit of data
	hasher := murmur3.New128()
	hasher.Write(data) // #nosec
	v1, v2 := hasher.Sum128()
	hasher.Write(a1) // #nosec
	v3, v4 := hasher.Sum128()
	return [4]uint64{
		v1, v2, v3, v4,
	}
}

func location(h [4]uint64, i uint, length uint) uint64 {
	ii := uint64(i)
	return (h[ii%2] + ii*h[2+(((ii+(ii%2))%4)/2)]) % uint64(length)
}
