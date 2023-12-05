package handler

import (
	"context"
	"encoding/json"
	"go-todolist-sber/internal/apperror"
	"go-todolist-sber/internal/entity"
	"go-todolist-sber/internal/usecase"
	"go-todolist-sber/pkg/logger"
	"net/http"
	"strconv"
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
	ID          int    `json:"id"`
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

func (t *taskHandler) DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	var data struct {
		ID int `json:"id"`
	}
	d := json.NewDecoder(r.Body)
	err := d.Decode(&data)
	if err != nil {
		t.log.Error("json.NewDecoder: %v", err)
		DecodingError(w)
		return
	}

	if err := t.taskUsecase.DeleteTask(context.Background(), data.ID); err != nil {
		t.log.Error("taskUsecase.DeleteTask: %v", err)
		HandleError(w, err, apperror.ParseHTTPErrStatusCode(err))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (t *taskHandler) UpdateTaskHandler(w http.ResponseWriter, r *http.Request) {
	data := new(TaskRequest)
	d := json.NewDecoder(r.Body)
	err := d.Decode(&data)
	if err != nil {
		t.log.Error("json.NewDecoder: %v", err)
		DecodingError(w)
		return
	}

	var parsedTime time.Time
	if data.StartDate != "" {
		parsedTime, err = time.Parse("02.01.2006 15:04", data.StartDate)
		if err != nil {
			t.log.Error("time.Parse: %v", err)
			ParseTimeError(w)
			return
		}
	}

	task := &entity.Task{
		Header:      data.Header,
		Description: data.Description,
		StartDate:   parsedTime,
		ID:          data.ID,
		UserID:      "53153c2c-1c10-4b92-b5ff-0cf67b116654",
	}
	updatedTask, err := t.taskUsecase.UpdateTask(context.Background(), task)
	if err != nil {
		t.log.Error("taskUsecase.UpdateTask: %v", err)
		HandleError(w, err, apperror.ParseHTTPErrStatusCode(err))
		return
	}

	w.WriteHeader(http.StatusOK)
	e := json.NewEncoder(w)
	e.Encode(updatedTask)
}

func (t *taskHandler) GetAllTasksHandler(w http.ResponseWriter, r *http.Request) {
	tasks, err := t.taskUsecase.GetAllTasks(context.Background())
	if err != nil {
		t.log.Error("taskUsecase.GetAllTasks: %v", err)
		HandleError(w, err, apperror.ParseHTTPErrStatusCode(err))
		return
	}

	w.WriteHeader(http.StatusOK)
	e := json.NewEncoder(w)
	e.Encode(tasks)
}

//localhost:8080/task/pagination?page=3&status=false
func (t *taskHandler) GetTaskWithPaginationHandler(w http.ResponseWriter, r *http.Request) {
	page, errPage := strconv.Atoi(r.URL.Query().Get("page"))
	status, errStatus := strconv.ParseBool(r.URL.Query().Get("status"))
	if errPage != nil || errStatus != nil {
		t.log.Error("Empty query result")
		QueryError(w)
		return
	}

	tasks, err := t.taskUsecase.PaginationTasks(context.Background(), "53153c2c-1c10-4b92-b5ff-0cf67b116654", status, page)
	if err != nil {
		t.log.Error("taskUsecase.PaginationTasks: %v", err)
		HandleError(w, err, apperror.ParseHTTPErrStatusCode(err))
		return
	}

	if len(tasks) == 0 {
		t.log.Error("tasks = 0")
		HandleError(w, apperror.ErrNotFound, http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	e := json.NewEncoder(w)
	e.Encode(tasks)
}

func (t *taskHandler) GetFilteredHandler(w http.ResponseWriter, r *http.Request) {
	date := r.URL.Query().Get("datetime")
	status, err := strconv.ParseBool(r.URL.Query().Get("status"))
	if err != nil || date == "" {
		t.log.Error("Not correct query result")
		QueryError(w)
		return
	}

	parsedDate, err := time.Parse("02.01.2006 15:04", date)
	if err != nil {
		t.log.Error("time.Parse: %v", err)
		ParseTimeError(w)
		return
	}

	tasks, err := t.taskUsecase.GetFilteredTasks(context.Background(), "53153c2c-1c10-4b92-b5ff-0cf67b116654", parsedDate, status)
	if err != nil {
		t.log.Error("taskUsecase.GetFilteredTasks: %v", err)
		HandleError(w, err, apperror.ParseHTTPErrStatusCode(err))
		return
	}

	w.WriteHeader(http.StatusOK)
	e := json.NewEncoder(w)
	e.Encode(tasks)
}

func getID(ctx context.Context) string {
	userID, _ := ctx.Value("userID").(string)

	return userID
}
