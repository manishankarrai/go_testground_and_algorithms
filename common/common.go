package common

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"

	"test/constant"
	connection "test/db/mongo"
	"test/db/mongoRepo"
	"test/models"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

// --- Structures ---

type Config struct{}

// --- Custom Logger Bridge ---

type MongoWriter struct {
	FileOut io.Writer
	StdOut  io.Writer
}

type CommonServiceStruct struct {
	LogRepo *mongoRepo.LogDetailRepository
	wg      sync.WaitGroup // Tracks active background writes
}

var Services CommonServiceStruct

func (mw *MongoWriter) Write(p []byte) (n int, err error) {
	n, err = mw.FileOut.Write(p)
	mw.StdOut.Write(p)

	if Services.LogRepo != nil {
		content := string(p)
		// Increment the WaitGroup counter
		Services.wg.Add(1)

		go func(msg string) {
			// Decrement the counter when the function finishes
			defer Services.wg.Done()

			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()

			_ = Services.LogRepo.Create(ctx, models.ExecutionLog{
				LogID:     uuid.New().String(),
				Level:     constant.LogTypeAuto,
				Message:   strings.TrimSpace(msg),
				Timestamp: time.Now(),
			})
		}(content)
	}
	return n, err
}

// Inside package common

func (s *CommonServiceStruct) Close() {
	// 1. Wait for all background 'go func' logs to finish
	//log.Println("Syncing logs to MongoDB...")
	s.wg.Wait()

	// 2. Get the client from the connection package and disconnect
	client := connection.GetClient()
	if client != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := client.Disconnect(ctx); err != nil {
			log.Printf("Error during MongoDB disconnect: %v", err)
		}
		// else {
		// 	log.Println("MongoDB connection closed gracefully.")
		// }
	}
}

// --- Core Functions ---

func (c *Config) SetUpProgram() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found")
	}

	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Failed to open log file:", err)
	}

	// Connect to Mongo
	connection.ConnectMongo()
	// Initialize the repository
	Services.LogRepo = mongoRepo.NewLogDetailRepository()

	customWriter := &MongoWriter{
		FileOut: file,
		StdOut:  os.Stdout,
	}
	log.SetOutput(customWriter)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	//log.Println("[System Ready] File and MongoDB logging linked")
}

// LogToAll is for manual, structured logging linked to a specific RunDetailId
func LogToAll(runID string, level string, message string) {
	// Adding (ManualLog) tag so the MongoWriter doesn't create a duplicate "AUTO" entry
	log.Printf("[%s] %s (RunID: %s) (ManualLog)", level, message, runID)

	if Services.LogRepo != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		err := Services.LogRepo.Create(ctx, models.ExecutionLog{
			LogID:       uuid.New().String(),
			RunDetailId: runID,
			Level:       level,
			Message:     message,
			Timestamp:   time.Now(),
		})
		if err != nil {
			// Use fmt here to avoid triggering the custom writer logic recursively on error
			fmt.Printf("CRITICAL ERROR: Failed to save to DB: %v\n", err)
		}
	}
}

// --- File Handling Utilities ---

func AddMyCodeIntoFile(filename string) {
	if strings.TrimSpace(filename) == "" {
		log.Fatal("file name is not correct")
	}
	filename = GiveFilenameByRemovingSpaces(filename)
	log.Printf("file is created with name %s", filename)
	destination := fmt.Sprintf("codehistory/%s", filename)

	err := AppendFile("testPlayground/testPlayground.go", destination)
	if err != nil {
		log.Fatal("Append failed:", err)
	}
	log.Println("File appended successfully")
}

func AppendFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer sourceFile.Close()

	destFile, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return fmt.Errorf("failed to open destination file: %w", err)
	}
	defer destFile.Close()

	if _, err := io.Copy(destFile, sourceFile); err != nil {
		return fmt.Errorf("failed to copy contents: %w", err)
	}
	return destFile.Sync()
}

func GiveFilenameByRemovingSpaces(filename string) string {
	filename = strings.TrimSpace(filename)
	if filename == "" {
		return ""
	}
	filename = regexp.MustCompile(`\s+`).ReplaceAllString(filename, "_")
	filename = regexp.MustCompile(`[\/\\:\*\?"<>\|]`).ReplaceAllString(filename, "")
	filename = regexp.MustCompile(`_+`).ReplaceAllString(filename, "_")
	return filename
}
func SaveRunDefaultToDB(detail models.RunDetail) {
	var runDetailService = mongoRepo.NewRunDetailRepository()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel() // Important to prevent context leaks
	// 2. Pass ctx and detail as two separate arguments
	err := runDetailService.Create(ctx, detail)
	if err != nil {
		log.Printf("Failed to save to MongoDB: %v", err)
		LogToAll(detail.RunDetailId, constant.LogTypeManual, fmt.Sprintf("Failed to save to MongoDB: %v", err))
	}
	// else {
	// 	log.Println("Successfully saved run detail to MongoDB")
	// }
}
