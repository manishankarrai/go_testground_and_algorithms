package models

import (
	"time"
)

// for saving logs
type RunDetail struct {
	RunDetailId    string                 `json:"runDetailId" bson:"runDetailId"`
	Problem        string                 `json:"problem" bson:"problem"`
	Input          map[string]interface{} `json:"input" bson:"input"`
	ExpectedResult interface{}            `json:"expectedResult" bson:"expectedResult"`
	Result         interface{}            `json:"result" bson:"result"`
	Error          string                 `json:"error" bson:"error"`
	StartAt        time.Time              `json:"startAt" bson:"startAt"`
	EndAt          time.Time              `json:"endAt" bson:"endAt"`
}

// ExecutionLog represents a single log entry in MongoDB
type ExecutionLog struct {
	LogID       string    `json:"logId" bson:"logId"`
	RunDetailId string    `json:"runDetailId,omitempty" bson:"runDetailId,omitempty"`
	Level       string    `json:"level" bson:"level"`
	Message     string    `json:"message" bson:"message"`
	Timestamp   time.Time `json:"timestamp" bson:"timestamp"`
}
