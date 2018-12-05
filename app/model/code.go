package model

import (
	"Server/app/constants"
	"Server/app/model/database"
	"context"
	"github.com/arangodb/go-driver"
)

type Code struct {
	Code      int   `json:"code" form:"code" binding:"required"`
	CreatedAt int64 `json:"created_at"`
}

func (code *Code) find(key string) (map[string]interface{}, error) {

	var c map[string]interface{}

	ctx := context.Background()
	_, err := database.Codes().ReadDocument(ctx, key, &c)

	if err != nil {
		return c, err
	}

	return c, nil
}

func (code *Code) findAll(count, page int64) ([]map[string]interface{}, error) {

	var u []map[string]interface{}
	return u, nil
}

func (code *Code) create() (driver.DocumentID, error) {
	ctx := context.Background()
	if meta, err := database.Codes().CreateDocument(ctx, code); err != nil {
		return "", err
	} else {
		return meta.ID, nil
	}
}

func (code *Code) destroy(key string) error {
	ctx := context.Background()
	_, err := database.Codes().RemoveDocument(ctx, key)
	if err != nil {
		return err
	}
	return nil
}

func (code *Code) update(key string) error {

	ctx := context.Background()
	_, err := database.Codes().UpdateDocument(ctx, key, code)

	if err != nil {
		return err
	}
	return nil
}

func (code *Code) count() (int64, error) {
	query := `FOR v IN ` + constants.Codes + ` FILTER RETURN v`
	ctx := driver.WithQueryCount(context.Background())
	cursor, err := database.DB().Query(ctx, query, nil)
	defer cursor.Close()
	if err != nil {
		return 0, err
	}

	return cursor.Count(), nil
}