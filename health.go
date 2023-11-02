package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func health(c *gin.Context) {
	c.JSON(http.StatusOK, "Service is UP!")
}
