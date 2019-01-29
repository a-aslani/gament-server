package model

import (
	"Server/app/constants"
	"Server/app/model/database"
	"Server/app/utility"
	"context"
	"github.com/arangodb/go-driver"
)

type Tournament struct {
	Name       string `json:"name" form:"name" binding:"required"`
	Type       int    `json:"type" form:"type" binding:"required,max=4,min=1"`
	GameKey    string `json:"game_key" form:"game_key" binding:"required"`
	RuleKey    string `json:"rule_key" form:"rule_key" binding:"required"`
	Ticket     int    `json:"ticket" form:"ticket"`
	Quantity   int    `json:"quantity" form:"quantity" binding:"required"`
	Platform   string `json:"platform" form:"platform" binding:"required"`
	Members    int    `json:"members"`
	Sum        int    `json:"sum"`
	Income     int    `json:"income"`
	Percentage string `json:"percentage" form:"percentage" binding:"required,numeric,max=100,min=1"`
	Award      int    `json:"award"`
	State      int    `json:"state"`
	Approved   bool   `json:"approved"`
	CreatedAt  int64  `json:"created_at"`
	UpdatedAt  int64  `json:"updated_at"`
}

func (tournament *Tournament) find(key string) (map[string]interface{}, error) {
	var t map[string]interface{}
	ctx := context.Background()
	_, err := database.Tournaments().ReadDocument(ctx, key, &t)
	if err != nil {
		return t, err
	}
	return t, nil
}

func (tournament *Tournament) findAll(count, page int64) ([]map[string]interface{}, error) {
	query := `FOR v IN ` + constants.Tournaments + ` FILTER v.approved == true LIMIT @offset, @count RETURN v`
	bindVars := map[string]interface{}{
		"offset": (page - 1) * count,
		"count":  count,
	}
	ctx := context.Background()
	cursor, err := database.DB().Query(ctx, query, bindVars)
	defer cursor.Close()

	var tournaments []map[string]interface{}
	if err != nil {
		return tournaments, err
	}
	for {
		var t map[string]interface{}
		_, err := cursor.ReadDocument(ctx, &t)
		if driver.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			utility.CheckErr(err)
		}
		tournaments = append(tournaments, t)
	}
	return tournaments, nil
}

func (tournament *Tournament) create() (driver.DocumentID, error) {
	ctx := context.Background()
	if meta, err := database.Tournaments().CreateDocument(ctx, tournament); err != nil {
		return "", err
	} else {
		return meta.ID, nil
	}
}

func (tournament *Tournament) destroy(key string) error {
	ctx := context.Background()
	_, err := database.Tournaments().RemoveDocument(ctx, key)
	if err != nil {
		return err
	}
	return nil
}

func (tournament *Tournament) update(key string) error {
	ctx := context.Background()
	_, err := database.Tournaments().UpdateDocument(ctx, key, tournament)
	if err != nil {
		return err
	}
	return nil
}

func (tournament *Tournament) count() (int64, error) {
	query := `FOR v IN ` + constants.Tournaments + ` FILTER v.approved == true RETURN v`
	ctx := driver.WithQueryCount(context.Background())
	cursor, err := database.DB().Query(ctx, query, nil)
	defer cursor.Close()
	if err != nil {
		return 0, err
	}
	return cursor.Count(), nil
}
