package common

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

// structure
type Config struct{}
type CommonServiceStruct struct{}

var systemServices CommonServiceStruct

type CommonFunc struct{}

func Callme() {
	log.Println("[Welcome]")
	log.Println("[HERE]")
}
func (c *Config) SetUpProgram() {
	// log setup
	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Failed to open log file:", err)
	}
	//defer file.Close() // don't close it
	// MultiWriter sends logs to both stdout and file
	multiWriter := io.MultiWriter(os.Stdout, file) // Configure the logger
	log.SetOutput(multiWriter)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	// dotenv load by godotenv otherwise os.Getenv lode system env
	errorInDotev := godotenv.Load() // godotenv.Load("config/.env")
	if errorInDotev != nil {
		log.Fatal("Error loading .env file")
	}
}

// AddMyCodeIntoFile appends the contents of testPlayground.go into a file under codehistory/
func AddMyCodeIntoFile(filename string) {
	if strings.TrimSpace(filename) == "" {
		log.Fatal("file name is not correct")
	}

	// Ensure destination path is valid
	destination := fmt.Sprintf("codehistory/%s", filename)

	// Append contents from source to destination
	err := AppendFile("testPlayground/testPlayground.go", destination)
	if err != nil {
		log.Fatal("Append failed:", err)
	}

	log.Println("File appended successfully")
}

// AppendFile copies everything from src file and appends it into dst file
func AppendFile(src, dst string) error {
	// Open source file for reading
	sourceFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer sourceFile.Close()

	// Open destination file in append mode (create if not exists)
	destFile, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return fmt.Errorf("failed to open destination file: %w", err)
	}
	defer destFile.Close()

	// Copy contents from src to dst (append mode)
	if _, err := io.Copy(destFile, sourceFile); err != nil {
		return fmt.Errorf("failed to copy contents: %w", err)
	}

	// Flush to disk
	if err := destFile.Sync(); err != nil {
		return fmt.Errorf("failed to sync destination file: %w", err)
	}

	return nil
}
