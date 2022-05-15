package cmdutils

import (
	"time"
)

// Recebe uma string com uma data e hora no formato padrão dd/mm/yyy hh:mm:ss
// e retorna uma string no formato mysql yyyy-mm-dd hh:mm:ss e um tipo error
func FormataDataHoraMySQL(DataHora string) (string, error) {
	formatoDataHoraEntrada := "02/01/2006 15:04:05"
	formatoDataHoraSaida := "2006-01-02 15:04:05"
	dataHoraConvertida, err := time.Parse(formatoDataHoraEntrada, DataHora)
	if err != nil {
		return "", err
	}

	return dataHoraConvertida.Format(formatoDataHoraSaida), nil
}

// Recebe uma string com uma data e hora no padrão mysql yyyy-mm-dd hh:mm:ss
// e retorna uma string apenas com a data no formato yyyy-mm-dd e um tipo error
func ExtrairDataMySQL(DataHora string) (string, error) {
	formatoDataHoraEntrada := "2006-01-02 15:04:05"
	formatoDataHoraSaida := "2006-01-02"
	dataHoraConvertida, err := time.Parse(formatoDataHoraEntrada, DataHora)
	if err != nil {
		return "", err
	}

	return dataHoraConvertida.Format(formatoDataHoraSaida), nil
}

// Recebe uma string com uma data e hora no padrão mysql yyyy-mm-dd hh:mm:ss
// e retorna uma string com a data e hora no formato yyyy-mm-dd hh:mm:ss e um tipo error
func ExtrairDataHoraMySQL(DataHora string) (string, error) {
	dataHoraConvertida, err := time.Parse("2006-01-02 15:04:05", DataHora)
	if err != nil {
		return "", err
	}

	return dataHoraConvertida.Format("2006-01-02 15:04:05"), nil
}
