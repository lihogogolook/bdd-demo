package main

import (
	"testing"

	"github.com/cucumber/godog"
)

func TestBDDFeatures(t *testing.T) {
	suite := godog.TestSuite{
		ScenarioInitializer: InitializeScenario,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"features"},
			TestingT: t,
		},
	}

	if suite.Run() != 0 {
		t.Fatal("非零狀態返回，BDD 測試失敗")
	}
}
