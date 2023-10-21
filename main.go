package main

import (
	"flag"
	"fmt"
	"io/fs"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

type Summary struct {
	fileName       string
	inputFileSize  int64
	outputFileSize int64
}

func main() {
	timeStart := time.Now()

	inputPtr := flag.String("input", "", "Input file or folder")
	outputPtr := flag.String("output", "", "Output file or folder")

	flag.Parse()

	inputValue := *inputPtr
	outputValue := *outputPtr

	if inputValue == "" {
		fmt.Println("Input file or folder is required. Example: `go run main.go -input /path/to/input -output /path/to/output`")
		return
	}

	if outputValue == "" {
		fmt.Println("Output file or folder is required. Example: `go run main.go -input /path/to/input -output /path/to/output`")
		return
	}

	var _, err = exec.LookPath("ffmpeg")
	if err != nil {
		fmt.Println("ffmpeg not found")
		return
	}

	fmt.Println("Input value:", inputValue)
	fmt.Println("Output value:", outputValue)

	if _, err := os.Stat(outputValue); os.IsNotExist(err) {
		fmt.Println("Creating output folder", outputValue)

		err := os.MkdirAll(outputValue, 0755)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	var summaries []Summary

	var walkErr = filepath.Walk(inputValue, func(path string, info fs.FileInfo, err error) error {
		ALLOWED_EXT := []string{".mp4", ".mkv", ".avi", ".mov", ".wmv", ".flv", ".webm", ".m4v"}

		if !arrayContains(ALLOWED_EXT, filepath.Ext(path)) {
			return nil
		}

		outputFile := filepath.Join(outputValue, filepath.Base(path))
		convertFile(path, outputFile)

		var summary = Summary{
			fileName:       filepath.Base(path),
			inputFileSize:  info.Size(),
			outputFileSize: 0,
		}

		if _, err := os.Stat(outputFile); err == nil {
			outputFileInfo, _ := os.Stat(outputFile)
			summary.outputFileSize = outputFileInfo.Size()
		}

		summaries = append(summaries, summary)

		return nil
	})

	if walkErr != nil {
		fmt.Println(walkErr)
		return
	}

	timeEnd := time.Now()
	duration := timeEnd.Sub(timeStart) / 1000000000

	fmt.Println()
	fmt.Printf("\033[1;32m")
	fmt.Println("Finished in", duration, "seconds")
	for _, summary := range summaries {
		fmt.Println("File:", summary.fileName)
		fmt.Println("Input file size:", convertToHumanFileSize(summary.inputFileSize))
		fmt.Println("Output file size:", convertToHumanFileSize(summary.outputFileSize))
		fmt.Println("Saved:", convertToHumanFileSize(summary.inputFileSize-summary.outputFileSize))
		fmt.Println()
	}
	fmt.Printf("\033[0m")
}

func convertFile(file string, outputFile string) {
	fmt.Println("Converting file", file)
	fmt.Println("Output file", outputFile)

	file = filepath.Clean(file)
	outputFile = filepath.Clean(outputFile)

	if _, err := os.Stat(outputFile); err == nil {
		fmt.Println("File already exists, skipping")
		return
	}

	cmdArgs := []string{"-i", file, "-vf", "scale=640:480", "-c:a", "aac", "-strict", "experimental", "-b:a", "128k", outputFile}

	cmd := exec.Command("ffmpeg", cmdArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()

	if err != nil {
		fmt.Println(err)
		return
	}
}

func arrayContains(array []string, value string) bool {
	for _, v := range array {
		if v == value {
			return true
		}
	}

	return false
}

func convertToHumanFileSize(size int64) string {
	if size < 1024 {
		return fmt.Sprintf("%d B", size)
	}

	const unit = 1024
	exp := int64(math.Log(float64(size)) / math.Log(float64(unit)))
	pre := "KMGTPE"[exp-1]

	return fmt.Sprintf("%.1f %ciB", float64(size)/math.Pow(float64(unit), float64(exp)), pre)
}
