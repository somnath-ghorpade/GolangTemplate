## Steps to run Go Template project on local

 - Server is a golang gin framework developed at Go v1.17.9
 - Before running update following files according to environment:
	 - Update db related config in - `mongodb.json`
	 - Server requires mongodb database (version 4.4.13)
 - To run:
	 - `go mod tidy`
	 - `go run main.go`
 - To Build:
	 - `go build`

## Test the gotemplate project
 - `GoTemplate.postman_collection.json` postman collection.

## Pprof usage
 - To check pprof uncomment pprof call from main file.
 - Use profiling for your server see the documentation  https://pkg.go.dev/net/http/pprof
 - also after starting pprof server use this command to check profiling 
    `go tool pprof http://localhost:YOUR_PPROF_PORT/debug/pprof/heap`

## JWT Usage
 - To use jwt authentication and its variables uncomment decodeToken method to access all token variables.

