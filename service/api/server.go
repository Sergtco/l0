package api

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"service/data"
	"service/data/cache"
)

func index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	tmpl, err := template.ParseFiles("./api/views/index.html")
	if err != nil {
		log.Println(err)
		http.Error(w, "InternalServerError", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, "")
}

func getData(w http.ResponseWriter, r *http.Request) {
	uid := r.URL.Query().Get("uid")
	order, err := data.GetOrder(uid)
	if err == cache.NotExistError {
		http.Error(w, "Order with such id does not exist", http.StatusBadRequest)
		return
	} else if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	data, _ := json.MarshalIndent(order, "", "    ")
	w.Write(data)
}

func RunServer() error {
	mx := http.NewServeMux()
	mx.HandleFunc("GET /", index)
	mx.HandleFunc("GET /data/", getData)
	server := http.Server{
		Handler: mx,
		Addr:    ":6969",
	}
	if err := server.ListenAndServe(); err != nil {
		return err
	}
	return nil
}
