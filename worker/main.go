package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gammazero/workerpool"
	"github.com/garyburd/redigo/redis"
)

var redisPool *redis.Pool

var (
	CommandPing      = "PING"
	CommandAuth      = "AUTH"
	CommandPublish   = "PUBLISH"
	CommandIncrement = "INCR"

	RedisURL = "redis:6379"
	// RedisPass = "nakama"
)

func Work() string {
	respRNG, err := http.Get("http://rng:11991/rng?max=1000000")
	if err != nil {
		log.Println(err)
		return ""
	}

	body, err := ioutil.ReadAll(respRNG.Body)
	if err != nil {
		log.Println(err)
		return ""
	}

	type tempBody struct {
		Rand int64 `json:"rand"`
	}
	var tb tempBody
	if err := json.Unmarshal(body, &tb); err != nil {
		log.Println(err)
		return ""
	}

	// count as failure
	if tb.Rand < 100000 {
		return ""
	}

	responseBody := bytes.NewBuffer([]byte(fmt.Sprintf("%d", tb.Rand)))
	respHash, err := http.Post("http://hash:11992/hash", "", responseBody)
	if err != nil {
		log.Println(err)
		return ""
	}

	bodyHash, err := ioutil.ReadAll(respHash.Body)
	if err != nil {
		log.Println(err)
		return ""
	}

	type tempBodyHash struct {
		Hash string `json:"hash"`
	}
	var tbh tempBodyHash
	if err := json.Unmarshal(bodyHash, &tbh); err != nil {
		log.Println(err)
		return ""
	}

	redisconn2 := redisPool.Get()
	rep, err := redisconn2.Do("HSET", "wallet", tbh.Hash, tb.Rand)
	log.Println(rep, err)
	redisconn2.Close()

	return tbh.Hash
}

func main() {
	countHash := 0

	redisPool = &redis.Pool{
		Dial: func() (conn redis.Conn, e error) {
			c, err := redis.Dial("tcp", RedisURL)
			if err != nil {
				return nil, err
			}

			// _, err = c.Do(CommandAuth, RedisPass)
			// if err != nil {
			// 	return nil, err
			// }

			return c, nil
		},
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
	}

	conn := redisPool.Get()
	defer conn.Close()

	_, err := conn.Do(CommandPing)
	if err != nil {
		log.Fatalln(err)
	}

	wp := workerpool.New(5)
	redisconn := redisPool.Get()
	defer redisconn.Close()

	for i := 0; i < 999999; i++ {
		// let worker work
		wp.Submit(func() {
			Work()
			countHash++
		})

		rep, err := redisconn.Do("INCRBY", "hashes", countHash)
		log.Println(rep, err)
		time.Sleep(time.Microsecond * time.Duration(rand.Intn(500)))
	}
	wp.StopWait()

	log.Println("HERE DONE")
}
