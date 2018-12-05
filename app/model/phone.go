package model

import (
	"Server/app/constants"
	"Server/app/model/database"
	"context"
	"github.com/arangodb/go-driver"
)

type Phone struct {
	Phone     string `json:"phone" form:"phone" binding:"required,numeric,max=10,min=10"`
	CreatedAt int64  `json:"created_at"`
}

func (phone *Phone) find(key string) (map[string]interface{}, error) {

	var p map[string]interface{}
	ctx := context.Background()
	_, err := database.Phones().ReadDocument(ctx, key, &p)
	if err != nil {
		return p, err
	}
	return p, nil
}

func (phone *Phone) findAll(count, page int64) ([]map[string]interface{}, error) {

	var u []map[string]interface{}
	return u, nil
}

func (phone *Phone) create() (driver.DocumentID, error) {
	ctx := context.Background()
	if meta, err := database.Phones().CreateDocument(ctx, &phone); err != nil {
		return "", err
	} else {
		return meta.ID, nil
	}
}

func (phone *Phone) destroy(key string) error {
	ctx := context.Background()
	_, err := database.Phones().RemoveDocument(ctx, key)
	if err != nil {
		return err
	}
	return nil
}

func (phone *Phone) update(key string) error {
	return nil
}

func (phone *Phone) count() (int64, error) {
	query := `FOR v IN ` + constants.Phones + ` RETURN v`
	ctx := driver.WithQueryCount(context.Background())
	cursor, err := database.DB().Query(ctx, query, nil)
	defer cursor.Close()
	if err != nil {
		return 0, err
	}

	return cursor.Count(), nil
}