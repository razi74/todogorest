package dao

import (
	"log"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	. "main/model"
)

type TasksDAO struct {
	Server   string
	Database string
}

var db *mgo.Database

const (
	COLLECTION_TASKS = "tasks"
)

func (m *TasksDAO) Connect() {
	session, err := mgo.Dial(m.Server)
	if err != nil {
		log.Fatal(err)
	}

	session.SetMode(mgo.Monotonic, true)

	c := session.DB(m.Database).C(COLLECTION_TASKS)

	index := mgo.Index{
		Key:        []string{"id"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
	errr := c.EnsureIndex(index)
	if errr != nil {
		panic(err)
	}

	db = session.DB(m.Database)
}

func (m *TasksDAO) FindAll() ([]Task, error) {
	var tasks []Task
	err := db.C(COLLECTION_TASKS).Find(bson.M{}).All(&tasks)
	return tasks, err
}

func (m *TasksDAO) FindById(id string) (Task, error) {
	var task Task
	err := db.C(COLLECTION_TASKS).FindId(bson.ObjectIdHex(id)).One(&task)
	return task, err
}

func (m *TasksDAO) Insert(task Task) error {
	err := db.C(COLLECTION_TASKS).Insert(&task)
	return err
}

func (m *TasksDAO) Delete(id string) error {
	err := db.C(COLLECTION_TASKS).RemoveId(bson.ObjectIdHex(id))
	return err
}

func (m *TasksDAO) Update(task Task) error {
	err := db.C(COLLECTION_TASKS).UpdateId(task.ID, &task)
	return err
}
