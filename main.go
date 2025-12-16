package main

/*
simple to-do list
fearutes to be added later
multi list
sort based on: creation time, how much urgent "taken add input", alphapeticly (ascending and decending)
all data saved to json
commads are add,start (change status to in progress), check (as done), remove, show, clear (the list), create (later for multilists), sort type_of_sort (default by creation time) order (default ascending)
maybe ui later if needed
*/
import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type task struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"Status"`
	Createdat   time.Time `json:"createdat"`
	Completedat time.Time `json:"completedat"`
}

var todos []task

func loadTasks() ([]task, error) {
	file, err := os.ReadFile("tasks.json")
	if err != nil {
		if os.IsNotExist(err) {
			return []task{}, nil
		}
		return nil, err
	}
	err = json.Unmarshal(file, &todos)
	return todos, err

}

func saveTasks(tasks []task) error {
	data, err := json.MarshalIndent(tasks, "", "	")
	if err != nil {
		return err
	}
	return os.WriteFile("tasks.json", data, 0644)
}
func addTask(title string) error {
	tasks, err := loadTasks()
	if err != nil {
		return err
	}
	toadd := task{
		ID:        len(tasks) + 1,
		Title:     title,
		Status:    "to-do",
		Createdat: time.Now(),
	}
	tasks = append(tasks, toadd)
	return saveTasks(tasks)
}
func showTasks() error {
	tasks, err := loadTasks()
	if err != nil {
		return err
	}
	if len(tasks) == 0 {
		fmt.Println("No tasks found")
		return nil
	}

}
func inprogress(taskID int) error {
	tasks, err := loadTasks()
	if err != nil {
		return err
	}
	if len(tasks) == 0 {
		return fmt.Errorf("there is no tasks in the list to check as in progress")
	}
	for i := range tasks {
		if tasks[i].ID == taskID {
			tasks[i].Status = "in-progress"
			return saveTasks(tasks)
		}
	}
	return fmt.Errorf("task not found, ID may be wrong")
}
func checkTask(taskid int) error {
	tasks, err := loadTasks()
	if err != nil {
		return err
	}
	if len(tasks) == 0 {
		return fmt.Errorf("there is no tasks in the list to check as done")
	}
	for i := range tasks {
		if tasks[i].ID == taskid {
			tasks[i].Status = "done"
			tasks[i].Completedat = time.Now()
			return saveTasks(tasks)
		}
	}
	return fmt.Errorf("task not found, ID may be wrong")

}
func removeTask(taskID int) error {
	tasks, err := loadTasks()
	if err != nil {
		return err
	}
	newTasks := []task{}
	for i := range tasks {
		if tasks[i].ID != taskID {
			newTasks = append(newTasks, tasks[i])
		}
	}
	return saveTasks(newTasks)
}
func clearList() error {
	return saveTasks([]task{})
}
func main() {

	var current_command string
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		current_command = scanner.Text()
		if current_command == "exit" {
			break
		}
		splited := strings.Fields(current_command)
		switch splited[0] {
		case "add":
			title := current_command[3:]
			addTask(title)
			fmt.Println(title, " add to your to-do list succesfully")
		case "show":
			tasks, err := loadTasks()
			fmt.Println(tasks, err)
		case "start":
			number_of_task, err := strconv.Atoi(splited[1])
			if err != nil {
				fmt.Printf("error in extracting task ID from your command %v\n", err)
				return
			}
			inprogress(number_of_task)
			fmt.Println("task no.", number_of_task, " in progress now, go finish the job")
		case "check":
			number_of_task, err := strconv.Atoi(splited[1])
			if err != nil {
				fmt.Printf("error in extracting task ID from your command %v\n", err)
				return
			}
			checkTask(number_of_task)
			fmt.Println("task no.", number_of_task, " checked as done succesfully")
		case "remove":
			number_of_task, err := strconv.Atoi(splited[1])
			if err != nil {
				fmt.Printf("error in extracting task ID from your command %v\n", err)
			}
			removeTask(number_of_task)
			fmt.Println("task no.", number_of_task, " removed succesfully from your list")
		case "clear":
			clearList()
			fmt.Println("list cleared succesfully")
		case "help":
			fmt.Println("This a to-do lists app where you can add tasks to do\n\nCommands\n\"add <task>\" to add new task as following\n\n\"show\" show table of all tasks at your list\n\n\"start <task ID>\" start doing the task, status be in porgress\n\n\"check <task ID>\" check the task as done\n\n\"remove <task ID>\" remove task from the list\n\n\"clear\" clear all data at the table\n\n\"help\" print this help message")
		}
	}
	if scanner.Err() != nil {
		panic(scanner.Err())
	}

}
