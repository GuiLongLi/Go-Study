package main

import (
	"fmt"
	"net/http"
)

type MyMux struct {

}

func (p *MyMux) ServeHTTP(w http.ResponseWriter,r *http.Request)  {
	if r.URL.Path == "/"{
		sayhelloName(w,r)
		return
	}
	http.NotFound(w,r)
	return
}

func sayhelloName(w http.ResponseWriter,r *http.Request)  {
	fmt.Fprintf(w,"hello testmux.go")
}

func main() {
	mux := &MyMux{}
	http.ListenAndServe(":6665",mux)
}
