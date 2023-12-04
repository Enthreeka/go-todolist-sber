package handler

import (
	"context"
	"encoding/json"
	"go-todolist-sber/internal/apperror"
	"go-todolist-sber/internal/usecase"
	"net/http"
)

type taskHandler struct {
	taskUsecase usecase.Task
}

func NewTaskHandler(taskUsecase usecase.Task) *taskHandler {
	return &taskHandler{
		taskUsecase: taskUsecase,
	}
}

func (t *taskHandler) GetTaskHandler(w http.ResponseWriter, r *http.Request) {
	userID := getID(r.Context())

	task, err := t.taskUsecase.GetUserTasks(context.Background(), userID)
	if err != nil {
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
