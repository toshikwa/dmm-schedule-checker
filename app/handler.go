package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/toshikwa/dmm-schedule-checker/app/dmm"
)

type Response struct {
	Message string `json:"message"`
}

func RespondJSON(ctx context.Context, w http.ResponseWriter, body any, status int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		rsp := Response{
			Message: http.StatusText(http.StatusInternalServerError),
		}
		if err := json.NewEncoder(w).Encode(rsp); err != nil {
			fmt.Printf("write error response error: %v", err)
		}
		return
	}
	w.WriteHeader(status)
	if _, err := fmt.Fprintf(w, "%s", bodyBytes); err != nil {
		fmt.Printf("write response error: %v", err)
	}
}

type AddTeacherHandler struct {
	api *dynamodb.Client
}

func (at *AddTeacherHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var body struct {
		Id string `json:"id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		RespondJSON(ctx, w, &Response{Message: err.Error()}, http.StatusInternalServerError)
		return
	}
	err := dmm.AddTeacher(ctx, at.api, body.Id)
	if err != nil {
		RespondJSON(ctx, w, &Response{Message: err.Error()}, http.StatusInternalServerError)
		return
	}
	RespondJSON(ctx, w, struct {
		Id string `json:"id"`
	}{Id: body.Id}, http.StatusOK)
}

type DeleteTeacherHandler struct {
	api *dynamodb.Client
}

func (dt *DeleteTeacherHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var body struct {
		Id string `json:"id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		RespondJSON(ctx, w, &Response{Message: err.Error()}, http.StatusInternalServerError)
		return
	}
	err := dmm.DeleteTeacher(ctx, dt.api, body.Id)
	if err != nil {
		RespondJSON(ctx, w, &Response{Message: err.Error()}, http.StatusInternalServerError)
		return
	}
	RespondJSON(ctx, w, struct {
		Id string `json:"id"`
	}{Id: body.Id}, http.StatusOK)
}

type CheckScheduleHandler struct {
	api *dynamodb.Client
}

func (fs *CheckScheduleHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	teachers, err := dmm.ListTeachers(ctx, fs.api)
	if err != nil {
		RespondJSON(ctx, w, &Response{Message: err.Error()}, http.StatusInternalServerError)
		return
	}
	for _, t := range teachers {
		// check schedule
		name, slots, err := dmm.CheckSchedule(t.Id)
		if err != nil {
			RespondJSON(ctx, w, &Response{Message: err.Error()}, http.StatusInternalServerError)
			return
		}
		if len(slots) == 0 {
			continue
		}
		// check dynamodb
		exists, err := dmm.ListExistingSlots(ctx, fs.api, t.Id, slots)
		if err != nil {
			RespondJSON(ctx, w, &Response{Message: err.Error()}, http.StatusInternalServerError)
			return
		}
		// find new slots
		news, err := dmm.FindSlotDiff(slots, exists)
		if err != nil {
			RespondJSON(ctx, w, &Response{Message: err.Error()}, http.StatusInternalServerError)
			return
		}
		// send message to line
		err = dmm.SendMessage(name, news)
		if err != nil {
			RespondJSON(ctx, w, &Response{Message: err.Error()}, http.StatusInternalServerError)
			return
		}
		// write to dynamodb
		err = dmm.AddNewSlots(ctx, fs.api, t.Id, news)
		if err != nil {
			RespondJSON(ctx, w, &Response{Message: err.Error()}, http.StatusInternalServerError)
			return
		}
	}
	RespondJSON(ctx, w, &Response{Message: "ok"}, http.StatusOK)
}
