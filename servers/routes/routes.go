package routes

import (
	"golangtemplate/servers/app/models"
	"golangtemplate/servers/app/modules/databasemdl"
	"net/http"

	"github.com/gin-gonic/contrib/jwt"

	"corelab.mkcl.org/MKCLOS/coredevelopmentplatform/corepkgv2/routebuildermdl"

	"github.com/gin-gonic/gin"
)

// Init - route
func Init(g *gin.Engine) {
	o := g.Group("/o")
	r := g.Group("/r")
	r.Use(jwt.Auth(models.JWTKey))
	c := r.Group("/c")

	r.POST("/readData", databasemdl.ReadData)
	r.POST("/insertData", databasemdl.InsertData)
	r.POST("/updateData", databasemdl.UpdateData)
	c.POST("/deleteData", databasemdl.DeleteData)
	g.GET("/checkServerStatus", func(c *gin.Context) {
		c.String(http.StatusOK, "Server is running - OK")
	})
	routebuildermdl.Init(o, r, c, models.JWTKey)
}
