package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/unrolled/render"
)

type msg struct {
	Message string `json:"message"`
}

func main() {
	r := mux.NewRouter()
	renderer := render.New()

	m := msg{
		Message: "pong",
	}

	r.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		renderer.JSON(w, http.StatusOK, m)
		// json.NewEncoder(w).Encode(m)
	})

	http.ListenAndServe(":8080", r)
	//install rest client in extension to run the request in tests.http
	// docker build . -t docker-containerised-api:latest
	// docker run -p 9000:8080 docker-containerised-api:latest
	//.  docker endpoint : which endpoint it needs to point to
	// the below command is not applicable here but there is an option in gin
	// where you can skip giving the port and later on it can be set as an env var as below
	// docker run -e PORT=9000 -p 9000:9000 docker-containerised-api:latest
	// see this video for more info https://www.youtube.com/watch?v=C5y-14YFs_8
}
