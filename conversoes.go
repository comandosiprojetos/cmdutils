package cmdutils

import (
	"fmt"
	"math"
	"strconv"
)

// Recebe um inteiro contendo um valor em segundos e converte ele para uma formato legível
// e retorna uma string com essa informação
// exemplo: valor informado -> 1368055 - valor convertido -> 2 weeks 1 day 20 hours 0 minute 55 seconds
func SecondsToHuman(input int) (result string) {
	years := math.Floor(float64(input) / 60 / 60 / 24 / 7 / 30 / 12)
	seconds := input % (60 * 60 * 24 * 7 * 30 * 12)
	months := math.Floor(float64(seconds) / 60 / 60 / 24 / 7 / 30)
	seconds = input % (60 * 60 * 24 * 7 * 30)
	weeks := math.Floor(float64(seconds) / 60 / 60 / 24 / 7)
	seconds = input % (60 * 60 * 24 * 7)
	days := math.Floor(float64(seconds) / 60 / 60 / 24)
	seconds = input % (60 * 60 * 24)
	hours := math.Floor(float64(seconds) / 60 / 60)
	seconds = input % (60 * 60)
	minutes := math.Floor(float64(seconds) / 60)
	seconds = input % 60

	if years > 0 {
		result =
			fmt.Sprint(pluralSecondsToHuman(int(years), "year"),
				pluralSecondsToHuman(int(months), "month"),
				pluralSecondsToHuman(int(weeks), "week"),
				pluralSecondsToHuman(int(days), "day"),
				pluralSecondsToHuman(int(hours), "hour"),
				pluralSecondsToHuman(int(minutes), "minute"),
				pluralSecondsToHuman(int(seconds), "second"))
	} else if months > 0 {
		result =
			fmt.Sprint(pluralSecondsToHuman(int(months), "month"),
				pluralSecondsToHuman(int(weeks), "week"),
				pluralSecondsToHuman(int(days), "day"),
				pluralSecondsToHuman(int(hours), "hour"),
				pluralSecondsToHuman(int(minutes), "minute"),
				pluralSecondsToHuman(int(seconds), "second"))
	} else if weeks > 0 {
		result =
			fmt.Sprint(pluralSecondsToHuman(int(weeks), "week"),
				pluralSecondsToHuman(int(days), "day"),
				pluralSecondsToHuman(int(hours), "hour"),
				pluralSecondsToHuman(int(minutes), "minute"),
				pluralSecondsToHuman(int(seconds), "second"))
	} else if days > 0 {
		result =
			fmt.Sprint(pluralSecondsToHuman(int(days), "day"),
				pluralSecondsToHuman(int(hours), "hour"),
				pluralSecondsToHuman(int(minutes), "minute"),
				pluralSecondsToHuman(int(seconds), "second"))
	} else if hours > 0 {
		result =
			fmt.Sprint(pluralSecondsToHuman(int(hours), "hour"),
				pluralSecondsToHuman(int(minutes), "minute"),
				pluralSecondsToHuman(int(seconds), "second"))
	} else if minutes > 0 {
		result =
			fmt.Sprint(pluralSecondsToHuman(int(minutes), "minute"),
				pluralSecondsToHuman(int(seconds), "second"))
	} else {
		result = pluralSecondsToHuman(int(seconds), "second")
	}

	return
}

func pluralSecondsToHuman(count int, singular string) (result string) {
	if (count == 1) || (count == 0) {
		result = fmt.Sprint(strconv.Itoa(count), " ", singular, " ")
	} else {
		result = fmt.Sprint(strconv.Itoa(count), " ", singular, "s ")
	}

	return
}
