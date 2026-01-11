package testplayground

import (
	"log"
	"test/common"
	"test/models"
	"time"

	"github.com/google/uuid"
	//"test/common"
)

// All substrings of a given String
func Run() {
	//common.AddMyCodeIntoFile("brute force  file.go") // use this line when you need to create file in codehistory
	var problem = "All substrings of a given String"
	str := "abc"
	startTime := time.Now()

	// call function here
	result := easyApproch(str)

	var detail = models.RunDetail{
		RunDetailId: uuid.New().String(),
		Problem:     problem,
		Input: map[string]interface{}{
			"str": str,
		},
		ExpectedResult: []string{"a", "ab", "abc", "b", "bc", "c"},
		Result:         result,
		StartAt:        startTime,
		EndAt:          time.Now(),
	}

	log.Printf("result: %+v", detail)
	common.SaveRunDefaultToDB(detail)
}
func easyApproch(str string) []string {
	var result = make([]string, 0)
	var substring string
	for i := 0; i < len(str)-1; i++ {
		for j := i + 1; j < len(str); j++ {
			substring = str[i:j]
			result = append(result, substring)
		}
	}
	return result
}
