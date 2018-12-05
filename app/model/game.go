package model

import (
	"Server/app/constants"
	"Server/app/model/database"
	"Server/app/utility"
	"context"
	"github.com/arangodb/go-driver"
)

type Game struct {
	Image       string `json:"image" form:"image"`
	Name        string `json:"name" form:"name" binding:"required"`
	Description string `json:"description"`
	PC          bool   `json:"pc"`
	PS          bool   `json:"ps"`
	Xbox        bool   `json:"xbox"`
	Mobile      bool   `json:"mobile"`
	Approved    bool   `json:"approved"`
	CreatedAt   int64  `json:"created_at"`
	UpdatedAt   int64  `json:"updated_at"`
}

func (game *Game) find(key string) (map[string]interface{}, error) {

	var g map[string]interface{}
	ctx := context.Background()
	_, err := database.Games().ReadDocument(ctx, key, &g)

	if err != nil {
		return g, err
	}
	return g, nil
}

func (game *Game) findAll(count, page int64) ([]map[string]interface{}, error) {

	query := `FOR v IN ` + constants.Games + ` FILTER v.approved == true SORT v.created_at DESC LIMIT @offset, @count RETURN 
		{key: v._key, image: v.image, name: v.name, description: v.description, pc: v.pc, ps: v.ps, xbox: v.xbox, mobile: v.mobile, created_at: v.created_at, updated_at: v.updated_at}`

	bindVars := map[string]interface{}{
		"offset": (page - 1) * count,
		"count":  count,
	}

	ctx := context.Background()

	cursor, err := database.DB().Query(ctx, query, bindVars)
	defer cursor.Close()

	var games []map[string]interface{}
	if err != nil {
		return games, err
	}

	for {
		var game map[string]interface{}
		_, err := cursor.ReadDocument(ctx, &game)
		if driver.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			utility.CheckErr(err)
		}
		games = append(games, game)
	}

	return games, nil
}

func (game *Game) create() (driver.DocumentID, error) {
	ctx := context.Background()
	if meta, err := database.Games().CreateDocument(ctx, game); err != nil {
		return "", err
	} else {
		return meta.ID, nil
	}
}

func (game *Game) destroy(key string) error {
	ctx := context.Background()
	_, err := database.Games().RemoveDocument(ctx, key)
	if err != nil {
		return err
	}
	return nil
}

func (game *Game) update(key string) error {

	ctx := context.Background()
	_, err := database.Games().UpdateDocument(ctx, key, game)

	if err != nil {
		return err
	}
	return nil
}

func (game *Game) count() (int64, error) {
	query := `FOR v IN ` + constants.Games + ` FILTER v.approved == true RETURN v`
	ctx := driver.WithQueryCount(context.Background())
	cursor, err := database.DB().Query(ctx, query, nil)
	defer cursor.Close()
	if err != nil {
		return 0, err
	}

	return cursor.Count(), nil
}
