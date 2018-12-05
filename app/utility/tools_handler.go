package utility

import "math"

func Pages(countRecords, count int64) float64 {
	var pages float64
	if countRecords%count == 0 {
		pages = math.Ceil(float64(countRecords / count))
	} else {
		pages = math.Ceil(float64(countRecords / count)) + 1
	}
	return pages
}