package common

import (
	"encoding/json"
	"net"
	"net/http"
	_ "net/http/pprof"

	"corelab.mkcl.org/MKCLOS/coredevelopmentplatform/corepkgv2/dalmdl/coremongo"
	"corelab.mkcl.org/MKCLOS/coredevelopmentplatform/corepkgv2/dalmdl/mysql"
	"corelab.mkcl.org/MKCLOS/coredevelopmentplatform/corepkgv2/errormdl"
	"corelab.mkcl.org/MKCLOS/coredevelopmentplatform/corepkgv2/loggermdl"
	"github.com/gocraft/dbr/v2"
)

func InitMongoDBConnectionUsingJson(jsonString string) error {
	var hosts []coremongo.MongoHost
	json.Unmarshal([]byte(jsonString), &hosts)

	err := coremongo.InitUsingJSON(hosts)
	if err != nil {
		loggermdl.LogError(err)
		return errormdl.Wrap(err.Error())
	}
	return err
}

func InitMysqlDBConnectionUsingJson(jsonString string) (mysqlConn *dbr.Connection, err error) {
	var hosts mysql.MySQLConnection
	json.Unmarshal([]byte(jsonString), &hosts)
	mysqlConnection, err := mysql.InitConnection(hosts)
	if err != nil {
		loggermdl.LogError(err)
		return nil, errormdl.Wrap(err.Error())
	}
	res, err := mysqlConnection.Query("Select * from Users;")
	if err != nil {
		loggermdl.LogError(err)
		return nil, errormdl.Wrap(err.Error())
	}
	loggermdl.LogError("res", res)
	return mysqlConnection, err
}

// GeneratePort - generate new availabe port
func GeneratePort() (net.Listener, string, error) {
	var ln net.Listener
	var err error
	// get new available port
	ln, err = net.Listen("tcp", ":0")
	if err != nil {
		loggermdl.LogError("Failed to generate new port - ", err)
		return ln, "", errormdl.Wrap(err.Error())
	}
	_, port, err := net.SplitHostPort(ln.Addr().String())
	if err != nil {
		loggermdl.LogError("Failed to get new port from address - ", err)
		return ln, "", errormdl.Wrap(err.Error())
	}
	return ln, port, nil
}

// startPprof - start the pprof server
func StartPprof(port string) (string, error) {
	if port == "0" || port == "" {
		loggermdl.LogError("Pporf port is empty")
		return "", errormdl.Wrap("Pporf port is empty")
	}
	go func(port string) {
		loggermdl.LogError(http.ListenAndServe("localhost:"+port, nil))
	}(port)
	return port, nil
}
