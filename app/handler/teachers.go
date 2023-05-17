package handler

import (
	"encoding/json"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/go-chi/chi/v5"
	"github.com/toshikwa/dmm-schedule-checker/app/dmm"
)

type PostTeacherHandler struct{ Api *dynamodb.Client }

func (h *PostTeacherHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// get id fron request body
	ctx := r.Context()
	var body struct {
		Id string `json:"id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		RespondMessage(ctx, w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !dmm.AssertTeacherId(body.Id) {
		RespondMessage(ctx, w, "Teacher ID must be 5-digit numeric string", http.StatusBadRequest)
		return
	}
	err := dmm.AddTeacher(ctx, h.Api, body.Id)
	if err != nil {
		RespondMessage(ctx, w, err.Error(), http.StatusInternalServerError)
		return
	}
	RespondMessage(ctx, w, "ok", http.StatusOK)
}

type DeleteTeacherHandler struct{ Api *dynamodb.Client }

func (h *DeleteTeacherHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// get id from the request uri path
	id := chi.URLParam(r, "id")
	ctx := r.Context()
	if !dmm.AssertTeacherId(id) {
		RespondMessage(ctx, w, "Teacher ID must be 5-digit numeric string", http.StatusBadRequest)
		return
	}
	err := dmm.DeleteTeacher(ctx, h.Api, id)
	if err != nil {
		RespondMessage(ctx, w, err.Error(), http.StatusInternalServerError)
		return
	}
	RespondMessage(ctx, w, "ok", http.StatusOK)
}
