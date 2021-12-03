package api

import (
	"reflect"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/go-playground/validator/v10"
)

type FilterSetup struct {
	validate *validator.Validate
}

func TestConvertToMapWithValidation(t *testing.T) {
	s := setupFilter(t)

	currentTime := time.Now()
	truncatedFor15Minutes := currentTime.Truncate(15 * time.Minute)
	testCases := map[string]*HistoryFilter{
		"Only rows filter":        {Rows: 50},
		"Rows + operation filter": {Rows: 50, Operation: "QW_RUB"},
		"Full filter":             {Rows: 50, Operation: "QW_RUB", StartDate: &truncatedFor15Minutes, EndDate: &currentTime},
	}

	for testName, test := range testCases {
		t.Logf("Test filters. Current test case: %s", testName)

		filterMap, err := test.ConvertToMapWithValidation(s.validate)
		assert.NoError(t, err)
		assert.Equal(t, len(filterMap), getNotZeroFieldsNum(test))
	}

}

func TestValidationErrorWhenInvalidStructConvertingToMap(t *testing.T) {
	s := setupFilter(t)

	currentTime := time.Now()
	truncatedFor15Minutes := currentTime.Truncate(15 * time.Minute)
	testCases := map[string]*HistoryFilter{
		"Rows are greater than limit(50)":    {Rows: 60},
		"Operation is invalid":               {Operation: "bla-bla"},
		"StartDate is greater than EndDate ": {StartDate: &currentTime, EndDate: &truncatedFor15Minutes},
	}

	for testName, test := range testCases {
		t.Logf("Test wrong filters. Current test case: %s", testName)

		_, err := test.ConvertToMapWithValidation(s.validate)
		if assert.Error(t, err) {
			_, ok := err.(validator.ValidationErrors)
			assert.Equal(t, true, ok)
		}
	}

}

func setupFilter(t *testing.T) *FilterSetup {
	t.Helper()

	s := &FilterSetup{}

	s.validate = validator.New()
	return s
}

func getNotZeroFieldsNum(s interface{}) int {
	v := reflect.ValueOf(s).Elem()
	var counter uint64
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		if field.IsZero() {
			continue
		}
		atomic.AddUint64(&counter, 1)

	}
	return int(counter)
}
