package controllers

import (
	"encoding/json"
	"net/http"
	model "todo/api/models"
	"todo/api/responses"
	"todo/api/utils"
)

func (app *App) registerUser(w http.ResponseWriter, r *http.Request) {
	var user = new(model.User)
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		responses.ErrorResponse(w, http.StatusBadRequest, err)
		return
	}
	defer r.Body.Close()

	user.PrepareData()

	err = user.ValidateFields("signup")
	if err != nil {
		responses.ErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	foundRecord, _ := model.GetUserByEmail(user.Email, app.DB)
	if foundRecord != nil && foundRecord.IsRegistered {
		responses.JsonResponse(w, http.StatusConflict, &responses.ErrorType{Detail: "User already exists", Status: "failed"})
		return
	} else if foundRecord != nil && !foundRecord.IsRegistered {
		user.IsRegistered = true
		user.ID = foundRecord.ID
		err = user.UpdateUser(app.DB)
		if err != nil {
			responses.ErrorResponse(w, http.StatusBadRequest, err)
			return
		} else {
			var resp = map[string]interface{}{"status": "success", "message": "Registered successfully", "user": user}
			responses.JsonResponse(w, http.StatusCreated, resp)
			return
		}
	}

	user.IsRegistered = true
	err = user.SaveUser(app.DB)
	if err != nil {
		responses.ErrorResponse(w, http.StatusBadRequest, err)
		return
	}
	var resp = map[string]interface{}{"status": "success", "message": "Registered successfully", "user": user}
	responses.JsonResponse(w, http.StatusCreated, resp)
}

func (app *App) loginUser(w http.ResponseWriter, r *http.Request) {
	var userInput = new(model.User)
	err := json.NewDecoder(r.Body).Decode(&userInput)
	if err != nil {
		responses.ErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	defer r.Body.Close()

	userInput.PrepareData()

	err = userInput.ValidateFields("login")
	if err != nil {
		responses.ErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	foundRecord, _ := model.GetUserByEmail(userInput.Email, app.DB)
	if foundRecord == nil || !foundRecord.IsRegistered {
		responses.JsonResponse(w, http.StatusUnauthorized, &responses.ErrorType{Detail: "Please create an account", Status: "failed"})
		return
	} else if err = model.CheckPasswordHash(userInput.Password, foundRecord.Password); foundRecord.Email != userInput.Email || err != nil {
		responses.JsonResponse(w, http.StatusUnauthorized, &responses.ErrorType{Detail: "Invalid Credentials", Status: "failed"})
		return
	}

	token, err := utils.EncodeAuthToken(foundRecord.ID)
	if err != nil {
		responses.ErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	var resp = map[string]interface{}{"status": "success", "message": "Logged in successfully", "accessToken": token}
	responses.JsonResponse(w, http.StatusOK, resp)

}

func (app *App) assignTaskToUser(w http.ResponseWriter, r *http.Request) {
	var assignedTask = new(model.AssignedTask)
	err := json.NewDecoder(r.Body).Decode(&assignedTask)
	if err != nil {
		responses.ErrorResponse(w, http.StatusBadRequest, err)
		return
	}
	defer r.Body.Close()
	foundRecord, _ := model.GetUserByEmail(assignedTask.Email, app.DB)
	if foundRecord == nil {
		user := model.User{Email: assignedTask.Email, IsRegistered: false}
		err = user.SaveUser(app.DB)
		if err != nil {
			responses.ErrorResponse(w, http.StatusBadRequest, err)
			return
		}
		assignedTask.Task.UserID = user.ID
	} else {
		assignedTask.Task.UserID = foundRecord.ID
	}
	err = assignedTask.Task.CreateTask(app.DB)
	if err != nil {
		responses.ErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	var resp = map[string]interface{}{"status": "success", "message": "Task Assigned Succesfully", "task": assignedTask.Task}
	responses.JsonResponse(w, http.StatusCreated, resp)
}
