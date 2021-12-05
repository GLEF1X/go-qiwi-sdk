package filters

import (
	"strconv"
	"time"

	"github.com/fatih/structs"
	"github.com/go-playground/validator/v10"
)

const (
	MaxTransactionsAPILimit = 50
)

type HistoryFilter struct {
	Rows      int        `validate:"omitempty,min=0,max=50" json:"rows"`
	Operation string     `validate:"omitempty,oneof=QW_RUB QW_USD QW_EUR CARD MK" json:"operation"`
	StartDate *time.Time `json:"startDate"`
	EndDate   *time.Time `validate:"omitempty,gtfield=StartDate" json:"endDate"`
}

func (filter *HistoryFilter) ConvertToMapWithValidation(validate *validator.Validate) (map[string]string, error) {
	err := validate.Struct(filter)

	if err != nil {
		return nil, err
	}

	formattedStartDate, formattedEndDate, err := filter.formatDates()
	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	if filter.Rows == 0 {
		filter.Rows = MaxTransactionsAPILimit
	}

	generatedMap := structs.Map(struct {
		Rows      string
		Operation string
		StartDate string
		EndDate   string
	}{
		Rows:      strconv.Itoa(filter.Rows),
		Operation: filter.Operation,
		StartDate: formattedStartDate,
		EndDate:   formattedEndDate,
	})

	stringMap := make(map[string]string, len(generatedMap))

	for key, value := range generatedMap {
		value := value.(string)
		if value == "" {
			continue
		}
		stringMap[lowerFirstLetter(key)] = value
	}

	return stringMap, nil
}

func (filter HistoryFilter) formatDates() (string, string, error) {
	moscowLocation, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		return "", "", nil
	}
	var formattedStartDate, formattedEndDate string
	if filter.StartDate != nil {
		formattedStartDate = filter.StartDate.In(moscowLocation).Format(time.RFC3339)
	}
	if filter.EndDate != nil {
		formattedEndDate = filter.EndDate.In(moscowLocation).Format(time.RFC3339)
	}
	return formattedStartDate, formattedEndDate, nil
}
