package main

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func GetIndex(c *gin.Context) {
	entries, _ := getStatusAll()

	flashMessage, _ := c.Cookie("flash-message")
	flashMessageType, _ := c.Cookie("flash-message-type")

	if flashMessage != "" || flashMessageType != "" {
		c.SetCookie("flash-message", "", -1440, "/", "", false, false)
		c.SetCookie("flash-message-type", "", -1440, "/", "", false, false)
	}
	if flashMessage != "" && flashMessageType == "" {
		flashMessageType = "info"
	}

	c.HTML(http.StatusOK, "index.go.html", gin.H{
		"entries":          entries,
		"flashMessage":     flashMessage,
		"flashMessageType": flashMessageType,
	})
}

func GetStd(pipe string) gin.HandlerFunc {
	return func(c *gin.Context) {
		program := c.Param("program")

		numBytesString := c.Query("size")
		var numBytes uint64 = 50000
		if numBytesString != "" {
			numBytesConv, _ := strconv.ParseUint(numBytesString, 10, 64)
			if numBytesConv > 0 {
				numBytes = numBytesConv
			}
		}

		std, _, _ := getTailRaw(program, pipe, numBytes)
		c.HTML(http.StatusOK, "std.go.html", gin.H{
			"pipe":     pipe,
			"std":      std,
			"program":  program,
			"numBytes": numBytes,
		})
	}
}
func ManageAction(action string) gin.HandlerFunc {
	return func(c *gin.Context) {
		program := c.PostForm("program")
		if program == "" {
			c.String(http.StatusBadRequest, "Missing parameter: program")
			return
		}
		stdout, stderr, err := manageProcess(program, action)

		stderr = strings.TrimSpace(stderr)
		if stderr != "" {
			c.SetCookie("flash-message", stderr, 0, "/", "", false, false)
			c.SetCookie("flash-message-type", "error", 0, "/", "", false, false)
			c.Redirect(http.StatusFound, "/")
			return
		}

		if err != nil {
			c.SetCookie("flash-message", err.Error(), 0, "/", "", false, false)
			c.SetCookie("flash-message-type", "error", 0, "/", "", false, false)
			c.Redirect(http.StatusFound, "/")
			return
		}

		stdout = strings.TrimSpace(stdout)
		if stdout != "" {
			c.SetCookie("flash-message", stdout, 0, "/", "", false, false)
			c.SetCookie("flash-message-type", "success", 0, "/", "", false, false)
			c.Redirect(http.StatusFound, "/")
			return
		}

		c.String(http.StatusOK, "Success")
	}
}
