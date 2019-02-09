package model

import (
	"Server/app/constants"
	"Server/app/model/database"
	"Server/app/utility"
	"context"
	"github.com/arangodb/go-driver"
)

type User struct {
	Image     string `json:"image"`
	Name      string `json:"name" form:"name" binding:"required"`
	Family    string `json:"family" form:"family" binding:"required"`
	Username  string `json:"username" form:"username" binding:"required,min=4"`
	Password  string `json:"password"`
	Approved  bool   `json:"approved"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}

func (user *User) find(key string) (map[string]interface{}, error) {

	var u map[string]interface{}
	ctx := context.Background()
	_, err := database.Users().ReadDocument(ctx, key, &u)

	if err != nil {
		return u, err
	}

	return u, nil
}

func (user *User) findAll(count, page int64) ([]map[string]interface{}, error) {

	query := `FOR v IN ` + constants.Users + ` FILTER v.approved == true LIMIT @offset, @count RETURN v`

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

func (user *User) create() (driver.DocumentID, error) {
	ctx := context.Background()
	if meta, err := database.Users().CreateDocument(ctx, &user); err != nil {
		return "", err
	} else {
		return meta.ID, nil
	}
}

func (user *User) destroy(key string) error {
	return nil
}

func (user *User) update(key string) error {
	return nil
}

func (user *User) count() (int64, error) {
	query := `FOR v IN ` + constants.Users + ` FILTER v.approved == true RETURN v`
	ctx := driver.WithQueryCount(context.Background())
	cursor, err := database.DB().Query(ctx, query, nil)
	defer cursor.Close()
	if err != nil {
		return 0, err
	}

	return cursor.Count(), nil
}
