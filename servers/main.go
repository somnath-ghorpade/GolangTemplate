//go:build !fasthttp
// +build !fasthttp

package main

import (
	"golangtemplate/servers/app/models"
	"golangtemplate/servers/common"
	"golangtemplate/servers/middleware"
	"golangtemplate/servers/routes"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"

	"corelab.mkcl.org/MKCLOS/coredevelopmentplatform/corepkgv2/filemdl"
	"go.uber.org/zap/zapcore"

	"github.com/gin-contrib/cors"
	"github.com/tidwall/gjson"

	"net/http"

	"corelab.mkcl.org/MKCLOS/coredevelopmentplatform/corepkgv2/errormdl"
	"corelab.mkcl.org/MKCLOS/coredevelopmentplatform/corepkgv2/loggermdl"

	"github.com/gin-gonic/gin"
)

func initLogger() {
	loggerLevel := zapcore.ErrorLevel
	loggermdl.Init("/logs/main.log", 0, 1, 0, loggerLevel)
}
func main() {

	gin.SetMode(gin.ReleaseMode)
	//Access log
	accessLog := filepath.Join("logs")
	if err := filemdl.CreateDirectoryRecursive(accessLog); err != nil {
		loggermdl.LogError(err)
		return
	}

	logFile, err := os.Create(filepath.Join(accessLog, "access.log"))
	if err != nil {
		loggermdl.LogError(err)
		return
	}
	g := gin.New()
	gin.DefaultErrorWriter = logFile
	gin.DefaultWriter = logFile
	g.Use(gin.LoggerWithWriter(logFile))
	g.Use(middleware.Recovery())
	md := cors.DefaultConfig()
	md.AllowAllOrigins = true
	md.AllowHeaders = []string{"*"}
	md.AllowMethods = []string{"*"}
	md.ExposeHeaders = []string{"Authorization"}
	g.Use(cors.New(md))

	middleware.Init(g)

	if err := initializeAll(g); err != nil {
		loggermdl.LogError(err)
		return
	}

	var netListen net.Listener
	if models.DefaultPort == "" {
		netListen, _, err = common.GeneratePort()
		if err != nil {
			loggermdl.LogError("Failed to generate new port - ", err)
			return
		}
	} else {
		netListen, err = net.Listen("tcp", ":"+models.DefaultPort)
		if err != nil {
			loggermdl.LogError("start server: failed to listen on provided port - ", err)
			return
		}
	}
	initLogger()
	loggermdl.LogError("Server started on ", netListen.Addr().String())
	err = http.Serve(netListen, g)
	if errormdl.CheckErr(err) != nil {
		errormdl.Wrap(err.Error())
		return
	}
}

func initializeAll(g *gin.Engine) error {
	routes.Init(g)

	// read file and initialise mongo connection
	content, err := ioutil.ReadFile(models.ConfigPath)
	if err != nil {
		loggermdl.LogError(err)
		return err
	}
	ConfigData := gjson.ParseBytes(content)
	err = common.InitMongoDBConnectionUsingJson(ConfigData.Get("MongoConfig").String())
	if err != nil {
		loggermdl.LogError(err)
		return err
	}

	// Decode Provided jwt token
	// err = models.DecodeToken("Token")
	// if err != nil {
	// 	return err
	// }

	// Initialise Mysql Connection
	cnn, err := common.InitMysqlDBConnectionUsingJson(ConfigData.Get("MysqlConfig").String())
	if err != nil {
		loggermdl.LogError(err)
		return err
	}
	// mysq
	// session := cnn.NewSession()
	loggermdl.LogError("session", cnn)
	// Pprof
	// pport, err := common.StartPprof(models.PprofPort)
	// if err != nil {
	// 	loggermdl.LogError(err)
	// 	return err
	// }
	// loggermdl.LogError("pprof server is running on port", pport)
	return nil
}
