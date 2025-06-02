package tests

import (
	"github.com/stretchr/testify/assert"
	"github.com/tebeka/selenium"
	"lab9/common"
	"lab9/config"
	"lab9/pages"
	"math/rand"
	"testing"
	"time"
)

func TestCalculateTextStatistics(t *testing.T) {
	testCases := []TestCase{
		{
			name: "Латинские буквы",
			text: config.Text{
				Value:   "abcdefXYZ",
				Country: "Франция",
			},
			expectedStatistics: config.Statistics{
				Rank:        "0",
				IsDuplicate: true,
			},
		},
		{
			name: "Кириллица",
			text: config.Text{
				Value:   "абвгде",
				Country: "Франция",
			},
			expectedStatistics: config.Statistics{
				Rank:        "0",
				IsDuplicate: true,
			},
		},
		{
			name: "Только цифры",
			text: config.Text{
				Value:   "012345",
				Country: "Франция",
			},
			expectedStatistics: config.Statistics{
				Rank:        "1",
				IsDuplicate: true,
			},
		},
		{
			name: "Только знаки/символы",
			text: config.Text{
				Value:   "!@#$%",
				Country: "Франция",
			},
			expectedStatistics: config.Statistics{
				Rank:        "1",
				IsDuplicate: true,
			},
		},
		{
			name: "Смешанный текст (буквы/цифры/символы)",
			text: config.Text{
				Value:   "a1!б2@",
				Country: "Франция",
			},
			expectedStatistics: config.Statistics{
				Rank:        "0.666666",
				IsDuplicate: true,
			},
		},
		{
			name: "Только пробелы и управляющие символы",
			text: config.Text{
				Value:   "   \t\n",
				Country: "Франция",
			},
			expectedStatistics: config.Statistics{
				Rank:        "1",
				IsDuplicate: true,
			},
		},
		{
			name: "Китайские иероглифы",
			text: config.Text{
				Value:   "漢字テスト",
				Country: "Франция",
			},
			expectedStatistics: config.Statistics{
				Rank:        "0",
				IsDuplicate: true,
			},
		},
	}

	uniqueTextCase := TestCase{
		name: "Уникальный текст только из пробелов",
		text: config.Text{
			Value:   randStringRandomLength(100, 200),
			Country: "Россия",
		},
		expectedStatistics: config.Statistics{
			Rank:        "0",
			IsDuplicate: false,
		},
	}

	testFunc := func(t *testing.T, driver selenium.WebDriver, test TestCase) {
		addTextPage := pages.AnalyzePage{}
		addTextPage.Init(driver)
		err := addTextPage.OpenPage("/text/add-form")
		assert.NoError(t, err, "Не удалось открыть страницу")

		err = addTextPage.TypeNewText(test.text.Value)
		assert.NoError(t, err, "Ошибка при вписывании текста")

		err = addTextPage.SelectCountry(test.text.Country)
		assert.NoError(t, err, "Ошибка при выборе страны")

		err = addTextPage.Analyze()
		assert.NoError(t, err, "Ошибка при нажатии на кнопку")

		statistics, err := addTextPage.GetStatistics()
		assert.NoError(t, err, "Ошибка при получении статистики")

		err = common.AssertEqualsStatistics(t, test.expectedStatistics, statistics)
		assert.NoError(t, err, "Ошибка при сравнении статистики")
	}

	for _, testCase := range testCases {
		runTestForBrowser(t, "chrome", testCase, testFunc)
		runTestForBrowser(t, "firefox", testCase, testFunc)
	}
	runTestForBrowser(t, "chrome", uniqueTextCase, testFunc)
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func randStringRandomLength(minLen, maxLen int) string {
	if minLen > maxLen || minLen < 0 {
		return ""
	}

	rand.Seed(time.Now().UnixNano())
	length := rand.Intn(maxLen-minLen+1) + minLen

	b := make([]byte, length)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
