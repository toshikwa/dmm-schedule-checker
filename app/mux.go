package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/go-chi/chi/v5"
	"github.com/toshikwa/dmm-schedule-checker/app/handler"
)

func NewMux(ctx context.Context) (http.Handler, error) {
	mux := chi.NewRouter()
	// add health endpoint
	mux.HandleFunc(
		"/health",
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			_, _ = w.Write([]byte(`{"status": "ok"}`))
		},
	)
	// create client
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("Failed to load config: %s", err)
	}
	ddbClient := dynamodb.NewFromConfig(cfg)
	// POST /teacher
	pt := &handler.PostTeacherHandler{Api: ddbClient}
	mux.Post("/teacher", pt.ServeHTTP)
	// DELETE /teacher
	dt := &handler.DeleteTeacherHandler{Api: ddbClient}
	mux.Delete("/teacher", dt.ServeHTTP)
	// GET /check
	check := &handler.GetCheckHandler{Api: ddbClient}
	mux.Get("/check", check.ServeHTTP)
	return mux, nil
}
