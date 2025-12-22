package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

type task struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	Priority    string    `json:"priority"`
	Createdat   time.Time `json:"createdat"`
	Completedat time.Time `json:"completedat"`
}

func loadTasks() ([]task, error) {
	var tasks []task
	file, err := os.ReadFile("tasks.json")
	if err != nil {
		if os.IsNotExist(err) {
			return []task{}, nil
		}
		return nil, err
	}
	err = json.Unmarshal(file, &tasks)
	return tasks, err
}

func saveTasks(tasks []task) error {
	data, err := json.MarshalIndent(tasks, "", "	")
	if err != nil {
		return err
	}
	return os.WriteFile("tasks.json", data, 0644)
}
func addTask(title string, priority string) error {
	tasks, err := loadTasks()
	if err != nil {
		return err
	}
	maxID := 0
	for i := range tasks {
		if tasks[i].ID > maxID {
			maxID = tasks[i].ID
		}
	}
	toadd := task{
		ID:        maxID + 1,
		Title:     title,
		Status:    "to-do",
		Priority:  priority,
		Createdat: time.Now(),
	}
	tasks = append(tasks, toadd)
	return saveTasks(tasks)
}
func timeAgo(tim time.Time) string {
	age := time.Since(tim)

	switch {
	case age < time.Minute:
		return fmt.Sprintf("%d seconds ago", int(age.Seconds()))

	case age < time.Hour:
		return fmt.Sprintf("%d minutes ago", int(age.Minutes()))

	case age < 24*time.Hour:
		return fmt.Sprintf("%d hours ago", int(age.Hours()))

	default:
		return fmt.Sprintf("%d days ago", int(age.Hours()/24))
	}
}
func showTasks(filter string) error {
	tasks, err := loadTasks()
	if err != nil {
		return err
	}
	if len(tasks) == 0 {
		fmt.Println("No tasks found")
		return nil
	}
	tabl := table.NewWriter()
	tabl.SetOutputMirror(os.Stdout)
	tabl.AppendHeader(table.Row{
		"ID", "Title", "Status", "Priority", "Age", "Created At", "Completed At",
	})
	for i := range tasks {
		if filter != "" {
			if tasks[i].Status != filter && tasks[i].Priority != filter {
				continue
			}
		}
		status := tasks[i].Status
		switch status {
		case "to-do":
			status = text.FgRed.Sprint("TO-DO")
		case "in-progress":
			status = text.FgYellow.Sprint("IN-PROGRESS")
		case "done":
			status = text.FgGreen.Sprint("DONE")
		}
		priority := tasks[i].Priority
		switch priority {
		case "routine":
			priority = text.FgGreen.Sprint("Routine")
		case "urgent":
			priority = text.FgRed.Sprint("Urgent")
		case "important":
			priority = text.FgYellow.Sprint("important")
		case "luxury":
			priority = text.FgBlue.Sprint("Luxury")
		default:
			priority = text.FgHiWhite.Sprint(priority)
		}
		title := tasks[i].Title
		if tasks[i].Status == "done" {
			title = text.FgHiBlack.Sprint(tasks[i].Title)
		}
		completedat := " "
		if !tasks[i].Completedat.IsZero() {
			completedat = tasks[i].Completedat.Format("2006-01-02 15:04:05")
		}
		age := timeAgo(tasks[i].Createdat)
		tabl.AppendRow(table.Row{tasks[i].ID, title, status, priority, age, tasks[i].Createdat.Format("2006-01-02 15:04:05"), completedat})
	}
	tabl.SetStyle(table.StyleDouble)
	tabl.Render()
	return nil
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

func changePriority(taskID int, newpriorrity string) error {
	tasks, err := loadTasks()
	if err != nil {
		return err
	}
	for i := range tasks {
		if tasks[i].ID == taskID {
			tasks[i].Priority = newpriorrity
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
			if len(splited) < 2 {
				fmt.Println("missing task content")
				continue
			}
			if splited[1] == "priority" {
				if splited[2] == "urgent" || splited[2] == "important" || splited[2] == "luxury" {
					title := strings.Join(splited[3:], " ")
					err := addTask(title, splited[2])
					if err != nil {
						fmt.Println("Error in adding your item: ", err)
						continue
					}
					fmt.Println(title, " added to your to-do list succesfully")
					continue
				}
			}
			title := current_command[3:]
			priority := "routine"
			err := addTask(title, priority)
			if err != nil {
				fmt.Println("Error in adding your item: ", err)
				continue
			}
			fmt.Println(title, " added to your to-do list succesfully")

		case "show":
			filter := ""
			if len(splited) > 1 {
				filter = strings.ToLower(splited[1])
			}
			err := showTasks(filter)
			if err != nil {
				fmt.Println("error: ", err)
				continue
			}

		case "start":
			if len(splited) < 2 {
				fmt.Println("missing task ID")
				continue
			}
			number_of_task, err := strconv.Atoi(splited[1])
			if err != nil {
				fmt.Printf("error in extracting task ID from your command %v\n", err)
				continue
			}
			err1 := inprogress(number_of_task)
			if err1 != nil {
				fmt.Println("error: ", err1)
				continue
			}
			fmt.Println("task no.", number_of_task, " in progress now, go finish the job")

		case "check":
			if len(splited) < 2 {
				fmt.Println("missing task ID")
				continue
			}
			number_of_task, err := strconv.Atoi(splited[1])
			if err != nil {
				fmt.Printf("error in extracting task ID from your command %v\n", err)
				continue
			}
			err1 := checkTask(number_of_task)
			if err1 != nil {
				fmt.Println("error: ", err1)
				continue
			}
			fmt.Println("task no.", number_of_task, " checked as done succesfully")

		case "priority":
			if len(splited) < 3 {
				fmt.Println("missing task ID or new priority")
				continue
			}

			number_off_task, err := strconv.Atoi(splited[1])
			if err != nil {
				fmt.Printf("error in extracting task ID from your command %v\n", err)
				continue
			}

			err1 := changePriority(number_off_task, splited[2])
			if err1 != nil {
				fmt.Println("error: ", err1)
				continue
			}
			fmt.Println("Priority of task no.", number_off_task, " has changed to ", splited[2])

		case "remove":
			if len(splited) < 2 {
				fmt.Println("missing task ID")
				continue
			}
			number_of_task, err := strconv.Atoi(splited[1])
			if err != nil {
				fmt.Printf("error in extracting task ID from your command %v\n", err)
				continue
			}
			err1 := removeTask(number_of_task)
			if err1 != nil {
				fmt.Println("error: ", err1)
				continue
			}
			fmt.Println("task no.", number_of_task, " removed succesfully from your list")

		case "clear":
			err := clearList()
			if err != nil {
				fmt.Println("error: ", err)
				continue
			}
			fmt.Println("list cleared succesfully")

		case "help":
			fmt.Println("This a to-do lists app where you can add tasks to do\n\nCommands\n\"add <task>\" to add new task as following \noptional set priority: \"add priority <priority level> <task>\" where prioty levels in order are:\nurgent\nimportant\nroutine\nluxury\nnote: the default is routine level where you don't have to set priority manually\n\n\"show\" show table of all tasks at your list\noptional filter by priority or status \"show <filter>\"\n\n\"start <task ID>\" start doing the task, status be in porgress\n\n\"check <task ID>\" check the task as done\n\n\"priority <task ID> <new value>\" change default priority that routine to new one from (urgent ,important, luxury)\n\n\"remove <task ID>\" remove task from the list\n\n\"clear\" clear all data at the table\n\n\"help\" print this help message")
		}
	}
	if scanner.Err() != nil {
		panic(scanner.Err())
	}

}
