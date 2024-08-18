package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
)

var automationClient *APIClient

func main() {
	automationClient = &APIClient{BaseURL: "http://localhost:8080/api"}

	luaInit()

	runCode("init.lua")

	model := flag.String("m", "l", "Set model (default: l)")
	temp := flag.Float64("p", 0.1, "Set temperature (default: 0.1)")
	prompt := flag.String("t", "ct", "Set prompt template (default: c)")
	cmd := flag.String("c", "", "Command string to execute")
	cf := flag.String("cf", "", "Context file")
	cd := flag.String("cd", "", "Context data")
	clean := flag.Bool("clean", false, "Clean output (default: false)")
	save := flag.String("o", "", "Save output to file")
	debug := flag.Bool("d", false, "Enable debug mode (default: false)")
	script := flag.String("s", "", "Lua script to run")

	flag.Parse()

	if *cmd == "" {
		flag.Usage()
		os.Exit(1)
	}

	var input string = ""

	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		reader := bufio.NewReader(os.Stdin)
		data, _ := reader.ReadString('\n')
		input = data
	}

	var contextData string = input

	if *cd != "" {
		contextData = *cd
	}

	if *cf != "" {
		contextData, _ = loadContext(*cf)
	}
	fmt.Println(*prompt)
	promptTemplateStr, _ := getPromptTemplate(*prompt)
	fmt.Println(promptTemplateStr)
	promptTemplateGen := generatePromptStr(promptTemplateStr, contextData, *cmd)
	fmt.Println(promptTemplateGen)

	if *debug {
		print(promptTemplateGen)
	}

	completion := runPrompt(*model, float32(*temp), promptTemplateGen)

	if *clean {
		completion = cleanOutputTxt(completion)
	}

	if *save != "" {
		saveToFile(completion, *save)
	}
	if *script != "" {
		setLuaPromptOutput(completion)
		runCode(*script)
	}
	if *debug {
		print(completion)
	}
}
func runPrompt(modelName string, temperature float32, prompt string) string {
	ctx := context.Background()
	modelName = getModelFunction(modelName) // Assuming this function gets the correct model
	llm, _ := ollama.New(ollama.WithModel(modelName))

	completion, err := llms.GenerateFromSinglePrompt(
		ctx,
		llm,
		prompt,
		llms.WithTemperature(float64(temperature)),
		llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
			fmt.Print(string(chunk))
			return nil
		}),
	)
	if err != nil {
		log.Fatal(err)
	}

	return completion
}
