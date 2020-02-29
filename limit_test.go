package main

import (
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/v7"
)

func Test_limitVisit(t *testing.T) {
	var redisTestServer redisServer
	s, err := miniredis.Run()
	if err != nil {
		log.Fatal(err)
	}
	cli := redis.NewClient(&redis.Options{
		Addr: s.Addr(),
	})
	redisTestServer.client = cli
	redisTestServer.maxIP = 10
	redisTestServer.timeout = 100 * time.Second
	ts := httptest.NewServer(limitVisit(hello, &redisTestServer))
	defer ts.Close()
	// test if the amount of visit is limited as expected
	for i := 0; i < redisTestServer.maxIP*2; i++ {
		res, err := http.Get(ts.URL)
		if err != nil {
			t.Fail()
		}
		if i < redisTestServer.maxIP && res.StatusCode != 200 {
			t.Fatal("before reach maximum visited", i, res.StatusCode, res.Header)
		} else if i >= redisTestServer.maxIP && res.StatusCode != 429 {
			t.Fatal("not sending 429 code correctly")
		}
	}
	// test if ttl works as expected
	s.FastForward(redisTestServer.timeout)
	res, _ := http.Get(ts.URL)
	val, err := strconv.Atoi(res.Header.Get("X-RateLimit-Remaining"))
	if err != nil {
		t.Fatal("X-RateLimit-Remaining is not a number!")
	}
	if val+1 != redisTestServer.maxIP {
		t.Fatal("TTL doesn't work as expected")
	}
}
