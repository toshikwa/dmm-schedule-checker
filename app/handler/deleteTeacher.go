package handler

import (
	"encoding/json"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/toshikwa/dmm-schedule-checker/app/dmm"
)

type DeleteTeacherHandler struct{ Api *dynamodb.Client }

func (h *DeleteTeacherHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var body struct {
		Id string `json:"id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		RespondMessage(ctx, w, err.Error(), http.StatusInternalServerError)
		return
	}
	err := dmm.DeleteTeacher(ctx, h.Api, body.Id)
	if err != nil {
		RespondMessage(ctx, w, err.Error(), http.StatusInternalServerError)
		return
	}
	RespondMessage(ctx, w, "ok", http.StatusOK)
}
