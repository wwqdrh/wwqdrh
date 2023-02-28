package main

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/catbugdemo/errors"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type CatUser struct {
	ID        uint      `gorm:"column:id" json:"id"`
	CreatedAt time.Time `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updatedAt"`
	UserId    int       `gorm:"column:user_id" json:"user_id"`
	Name      string    `gorm:"column:name" json:"name"`
	Password  string    `gorm:"cloumn:password" json:"password"`
}

// RedisKey 缓存key
func (o *CatUser) RedisKey() string {
	// 这里推荐: 1.用:分隔 1.如果有能够识别唯一标识的 id ,用它 -- (用id也行)
	// user_id 能唯一标识该数据 -- 同 id 类似
	return fmt.Sprintf("cat_user:%d", o.UserId)
}

func (o *CatUser) ArrayRedisKey() string {
	return fmt.Sprintf("cat_user")
}

// 缓存时间
func (o *CatUser) RedisDuration() time.Duration {
	// 这个时候可以用随机时间 解决缓存雪崩问题
	// 设置 30 ~ 60 分钟  -- 这里记得不要设置  0 ~ n 时间，因为万一是 0 相当于没有设置
	return time.Duration((rand.Intn(60-30) + 30)) * time.Minute
}

// SyncToRedis 添加缓存
// 使用 序列化后的 string 类型存储 缓存
func (o *CatUser) SyncToRedis(conn *redis.Conn) error {
	if o.RedisKey() == "" {
		return errors.New("not set redis key")
	}
	buf, err := json.Marshal(o)
	if err != nil {
		return errors.WithStack(err)
	}
	if err = conn.SetEX(context.Background(), o.RedisKey(), string(buf), o.RedisDuration()).Err(); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// 获取缓存
// 1.判断是否存在 key
// 2.获取是否为空
// 3.判断是否缓存穿透
// 3.获取后反序列化

// GetFromRedis 获取缓存
func (o *CatUser) GetFromRedis(conn *redis.Conn) error {
	if o.RedisKey() == "" {
		return errors.New("not set redis key")
	}
	buf, err := conn.Get(context.Background(), o.RedisKey()).Bytes()
	if err != nil {
		if err == redis.Nil {
			return redis.Nil
		}
		return errors.WithStack(err)
	}
	// 是否出现过缓存穿透
	if string(buf) == "DISABLE" {
		return errors.New("not found data in redis nor db")
	}

	if err = json.Unmarshal(buf, o); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (o *CatUser) DeleteFromRedis(conn *redis.Conn) error {
	if o.RedisKey() != "" {
		if err := conn.Del(context.Background(), o.RedisKey()).Err(); err != nil {
			return errors.WithStack(err)
		}
	}
	// 同时删除数组缓存
	if o.ArrayRedisKey() != "" {
		if err := conn.Del(context.Background(), o.ArrayRedisKey()).Err(); err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}

// MustGet 获取数据
// 1.先从缓存中获取
// 2.如果没找到 --找数据库 (也没找到--设置DISABLE 防止缓存穿透)
func (o *CatUser) MustGet(engine *gorm.DB, conn *redis.Conn) error {
	err := o.GetFromRedis(conn)
	// 如果为空证明找到了，提前返回不考虑后续操作
	if err == nil {
		return nil
	}

	if err != nil && err != redis.Nil {
		return errors.WithStack(err)
	}
	// 在缓存中没有找到这条数据，则从数据库中找
	var count int64
	if err = engine.Count(&count).Error; err != nil {
		return errors.WithStack(err)
	}
	// 如果 为 count =0  设置 DISABLE 防止缓存穿透
	if count == 0 {
		if err = conn.SetNX(context.Background(), o.RedisKey(), "DISABLE", o.RedisDuration()).Err(); err != nil {
			return errors.WithStack(err)
		}
		return errors.New("not found data in redis nor db")
	}

	// 这个时候找到了 -- 并且数据库中存在数据 --加锁防止缓存击穿
	// 设置 5 秒的互斥锁锁时间
	var mutex = o.RedisKey() + "_MUTEX"
	if err = conn.Get(context.Background(), mutex).Err(); err != nil {
		// 非 缓存为空 异常错误，提前报错
		if err != redis.Nil {
			return errors.WithStack(err)
		}
		// err == redis.Nil
		// 设置 5 s 的互斥锁时间
		if err = conn.SetNX(context.Background(), mutex, 1, 3*time.Second).Err(); err != nil {
			return errors.WithStack(err)
		}
		// 从数据库中查找
		if err = engine.First(&o).Error; err != nil {
			return errors.WithStack(err)
		}
		// 同步缓存
		if err = o.SyncToRedis(conn); err != nil {
			return errors.WithStack(err)
		}
		// 删除锁
		if err = conn.Del(context.Background(), mutex).Err(); err != nil {
			return errors.WithStack(err)
		}
	} else {
		// 这个时候不为空,加了锁 -- 进行循环等等待
		var index int
		for {
			if index > 10 {
				return errors.New(mutex + " lock error")
			}
			if err2 := conn.Get(context.Background(), mutex).Err(); err2 != nil {
				break
			} else {
				time.Sleep(30 * time.Millisecond)
				index++
				continue
			}
		}
		if err = o.MustGet(engine, conn); err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}

func (o *CatUser) ArraySyncToRedis(list []CatUser, conn *redis.Conn) error {
	if o.ArrayRedisKey() == "" {
		return errors.New("not set redis key")
	}
	buf, err := json.Marshal(list)
	if err != nil {
		return errors.WithStack(err)
	}
	if err = conn.SetEX(context.Background(), o.ArrayRedisKey(), string(buf), o.RedisDuration()).Err(); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (o *CatUser) ArrayGetFromRedis(conn *redis.Conn) ([]CatUser, error) {
	if o.RedisKey() == "" {
		return nil, errors.New("not set redis key")
	}
	buf, err := conn.Get(context.Background(), o.ArrayRedisKey()).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, redis.Nil
		}
		return nil, errors.WithStack(err)
	}
	if string(buf) == "DISABLE" {
		return nil, errors.New("not found data in redis nor db")
	}

	var list []CatUser
	if err = json.Unmarshal(buf, &list); err != nil {
		return nil, errors.WithStack(err)
	}
	return list, nil
}

func (o *CatUser) ArrayDeleteFromRedis(conn *redis.Conn) error {
	return o.DeleteFromRedis(conn)
}

// ArrayMustGet
func (o *CatUser) ArrayMustGet(engine *gorm.DB, conn *redis.Conn) ([]CatUser, error) {
	list, err := o.ArrayGetFromRedis(conn)
	if err == nil {
		return list, nil
	}
	if err != nil && err != redis.Nil {
		return nil, errors.WithStack(err)
	}

	// not found in redis
	var count int64
	if err = engine.Count(&count).Error; err != nil {
		return nil, errors.WithStack(err)
	}
	if count == 0 {
		if err = conn.SetNX(context.Background(), o.ArrayRedisKey(), "DISABLE", o.RedisDuration()).Err(); err != nil {
			return nil, errors.WithStack(err)
		}
		return nil, errors.New("not found data in redis nor db")
	}

	var mutex = o.ArrayRedisKey() + "_MUTEX"
	if err = conn.Get(context.Background(), mutex).Err(); err != nil {
		if err != redis.Nil {
			return nil, errors.WithStack(err)
		}
		// err = redis.Nil
		if err = conn.SetNX(context.Background(), mutex, 1, 3*time.Second).Err(); err != nil {
			return nil, errors.WithStack(err)
		}
		if err = engine.Find(&list).Error; err != nil {
			return nil, errors.WithStack(err)
		}
		if err = o.ArraySyncToRedis(list, conn); err != nil {
			return nil, errors.WithStack(err)
		}
		if err = conn.Del(context.Background(), mutex).Err(); err != nil {
			return nil, errors.WithStack(err)
		}
	} else {
		var index int
		for {
			if index > 10 {
				return nil, errors.New(mutex + " lock error")
			}
			if err2 := conn.Get(context.Background(), mutex).Err(); err2 != nil {
				break
			} else {
				time.Sleep(50 * time.Millisecond)
				index++
				continue
			}
		}
		list, err = o.ArrayMustGet(engine, conn)
		if err != nil {
			return nil, errors.WithStack(err)
		}
	}
	return list, nil
}
