package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
	model "todo/api/models"
	"todo/api/responses"

	"github.com/gorilla/mux"
)

func (app *App) AddTask(w http.ResponseWriter, r *http.Request) {
	requestedTask := new(model.Task)
	userId := r.Context().Value("userID").(float64)
	err := json.NewDecoder(r.Body).Decode(&requestedTask)
	if err != nil {
		responses.ErrorResponse(w, http.StatusBadRequest, err)
		return
	}
	requestedTask.UserID = uint(userId)
	requestedTask.PrepareData()
	if err = requestedTask.ValidateFields(); err != nil {
		responses.ErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	err = requestedTask.CreateTask(app.DB)
	if err != nil {
		responses.ErrorResponse(w, http.StatusBadRequest, err)
		return
	} else {
		var resp = map[string]interface{}{"status": "success", "message": "Task Created Successfully", "task": requestedTask}
		responses.JsonResponse(w, http.StatusCreated, resp)
	}

}

func (app *App) GetTasksForUser(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("userID").(float64)
	tasksList := model.GetTasksByUserId(uint(userId), app.DB)
	resp := map[string]interface{}{"status": "success", "message": "tasks found", "tasks": tasksList}
	if tasksList == nil {
		resp["message"] = "No task found"
	}
	responses.JsonResponse(w, http.StatusOK, resp)

}

func (app *App) DeleteTask(w http.ResponseWriter, r *http.Request) {
	requestParameters := mux.Vars(r)
	taskId, err := strconv.Atoi(requestParameters["taskId"])
	if err != nil {
		responses.ErrorResponse(w, http.StatusBadRequest, err)
		return
	}
	userId := r.Context().Value("userID").(float64)
	recordTask := model.GetTasksByTaskId(taskId, app.DB)
	if recordTask.ID == 0 {
		responses.JsonResponse(w, http.StatusBadRequest, responses.ErrorType{Detail: "No record found", Status: "failed"})
		return
	} else if recordTask.UserID != uint(userId) {
		responses.JsonResponse(w, http.StatusUnauthorized, responses.ErrorType{Detail: "You don't have permission to delete this task", Status: "failed"})
		return
	}

	err = recordTask.DeleteTask(app.DB)
	if err != nil {
		responses.JsonResponse(w, http.StatusInternalServerError, err)
		return
	}
	resp := map[string]interface{}{"status": "success", "message": "tasks deleted successfully"}
	responses.JsonResponse(w, http.StatusOK, resp)

}

func (app *App) UpdateTask(w http.ResponseWriter, r *http.Request) {
	requestParameters := mux.Vars(r)
	taskId, err := strconv.Atoi(requestParameters["taskId"])
	if err != nil {
		responses.ErrorResponse(w, http.StatusBadRequest, err)
		return
	}
	userId := r.Context().Value("userID").(float64)
	recordTask := model.GetTasksByTaskId(taskId, app.DB)
	if recordTask.ID == 0 {
		responses.JsonResponse(w, http.StatusBadRequest, responses.ErrorType{Detail: "No record found", Status: "failed"})
		return
	} else if recordTask.UserID != uint(userId) {
		responses.JsonResponse(w, http.StatusUnauthorized, responses.ErrorType{Detail: "You don't have permission to update this task", Status: "failed"})
		return
	}

	requestedTask := new(model.Task)

	err = json.NewDecoder(r.Body).Decode(&requestedTask)
	if err != nil {
		responses.ErrorResponse(w, http.StatusBadRequest, err)
		return
	}
	requestedTask.PrepareData()

	err = requestedTask.UpdateTask(taskId, app.DB)
	if err != nil {
		responses.JsonResponse(w, http.StatusInternalServerError, err)
		return
	}
	resp := map[string]interface{}{"status": "success", "message": "tasks updated successfully"}
	responses.JsonResponse(w, http.StatusOK, resp)

}
