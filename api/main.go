package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"stefma.guru/valueStorage/storage"
)

type bodyData struct {
	Value string
}

const authorizationToken = "Token ABC123"

func HandleFunc(w http.ResponseWriter, r *http.Request) {
	authorization := r.Header.Get("Authorization")
	if authorization != authorizationToken {
		http.Error(w, errors.New("Authorization failed").Error(), http.StatusBadRequest)
		return
	}

	var jsonData bodyData
	err := json.NewDecoder(r.Body).Decode(&jsonData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	key := r.URL.Query().Get("key")

	storage, err := storage.CreateStorage()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer storage.Close()

	err = storage.Add(key, jsonData.Value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Fprint(w, "Success!")
}
