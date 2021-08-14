package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pscompsci/vmmm/internal/explorer"
)

func main() {
	db, err := sql.Open("sqlite3", "./vmmm.db")
	if err != nil {
		log.Fatalf("Could not connect to database: %v\n", err)
	}

	r := mux.NewRouter()

	s := server{router: r, db: db}
	s.routes()

	log.Fatal(http.ListenAndServe(":3001", s.router))
}

type server struct {
	router *mux.Router
	db     *sql.DB
}

func withCORS(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Expose-Headers", "Location")
		fn(w, r)
	}
}

func (s *server) routes() {
	s.router.HandleFunc("/", withCORS(s.handleIndex()))
}

func (s *server) handleIndex() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		e := explorer.Explorer{}
		vms, err := e.GetVMListFromHost(context.Background(), "https://root:password@192.168.0.150/sdk")
		if err != nil {
			log.Fatalf("Could not obtain vm list: %v\n", err)
		}
		// TODO: Remove this and put into debug logging
		for _, vm := range *vms {
			fmt.Printf("%s, %s, %s\n", vm.Name, vm.Parent, vm.OverallStatus)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(*vms)
	}
}
