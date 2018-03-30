package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"goji.io"
	"goji.io/pat"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	. "main/dao"
	. "main/model"
)

const (
	DB     = "tododb"
	SERVER = "todogomongo"
)

var dao = TasksDAO{}

func init() {
	dao.Server = SERVER
	dao.Database = DB
	dao.Connect()
}

func main() {
	mux := goji.NewMux()
	mux.HandleFunc(pat.Get("/tasks"), AllTasks())
	mux.HandleFunc(pat.Post("/tasks"), AddTask())
	mux.HandleFunc(pat.Get("/tasks/:id"), FindTaskById())
	mux.HandleFunc(pat.Put("/tasks"), UpdateTask())
	mux.HandleFunc(pat.Delete("/tasks/:id"), DeleteTask())
	log.Println("Listening...")

	if err := http.ListenAndServe(":3000", mux); err != nil {
		log.Fatal(err)
	}
}

func AllTasks() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var tasks []Task
		tasks, err := dao.FindAll()

		if err != nil {
			errorWithJSON(w, "Database error", http.StatusInternalServerError)
			log.Println("Failed get all tasks: ", err)
			return
		}
		fmt.Println(tasks)
		responseWithJSON(w, tasks, http.StatusOK)
	}
}

func AddTask() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var task Task
		defer r.Body.Close()
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&task)
		if err != nil {
			errorWithJSON(w, "Incorrect body", http.StatusBadRequest)
			return
		}
		task.ID = bson.NewObjectId()
		fmt.Println(task)
		err = dao.Insert(task)
		if err != nil {
			if mgo.IsDup(err) {
				errorWithJSON(w, "Task with this taskname already exists", http.StatusBadRequest)
				return
			}
			errorWithJSON(w, "Database error", http.StatusInternalServerError)
			log.Println("Failed insert task: ", err)
			return
		}
		responseWithJSON(w, task, http.StatusCreated)
	}
}

func DeleteTask() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		id := pat.Param(r, "id")
		fmt.Println(id)
		if err := dao.Delete(id); err != nil {
			errorWithJSON(w, err.Error(), http.StatusInternalServerError)
			fmt.Println(err.Error())
			return
		}
		responseWithJSON(w, map[string]string{"result": "success"}, http.StatusOK)
	}
}

func FindTaskById() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id := pat.Param(r, "id")
		task, err := dao.FindById(id)
		if err != nil {
			errorWithJSON(w, "Invalid Task ID", http.StatusBadRequest)
			return
		}
		responseWithJSON(w, task, http.StatusOK)
	}
}

func UpdateTask() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		var task Task
		if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
			errorWithJSON(w, "Invalid request payload", http.StatusBadRequest)
			return
		}
		if err := dao.Update(task); err != nil {
			errorWithJSON(w, err.Error(), http.StatusInternalServerError)
			return
		}
		responseWithJSON(w, map[string]string{"result": "success"}, http.StatusOK)
	}
}

func errorWithJSON(w http.ResponseWriter, message string, code int) {
	fmt.Fprintf(w, "{message: %q}", message)
	responseWithJSON(w, map[string]string{"error": message}, code)
}

func responseWithJSON(w http.ResponseWriter, payload interface{}, code int) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	w.Write(response)
}
