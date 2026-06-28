package main

import (
	"fmt"
	"strings"

	prompt "github.com/c-bata/go-prompt"
)

func executor(input string) {
	input = strings.TrimSpace(input)

	switch input {
	case "login":
		fmt.Println("Logging in...")

	case "register":
		fmt.Println("Registering user...")

	case "logout":
		fmt.Println("Logging out...")

	case "help":
		fmt.Println("Available commands:")
		fmt.Println("  login")
		fmt.Println("  register")
		fmt.Println("  logout")
		fmt.Println("  exit")

	case "exit":
		fmt.Println("Goodbye!")
		panic("exit") // Simple way to stop the example

	default:
		if input != "" {
			fmt.Printf("Unknown command: %s\n", input)
		}
	}
}

func completer(d prompt.Document) []prompt.Suggest {
	suggestions := []prompt.Suggest{
		{Text: "login", Description: "Login to your account"},
		{Text: "register", Description: "Create a new account"},
		{Text: "logout", Description: "Logout of your account"},
		{Text: "help", Description: "Show available commands"},
		{Text: "exit", Description: "Exit application"},
	}

	return prompt.FilterHasPrefix(
		suggestions,
		d.GetWordBeforeCursor(),
		true,
	)
}

func main() {
	fmt.Println("Simple CLI")
	fmt.Println("Press TAB for autocomplete")
	fmt.Println("Type 'help' for commands")

	p := prompt.New(
		executor,
		completer,
		prompt.OptionPrefix(">>> "),
		prompt.OptionTitle("CLI Authentication"),
	)

	p.Run()
}