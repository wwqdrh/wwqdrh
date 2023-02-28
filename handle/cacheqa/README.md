在引入缓存组件的时候最常见需要处理的三个问题如下

## 缓存雪崩

大量key同时失效

解决方案：设置缓存时间为一定范围的随机数

```go
// 缓存时间
func (o *CatUser) RedisDuration() time.Duration {
	// 这个时候可以用随机时间 解决缓存雪崩问题
	// 设置 30 ~ 60 分钟  -- 这里记得不要设置  0 ~ n 时间，因为万一是 0 相当于没有设置
	return time.Duration((rand.Intn(60-30) + 30)) * time.Minute
}
```

## 缓存穿透

缓存和数据库中都不存在（请求数据没有被缓存拦截，一直都在找数据库，但是数据库没有，所以一直找）

解决方法：当第一次命中时设置 该缓存 value 为 DISABLE ，之后每次都只会打到该缓存上，或者使用布隆过滤器进行过滤

```go
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
```

## 缓存击穿

缓存失效后，有某些 key 被超高并发地访问

解决方法：使用互斥锁，有锁时，等待获取

```go
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
```
