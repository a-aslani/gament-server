package database

import "github.com/arangodb/go-driver"

func FindItemByCondition(collection ,field, condition, target string) (map[string]interface{}, bool) {
	return Builder().Collection(collection).Filter(field, condition, target).Build().findItem()
}

func FindItemInGraph(graphName string, fromId driver.DocumentID, graphType string) (map[string]interface{}, bool) {
	return Builder().Graph(graphName).From(fromId).Type(graphType).Build().findItemInGraph()
}

func FindItemsInGraph(graphName string, fromId driver.DocumentID, graphType string, page, count int64) ([]map[string]interface{}, bool) {
	return Builder().Graph(graphName).From(fromId).Type(graphType).Page(page).Count(count).Build().findItemsInGraph()
}

func RemoveItemInEdge(graphName, collection string, fromId driver.DocumentID, graphType string) error {
	return Builder().Graph(graphName).Collection(collection).From(fromId).Type(graphType).Build().removeItemInEdge()
}

func TotalCount(graphName string, fromId driver.DocumentID, graphType string) (int64, error) {
	return Builder().Graph(graphName).From(fromId).Type(graphType).Build().totalCount()
}