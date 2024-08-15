package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strconv"
	"time"
)

type Task struct {
	Id          int    `json:"id"`
	Description string `json:"description"`
	Status      string `json:"status"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}

type DBJson map[string]Task

var dbKeys []int

/*
Initialize the internal database to save the tasks,
if the file not exists then create a new one.
After that load the file into an struct
*/
func (dbj DBJson) Initialize() error {
	//if not found the file db.json create a new one
	checkDir, err := os.ReadDir("/home/fornazierr/Documents/projects/golang/roadmap-projects/task-tracker")
	if err != nil {
		return err
	}
	find := false
	for i := 0; i < len(checkDir); i++ {
		if checkDir[i].Name() == "db.json" {
			find = true
		}
	}
	if !find {
		f, err := os.Create("db.json")
		if err != nil {
			return err
		}
		defer f.Close()
	} else {
		byt, err := os.ReadFile("db.json")
		if err != nil {
			return err
		}
		if err := json.Unmarshal(byt, &dbj); err != nil {
			return err
		}
	}
	fmt.Printf("Dados: %v\n", dbj)

	dbj.organizeKeys()
	return nil
}

/*
Organize the cli key list
*/
func (dbj DBJson) organizeKeys() {
	for k := range dbj {
		i, err := strconv.Atoi(k)
		if err != nil {
			panic(err)
		}
		dbKeys = append(dbKeys, i)
	}
	sort.Ints(dbKeys)
	fmt.Printf("OrganizeKeys: %v\n", dbKeys)
}

/*
Return a boolean, true if the task exist or false if not.
*/
func (dbj DBJson) exists(id string) bool {
	_, exists := dbj[id]
	return exists
}

/*
Save/write the current map in to a json file
*/
func (dbj DBJson) saveToFile() {
	byt, err := json.Marshal(dbj)
	if err != nil {
		panic(err)
	}
	err = os.WriteFile("db.json", byt, 0664)
	if err != nil {
		panic(err)
	}
}

/*
Add a new task
*/
func (dbj DBJson) Add(description string) {
	var idTask int

	if len(dbKeys) == 0 {
		idTask = 1
	} else {
		idTask = dbKeys[len(dbKeys)-1] + 1
	}
	dbKeys = append(dbKeys, idTask)

	t := time.Now()
	task := Task{
		Id:          idTask,
		Description: description,
		Status:      "created",
		CreatedAt: fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d-00:00",
			t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second()),
		UpdatedAt: fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d-00:00",
			t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second()),
	}

	dbj[strconv.Itoa(idTask)] = task

	fmt.Printf("Task added successfully (ID: %v)\n", idTask)

	dbj.saveToFile()
}

func main() {
	//initializing the database
	var dbJson DBJson
	dbJson = make(map[string]Task, 0)
	err := dbJson.Initialize()
	if err != nil {
		fmt.Errorf("Failed initializing the database: %v\n", err)
	}

	//gattering the args
	args := os.Args[1:]
	//filtering the command
	command := args[0]

	//check if the argument is a valid one listed above
	// add update delete list mark-in-progress mark-done
	if command != "add" && command != "update" && command != "delete" && command != "list" && command != "mark-in-progress" && command != "mark-done" {
		fmt.Println("\nCommand not found. Check if are one of the following above:\n   add\n   update\n   delete\n   mark-in-progress\n   mark-done\n   list")
	}

	// ADD COMMAND
	if command == "add" {
		//check for enough args to perform this action
		if len(args) < 2 {
			fmt.Println("Not enough args to perform an add action. Example:\n     ./task-cli add \"My task here\"")
			os.Exit(1)
		}
		task := args[1]
		dbJson.Add(task)
	}
}
