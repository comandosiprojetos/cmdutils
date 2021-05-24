package cmdutils

import (
	"time"
)

// Recebe uma string com uma data e hora no formato padrão dd/mm/yyy hh:mm:ss 
// e retorna uma string no formato mysql yyyy-mm-dd hh:mm:ss e um tipo error
func FormataDataHoraMySQL(DataHora string) (string, error) {
	dataHoraConvertida, err := time.Parse("2006-01-02T15:04:05.000-0700", DataHora)
	if err != nil {
		return "", err
	}

	return dataHoraConvertida.Format("2006-01-02 03:04:05"), nil
}

// Recebe uma string com uma data e hora no padrão mysql yyyy-mm-dd hh:mm:ss 
// e retorna uma string apenas com a data no formato yyyy-mm-dd e um tipo error
func ExtrairDataMySQL(DataHora string) (string, error) {
	dataHoraConvertida, err := time.Parse("2006-01-02 15:04:05", DataHora)
	if err != nil {
		return "", err
	}

	return dataHoraConvertida.Format("2006-01-02"), nil
}
