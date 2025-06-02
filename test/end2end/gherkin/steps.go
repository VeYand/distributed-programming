package gherkin

import (
	"context"
	"github.com/stretchr/testify/assert"
	"lab9/common"
	"lab9/config"
	"testing"

	"github.com/cucumber/godog"
	"github.com/tebeka/selenium"
	"lab9/pages"
)

type testContext struct {
	driver      selenium.WebDriver
	addTextPage pages.AnalyzePage
	t           *testing.T
}

func (ctx *testContext) openPage() error {
	ctx.addTextPage = pages.AnalyzePage{}
	ctx.addTextPage.Init(ctx.driver)
	return ctx.addTextPage.OpenPage("/text/add-form")
}

func (ctx *testContext) enterText(text string) error {
	return ctx.addTextPage.TypeNewText(text)
}

func (ctx *testContext) selectRegion(region string) error {
	return ctx.addTextPage.SelectCountry(region)
}

func (ctx *testContext) submitForm() error {
	return ctx.addTextPage.Analyze()
}

func (ctx *testContext) seeResults(table *godog.Table) error {
	actualStatistics, err := ctx.addTextPage.GetStatistics()
	assert.NoError(ctx.t, err)

	results := make(map[string]string)
	for _, row := range table.Rows[0:] {
		param := row.Cells[0].Value
		value := row.Cells[1].Value
		results[param] = value
	}
	return common.AssertEqualsStatistics(ctx.t, config.Statistics{
		Rank:        results["Rank"],
		IsDuplicate: results["Similarity"] == "1",
	}, actualStatistics)
}

func InitializeScenario(ctx *godog.ScenarioContext, t *testing.T) {
	tc := &testContext{t: t}

	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		caps := selenium.Capabilities{"browserName": "chrome"}
		driver, err := selenium.NewRemote(caps, "http://localhost:4444/wd/hub")
		if err != nil {
			t.Fatal(err)
		}
		tc.driver = driver
		return ctx, nil
	})

	ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		if tc.driver != nil {
			err := tc.driver.Quit()
			if err != nil {
				return nil, err
			}
		}
		return ctx, nil
	})

	ctx.Step(`^открываю главную страницу$`, tc.openPage)
	ctx.Step(`^ввожу текст "([^"]*)"$`, tc.enterText)
	ctx.Step(`^выбираю регион "([^"]*)"$`, tc.selectRegion)
	ctx.Step(`^отправляю текст на анализ$`, tc.submitForm)
	ctx.Step(`^вижу результаты$`, tc.seeResults)
}
