package handler

import (
	"context"
	"encoding/json"
	"go-todolist-sber/internal/apperror"
	"go-todolist-sber/internal/usecase"
	"go-todolist-sber/pkg/logger"
	"net/http"
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

func (t *taskHandler) GetTaskHandler(w http.ResponseWriter, r *http.Request) {
	//userID := getID(r.Context())

	task, err := t.taskUsecase.GetUserTasks(context.Background(), "4bd6d3af-da49-4bef-8caa-0030fc30e93b")
	if err != nil {
		t.log.Error("taskUsecase.GetUserTasks: %v", err)
		HandleError(w, err, apperror.ParseHTTPErrStatusCode(err))
		return
	}

	w.WriteHeader(http.StatusOK)
	e := json.NewEncoder(w)
	e.Encode(task)
}

func getID(ctx context.Context) string {
	userID, _ := ctx.Value("userID").(string)

	return userID
}
