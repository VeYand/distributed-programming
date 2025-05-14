package pages

import (
	"fmt"
	"lab9/config"
	"time"

	"github.com/tebeka/selenium"
)

type Page struct {
	Driver  selenium.WebDriver
	BaseURL string
}

func (c *Page) Init(driver selenium.WebDriver) {
	c.Driver = driver
	c.BaseURL = config.BaseURL
}

func (c *Page) OpenPage(url string) error {
	return c.Driver.Get(c.BaseURL + url)
}

func (c *Page) WaitForElement(by, value string, timeout time.Duration) (selenium.WebElement, error) {
	var elem selenium.WebElement
	var err error

	startTime := time.Now()
	for time.Since(startTime) < timeout {
		elem, err = c.Driver.FindElement(by, value)
		if err == nil {
			return elem, nil
		}
		time.Sleep(500 * time.Millisecond)
	}

	return nil, fmt.Errorf("element not found after %v: %v", timeout, err)
}

func (c *Page) WaitForElements(by, value string, timeout time.Duration) ([]selenium.WebElement, error) {
	var elems []selenium.WebElement
	var err error

	startTime := time.Now()
	for time.Since(startTime) < timeout {
		elems, err = c.Driver.FindElements(by, value)
		if err == nil && len(elems) > 0 {
			return elems, nil
		}
		time.Sleep(500 * time.Millisecond)
	}
	return nil, fmt.Errorf("elements not found after %v: %v", timeout, err)
}

func (c *Page) FindElement(by, value string) (selenium.WebElement, error) {
	return c.WaitForElement(by, value, 10*time.Second)
}

func (c *Page) FindElements(by, value string) ([]selenium.WebElement, error) {
	return c.WaitForElements(by, value, 10*time.Second)
}
