package cmdutils

import (
	"time"
)

func FormataDataHoraMySQL(DataHora string) (string, error) {
	dataHoraConvertida, err := time.Parse("2006-01-02T15:04:05.000-0700", DataHora)
	if err != nil {
		return "", err
	}

	return dataHoraConvertida.Format("2006-01-02 03:04:05"), nil
}

func ExtrairDataMySQL(DataHora string) (string, error) {
	dataHoraConvertida, err := time.Parse("2006-01-02 15:04:05", DataHora)
	if err != nil {
		return "", err
	}

	return dataHoraConvertida.Format("2006-01-02"), nil
}
