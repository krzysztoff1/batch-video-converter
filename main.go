package main

import (
	"flag"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
)


func main() {
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

	fmt.Println("Input value:", inputValue)
	fmt.Println("Output value:", outputValue)

	if _, err := os.Stat(outputValue); os.IsNotExist(err) {
		fmt.Println("Creating output folder", outputValue)
		os.Mkdir(outputValue, 0755)
	}

	filepath.Walk(inputValue, func(path string, info fs.FileInfo, err error) error {
		fmt.Println("Input file:", path)
		if filepath.Ext(path) != ".mp4" {
			return nil
		}

		outputFile := filepath.Join(outputValue, filepath.Base(path))
		fmt.Println("Output file:", filepath.Base(path))

		convertFile(path, outputFile)

		return nil
	})

	fmt.Println("Done")
}

func convertFile(file string, outputFile string) {
	fmt.Println("Converting file", file)
	fmt.Println("Output file", outputFile)

	var _, err = exec.LookPath("ffmpeg")
	fmt.Println("ffmpeg path", err)

	if err != nil {
		fmt.Println("ffmpeg not found")
		return
	}

	file = filepath.Clean(file)
	file = filepath.ToSlash(file)

	cmdArgs := []string{"-i", file, "-vf", "scale=640:480", "-c:a", "aac", "-strict", "experimental", "-b:a", "128k", outputFile}

	cmd := exec.Command("ffmpeg", cmdArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}
