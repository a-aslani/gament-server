package model

import "github.com/arangodb/go-driver"

type Repository interface {
	find(string) (map[string]interface{}, error)
	findAll(count, page int64) ([]map[string]interface{}, error)
	create() (driver.DocumentID, error)
	update(string) error
	destroy(string) error
	count() (int64, error)
}

//Find one document from target collection by _key
func Find(key string, r Repository) (map[string]interface{}, error) {
	return r.find(key)
}

//Find all documents from target collection
func FindAll(count, page int64,r Repository) ([]map[string]interface{}, error) {
	return r.findAll(count, page)
}

//Create new document in target collection
func Create(r Repository) (driver.DocumentID, error) {
	return r.create()
}

//Update document from target collection by _key
func Update(key string, r Repository) error {
	return r.update(key)
}

//Remove document from target collection by _key
func Destroy(key string, r Repository) error {
	return r.destroy(key)
}

//Get count of documents in target collection
func Count(r Repository) (int64, error) {
	return r.count()
}
