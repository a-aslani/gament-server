package model

import (
	"Server/app/constants"
	"Server/app/model/database"
	"Server/app/utility"
	"context"
	"github.com/arangodb/go-driver"
)

type Rule struct {
	Description string `json:"description" form:"description" binding:"required"`
	CreatedAt   int64  `json:"created_at"`
	UpdatedAt   int64  `json:"updated_at"`
}

func (rule *Rule) find(key string) (map[string]interface{}, error) {

	var r map[string]interface{}
	ctx := context.Background()
	_, err := database.Rules().ReadDocument(ctx, key, &r)

	if err != nil {
		return r, err
	}
	return r, nil
}

func (rule *Rule) findAll(count, page int64) ([]map[string]interface{}, error) {

	query := `FOR v IN ` + constants.Rules + ` FILTER LIMIT @offset, @count RETURN v`

	bindVars := map[string]interface{}{
		"offset": (page - 1) * count,
		"count":  count,
	}

	ctx := context.Background()

	cursor, err := database.DB().Query(ctx, query, bindVars)
	defer cursor.Close()

	var rules []map[string]interface{}
	if err != nil {
		return rules, err
	}

	for {
		var rule map[string]interface{}
		_, err := cursor.ReadDocument(ctx, &rule)
		if driver.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			utility.CheckErr(err)
		}
		rules = append(rules, rule)
	}

	return rules, nil
}

func (rule *Rule) create() (driver.DocumentID, error) {
	ctx := context.Background()
	if meta, err := database.Rules().CreateDocument(ctx, rule); err != nil {
		return "", err
	} else {
		return meta.ID, nil
	}
}

func (rule *Rule) destroy(key string) error {
	ctx := context.Background()
	_, err := database.Rules().RemoveDocument(ctx, key)
	if err != nil {
		return err
	}
	return nil
}

func (rule *Rule) update(key string) error {

	ctx := context.Background()
	_, err := database.Rules().UpdateDocument(ctx, key, rule)

	if err != nil {
		return err
	}
	return nil
}

func (rule *Rule) count() (int64, error) {
	query := `FOR v IN ` + constants.Rules + ` FILTER RETURN v`
	ctx := driver.WithQueryCount(context.Background())
	cursor, err := database.DB().Query(ctx, query, nil)
	defer cursor.Close()
	if err != nil {
		return 0, err
	}

	return cursor.Count(), nil
}
