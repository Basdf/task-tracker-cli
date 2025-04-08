package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"task-tracker-cli/models"
	"task-tracker-cli/services"
)

func processArgs(args []string) []string {
	if len(args) < 2 {
		return args
	}

	subcommand := args[1]
	subcommandArgs := args[2:]

	if len(subcommandArgs) == 0 {
		return args
	}

	if strings.HasPrefix(subcommandArgs[0], "-") {
		return args
	}

	switch subcommand {
	case "add":
		if len(subcommandArgs) > 0 {
			return []string{args[0], subcommand, "-desc", subcommandArgs[0]}
		}
	case "update":
		if len(subcommandArgs) >= 2 {
			return []string{args[0], subcommand, "-id", subcommandArgs[0], "-desc", subcommandArgs[1]}
		}
	case "delete", "mark-in-progress", "mark-done":
		if len(subcommandArgs) > 0 {
			return []string{args[0], subcommand, "-id", subcommandArgs[0]}
		}
	case "list":
		if len(subcommandArgs) > 0 {
			return []string{args[0], subcommand, "-status", subcommandArgs[0]}
		}
	case "help":
		if len(subcommandArgs) > 0 {
			return []string{args[0], subcommand, "-flag", subcommandArgs[0]}
		}
	}

	return args
}

func main() {
	// Subcommands
	addCmd := flag.NewFlagSet("add", flag.ExitOnError)
	addDescription := addCmd.String("desc", "", "Task description")
	addCmd.StringVar(addDescription, "d", "", "Task description (shorthand)")

	updateCmd := flag.NewFlagSet("update", flag.ExitOnError)
	updateID := updateCmd.Int("id", 0, "Task ID")
	updateDescription := updateCmd.String("desc", "", "New task description")
	updateCmd.IntVar(updateID, "i", 0, "Task ID (shorthand)")
	updateCmd.StringVar(updateDescription, "d", "", "New task description (shorthand)")

	deleteCmd := flag.NewFlagSet("delete", flag.ExitOnError)
	deleteID := deleteCmd.Int("id", 0, "Task ID")
	deleteCmd.IntVar(deleteID, "i", 0, "Task ID (shorthand)")

	markInProgressCmd := flag.NewFlagSet("mark-in-progress", flag.ExitOnError)
	markInProgressID := markInProgressCmd.Int("id", 0, "Task ID")
	markInProgressCmd.IntVar(deleteID, "i", 0, "Task ID (shorthand)")

	markDoneCmd := flag.NewFlagSet("mark-done", flag.ExitOnError)
	markDoneID := markDoneCmd.Int("id", 0, "Task ID")
	markDoneCmd.IntVar(deleteID, "i", 0, "Task ID (shorthand)")

	listCmd := flag.NewFlagSet("list", flag.ExitOnError)
	listStatus := listCmd.String("status", "", "Filter tasks by status (todo, in-progress, done)")
	listCmd.StringVar(listStatus, "s", "", "Filter tasks by status (shorthand)")

	helpCmd := flag.NewFlagSet("help", flag.ExitOnError)
	helpFlag := helpCmd.String("flag", "", "Command to display help for")
	helpCmd.StringVar(helpFlag, "f", "", "Command to display help for (shorthand)")

	os.Args = processArgs(os.Args)

	if len(os.Args) < 2 {
		fmt.Println("Expected subcommands: add, update, delete, mark-in-progress, mark-done, list, help")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "add":
		addCmd.Parse(os.Args[2:])
		if *addDescription == "" {
			fmt.Println("Error: Task description is required")
			addCmd.PrintDefaults()
			os.Exit(1)
		}
		services.AddTask(*addDescription)

	case "update":
		updateCmd.Parse(os.Args[2:])
		if *updateID == 0 || *updateDescription == "" {
			fmt.Println("Error: Task ID and description are required")
			updateCmd.PrintDefaults()
			os.Exit(1)
		}
		services.UpdateTask(fmt.Sprintf("%d", *updateID), *updateDescription)

	case "delete":
		deleteCmd.Parse(os.Args[2:])
		if *deleteID == 0 {
			fmt.Println("Error: Task ID is required")
			deleteCmd.PrintDefaults()
			os.Exit(1)
		}
		services.DeleteTask(fmt.Sprintf("%d", *deleteID))

	case "mark-in-progress":
		markInProgressCmd.Parse(os.Args[2:])
		if *markInProgressID == 0 {
			fmt.Println("Error: Task ID is required")
			markInProgressCmd.PrintDefaults()
			os.Exit(1)
		}
		services.UpdateTaskStatus(fmt.Sprintf("%d", *markInProgressID), models.StatusInProgress)

	case "mark-done":
		markDoneCmd.Parse(os.Args[2:])
		if *markDoneID == 0 {
			fmt.Println("Error: Task ID is required")
			markDoneCmd.PrintDefaults()
			os.Exit(1)
		}
		services.UpdateTaskStatus(fmt.Sprintf("%d", *markDoneID), models.StatusDone)

	case "list":
		listCmd.Parse(os.Args[2:])
		services.ListTasks(*listStatus)
	case "help":
		helpCmd.Parse(os.Args[2:])
		switch *helpFlag {
		case "add":
			fmt.Println("Add a new task to the tracking system")
			fmt.Println("\nUsage:")
			fmt.Println("  task add -desc <description>     # Using full flag")
			fmt.Println("  task add -d <description>        # Using short flag")
			fmt.Println("  task add <description>          # Direct usage")
			fmt.Println("\nExamples:")
			fmt.Println("  task add -desc \"Complete monthly report\"")
			fmt.Println("  task add -d \"Complete monthly report\"")
			fmt.Println("  task add \"Complete monthly report\"")
			addCmd.PrintDefaults()
		case "update":
			fmt.Println("Update an existing task's description")
			fmt.Println("\nUsage:")
			fmt.Println("  task update -id <task_id> -desc <new_description>      # Using full flags")
			fmt.Println("  task update -i <task_id> -d <new_description>         # Using short flags")
			fmt.Println("  task update <task_id> <new_description>             # Direct usage")
			fmt.Println("\nExamples:")
			fmt.Println("  task update -id 1 -desc \"Review monthly report\"")
			fmt.Println("  task update -i 1 -d \"Review monthly report\"")
			fmt.Println("  task update 1 \"Review monthly report\"")
			updateCmd.PrintDefaults()
		case "delete":
			fmt.Println("Delete a task from the system")
			fmt.Println("\nUsage:")
			fmt.Println("  task delete -id <task_id>      # Using full flag")
			fmt.Println("  task delete -i <task_id>       # Using short flag")
			fmt.Println("  task delete <task_id>          # Direct usage")
			fmt.Println("\nExamples:")
			fmt.Println("  task delete -id 1")
			fmt.Println("  task delete -i 1")
			fmt.Println("  task delete 1")
			deleteCmd.PrintDefaults()
		case "mark-in-progress":
			fmt.Println("Mark a task as in progress")
			fmt.Println("\nUsage:")
			fmt.Println("  task mark-in-progress -id <task_id>      # Using full flag")
			fmt.Println("  task mark-in-progress -i <task_id>       # Using short flag")
			fmt.Println("  task mark-in-progress <task_id>          # Direct usage")
			fmt.Println("\nExamples:")
			fmt.Println("  task mark-in-progress -id 1")
			fmt.Println("  task mark-in-progress -i 1")
			fmt.Println("  task mark-in-progress 1")
			markInProgressCmd.PrintDefaults()
		case "mark-done":
			fmt.Println("Mark a task as completed")
			fmt.Println("\nUsage:")
			fmt.Println("  task mark-done -id <task_id>      # Using full flag")
			fmt.Println("  task mark-done -i <task_id>       # Using short flag")
			fmt.Println("  task mark-done <task_id>          # Direct usage")
			fmt.Println("\nExamples:")
			fmt.Println("  task mark-done -id 1")
			fmt.Println("  task mark-done -i 1")
			fmt.Println("  task mark-done 1")
			markDoneCmd.PrintDefaults()
		case "list":
			fmt.Println("List all tasks or filter by status")
			fmt.Println("\nUsage:")
			fmt.Println("  task list -status <status>      # Using full flag")
			fmt.Println("  task list -s <status>           # Using short flag")
			fmt.Println("  task list <status>              # Direct usage")
			fmt.Println("\nExamples:")
			fmt.Println("  task list                       # Show all tasks")
			fmt.Println("  task list -status todo          # Using full flag")
			fmt.Println("  task list -s in-progress        # Using short flag")
			fmt.Println("  task list done                  # Direct usage")
			listCmd.PrintDefaults()
		case "help":
			fmt.Println("Display help for specific commands")
			fmt.Println("\nUsage:")
			fmt.Println("  task help -flag <command>      # Using full flag")
			fmt.Println("  task help -f <command>         # Using short flag")
			fmt.Println("  task help <command>            # Direct usage")
			fmt.Println("\nExamples:")
			fmt.Println("  task help -flag add")
			fmt.Println("  task help -f add")
			fmt.Println("  task help add")
			helpCmd.PrintDefaults()
		default:
			fmt.Println("Task Tracking System - Available Commands:")
			fmt.Println("\nCommands:")
			fmt.Println("  add             Add a new task")
			fmt.Println("  update          Update an existing task")
			fmt.Println("  delete          Delete a task")
			fmt.Println("  mark-in-progress Mark a task as in progress")
			fmt.Println("  mark-done       Mark a task as completed")
			fmt.Println("  list            List tasks")
			fmt.Println("  help            Show this help or specific command help")
			fmt.Println("\nTo see detailed help for a specific command:")
			fmt.Println("  task help -flag <command>")
			fmt.Println("\nExample: task help -flag add")
		}

	default:
		fmt.Println("Error: Unknown command")
		os.Exit(1)
	}
}
