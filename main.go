package main

import (
	"context"
	"flag"
	"fmt"
	gogpt "github.com/sashabaranov/go-gpt3"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

func main() {

	token := os.Getenv("CHATGPT_TOKEN")
	if token == "" {
		fmt.Println("You must have an environmental variable 'CHATGPT_TOKEN' with a ChatGPT token.")
		os.Exit(1)
	}

	var describe bool
	flag.BoolVar(&describe, "describe", false, "Describe the generated commands. (Not implemented)")

	flag.Parse()

	queryString := strings.Join(flag.Args(), " ")
	if queryString == "" {
		fmt.Println("You need to ask a question.")
		os.Exit(1)
	}

	osName := runtime.GOOS
	archName := runtime.GOARCH
	shellName := os.Getenv("SHELL")
	cwd, _ := os.Getwd()

	queryPrefix := fmt.Sprintf("Only provide suitable shell command in plain text, on one line and no description. OS=%s ARCH=%s SHELL=%s CWD=%s.", osName, archName, shellName, cwd)

	c := gogpt.NewClient(token)
	ctx := context.Background()

	req := gogpt.ChatCompletionRequest{
		Model:     gogpt.GPT3Dot5Turbo,
		MaxTokens: 10,
		Messages: []gogpt.ChatCompletionMessage{
			{
				Role:    "system",
				Content: queryPrefix,
			},
			{
				Role:    "user",
				Content: queryString,
			},
		},
	}
	resp, err := c.CreateChatCompletion(ctx, req)
	if err != nil {
		log.Println(err)
		os.Exit(2)
	}
	command := resp.Choices[0].Message.Content
	fmt.Println("=>", command)

	var input string
	fmt.Print("Run? y/(n): ")
	_, _ = fmt.Scanln(&input)

	if input != "y" && input != "Y" {
		return
	}

	cmd := exec.Command(shellName, "-c", command)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
