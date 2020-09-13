package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	redis "github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
)

// Todo represents an todo model.
type Todo struct {
	ID        int       `json:"id,omitempty"`
	Name      string    `json:"name,omitempty"`
	Done      byte      `json:"done,omitempty"`
	CreatedAt time.Time `json:"createdAt,omitempty"`
	UpdatedAt time.Time `json:"updatedAt,omitempty"`
}

func newMysqlDB(dbName string) *sql.DB {
	usr := os.Getenv("MYSQL_APP_USER")
	pwd := os.Getenv("MYSQL_APP_PASSWORD")
	connStr := fmt.Sprintf(
		"%v:%v@tcp(%v)/%v?charset=utf8mb4&parseTime=true&loc=Local",
		usr,
		pwd,
		"mysql",
		dbName,
	)
	db, err := sql.Open("mysql", connStr)
	if err != nil {
		log.Fatal("fail to connect to mysql")
	}
	return db
}

func newRedis() *redis.Client {
	addr := os.Getenv("REDIS_ADDR")
	if len(addr) < 1 {
		log.Printf("empty redis address. set %v env value", "REDIS_ADDR")
	}

	password := os.Getenv("REDIS_PASSWORD")
	c := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0,
	})

	if c == nil {
		log.Fatal("fail to get a redis client")
	}
	_, err := c.Ping(context.Background()).Result()
	if err != nil {
		log.Fatal("Ping to redis server failed")
	}
	return c
}

func main() {
	db := newMysqlDB("test")
	rc := newRedis()
	r := gin.Default()

	r.GET("/todos", func(c *gin.Context) {
		qr := c.Query("cache")
		if qr == "true" {
			val, err := rc.Get(context.Background(), "todos").Result()
			if err != nil && err != redis.Nil {
				log.Printf("get error from redis. %v, %v", err, val)
				c.Status(http.StatusInternalServerError)
				return
			}
			var todos []Todo
			err = json.Unmarshal([]byte(val), &todos)
			if err != nil {
				log.Printf("unmarshal error. %v", err)
				c.Status(http.StatusInternalServerError)
				return
			}
			c.JSON(http.StatusOK, gin.H{
				"todos": todos,
			})
			return
		}
		rows, err := db.Query("SELECT id, name, done, created_at, updated_at FROM todos LIMIT 10")
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}

		var (
			id        int
			name      string
			done      byte
			createdAt time.Time
			updatedAt time.Time
			todos     []Todo
		)

		for rows.Next() {
			err := rows.Scan(&id, &name, &done, &createdAt, &updatedAt)
			if err != nil {
				log.Printf("scan error - %v", err)
				continue
			}
			todos = append(todos, Todo{ID: id, Name: name, Done: done, CreatedAt: createdAt, UpdatedAt: updatedAt})
		}
		_, err = rc.Get(context.Background(), "todos").Result()
		if err == redis.Nil {
			log.Println("empty value in redis")
			serialized, _ := json.Marshal(todos)
			_, err = rc.Set(context.Background(), "todos", serialized, 0).Result()
			if err != nil {
				log.Printf("set error. %v", err)
			}
		} else {
			log.Printf("fail to get value from redis. %v", err)
		}
		c.JSON(http.StatusOK, gin.H{
			"todos": todos,
		})
	})

	r.Run()
}
