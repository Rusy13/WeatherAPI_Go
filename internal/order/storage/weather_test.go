package storage

import (
	"WbTest/internal/infrastructure/database/postgres/mock_database"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

// MockRow is a mock implementation of pgx.Row
type MockRow struct {
	values []interface{}
	err    error
}

func (r *MockRow) Scan(dest ...interface{}) error {
	if r.err != nil {
		return r.err
	}
	for i := range dest {
		switch dest[i].(type) {
		case *string:
			*dest[i].(*string) = r.values[i].(string)
		case *float64:
			*dest[i].(*float64) = r.values[i].(float64)
		}
	}
	return nil
}

func TestGetWeather(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mock_database.NewMockDatabase(ctrl)

	cityName := "Moscow"
	date := "2021-07-20 10:00:00"
	expectedWeather := Weather{
		CityName:    cityName,
		Temperature: 30.5, // Изменяем ожидаемую температуру, чтобы вызвать ошибку
		DateTime:    date,
		Data:        `{"dt":1626817200,"main":{"temp":25.5}}`,
	}

	mockRow := &MockRow{
		values: []interface{}{
			cityName,
			25.5, // Реальное значение температуры
			date,
			`{"dt":1626817200,"main":{"temp":25.5}}`,
		},
		err: nil,
	}

	mockDB.EXPECT().
		QueryRow(gomock.Any(), gomock.Any(), cityName, date).
		Return(mockRow)

	weather, err := GetWeather(mockDB, cityName, date)
	assert.NoError(t, err)
	assert.Equal(t, expectedWeather, weather)
}
