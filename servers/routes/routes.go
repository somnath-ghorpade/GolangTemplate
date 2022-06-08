package routes

import (
	"golangtemplate/servers/app/models"
	"golangtemplate/servers/app/modules/localmongo"
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
	// use in case of mongo DB
	r.POST("/readData", localmongo.ReadData)
	r.POST("/insertData", localmongo.InsertData)
	r.POST("/updateData", localmongo.UpdateData)
	c.POST("/deleteData", localmongo.DeleteData)

	// use in case of mysql DB
	r.POST("/readMysqlData", localmongo.ReadData)

	g.GET("/checkServerStatus", func(c *gin.Context) {
		c.String(http.StatusOK, "Server is running - OK")
	})
	routebuildermdl.Init(o, r, c, models.JWTKey)
}
