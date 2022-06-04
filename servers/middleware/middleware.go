package middleware

import (
	"runtime/debug"

	"corelab.mkcl.org/MKCLOS/coredevelopmentplatform/corepkgv2/loggermdl"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Init -Init
func Init(g *gin.Engine) {
	g.Use(cors.Default())
	g.Use(Recovery())
}

// Recovery - a recoverymdl for gin to log the panic and recover
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				loggermdl.LogError("recovered:", err)
				loggermdl.LogError(string(debug.Stack()))
				// write to access log
				gin.DefaultErrorWriter.Write(debug.Stack())
			}
		}()
		c.Next()
	}
}
