package handler

import (
	"net/http"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/toshikwa/dmm-schedule-checker/app/dmm"
)

type GetCheckHandler struct{ Api *dynamodb.Client }

func (h *GetCheckHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	teachers, err := dmm.ListTeachers(ctx, h.Api)
	if err != nil {
		RespondJSON(ctx, w, &Response{Message: err.Error()}, http.StatusInternalServerError)
		return
	}
	for _, t := range teachers {
		// check schedule
		name, news, err := dmm.CheckSchedule(t.Id)
		if err != nil {
			RespondMessage(ctx, w, err.Error(), http.StatusInternalServerError)
			return
		}
		// check dynamodb
		exists, err := dmm.ListSlots(ctx, h.Api, t.Id)
		if err != nil {
			RespondMessage(ctx, w, err.Error(), http.StatusInternalServerError)
			return
		}
		// find new slots
		adds, dels, err := dmm.DiffSlots(news, exists)
		if err != nil {
			RespondMessage(ctx, w, err.Error(), http.StatusInternalServerError)
			return
		}
		if len(adds) != 0 {
			// send message to line
			err = dmm.SendMessage(t.Id, name, adds)
			if err != nil {
				RespondMessage(ctx, w, err.Error(), http.StatusInternalServerError)
				return
			}
			// add slots
			err = dmm.AddSlots(ctx, h.Api, t.Id, adds)
			if err != nil {
				RespondMessage(ctx, w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		if len(dels) != 0 {
			// delete slots
			err = dmm.DeleteSlots(ctx, h.Api, t.Id, dels)
			if err != nil {
				RespondMessage(ctx, w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}
	RespondMessage(ctx, w, "ok", http.StatusOK)
}
