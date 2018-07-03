package util

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

func QueryInteger(c *gin.Context, name string) (int, error) {
	v := c.Query(name)
	if v == "" {
		return -1, fmt.Errorf("not set")
	}

	n, err := strconv.Atoi(v)
	if err != nil {
		return -1, err
	}

	return n, nil
}
