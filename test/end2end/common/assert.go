package common

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"lab9/config"
	"math"
	"strconv"
	"testing"
)

const epsilon = 0.001

func AssertEqualsStatistics(t *testing.T, expected, actual config.Statistics) error {
	equals, err := compareFloatsFromStrings(expected.Rank, actual.Rank)
	if err != nil {
		return err
	}

	assert.True(t, equals, fmt.Sprintf("Rank (%s) не соответствует ожидаемому (%s)", actual.Rank, expected.Rank))
	assert.Equal(t, expected.IsDuplicate, actual.IsDuplicate, "IsDuplicate не соответствуют")
	return nil
}

func compareFloatsFromStrings(aStr, bStr string) (bool, error) {
	a, err := strconv.ParseFloat(aStr, 64)
	if err != nil {
		return false, fmt.Errorf("failed to parse first float '%s': %w", aStr, err)
	}

	b, err := strconv.ParseFloat(bStr, 64)
	if err != nil {
		return false, fmt.Errorf("failed to parse second float '%s': %w", bStr, err)
	}

	diff := math.Abs(a - b)
	return diff <= epsilon, nil
}
