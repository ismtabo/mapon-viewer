package main

import (
	"context"
	"testing"

	"github.com/Telefonica/golium"
	"github.com/Telefonica/golium/steps/common"
	"github.com/Telefonica/golium/steps/http"
	"github.com/cucumber/godog"
	"github.com/ismtabo/mapon-viewer/test/acceptance/steps"
)

func TestMain(m *testing.M) {
	launcher := golium.NewLauncher()
	launcher.Launch(InitializeTestSuite, InitializeScenario)
}

func InitializeTestSuite(ctx context.Context, suiteCtx *godog.TestSuiteContext) {
}

func InitializeScenario(ctx context.Context, scenarioCtx *godog.ScenarioContext) {
	stepsInitializers := []golium.StepsInitializer{
		common.Steps{},
		http.Steps{},
		steps.ServerSteps{},
		steps.MongoSteps{},
	}
	for _, stepsInitializer := range stepsInitializers {
		ctx = stepsInitializer.InitializeSteps(ctx, scenarioCtx)
	}
}
