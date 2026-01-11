package main

import (
	"test/common"
	testplayground "test/testPlayground"
	"time"
)

var config common.Config

func SetUpConfiguration() {
	config.SetUpProgram()

}

// main fun
// code start from where
func main() {
	SetUpConfiguration()
	// all testing and logic are written in testlogics
	testplayground.Run()
	// Give the background goroutine a moment to send the data to MongoDB
	time.Sleep(2 * time.Second)
	defer common.Services.Close()

}
