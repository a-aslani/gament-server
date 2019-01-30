package utility

import (
	"Server/app/utility/jalali"
	"math"
	"time"
)

func Pages(countRecords, count int64) float64 {
	var pages float64
	if countRecords%count == 0 {
		pages = math.Ceil(float64(countRecords / count))
	} else {
		pages = math.Ceil(float64(countRecords/count)) + 1
	}
	return pages
}

func RefactorResponse(docs []map[string]interface{}) []map[string]interface{} {

	var newDocs []map[string]interface{}

	for i, v := range docs {

		delete(docs[i], "_id")
		delete(docs[i], "_rev")
		docs[i]["key"] = v["_key"]
		delete(docs[i], "_key")

		if _, ok := docs[i]["approved"]; ok {
			delete(docs[i], "approved")
		}

		utc, _ := time.LoadLocation("UTC")
		t := time.Unix(int64(v["created_at"].(float64)), 0)
		t = t.In(utc)
		docs[i]["date"] = jalali.Strftime("%A, %e %b %Y %H:%M", t)

		newDocs = append(newDocs, docs[i])
	}

	return newDocs
}
