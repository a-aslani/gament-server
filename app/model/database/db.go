package database

import (
	"Server/app/utility"
	"context"
	"github.com/arangodb/go-driver"
	"github.com/arangodb/go-driver/http"
	"sync"
)

var (
	onceDb   sync.Once
	database driver.Database
)

// Get instance of game Database
func DB() driver.Database {
	onceDb.Do(func() {
		conn, err := http.NewConnection(http.ConnectionConfig{
			Endpoints: []string{"http://localhost:8529"},
		})
		utility.CheckErr(err)
		client, err := driver.NewClient(driver.ClientConfig{
			Connection:     conn,
			Authentication: driver.BasicAuthentication("root", "42212"),
		})
		utility.CheckErr(err)
		ctx := context.Background()
		db, err := client.Database(ctx, "gament")
		utility.CheckErr(err)
		database = db
	})
	return database
}
