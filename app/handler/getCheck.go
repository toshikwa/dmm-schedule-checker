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
		name, slots, err := dmm.CheckSchedule(t.Id)
		if err != nil {
			RespondMessage(ctx, w, err.Error(), http.StatusInternalServerError)
			return
		}
		if len(slots) == 0 {
			continue
		}
		// check dynamodb
		exists, err := dmm.ListExistingSlots(ctx, h.Api, t.Id, slots)
		if err != nil {
			RespondMessage(ctx, w, err.Error(), http.StatusInternalServerError)
			return
		}
		// find new slots
		news, err := dmm.FindSlotDiff(slots, exists)
		if err != nil {
			RespondMessage(ctx, w, err.Error(), http.StatusInternalServerError)
			return
		}
		// send message to line
		err = dmm.SendMessage(name, news)
		if err != nil {
			RespondMessage(ctx, w, err.Error(), http.StatusInternalServerError)
			return
		}
		// write to dynamodb
		err = dmm.AddNewSlots(ctx, h.Api, t.Id, news)
		if err != nil {
			RespondMessage(ctx, w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	RespondMessage(ctx, w, "ok", http.StatusOK)
}
