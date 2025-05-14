package pages

import (
	"fmt"
	"github.com/tebeka/selenium"
	"lab9/config"
	"strings"
)

const (
	Textarea            = "//*[@id=\"addForm\"]/label[1]/textarea"
	CountrySelect       = "//*[@id=\"country\"]"
	CountryOptionByText = ".//option[normalize-space(text())='%s']"
	AnalyzeButton       = "//*[@id=\"addForm\"]/input"
	StatisticsElem      = "//*[@id=\"results\"]"
	StatIdentifier      = "//*[@id=\"results\"]/p[1]"
	StatRank            = "//*[@id=\"results\"]/p[2]"
	StatDuplicate       = "//*[@id=\"results\"]/p[3]"
)

type AnalyzePage struct {
	Page
}

func (a *AnalyzePage) TypeNewText(text string) error {
	input, err := a.FindElement(selenium.ByXPATH, Textarea)
	if err != nil {
		return fmt.Errorf("failed to find input: %v", err)
	}

	if err := input.Clear(); err != nil {
		return fmt.Errorf("failed to clear input: %v", err)
	}

	return input.SendKeys(text)
}

func (a *AnalyzePage) SelectCountry(countryValue string) error {
	selectElem, err := a.FindElement(selenium.ByXPATH, CountrySelect)
	if err != nil {
		return fmt.Errorf("failed to find country select: %v", err)
	}

	optionXPath := fmt.Sprintf(CountryOptionByText, countryValue)
	optionElem, err := selectElem.FindElement(selenium.ByXPATH, optionXPath)
	if err != nil {
		return fmt.Errorf("failed to find option with value '%s': %v", countryValue, err)
	}

	if err := optionElem.Click(); err != nil {
		return fmt.Errorf("failed to click option '%s': %v", countryValue, err)
	}

	return nil
}

func (a *AnalyzePage) Analyze() error {
	butt, err := a.FindElement(selenium.ByXPATH, AnalyzeButton)
	if err != nil {
		return fmt.Errorf("failed to find button: %v", err)
	}

	return butt.Click()
}

func (a *AnalyzePage) GetStatistics() (config.Statistics, error) {
	_, err := a.FindElement(selenium.ByXPATH, StatisticsElem)
	if err != nil {
		return config.Statistics{}, fmt.Errorf("failed to find statistics: %v", err)
	}

	identifierElem, err := a.FindElement(selenium.ByXPATH, StatIdentifier)
	if err != nil {
		return config.Statistics{}, fmt.Errorf("failed to find identifier element: %v", err)
	}
	identifierText, err := identifierElem.Text()
	if err != nil {
		return config.Statistics{}, fmt.Errorf("failed to get identifier text: %v", err)
	}
	identifier := extractValueAfterColon(identifierText)

	rankElem, err := a.FindElement(selenium.ByXPATH, StatRank)
	if err != nil {
		return config.Statistics{}, fmt.Errorf("failed to find rank element: %v", err)
	}
	rankText, err := rankElem.Text()
	if err != nil {
		return config.Statistics{}, fmt.Errorf("failed to get rank text: %v", err)
	}
	rank := extractValueAfterColon(rankText)

	duplicateElem, err := a.FindElement(selenium.ByXPATH, StatDuplicate)
	if err != nil {
		return config.Statistics{}, fmt.Errorf("failed to find duplicate element: %v", err)
	}
	duplicateText, err := duplicateElem.Text()
	if err != nil {
		return config.Statistics{}, fmt.Errorf("failed to get duplicate text: %v", err)
	}
	duplicate := extractValueAfterColon(duplicateText)

	if duplicate != "Найден дубликат" && duplicate != "Дубликатов не найдено" {
		return config.Statistics{}, fmt.Errorf("invalid duplicate value: %s", duplicate)
	}

	return config.Statistics{
		ID:          identifier,
		Rank:        rank,
		IsDuplicate: duplicate == "Найден дубликат",
	}, nil
}

func extractValueAfterColon(s string) string {
	for i, ch := range s {
		if ch == ':' {
			return strings.TrimSpace(s[i+1:])
		}
	}
	return s
}
