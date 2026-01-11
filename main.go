package main

import (
	"test/common"
	testplayground "test/testPlayground"
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

}
