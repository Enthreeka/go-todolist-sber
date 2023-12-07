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

type StatusRequest struct {
	Status bool `json:"status"`
}

// GetTaskHandler godoc
// @Summary Get user task with filter
// @Tags Task
// @Description get user task with pagination and filter, by default without parameters return first page
// @Accept json
// @Produce json
// @Param page query int false "page number" Format(page)
// @Param datetime query string false "date and time required tasks" Format(datetime)
// @Param status query boolean false "task status" Format(status)
// @Success 200 {object} []entity.Task
// @Failure 400 {object} JSONError
// @Failure 404 {object} JSONError
// @Failure 500 {object} JSONError
// @Router /tasks [get]
func (t *taskHandler) GetTaskHandler(w http.ResponseWriter, r *http.Request) {
	opt := new(entity.ParamOption)

	pageString := r.URL.Query().Get("page")
	if pageString != "" {
		page, err := strconv.Atoi(pageString)
		if err != nil {
			t.log.Error("Not correct query result")
			QueryError(w)
			return
		}
		opt.Page = page
	}

	datetime := r.URL.Query().Get("datetime")
	if datetime != "" {
		parsedDate, err := time.Parse("02.01.2006 15:04", datetime)
		if err != nil {
			t.log.Error("time.Parse: %v", err)
			ParseTimeError(w)
			return
		}
		opt.DateTime = parsedDate
	}

	statusString := r.URL.Query().Get("status")
	if statusString != "" {
		status, err := strconv.ParseBool(statusString)
		if err != nil || statusString == "" {
			t.log.Error("Not correct query result")
			QueryError(w)
			return
		}
		opt.Status = &status
	}

	userID := getUserID(r.Context())

	tasks, err := t.taskUsecase.GetTaskWithPaginationAndFilter(context.Background(), userID, opt)
	if err != nil {
		t.log.Error("taskUsecase.GetTaskWithPaginationAndFilter: %v", err)
		HandleError(w, err, apperror.ParseHTTPErrStatusCode(err))
		return
	}

	w.WriteHeader(http.StatusOK)
	e := json.NewEncoder(w)
	e.SetIndent(" ", " ")
	e.Encode(tasks)
}

// GetUserTaskHandler godoc
// @Summary Get user task
// @Tags Task
// @Description get user task by userID from context, you can also make a filter for time and status, return tasks
// @Accept json
// @Produce json
// @Param datetime query string false "date and time required tasks" Format(datetime)
// @Param status query boolean false "task status" Format(status)
// @Success 200 {object} []entity.Task
// @Failure 400 {object} JSONError
// @Failure 404 {object} JSONError
// @Failure 500 {object} JSONError
// @Router /tasks/list [get]
func (t *taskHandler) GetUserTaskHandler(w http.ResponseWriter, r *http.Request) {
	opt := new(entity.ParamOption)

	datetime := r.URL.Query().Get("datetime")
	if datetime != "" {
		parsedDate, err := time.Parse("02.01.2006 15:04", datetime)
		if err != nil {
			t.log.Error("time.Parse: %v", err)
			ParseTimeError(w)
			return
		}
		opt.DateTime = parsedDate
	}

	statusString := r.URL.Query().Get("status")
	if statusString != "" {
		status, err := strconv.ParseBool(statusString)
		if err != nil || statusString == "" {
			t.log.Error("Not correct query result")
			QueryError(w)
			return
		}
		opt.Status = &status
	}

	if datetime == "" || statusString == "" {
		t.log.Error("Not correct query result")
		QueryError(w)
		return
	}

	userID := getUserID(r.Context())

	tasks, err := t.taskUsecase.GetUserTasks(context.Background(), userID, opt)
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
// @Router /tasks/add [post]
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
// @Failure 403 {object} JSONError
// @Failure 404 {object} JSONError
// @Failure 500 {object} JSONError
// @Router /tasks/{id} [delete]
func (t *taskHandler) DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	param := chi.URLParam(r, "id")
	taskID, err := strconv.Atoi(param)
	if err != nil {
		t.log.Error("strconv.Atoi: %v", err)
		DecodingError(w)
		return
	}

	userID := getUserID(r.Context())

	equal, err := t.taskUsecase.IsEqualUserID(context.Background(), userID, taskID)
	if err != nil {
		t.log.Error("taskUsecase.IsEqualUserID: %v", err)
		HandleError(w, err, apperror.ParseHTTPErrStatusCode(err))
		return
	}

	if !equal {
		AccessError(w)
		return
	}

	if err := t.taskUsecase.DeleteTask(context.Background(), taskID); err != nil {
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
// @Param id path int true "task id"
// @Param input body TaskRequest true "task attribute"
// @Success 200 {object} entity.Task
// @Failure 400 {object} JSONError
// @Failure 403 {object} JSONError
// @Failure 404 {object} JSONError
// @Failure 500 {object} JSONError
// @Router /tasks/{id} [put]
func (t *taskHandler) UpdateTaskHandler(w http.ResponseWriter, r *http.Request) {
	param := chi.URLParam(r, "id")
	taskID, err := strconv.Atoi(param)
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

	equal, err := t.taskUsecase.IsEqualUserID(context.Background(), userID, taskID)
	if err != nil {
		t.log.Error("taskUsecase.IsEqualUserID: %v", err)
		HandleError(w, err, apperror.ParseHTTPErrStatusCode(err))
		return
	}

	if !equal {
		AccessError(w)
		return
	}

	task := &entity.Task{
		Header:      data.Header,
		Description: data.Description,
		StartDate:   data.StartDate,
		ID:          taskID,
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
// @Summary Get all users task
// @Tags Task
// @Description Get all users tasks, available only to the admin
// @Accept json
// @Produce json
// @Success 200 {object} []entity.Task
// @Failure 403 {object} JSONError
// @Failure 404 {object} JSONError
// @Failure 500 {object} JSONError
// @Router /tasks/all [get]
func (t *taskHandler) GetAllTasksHandler(w http.ResponseWriter, r *http.Request) {
	role := getRole(r.Context())
	if role != "admin" {
		AccessError(w)
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

// UpdateStatusHandler godoc
// @Summary Set status
// @Tags Task
// @Description Set to task completed or not by userID from context, return updated task
// @Accept json
// @Produce json
// @Param id path int true "task id"
// @Param input body StatusRequest true "task attribute"
// @Success 200 {object} entity.Task
// @Failure 400 {object} JSONError
// @Failure 403 {object} JSONError
// @Failure 404 {object} JSONError
// @Failure 500 {object} JSONError
// @Router /tasks/{id}/status [put]
func (t *taskHandler) UpdateStatusHandler(w http.ResponseWriter, r *http.Request) {
	param := chi.URLParam(r, "id")
	taskID, err := strconv.Atoi(param)
	if err != nil {
		t.log.Error("strconv.Atoi: %v", err)
		DecodingError(w)
		return
	}

	data := new(StatusRequest)
	d := json.NewDecoder(r.Body)
	err = d.Decode(&data)
	if err != nil {
		t.log.Error("json.NewDecoder: %v", err)
		DecodingError(w)
		return
	}

	userID := getUserID(r.Context())

	equal, err := t.taskUsecase.IsEqualUserID(context.Background(), userID, taskID)
	if err != nil {
		t.log.Error("taskUsecase.IsEqualUserID: %v", err)
		HandleError(w, err, apperror.ParseHTTPErrStatusCode(err))
		return
	}

	if !equal {
		AccessError(w)
		return
	}

	updatedTask, err := t.taskUsecase.UpdateTaskStatus(context.Background(), data.Status, taskID)
	if err != nil {
		t.log.Error("taskUsecase.UpdateTaskStatus: %v", err)
		HandleError(w, err, apperror.ParseHTTPErrStatusCode(err))
		return
	}

	w.WriteHeader(http.StatusOK)
	e := json.NewEncoder(w)
	e.SetIndent(" ", " ")
	e.Encode(updatedTask)
}

func getUserID(ctx context.Context) string {
	userID, _ := ctx.Value("userID").(string)

	return userID
}

func getRole(ctx context.Context) string {
	role, _ := ctx.Value("role").(string)

	return role
}
