package main

import (
	"encoding/base64"
	"fmt"
	"image/png"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var doDebug bool = true

func getPromptTemplate(prompt string) (string, error) {
	filePath := filepath.Join("prompts", prompt+".txt")
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func Debug(message string) {
	if doDebug {
		fmt.Println(message)
	}
}

func saveToFile(content, filename string) error {
	return ioutil.WriteFile(filename, []byte(content), 0644)
}

func cleanOutputTxt(input string) string {
	lines := strings.Split(input, "\n")
	var cleanedLines []string
	for _, line := range lines {
		if !strings.HasPrefix(line, "```") {
			cleanedLines = append(cleanedLines, line)
		}
	}
	return strings.Join(cleanedLines, "\n")
}
func getModelFunction(model string) string {
	switch model {
	case "l":
		return "llama3.1"
	case "d":
		return "deepseek-coder-v2"
	case "c":
		return "codestral"
	case "s":
		return "sqlcoder"
	case "cl":
		return "codellama"
	default:
		return "llama3.1"
	}
}
func generatePromptStr(promptTemplate, context, query string) string {
	return fmt.Sprintf(promptTemplate, context, query)
}
func loadContext(filename string) (string, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func SaveBase64ImageToPNG(base64Image, filePath string) error {
	// Remove the data URL prefix if present
	if strings.Contains(base64Image, "data:image/png;base64,") {
		base64Image = strings.Replace(base64Image, "data:image/png;base64,", "", 1)
	}

	// Decode the base64 string to a byte slice
	imageData, err := base64.StdEncoding.DecodeString(base64Image)
	if err != nil {
		return fmt.Errorf("failed to decode base64 image: %v", err)
	}

	// Create a file to save the PNG image
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer file.Close()

	// Decode the byte slice into an image
	img, err := png.Decode(strings.NewReader(string(imageData)))
	if err != nil {
		return fmt.Errorf("failed to decode PNG image: %v", err)
	}

	// Encode the image to the file as PNG
	if err := png.Encode(file, img); err != nil {
		return fmt.Errorf("failed to encode image to PNG: %v", err)
	}

	return nil
}

func readFileAsBase64(filePath string) (string, error) {
	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	// Read the file content into a byte slice
	fileData, err := ioutil.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %v", err)
	}

	// Encode the file data to base64
	base64Data := base64.StdEncoding.EncodeToString(fileData)
	return base64Data, nil
}
