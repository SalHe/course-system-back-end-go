package middleware

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
	"time"
)

func RecoveryWithZerolog(logger *zerolog.Logger, stack bool, goLogStack bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Check for a broken connection, as it is not really a
				// condition that warrants a panic stack trace.
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				if brokenPipe {
					event := logger.Error().
						Str("request", string(httpRequest))
					if errErr, ok := err.(error); ok {
						event = event.Err(errErr)
					} else {
						event = event.Err(errors.New(fmt.Sprint(err)))
					}
					event.Msg(c.Request.URL.Path)

					// If the connection is dead, we can't write a status to it.
					c.Error(err.(error)) // nolint: errcheck
					c.Abort()
					return
				}

				if goLogStack {
					log.Printf("[Recovery] %s panic recovered:\n%s\n%s\n",
						c.Request.RemoteAddr, string(httpRequest), debug.Stack())
				}

				event := logger.Error().
					Time("time", time.Now()).
					Str("request", string(httpRequest)).
					Interface("error", err)
				if stack {
					event = event.Str("stack", string(debug.Stack()))
				}
				event.Msg("[Recovery from panic]")
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}

func LoggerWithZerolog(logger *zerolog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		latency := time.Since(start)

		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()

		comment := c.Errors.ByType(gin.ErrorTypePrivate).String()

		if raw != "" {
			path = path + "?" + raw
		}

		event := logger.Info()
		if comment != "" {
			event = logger.Error()
		}

		event.
			Int("statusCode", statusCode).
			Str("path", path).
			Dur("latency", latency).
			Str("clientIP", clientIP).
			Str("method", method).
			Msg(comment)
	}
}
