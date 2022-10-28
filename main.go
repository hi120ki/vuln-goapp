package main

import (
	"net/http"
	"os/exec"

	"github.com/gin-gonic/gin"
	"github.com/mattn/go-shellwords"
)

type Argument struct {
	Arg string `json:"arg" binding:"required"`
}

func main() {
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.POST("/ace", func(c *gin.Context) {
		var form Argument
		if err := c.Bind(&form); err != nil {
			c.JSON(400, gin.H{
				"error": "failed to parse param",
			})
			return
		}
		cmd, err := shellwords.Parse(form.Arg)
		if err != nil {
			c.JSON(500, gin.H{
				"arg":   form.Arg,
				"error": err,
			})
			return
		}
		switch len(cmd) {
		case 0:
			c.JSON(400, gin.H{
				"arg":   form.Arg,
				"error": "command length is 0",
			})
		case 1:
			out, err := exec.Command(cmd[0]).Output()
			if err != nil {
				c.JSON(500, gin.H{
					"arg":   form.Arg,
					"error": err,
				})
				return
			}
			c.JSON(200, gin.H{
				"arg":    form.Arg,
				"result": string(out),
			})
		default:
			out, err := exec.Command(cmd[0], cmd[1:]...).Output()
			if err != nil {
				c.JSON(500, gin.H{
					"arg":   form.Arg,
					"error": err,
				})
				return
			}
			c.JSON(200, gin.H{
				"arg":    form.Arg,
				"result": string(out),
			})
		}
	})

	if err := r.Run(); err != nil {
		panic(err)
	}
}
