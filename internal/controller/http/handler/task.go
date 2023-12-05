package handler

import (
	"context"
	"encoding/json"
	"go-todolist-sber/internal/apperror"
	"go-todolist-sber/internal/entity"
	"go-todolist-sber/internal/usecase"
	"go-todolist-sber/pkg/logger"
	"net/http"
	"time"
)

type taskHandler struct {
	taskUsecase usecase.Task
	log         *logger.Logger
}

func NewTaskHandler(taskUsecase usecase.Task, log *logger.Logger) *taskHandler {
	return &taskHandler{
		taskUsecase: taskUsecase,
		log:         log,
	}
}

type TaskRequest struct {
	Header      string `json:"header"`
	Description string `json:"description"`
	StartDate   string `json:"start_date"`
}

func (t *taskHandler) GetTaskHandler(w http.ResponseWriter, r *http.Request) {
	//userID := getID(r.Context())

	task, err := t.taskUsecase.GetUserTasks(context.Background(), "53153c2c-1c10-4b92-b5ff-0cf67b116654")
	if err != nil {
		t.log.Error("taskUsecase.GetUserTasks: %v", err)
		HandleError(w, err, apperror.ParseHTTPErrStatusCode(err))
		return
	}

	w.WriteHeader(http.StatusOK)
	e := json.NewEncoder(w)
	e.Encode(task)
}

func (t *taskHandler) CreateTaskHandler(w http.ResponseWriter, r *http.Request) {
	//userID := getID(r.Context())

	data := new(TaskRequest)
	d := json.NewDecoder(r.Body)
	err := d.Decode(&data)
	if err != nil {
		t.log.Error("json.NewDecoder: %v", err)
		DecodingError(w)
		return
	}

	parsedTime, err := time.Parse("02.01.2006 15:04", data.StartDate)
	if err != nil {
		t.log.Error("time.Parse: %v", err)
		ParseTimeError(w)
		return
	}

	task := &entity.Task{
		Header:      data.Header,
		Description: data.Description,
		StartDate:   parsedTime,
		UserID:      "53153c2c-1c10-4b92-b5ff-0cf67b116654",
	}
	createdTask, err := t.taskUsecase.CreateTask(context.Background(), task)
	if err != nil {
		t.log.Error("taskUsecase.CreateTask: %v", err)
		HandleError(w, err, apperror.ParseHTTPErrStatusCode(err))
		return
	}

	w.WriteHeader(http.StatusCreated)
	e := json.NewEncoder(w)
	e.Encode(createdTask)
}

func getID(ctx context.Context) string {
	userID, _ := ctx.Value("userID").(string)

	return userID
}
