package mongoRepo

import (
	"context"
	connection "test/db/mongo"
	"test/models"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type RunDetailRepository struct {
	col *mongo.Collection
}

func NewRunDetailRepository() *RunDetailRepository {
	db := connection.GetDatabase() // Uses the helper above
	return &RunDetailRepository{
		col: db.Collection("run_details"),
	}
}

// Create stores a new run detail
func (r *RunDetailRepository) Create(ctx context.Context, detail models.RunDetail) error {
	if detail.RunDetailId == uuid.Nil.String() {
		detail.RunDetailId = uuid.New().String() // Generate UUID if not provided
	}
	_, err := r.col.InsertOne(ctx, detail)
	return err
}

// GetByID finds a detail by its UUID string or object
func (r *RunDetailRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.RunDetail, error) {
	var detail models.RunDetail
	err := r.col.FindOne(ctx, bson.M{"runDetailId": id}).Decode(&detail)
	if err != nil {
		return nil, err
	}
	return &detail, nil
}

// ListByProblem fetches all runs for a specific problem name
func (r *RunDetailRepository) ListByProblem(ctx context.Context, problemName string) ([]models.RunDetail, error) {
	cursor, err := r.col.Find(ctx, bson.M{"problem": problemName})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []models.RunDetail
	err = cursor.All(ctx, &results)
	return results, err
}
