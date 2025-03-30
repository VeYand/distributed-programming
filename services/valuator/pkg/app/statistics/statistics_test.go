package statistics

import (
	"crypto/sha256"
	"encoding/hex"
	"testing"
	"valuator/pkg/app/errors"

	"github.com/mono83/maybe"
	"github.com/stretchr/testify/assert"

	"valuator/pkg/app/model"
	"valuator/pkg/app/query"
)

type mockTextRepository struct {
	Data map[string]model.Text
}

func (m *mockTextRepository) Find(id model.TextID) (maybe.Maybe[model.Text], error) {
	if text, ok := m.Data[string(id)]; ok {
		return maybe.Just(text), nil
	}
	return maybe.Nothing[model.Text](), nil
}

func (m *mockTextRepository) ListAll() ([]model.Text, error) {
	result := make([]model.Text, 0, len(m.Data))
	for _, text := range m.Data {
		result = append(result, text)
	}
	return result, nil
}

func (m *mockTextRepository) Store(text model.Text) error {
	if m.Data == nil {
		m.Data = make(map[string]model.Text)
	}
	m.Data[string(text.ID)] = text
	return nil
}

func (m *mockTextRepository) Remove(text model.Text) error {
	delete(m.Data, string(text.ID))
	return nil
}

func TestGetSummaryFound(t *testing.T) {
	textRepository := &mockTextRepository{}
	statisticsQS := NewStatisticsQueryService(query.NewTextQueryService(textRepository))

	text := createText("Hello world! This is a test text.")
	err := textRepository.Store(text)
	assert.NoError(t, err)

	summary, err := statisticsQS.GetSummary(string(text.ID))
	assert.NoError(t, err)

	assert.Equal(t, 25, summary.SymbolStatistics.AlphabetSymbolsCount)
	assert.Equal(t, 33, summary.SymbolStatistics.AllSymbolsCount)
	assert.Equal(t, false, summary.UniqueStatistics.IsDuplicate)
}

func TestGetSummaryNotFound(t *testing.T) {
	textRepository := &mockTextRepository{}
	statisticsQS := NewStatisticsQueryService(query.NewTextQueryService(textRepository))
	text := createText("Hello world! This is a test text.")

	_, err := statisticsQS.GetSummary(string(text.ID))

	assert.ErrorIs(t, err, errors.ErrTextNotFound)
}

func createText(value string) model.Text {
	newID := hashText(value)
	return model.Text{
		ID:    model.TextID(newID),
		Value: value,
	}
}

func hashText(text string) string {
	hash := sha256.New()
	hash.Write([]byte(text))
	return hex.EncodeToString(hash.Sum(nil))
}
