package main

import (
	"errors"
	"net/http"

	"corelab.mkcl.org/MKCLOS/coredevelopmentplatform/corepkgv2/loggermdl"
)

// startPprof - start the pprof server
func StartPprof(port string) (string, error) {
	if port == "0" || port == "" {
		loggermdl.LogError("Pporf port is empty")
		return "", errors.New("Pporf port is empty")
	}
	go func(port string) {
		loggermdl.LogError(http.ListenAndServe(":"+port, nil))
	}(port)
	return port, nil
}
