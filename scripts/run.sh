#/bin/bash
FILE_ROOT="resources" FILE_SERVER="http://localhost:8082/" go run cmd/server/server.go & go run cmd/fileserver/fileserver.go & go run sessions/main.go
