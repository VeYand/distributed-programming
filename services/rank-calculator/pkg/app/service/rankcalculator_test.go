package service_test

import (
	"github.com/mono83/maybe"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"

	"rankcalculator/pkg/app/data"
	"rankcalculator/pkg/app/event"
	"rankcalculator/pkg/app/model"
	"rankcalculator/pkg/app/service"
)

type mockRepo struct {
	mock.Mock
}

func (m *mockRepo) Find(ID model.TextID) (maybe.Maybe[model.Statistics], error) {
	args := m.Called(ID)

	if args.Get(0) == nil {
		return maybe.Nothing[model.Statistics](), args.Error(1)
	}

	stat := args.Get(0).(model.Statistics)
	return maybe.Just(stat), args.Error(1)
}

func (m *mockRepo) Store(stat model.Statistics) error {
	return m.Called(stat).Error(0)
}

type mockDispatcher struct {
	mock.Mock
}

func (m *mockDispatcher) Dispatch(evt event.Event) error {
	return m.Called(evt).Error(0)
}

type mockCentrifugo struct {
	mock.Mock
}

func (m *mockCentrifugo) Publish(channel string, data interface{}) error {
	return m.Called(channel, data).Error(0)
}

func TestRankCalculator_Calculate_TableDriven(t *testing.T) {
	cases := []struct {
		name      string
		value     string
		count     int
		wantAll   int
		wantAlpha int
		wantDup   bool
	}{
		{"Латинские буквы", "abcdefXYZ", 1, 9, 9, false},
		{"Кириллица", "абвгде", 2, 6, 6, true},
		{"Только знаки/символы", "!@#$%", 1, 5, 0, false},
		{"Смешанный текст (буквы/цифры/символы)", "a1!б2@", 3, 6, 2, true},
		{"Пустая строка", "", 1, 0, 0, false},
		{"Только пробелы и управляющие символы", "   \t\n", 1, 5, 0, false},
		{"Emoji и буквы", "a😀b😂", 1, 4, 2, false},
		{"Китайские иероглифы", "漢字テスト", 1, 5, 5, false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			textID := model.TextID("id_" + tc.name)
			input := data.Text{ID: textID, Value: tc.value, Count: tc.count}

			expectedStats := model.Statistics{
				TextID:               textID,
				AllSymbolsCount:      tc.wantAll,
				AlphabetSymbolsCount: tc.wantAlpha,
				IsDuplicate:          tc.wantDup,
			}

			mr := new(mockRepo)
			md := new(mockDispatcher)
			mc := new(mockCentrifugo)

			mr.On("Store", expectedStats).Return(nil).Once()

			channel := "statistics#" + string(textID)
			mc.On("Publish", channel, mock.MatchedBy(func(data interface{}) bool {
				m, ok := data.(map[string]interface{})
				return ok && m["text_id"] == string(textID)
			})).Return(nil).Once()

			wantEvt := event.CreateRankCalculatedEvent(textID, service.CalculateRank(expectedStats))
			md.On("Dispatch", wantEvt).Return(nil).Once()

			svc := service.NewRankCalculator(mr, md, mc)

			err := svc.Calculate(input)

			assert.NoError(t, err)
			mr.AssertExpectations(t)
			mc.AssertExpectations(t)
			md.AssertExpectations(t)
		})
	}
}
