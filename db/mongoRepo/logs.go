package mongoRepo

import (
	"context"
	connection "test/db/mongo"
	"test/models"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type LogDetailRepository struct {
	col *mongo.Collection
}

func NewLogDetailRepository() *LogDetailRepository {
	// GetDatabase() handles the singleton connection
	db := connection.GetDatabase()
	return &LogDetailRepository{
		col: db.Collection("log_details"),
	}
}

// Create stores a new execution log
func (r *LogDetailRepository) Create(ctx context.Context, logEntry models.ExecutionLog) error {
	// Ensure the LogID (primary ID for this entry) is set

	// Timestamp should be set if missing
	if logEntry.Timestamp.IsZero() {
		logEntry.Timestamp = time.Now()
	}

	_, err := r.col.InsertOne(ctx, logEntry)
	return err
}

// GetByRunID finds all logs associated with a specific problem run
func (r *LogDetailRepository) GetByRunID(ctx context.Context, runID uuid.UUID) ([]models.ExecutionLog, error) {
	cursor, err := r.col.Find(ctx, bson.M{"runDetailId": runID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []models.ExecutionLog
	err = cursor.All(ctx, &results)
	return results, err
}

// GetByID finds a specific single log entry by its own LogID
func (r *LogDetailRepository) GetByID(ctx context.Context, logID uuid.UUID) (*models.ExecutionLog, error) {
	var logEntry models.ExecutionLog
	err := r.col.FindOne(ctx, bson.M{"logId": logID}).Decode(&logEntry)
	if err != nil {
		return nil, err
	}
	return &logEntry, nil
}
