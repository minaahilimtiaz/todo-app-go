package controllers

import (
	"fmt"
	"log"
	"net/http"
	"todo/api/middlewares"
	model "todo/api/models"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type App struct {
	Router *mux.Router
	DB     *gorm.DB
}

func (app *App) Initialize(dbUser, dbName, dbPassword string) {
	connectionString := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", dbUser, dbPassword, dbName)

	db, connectionError := gorm.Open("postgres", connectionString)
	if connectionError != nil {
		fmt.Printf("\n Cannot connect to database %s", dbName)
		log.Fatal("This is the error:", connectionError)
	} else {
		fmt.Println("We are connected to the database")
	}
	app.DB = db
	app.DB.Debug().AutoMigrate(&model.User{}, &model.Task{})
	db.Model(&model.Task{}).AddForeignKey("user_id", "users(id)", "SET NULL", "CASCADE")

	app.Router = &mux.Router{}
	app.initializeRoutes()
	http.ListenAndServe(":8080", app.Router)

	defer db.Close()

}

func (app *App) initializeRoutes() {
	app.Router.Use(middlewares.SetContentTypeMiddleware)
	app.Router.HandleFunc("/user/register/", app.registerUser).Methods("POST")
	app.Router.HandleFunc("/user/login/", app.loginUser).Methods("GET")

	subRouter := app.Router.PathPrefix("/api/v1/task/").Subrouter()
	subRouter.Use(middlewares.AuthJwtVerify)
	subRouter.HandleFunc("/add/", app.AddTask).Methods("POST")
	subRouter.HandleFunc("/get/", app.GetTasksForUser).Methods("GET")
	subRouter.HandleFunc("/delete/{taskId:[0-9]+}/", app.DeleteTask).Methods("DELETE")
	subRouter.HandleFunc("/update/{taskId:[0-9]+}/", app.UpdateTask).Methods("PUT")
	subRouter.HandleFunc("/assign/", app.assignTaskToUser).Methods("POST")
}
