package logic

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (s *httpServer) initRouter() {
	router := gin.Default()

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": "pong",
		})
	})

	router.GET("/node-instance", s.nodeInstance)

	s.router = router
}

func (s *httpServer) nodeInstance(c *gin.Context) {
	addr, err := s.lh.NodeInstance(s.ctx)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"msg":  err.Error(),
			"data": "",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg":  "success",
		"data": addr,
	})
}
