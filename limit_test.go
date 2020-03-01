package main

import (
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/gin-gonic/gin"
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
	gin.SetMode(gin.TestMode)
	router := setupRouter(&redisTestServer)
	getResponse := func(router *gin.Engine) *httptest.ResponseRecorder {
		w := httptest.NewRecorder()
		r, _ := http.NewRequest("GET", "/", nil)
		router.ServeHTTP(w, r)
		return w
	}
	// test if the amount of visit is limited as expected
	// router.Run()
	for i := 0; i < redisTestServer.maxIP*2; i++ {
		w := getResponse(router)
		if i < redisTestServer.maxIP && w.Code != 200 {
			t.Fatal("before reach maximum visited", i, w.Code, w.Result().Header)
		} else if i >= redisTestServer.maxIP && w.Code != 429 {
			t.Fatal("not sending 429 code correctly")
		}
	}
	// test if ttl works as expected
	s.FastForward(redisTestServer.timeout)
	w := getResponse(router)
	val, err := strconv.Atoi(w.Result().Header.Get("X-RateLimit-Remaining"))
	if err != nil {
		t.Fatal("X-RateLimit-Remaining is not a number")
	}
	if val+1 != redisTestServer.maxIP {
		t.Fatal("TTL doesn't work as expected")
	}
}
