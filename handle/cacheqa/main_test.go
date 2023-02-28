package main

import (
	"context"
	"fmt"
	"testing"

	"github.com/go-redis/redis/v8"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitDb() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("./testdata/test.db"))
	if err != nil {
		return nil, err
	}
	db.AutoMigrate(&CatUser{})
	return db, nil
}

func InitRedis() *redis.Conn {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "redis6:6379",
		Password: "123456", // no password set
		DB:       0,        // use default DB
	})
	return rdb.Conn(context.TODO())
}

func TestCatUser(t *testing.T) {
	db, err := InitDb()
	if err != nil {
		t.Error(err)
		return
	}

	conn := InitRedis()

	t.Run("single", func(t *testing.T) {
		var cu CatUser
		engine := db.Model(&CatUser{}).Where("user_id=?", 1)
		cu.UserId = 1
		cu.DeleteFromRedis(conn)
		if err := cu.MustGet(engine, conn); err != nil {
			fmt.Printf("%+v", err)
			panic(err)
		}
		fmt.Println(cu)
	})

	t.Run("list", func(t *testing.T) {
		var cu CatUser
		engine := db.Model(&CatUser{})
		list, err := cu.ArrayMustGet(engine, conn)
		if err != nil {
			panic(err)
		}
		fmt.Println(list)
	})
}
