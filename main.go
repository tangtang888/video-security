package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"./auth"
	"./task-dispatcher/autofind"
	"./task-dispatcher/tasks"
	"github.com/gorilla/mux"
)

//autoFindHandler
var autoFindHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	bizTask := autofind.DeviceTask{
		Task: tasks.BizTask{
			ID:   "1",
			Name: "test",
		},
	}

	dispatcher := tasks.GetInstance()
	dispatcher.RunTask(&bizTask)

	devices := bizTask.Result.Devices
	for _, device := range devices {
		fmt.Fprintf(w, "device xaddres=", device.Xaddr)
	}

	fmt.Fprintf(w, "Find!")
})

//helloHandler home handler
var helloHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello!")
})

func main() {

	r := mux.NewRouter()

	v1 := r.PathPrefix("/v1/").Subrouter()
	v1.Handle("/get-token", auth.GetTokenHandler).Methods("GET")
	v1.Handle("/home", auth.MiddlewareHandler(helloHandler)).Methods("GET")
	//v1.Handle("/find", auth.MiddlewareHandler(autoFindHandler)).Methods("GET")
	v1.Handle("/find", autoFindHandler).Methods("GET")

	srv := &http.Server{
		Addr: ":8001",
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r,
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
