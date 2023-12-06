package handler

import (
	"context"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"go-todolist-sber/internal/apperror"
	"go-todolist-sber/internal/entity"
	"go-todolist-sber/internal/task"
	"go-todolist-sber/pkg/logger"
	"net/http"
	"strconv"
	"time"
)

type taskHandler struct {
	taskUsecase task.TaskUsecase
	log         *logger.Logger
}

func NewTaskHandler(taskUsecase task.TaskUsecase, log *logger.Logger) *taskHandler {
	return &taskHandler{
		taskUsecase: taskUsecase,
		log:         log,
	}
}

type TaskRequest struct {
	Header      string    `json:"header"`
	Description string    `json:"description"`
	StartDate   time.Time `json:"start_date"`
}

// GetTaskHandler godoc
// @Summary Get user task
// @Tags Task
// @Description get user task by userID from context, return tasks
// @Accept json
// @Produce json
// @Success 200 {object} []entity.Task
// @Failure 404 {object} JSONError
// @Failure 500 {object} JSONError
// @Router /task/list [get]
func (t *taskHandler) GetTaskHandler(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r.Context())

	tasks, err := t.taskUsecase.GetUserTasks(context.Background(), userID)
	if err != nil {
		t.log.Error("taskUsecase.GetUserTasks: %v", err)
		HandleError(w, err, apperror.ParseHTTPErrStatusCode(err))
		return
	}

	w.WriteHeader(http.StatusOK)
	e := json.NewEncoder(w)
	e.SetIndent(" ", " ")
	e.Encode(tasks)
}

// CreateTaskHandler godoc
// @Summary Create new task
// @Tags Task
// @Description create new user task by userID from context, return created task
// @Accept json
// @Produce json
// @Param input body TaskRequest true "task attribute"
// @Success 201 {object} entity.Task
// @Failure 400 {object} JSONError
// @Failure 404 {object} JSONError
// @Failure 500 {object} JSONError
// @Router /task/add [post]
func (t *taskHandler) CreateTaskHandler(w http.ResponseWriter, r *http.Request) {
	data := new(TaskRequest)
	d := json.NewDecoder(r.Body)
	err := d.Decode(&data)
	if err != nil {
		t.log.Error("json.NewDecoder: %v", err)
		DecodingError(w)
		return
	}

	userID := getUserID(r.Context())

	task := &entity.Task{
		Header:      data.Header,
		Description: data.Description,
		StartDate:   data.StartDate,
		UserID:      userID,
	}
	createdTask, err := t.taskUsecase.CreateTask(context.Background(), task)
	if err != nil {
		t.log.Error("taskUsecase.CreateTask: %v", err)
		HandleError(w, err, apperror.ParseHTTPErrStatusCode(err))
		return
	}

	w.WriteHeader(http.StatusCreated)
	e := json.NewEncoder(w)
	e.SetIndent(" ", " ")
	e.Encode(createdTask)
}

// DeleteTaskHandler godoc
// @Summary Delete task
// @Tags Task
// @Description delete task by id
// @Accept json
// @Produce json
// @Param id path int true "task id"
// @Success 204
// @Failure 404 {object} JSONError
// @Failure 500 {object} JSONError
// @Router /task/{id} [delete]
func (t *taskHandler) DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	param := chi.URLParam(r, "id")
	id, err := strconv.Atoi(param)
	if err != nil {
		t.log.Error("strconv.Atoi: %v", err)
		DecodingError(w)
		return
	}

	if err := t.taskUsecase.DeleteTask(context.Background(), id); err != nil {
		t.log.Error("taskUsecase.DeleteTask: %v", err)
		HandleError(w, err, apperror.ParseHTTPErrStatusCode(err))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// UpdateTaskHandler godoc
// @Summary Update task
// @Tags Task
// @Description Update header, description, datetime, task status by userID from context, return updated task
// @Accept json
// @Produce json
// @Param input body TaskRequest true "task attribute"
// @Success 200 {object} entity.Task
// @Failure 400 {object} JSONError
// @Failure 404 {object} JSONError
// @Failure 500 {object} JSONError
// @Router /task/delete [put]
func (t *taskHandler) UpdateTaskHandler(w http.ResponseWriter, r *http.Request) {
	param := chi.URLParam(r, "id")
	id, err := strconv.Atoi(param)
	if err != nil {
		t.log.Error("strconv.Atoi: %v", err)
		DecodingError(w)
		return
	}

	data := new(TaskRequest)
	d := json.NewDecoder(r.Body)
	err = d.Decode(&data)
	if err != nil {
		t.log.Error("json.NewDecoder: %v", err)
		DecodingError(w)
		return
	}

	userID := getUserID(r.Context())

	task := &entity.Task{
		Header:      data.Header,
		Description: data.Description,
		StartDate:   data.StartDate,
		ID:          id,
		UserID:      userID,
	}
	updatedTask, err := t.taskUsecase.UpdateTask(context.Background(), task)
	if err != nil {
		t.log.Error("taskUsecase.UpdateTask: %v", err)
		HandleError(w, err, apperror.ParseHTTPErrStatusCode(err))
		return
	}

	w.WriteHeader(http.StatusOK)
	e := json.NewEncoder(w)
	e.SetIndent(" ", " ")
	e.Encode(updatedTask)
}

// GetAllTasksHandler godoc
// @Summary Get all task
// @Tags Task
// @Description Get all users tasks
// @Accept json
// @Produce json
// @Success 200 {object} []entity.Task
// @Failure 403 {object} JSONError
// @Failure 404 {object} JSONError
// @Failure 500 {object} JSONError
// @Router /task/all [get]
func (t *taskHandler) GetAllTasksHandler(w http.ResponseWriter, r *http.Request) {
	role := getRole(r.Context())
	if role != "admin" {
		ErrorJSON(w, "access denied", http.StatusForbidden)
		return
	}

	tasks, err := t.taskUsecase.GetAllTasks(context.Background())
	if err != nil {
		t.log.Error("taskUsecase.GetAllTasks: %v", err)
		HandleError(w, err, apperror.ParseHTTPErrStatusCode(err))
		return
	}

	w.WriteHeader(http.StatusOK)
	e := json.NewEncoder(w)
	e.SetIndent(" ", " ")
	e.Encode(tasks)
}

// GetTaskWithPaginationHandler godoc
// @Summary Get tasks
// @Tags Task
// @Description Get user task with pagination by userID from context
// @Accept json
// @Produce json
// @Param page query int false "page number" Format(page)
// @Param status query boolean false "task status" Format(status)
// @Success 200 {object} []entity.Task
// @Failure 400 {object} JSONError
// @Failure 404 {object} JSONError
// @Failure 500 {object} JSONError
// @Router /task/pagination [get]
func (t *taskHandler) GetTaskWithPaginationHandler(w http.ResponseWriter, r *http.Request) {
	page, errPage := strconv.Atoi(r.URL.Query().Get("page"))
	status, errStatus := strconv.ParseBool(r.URL.Query().Get("status"))
	if errPage != nil || errStatus != nil {
		t.log.Error("Empty query result")
		QueryError(w)
		return
	}

	userID := getUserID(r.Context())

	tasks, err := t.taskUsecase.PaginationTasks(context.Background(), userID, status, page)
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
	e.SetIndent(" ", " ")
	e.Encode(tasks)
}

// GetFilteredHandler godoc
// @Summary Get task
// @Tags Task
// @Description Get user task with filter
// @Accept json
// @Produce json
// @Param datetime query string false "date and time required tasks" Format(datetime)
// @Param status query boolean false "task status" Format(status)
// @Success 200 {object} []entity.Task
// @Failure 400 {object} JSONError
// @Failure 404 {object} JSONError
// @Failure 500 {object} JSONError
// @Router /task/filter [get]
func (t *taskHandler) GetFilteredHandler(w http.ResponseWriter, r *http.Request) {
	datetime := r.URL.Query().Get("datetime")
	status, err := strconv.ParseBool(r.URL.Query().Get("status"))
	if err != nil || datetime == "" {
		t.log.Error("Not correct query result")
		QueryError(w)
		return
	}

	parsedDate, err := time.Parse("02.01.2006 15:04", datetime)
	if err != nil {
		t.log.Error("time.Parse: %v", err)
		ParseTimeError(w)
		return
	}

	userID := getUserID(r.Context())

	tasks, err := t.taskUsecase.GetFilteredTasks(context.Background(), userID, parsedDate, status)
	if err != nil {
		t.log.Error("taskUsecase.GetFilteredTasks: %v", err)
		HandleError(w, err, apperror.ParseHTTPErrStatusCode(err))
		return
	}

	w.WriteHeader(http.StatusOK)
	e := json.NewEncoder(w)
	e.SetIndent(" ", " ")
	e.Encode(tasks)
}

func getUserID(ctx context.Context) string {
	userID, _ := ctx.Value("userID").(string)

	return userID
}

func getRole(ctx context.Context) string {
	role, _ := ctx.Value("role").(string)

	return role
}
