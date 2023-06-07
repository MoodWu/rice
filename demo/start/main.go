package main

import (
	"fmt"
	"net/http"
	"os"
	"rice/router"

)

func main() {

	r, err := router.Register()
	if err != nil {
		fmt.Printf("router register with errors [%s]", err)
		os.Exit(-1)
	}
	fmt.Printf("service listen on [%s]\r\n","127.0.0.1:5250")
	fmt.Println(http.ListenAndServe("127.0.0.1:5250", r))
}
