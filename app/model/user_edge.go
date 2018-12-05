package model

import (
	"Server/app/constants"
	"Server/app/model/database"
	"context"
	"github.com/arangodb/go-driver"
)

type UsersEdge struct {
	From driver.DocumentID `json:"_from"`
	To   driver.DocumentID `json:"_to"`
	Type string            `json:"type"`
}

func (UE *UsersEdge) find(key string) (map[string]interface{}, error) {

	var u map[string]interface{}
	ctx := context.Background()
	_, err := database.UsersEdge().ReadDocument(ctx, key, &u)

	if err != nil {
		return u, err
	}

	return u, nil
}

func (UE *UsersEdge) findAll(count, page int64) ([]map[string]interface{}, error) {

	var u []map[string]interface{}
	return u, nil
}

func (UE *UsersEdge) create() (driver.DocumentID, error) {
	ctx := context.Background()
	if meta, err := database.UsersEdge().CreateDocument(ctx, &UE); err != nil {
		return "", err
	} else {
		return meta.ID, nil
	}
}

func (UE *UsersEdge) destroy(key string) error {
	return nil
}

func (UE *UsersEdge) update(key string) error {

	ctx := context.Background()

	_, err := database.UsersEdge().UpdateDocument(ctx, key, &UE)

	if err != nil {
		return err
	}

	return nil
}

func (UE *UsersEdge) count() (int64, error) {
	query := `FOR v IN ` + constants.UsersEdge + ` RETURN v`
	ctx := driver.WithQueryCount(context.Background())
	cursor, err := database.DB().Query(ctx, query, nil)
	defer cursor.Close()
	if err != nil {
		return 0, err
	}

	return cursor.Count(), nil
}