package main

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/basdf/tast-tracker-cli/cmd"
)

func main() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Enter the filename to save tasks")
	fmt.Print("> ")
	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println(err)
	}
	cmdClient := cmd.New(strings.TrimSpace(input))
	defer func(cmdClient *cmd.CMD) {
		if err := recover(); err != nil {
			fmt.Println("Error:", err)
		}
		fmt.Println("Exiting...")
		cmdClient.Close()
	}(cmdClient)

	go func() {
		for {
			fmt.Println("Enter command")
			fmt.Print("> ")
			input, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println("Error reading input:", err)
			}
			terminated := cmdClient.ExecuteCommand(input, sigChan)
			if terminated {
				return
			}
		}
	}()

	<-sigChan
}
