package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/mirsafari/meta-oci-sync/config"

	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Controller struct {
	ApiVersion      string `json:"apiVersion"`
	Kind            string `json:"kind"`
	meta.ObjectMeta `json:"metadata"`
	Spec            MirrorSpec `json:"spec"`
}

type SyncRequest struct {
	Parent Controller `json:"parent"`
}

type MirrorSpec struct {
	Source      ImageRegistry `json:"source"`
	Destination ImageRegistry `json:"destination"`
}
type ImageRegistry struct {
	Registry string   `json:"registry"`
	Image    string   `json:"image"`
	Tags     []string `json:"tags"`
}

func main() {

	err := config.Clients.Load()

	if err != nil {
		slog.Error("Failed creating application config. Exiting...", err)
		os.Exit(1)
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Post("/sync", SyncHandler)

	slog.Info("Server starting on port " + strconv.Itoa(config.Clients.HTTPPort))
	http.ListenAndServe(":"+strconv.Itoa(config.Clients.HTTPPort), r)
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

	fmt.Println("Creating a mirror job for named " + request.Parent.Name + " of kind " + request.Parent.Kind + "." + request.Parent.ApiVersion)
	fmt.Println("Images to sync: " + request.Parent.Spec.Source.Image)

	err = CreatePod(*request)
	if err != nil {
		slog.Error("Pod couldn't be created", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func CreatePod(SyncRequest) error {

	pods, err := config.Clients.K8S.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}
	fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))

	return nil
}
