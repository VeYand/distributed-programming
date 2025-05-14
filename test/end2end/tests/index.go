package tests

import (
	"github.com/stretchr/testify/assert"
	"github.com/tebeka/selenium"
	"lab9/config"
	"testing"
)

type TestCase struct {
	name               string
	text               config.Text
	expectedStatistics config.Statistics
}

func runTestForBrowser(t *testing.T, browserName string, test TestCase, testFunc func(*testing.T, selenium.WebDriver, TestCase)) {
	t.Helper()
	t.Run(test.name+"/"+browserName, func(t *testing.T) {
		caps := selenium.Capabilities{"browserName": browserName}
		driver, err := selenium.NewRemote(caps, "http://localhost:4444/wd/hub")
		if !assert.NoError(t, err, "Failed to start "+browserName+" session") {
			return
		}
		defer driver.Quit()

		testFunc(t, driver, test)
	})
}
