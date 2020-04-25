package main

import (
	"fmt"
	"log"
	"net/http"
)

func main(){
	fmt.Println("Starts server at 8082")
	err := http.ListenAndServeTLS(":8082","fullchain.pem","privkey.pem", http.FileServer(http.Dir("./resources")))
	if err != nil {
		log.Println(err)
		return
	}
}

