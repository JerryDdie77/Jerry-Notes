package httpserver

import "github.com/gin-gonic/gin"

func getUserID(c *gin.Context) (int64, bool) {
	v, ok := c.Get("userID")
	if !ok {
		return 0, false
	}

	id, ok := v.(int64)
	return id, ok
}
