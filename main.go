package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Controller struct {
	ApiVersion      string `json:"apiVersion"`
	Kind            string `json:"kind"`
	meta.ObjectMeta `json:"metadata"`
}

type SyncRequest struct {
	Parent Controller `json:"parent"`
}

func main() {
	port := 8080

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Post("/sync", SyncHandler)

	slog.Info("Server starting on port " + strconv.Itoa(port))
	http.ListenAndServe(":"+strconv.Itoa(port), r)
}

// Main handler for requests that come from metacontroller
func SyncHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		msg := "JSON could not be retrieved"
		slog.Error(msg, err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(msg))
		return
	}

	request := &SyncRequest{}
	// Unmarshal the JSON into the struct.
	err = json.Unmarshal(body, request)
	if err != nil {
		msg := "JSON could not be unmarshalled"
		slog.Error(msg, err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(msg))
		return
	}

	fmt.Println(request)

}
