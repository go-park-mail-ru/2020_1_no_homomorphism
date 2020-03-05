package middleware

import (
	"encoding/json"
	"log"
	"net/http"
)

func MarshallAndSendError(err error, w http.ResponseWriter) {
	errJson, e := json.Marshal(err.Error())
	if e != nil {
		log.Println(e)
		return
	}
	w.Header().Set("content-type", "application/json")
	_, e = w.Write(errJson)
	if e != nil {
		log.Println(e)
		return
	}

}
