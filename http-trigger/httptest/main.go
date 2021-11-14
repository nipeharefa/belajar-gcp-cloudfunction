package main

import (
	"net/http"

	p "github.com/nipeharefa/belajar-cloudfunction"
)

func main() {

	r := http.NewServeMux()
	r.HandleFunc("/", p.HelloWorld)
	http.ListenAndServe(":8080", r)
}
