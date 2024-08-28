package main

import (
	"cmp"
	"encoding/json"
	"fmt"
	"os"
	"slices"
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
	checkDir, err := os.ReadDir("./")
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
	// fmt.Printf("Dados: %v\n", dbj)

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
}

/*
Return a boolean, true if the task exist or false if not.
*/
func (dbj DBJson) exists(id string) (Task, bool) {
	task, exists := dbj[id]
	return task, exists
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
Filter the map by your status [done | todo | in-progress]
*/
func (dbj DBJson) filterByStatus(status string) DBJson {
	d := make(DBJson, 0)
	for k, v := range dbj {
		if v.Status == status {
			d[k] = v
		}
	}
	return d
}

/*
Print the list of tasks in a AESTHETIC way (*.*)
*/
func (dbj DBJson) aestheticPrint(tasks []Task) {
	fmt.Printf("|_%6v|_%50v|_%11v|_%40v|_%40v|\n",
		"ID", "Description", "Status", "Created at", "Updated at")
	for _, v := range tasks {
		fmt.Printf("|_%6v|_%50v|_%11v|_%40v|_%40v|\n",
			v.Id, v.Description, v.Status, v.CreatedAt, v.UpdatedAt)
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
		Status:      "todo",
		CreatedAt: fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d-00:00",
			t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second()),
		UpdatedAt: fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d-00:00",
			t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second()),
	}
	dbj[strconv.Itoa(idTask)] = task
	dbj.saveToFile()
	fmt.Printf("Task added successfully (ID: %v)\n", idTask)
}

/*
Update a task description by your id.
*/
func (dbj DBJson) Update(idTask string, description string) {
	task, exists := dbj.exists(idTask)
	if !exists {
		fmt.Println("Task not found. ID: ", idTask)
		os.Exit(1)
	}
	t := time.Now()
	task.Description = description
	task.UpdatedAt = fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d-00:00",
		t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
	dbj[idTask] = task
	dbj.saveToFile()
	fmt.Printf("Task updated successfully (ID: %v)\n", idTask)
}

/*
Delete a task by your id.
*/
func (dbj DBJson) Delete(idTask string) {
	_, exists := dbj.exists(idTask)
	if !exists {
		fmt.Println("Task not found, ID: ", idTask)
		os.Exit(1)
	}
	delete(dbj, idTask)
	dbj.saveToFile()
	fmt.Printf("Task deleted successfully (ID: %v)\n", idTask)
}

/*
Mark in-progress a task by your id
*/
func (dbj DBJson) MarkInProgress(idTask string) {
	task, exists := dbj.exists(idTask)
	if !exists {
		fmt.Println("Task not found, ID: ", idTask)
		os.Exit(1)
	}
	t := time.Now()
	task.Status = "in-progress"
	task.UpdatedAt = fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d-00:00",
		t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
	dbj[idTask] = task
	dbj.saveToFile()
	fmt.Printf("Task marked in-progress successfully (ID: %v)\n", idTask)
}

/*
Mark done a task by your id
*/
func (dbj DBJson) MarkDone(idTask string) {
	task, exists := dbj.exists(idTask)
	if !exists {
		fmt.Println("Task not found, ID: ", idTask)
		os.Exit(1)
	}
	t := time.Now()
	task.Status = "done"
	task.UpdatedAt = fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d-00:00",
		t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
	dbj[idTask] = task
	dbj.saveToFile()
	fmt.Printf("Task marked done successfully (ID: %v)\n", idTask)
}

/*
List all tasks, optionaly by your status
*/
func (dbj DBJson) List(status string) {
	if status != "" {
		dbj = dbj.filterByStatus(status)
	}
	//creating a list ok Task to apply SortFunc further
	tasks := make([]Task, 0)
	for _, v := range dbj {
		tasks = append(tasks, v)
	}

	slices.SortFunc(tasks,
		func(a, b Task) int {
			return cmp.Compare(a.Id, b.Id)
		})

	dbj.aestheticPrint(tasks)
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
	//	task-cli add "Buy groceries"
	if command == "add" {
		//check for enough args to perform this action
		if len(args) < 2 {
			fmt.Println("Not enough args to perform an 'add' action. Example:\n     ./task-cli add \"My task here\"")
			os.Exit(1)
		}
		task := args[1]
		dbJson.Add(task)
	}

	// UPDATE COMMAND
	//	task-cli update 1 "Buy groceries and cook dinner"
	if command == "update" {
		if len(args) < 3 {
			fmt.Println("Not enough args to perform an 'update' action. Example:\n     ./task-cli update 1 \"Buy groceries and cook dinner\"")
			os.Exit(1)
		}
		idTask := args[1]
		description := args[2]
		dbJson.Update(idTask, description)
	}

	// DELETE COMMAND
	//	task-cli delete 1
	if command == "delete" {
		if len(args) < 2 {
			fmt.Println("Not enough args to perform an 'delete' action. Example:\n     ./task-cli delete 1")
			os.Exit(1)
		}
		idTask := args[1]
		dbJson.Delete(idTask)
	}

	// mark-in-progress COMMAND
	// task-cli mark-in-progress 1
	if command == "mark-in-progress" {
		if len(args) < 2 {
			fmt.Println("Not enough args to perform an 'mark-in-progress' action. Example:\n     ./task-cli mark-in-progress 1")
			os.Exit(1)
		}
		idTask := args[1]
		dbJson.MarkInProgress(idTask)
	}

	// mark-done COMMAND
	// task-cli mark-done 1
	if command == "mark-done" {
		if len(args) < 2 {
			fmt.Println("Not enough args to perform an 'mark-done' action. Example:\n     ./task-cli mark-done 1")
			os.Exit(1)
		}
		idTask := args[1]
		dbJson.MarkDone(idTask)
	}

	// list COMMAND
	// task-cli list [done | todo | in-progress]
	if command == "list" {
		status := ""
		if len(args) > 1 {
			status = args[1]
		}
		if status != "" && status != "done" && status != "todo" && status != "in-progress" {
			fmt.Println("Not enough args to perform an 'list' action by status, please follow the example ahead. Example:\n     ./task-cli list [done | todo | in-progress]")
			os.Exit(1)
		}
		dbJson.List(status)
	}
}
