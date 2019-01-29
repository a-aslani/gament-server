package database

import (
	"Server/app/utility"
	"context"
	"github.com/arangodb/go-driver"
)

type aql struct {
	collection string
	target     string
	field      string
	condition  string
	from       driver.DocumentID
	graph      string
	count      int64
	page       int64
	graphType  string
}

type AQLBuilder interface {
	Collection(string) AQLBuilder
	Filter(field, condition, target string) AQLBuilder
	From(from driver.DocumentID) AQLBuilder
	Graph(name string) AQLBuilder
	Type(graphType string) AQLBuilder
	Count(count int64) AQLBuilder
	Page(page int64) AQLBuilder
	Build() Create
}

type aqlBuilder struct {
	collection string
	target     string
	field      string
	condition  string
	from       driver.DocumentID
	graph      string
	count      int64
	page       int64
	graphType  string
}

type Create interface {
	findItem() (map[string]interface{}, bool)
	findItemInGraph() (map[string]interface{}, bool)
	findItemsInGraph() ([]map[string]interface{}, bool)
	removeItemInEdge() error
	totalCount() (int64, error)
}

func Builder() AQLBuilder {
	return &aqlBuilder{}
}

func (a *aqlBuilder) Collection(collection string) AQLBuilder {
	a.collection = collection
	return a
}

func (a *aqlBuilder) Filter(field, condition, target string) AQLBuilder {
	a.field = field
	a.condition = condition
	a.target = target
	return a
}

func (a *aqlBuilder) From(from driver.DocumentID) AQLBuilder {
	a.from = from
	return a
}

func (a *aqlBuilder) Graph(name string) AQLBuilder {
	a.graph = name
	return a
}

func (a *aqlBuilder) Type(graphType string) AQLBuilder {
	a.graphType = graphType
	return a
}

func (a *aqlBuilder) Count(count int64) AQLBuilder {
	a.count = count
	return a
}

func (a *aqlBuilder) Page(page int64) AQLBuilder {
	a.page = page
	return a
}

func (a *aqlBuilder) Build() Create {
	return &aql{
		collection: a.collection,
		field:      a.field,
		condition:  a.condition,
		target:     a.target,
		from:       a.from,
		graph:      a.graph,
		graphType:  a.graphType,
		count:      a.count,
		page:       a.page,
	}
}

func (a *aql) findItem() (map[string]interface{}, bool) {

	query := `LET c = (FOR v IN ` + a.collection + ` FILTER v.` + a.field + ` ` + a.condition + ` @target RETURN v) RETURN first(c)`

	bindVars := map[string]interface{}{
		"target": a.target,
	}

	ctx := context.Background()

	cursor, err := DB().Query(ctx, query, bindVars)
	defer cursor.Close()
	utility.CheckErr(err)

	var doc map[string]interface{}

	_, err = cursor.ReadDocument(ctx, &doc)
	utility.CheckErr(err)

	if doc == nil {
		return doc, false
	}

	return doc, true
}

func (a *aql) findItemInGraph() (map[string]interface{}, bool) {

	query := `LET c = (FOR v, e, p IN OUTBOUND @from GRAPH @graph FILTER e.type == @type RETURN v) RETURN first(c)`

	bindVars := map[string]interface{}{
		"from":  a.from.String(),
		"graph": a.graph,
		"type":  a.graphType,
	}

	ctx := context.Background()
	cursor, err := DB().Query(ctx, query, bindVars)
	defer cursor.Close()
	utility.CheckErr(err)

	var doc map[string]interface{}

	_, err = cursor.ReadDocument(ctx, &doc)
	utility.CheckErr(err)

	if doc == nil {
		return doc, false
	}

	return doc, true
}

func (a *aql) findItemsInGraph() ([]map[string]interface{}, bool) {

	query := `LET c = (FOR v, e, p IN OUTBOUND @from GRAPH @graph FILTER e.type == @type SORT v.created_at DESC LIMIT @offset, @count RETURN v) RETURN c`

	bindVars := map[string]interface{}{
		"from":   a.from.String(),
		"graph":  a.graph,
		"type":   a.graphType,
		"offset": (a.page - 1) * a.count,
		"count":  a.count,
	}

	ctx := context.Background()
	cursor, err := DB().Query(ctx, query, bindVars)

	defer cursor.Close()
	utility.CheckErr(err)

	var docs []map[string]interface{}
	_, err = cursor.ReadDocument(ctx, &docs)
	utility.CheckErr(err)

	if docs == nil {
		return docs, false
	}

	return docs, true
}

func (a *aql) removeItemInEdge() error {
	query := `LET c = (FOR v, e, p IN 1..1 ANY @from GRAPH @graph FILTER e.type == @type RETURN e._key) REMOVE first(c) IN ` + a.collection
	bindVars := map[string]interface{}{
		"from":  a.from.String(),
		"graph": a.graph,
		"type":  a.graphType,
	}
	ctx := context.Background()
	cursor, err := DB().Query(ctx, query, bindVars)
	defer cursor.Close()
	if err != nil {
		return err
	}
	return nil
}

func (a *aql) totalCount() (int64, error) {

	query := `LET c = (FOR v, e, p IN OUTBOUND @from GRAPH @graph FILTER e.type == @type RETURN v) RETURN LENGTH(c)`

	bindVars := map[string]interface{}{
		"from":   a.from.String(),
		"graph":  a.graph,
		"type":   a.graphType,
	}

	ctx := context.Background()
	cursor, err := DB().Query(ctx, query, bindVars)
	defer cursor.Close()
	utility.CheckErr(err)

	if err != nil {
		return 0, err
	}

	var doc int64

	_, err = cursor.ReadDocument(ctx, &doc)
	utility.CheckErr(err)

	return doc, nil
}
