package model

import (
	"Server/app/constants"
	"Server/app/model/database"
	"Server/app/utility"
	"context"
	"github.com/arangodb/go-driver"
)

type TeamsEdge struct {
	From driver.DocumentID `json:"_from"`
	To   driver.DocumentID `json:"_to"`
	Type string            `json:"type"`
}

func (TE *TeamsEdge) find(key string) (map[string]interface{}, error) {

	var doc map[string]interface{}
	ctx := context.Background()
	_, err := database.TeamsEdge().ReadDocument(ctx, key, &doc)

	if err != nil {
		return doc, err
	}
	return doc, nil
}

func (TE *TeamsEdge) findAll(count, page int64) ([]map[string]interface{}, error) {

	query := `FOR v IN ` + constants.TeamsEdge + ` FILTER LIMIT @offset, @count RETURN v`

	bindVars := map[string]interface{}{
		"offset": (page - 1) * count,
		"count":  count,
	}

	ctx := context.Background()

	cursor, err := database.DB().Query(ctx, query, bindVars)
	defer cursor.Close()

	var docs []map[string]interface{}
	if err != nil {
		return docs, err
	}

	for {
		var doc map[string]interface{}
		_, err := cursor.ReadDocument(ctx, &doc)
		if driver.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			utility.CheckErr(err)
		}
		docs = append(docs, doc)
	}

	return docs, nil
}

func (TE *TeamsEdge) create() (driver.DocumentID, error) {
	ctx := context.Background()
	if meta, err := database.TeamsEdge().CreateDocument(ctx, TE); err != nil {
		return "", err
	} else {
		return meta.ID, nil
	}
}

func (TE *TeamsEdge) destroy(key string) error {
	ctx := context.Background()
	_, err := database.TeamsEdge().RemoveDocument(ctx, key)
	if err != nil {
		return err
	}
	return nil
}

func (TE *TeamsEdge) update(key string) error {

	ctx := context.Background()
	_, err := database.TeamsEdge().UpdateDocument(ctx, key, TE)

	if err != nil {
		return err
	}
	return nil
}

func (TE *TeamsEdge) count() (int64, error) {
	query := `FOR v IN ` + constants.TeamsEdge + ` FILTER RETURN v`
	ctx := driver.WithQueryCount(context.Background())
	cursor, err := database.DB().Query(ctx, query, nil)
	defer cursor.Close()
	if err != nil {
		return 0, err
	}

	return cursor.Count(), nil
}
