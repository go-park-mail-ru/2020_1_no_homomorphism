package main

import (
	"fmt"
	"log"
	"net/http"
)

func main(){
	fmt.Println("Starts server at 8082")
	err := http.ListenAndServe(":8082", http.FileServer(http.Dir("./resources")))
	if err != nil {
			log.Println(err)
			return
	}
}

