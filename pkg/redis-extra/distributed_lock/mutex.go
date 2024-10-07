package distributed_lock

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"github.com/hashicorp/go-multierror"
	"golang.org/x/sync/singleflight"
	"time"
)

type Mutex struct {
	name         string
	expiration   time.Duration
	reTry        RetryStrategy
	quorum       int
	genValueFunc func() (string, error)
	value        string    // 锁持有者的标识符
	until        time.Time // 锁的过期时间
	factor       float64   // 飘逸因子，0~1，用于调整 计算锁过期时间

	g       singleflight.Group
	clients []*Client
}

// 检查 KEYS[1]（即 Redis 键）的值是否等于 ARGV[1]（提供的值）
// 如果相等，则更新键的过期时间为 ARGV[2]（以毫秒为单位）；否则返回 0
var touchScript = `
	if redis.call("GET", KEYS[1]) == ARGV[1] then
		return redis.call("pexpire", KEYS[1], ARGV[2])
	else
		return 0
	end
`

var deleteScript = `
	if redis.call("GET", KEYS[1]) == ARGV[1] then
		return redis.call("DEL", KEYS[1])
	else
		return 0
	end
`

func (m *Mutex) SingleflightLock(ctx context.Context, key string) error {
	for {
		flag := false
		result := m.g.DoChan(key, func() (interface{}, error) {
			flag = true
			err := m.Lock(ctx)
			if err != nil {
				return nil, err
			}
			return nil, nil
		})
		select {
		case res := <-result:
			if flag { // 排队加锁
				m.g.Forget(key)
				if res.Err != nil {
					return res.Err
				}
				return nil // 抢到走人
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (m *Mutex) Lock(ctx context.Context) error {
	value, err := genValue()
	if err != nil {
		return err
	}

	for interval, goon := m.reTry.Next(); goon; interval, goon = m.reTry.Next() {
		start := time.Now()

		// 申请加锁
		n, err := m.actOnClientsAsync(func(c *Client) (bool, error) {
			return m.acquire(ctx, c, value)
		})
		if n == 0 && err != nil {
			return err
		}

		now := time.Now()
		remainingExpiry := m.expiration - time.Since(start)
		drift := time.Duration(float64(m.expiration) * m.factor)
		until := now.Add(remainingExpiry - drift)
		if n >= m.quorum && now.Before(until) { // 加锁成功
			m.value = value
			m.until = until
			return nil
		}

		// 失败
		m.actOnClientsAsync(func(c *Client) (bool, error) {
			return m.release(ctx, c, value)
		})
		time.Sleep(interval)
	}
	return ErrFailedToGetLock
}

func (m *Mutex) Unlock(ctx context.Context) (bool, error) {
	n, err := m.actOnClientsAsync(func(c *Client) (bool, error) {
		return m.release(ctx, c, m.value)
	})
	return n >= m.quorum, err
}

// Extend 尝试延长 Mutex（分布式锁）的过期时间
func (m *Mutex) Extend(ctx context.Context) (bool, error) {
	n, err := m.actOnClientsAsync(func(c *Client) (bool, error) {
		return m.touch(ctx, c, m.name, int(m.expiration/time.Millisecond))
	})
	if n < m.quorum {
		return false, err
	}
	return true, err
}

func (m *Mutex) touch(ctx context.Context, c *Client, value string, expiration int) (bool, error) {
	status, err := c.client.Eval(ctx, touchScript, []string{m.name}, value, expiration).Result()

	return err == nil && status.(int64) != 0, err
}

func (m *Mutex) acquire(ctx context.Context, client *Client, value string) (bool, error) {
	reply, err := client.client.SetNX(ctx, m.name, value, m.expiration).Result()
	if err != nil {
		return false, err
	}
	return reply, nil
}

func (m *Mutex) release(ctx context.Context, c *Client, value string) (bool, error) {
	result, err := c.client.Eval(ctx, deleteScript, []string{m.name}, value).Result()

	return err == nil && result.(int64) != 0, err
}

// 多个 Redis 实例或集群进行并发锁操作，如释放或加锁等。
func (m *Mutex) actOnClientsAsync(actFn func(*Client) (bool, error)) (int, error) {
	type res struct {
		Status bool
		err    error
	}

	ch := make(chan res)
	for _, c := range m.clients {
		go func(c *Client) {
			status, err := actFn(c)
			ch <- res{
				Status: status,
				err:    err,
			}
		}(c)
	}

	n := 0
	var err error
	for range m.clients {
		r := <-ch
		if r.Status {
			n++
		} else if r.err != nil {
			err = multierror.Append(err, r.err)
		}
	}
	return n, err
}

// 可替换
// 生成一个长度为 16 字节的随机字节数组，并将其编码为 Base64 字符串
func genValue() (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(b), nil
}
