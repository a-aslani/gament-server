package model

import (
	"Server/app/constants"
	"Server/app/model/database"
	"Server/app/utility"
	"context"
	"github.com/arangodb/go-driver"
)

type GamesEdge struct {
	From driver.DocumentID `json:"_from"`
	To   driver.DocumentID `json:"_to"`
	Type string            `json:"type"`
}

func (GE *GamesEdge) find(key string) (map[string]interface{}, error) {

	var g map[string]interface{}
	ctx := context.Background()
	_, err := database.GamesEdge().ReadDocument(ctx, key, &g)

	if err != nil {
		return g, err
	}
	return g, nil
}

func (GE *GamesEdge) findAll(count, page int64) ([]map[string]interface{}, error) {

	query := `FOR v IN ` + constants.GamesEdge + ` FILTER LIMIT @offset, @count RETURN v`

	bindVars := map[string]interface{}{
		"offset": (page - 1) * count,
		"count":  count,
	}

	ctx := context.Background()

	cursor, err := database.DB().Query(ctx, query, bindVars)
	defer cursor.Close()

	var gameEdge []map[string]interface{}
	if err != nil {
		return gameEdge, err
	}

	for {
		var g map[string]interface{}
		_, err := cursor.ReadDocument(ctx, &g)
		if driver.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			utility.CheckErr(err)
		}
		gameEdge = append(gameEdge, g)
	}

	return gameEdge, nil
}

func (GE *GamesEdge) create() (driver.DocumentID, error) {
	ctx := context.Background()
	if meta, err := database.GamesEdge().CreateDocument(ctx, GE); err != nil {
		return "", err
	} else {
		return meta.ID, nil
	}
}

func (GE *GamesEdge) destroy(key string) error {
	ctx := context.Background()
	_, err := database.GamesEdge().RemoveDocument(ctx, key)
	if err != nil {
		return err
	}
	return nil
}

func (GE *GamesEdge) update(key string) error {

	ctx := context.Background()
	_, err := database.GamesEdge().UpdateDocument(ctx, key, GE)

	if err != nil {
		return err
	}
	return nil
}

func (GE *GamesEdge) count() (int64, error) {
	query := `FOR v IN ` + constants.GamesEdge + ` FILTER RETURN v`
	ctx := driver.WithQueryCount(context.Background())
	cursor, err := database.DB().Query(ctx, query, nil)
	defer cursor.Close()
	if err != nil {
		return 0, err
	}

	return cursor.Count(), nil
}
