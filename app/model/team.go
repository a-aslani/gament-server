package model

import (
	"Server/app/constants"
	"Server/app/model/database"
	"Server/app/utility"
	"context"
	"github.com/arangodb/go-driver"
)

type Team struct {
	Image       string            `json:"image" form:"image"`
	Name        string            `json:"name" form:"name" binding:"required"`
	Description string            `json:"description" form:"description" binding:"required"`
	Admin       driver.DocumentID `json:"admin"`
	GameKey     string            `json:"game_key" form:"game_key" binding:"required"`
	CreatedAt   int64             `json:"created_at"`
	UpdatedAt   int64             `json:"updated_at"`
}

func (team *Team) find(key string) (map[string]interface{}, error) {

	var doc map[string]interface{}
	ctx := context.Background()
	_, err := database.Teams().ReadDocument(ctx, key, &doc)

	if err != nil {
		return doc, err
	}
	return doc, nil
}

func (team *Team) findAll(count, page int64) ([]map[string]interface{}, error) {

	query := `FOR v IN ` + constants.Teams + ` FILTER LIMIT @offset, @count RETURN v`

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

func (team *Team) create() (driver.DocumentID, error) {
	ctx := context.Background()
	if meta, err := database.Teams().CreateDocument(ctx, team); err != nil {
		return "", err
	} else {
		return meta.ID, nil
	}
}

func (team *Team) destroy(key string) error {
	ctx := context.Background()
	_, err := database.Teams().RemoveDocument(ctx, key)
	if err != nil {
		return err
	}
	return nil
}

func (team *Team) update(key string) error {

	ctx := context.Background()
	_, err := database.Teams().UpdateDocument(ctx, key, team)

	if err != nil {
		return err
	}
	return nil
}

func (team *Team) count() (int64, error) {
	query := `FOR v IN ` + constants.Teams + ` FILTER RETURN v`
	ctx := driver.WithQueryCount(context.Background())
	cursor, err := database.DB().Query(ctx, query, nil)
	defer cursor.Close()
	if err != nil {
		return 0, err
	}

	return cursor.Count(), nil
}
