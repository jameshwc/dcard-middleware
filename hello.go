package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func hello(c *gin.Context) {
	c.String(http.StatusOK, "Hello, world!")
}
