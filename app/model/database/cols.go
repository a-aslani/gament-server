package database

import (
	"Server/app/constants"
	"Server/app/utility"
	"context"
	"github.com/arangodb/go-driver"
	"sync"
)

var (
	onceUsersCollection       sync.Once
	usersCollection           driver.Collection
	oncePhonesCollection      sync.Once
	phonesCollection          driver.Collection
	onceCodesCollection       sync.Once
	codesCollection           driver.Collection
	onceUsersEdgeCollection   sync.Once
	usersEdgeCollection       driver.Collection
	onceGamesCollection       sync.Once
	gamesCollection           driver.Collection
	onceTournamentsCollection sync.Once
	tournamentsCollection     driver.Collection
	onceGamesEdgeCollection   sync.Once
	gamesEdgeCollection       driver.Collection
	onceRulesCollection       sync.Once
	rulesCollection           driver.Collection
	onceTeamsCollection       sync.Once
	teamsCollection           driver.Collection
	onceTeamsEdgeCollection   sync.Once
	teamsEdgeCollection       driver.Collection
)

// Get users collection
func Users() driver.Collection {
	onceUsersCollection.Do(func() {
		ctx := context.Background()
		collection, err := DB().Collection(ctx, constants.Users)
		utility.CheckErr(err)
		usersCollection = collection
	})
	return usersCollection
}

// Get phones collection
func Phones() driver.Collection {
	oncePhonesCollection.Do(func() {
		ctx := context.Background()
		collection, err := DB().Collection(ctx, constants.Phones)
		utility.CheckErr(err)
		phonesCollection = collection
	})
	return phonesCollection
}

// Get codes collection
func Codes() driver.Collection {
	onceCodesCollection.Do(func() {
		ctx := context.Background()
		collection, err := DB().Collection(ctx, constants.Codes)
		utility.CheckErr(err)
		codesCollection = collection
	})
	return codesCollection
}

//Get users_edge edge
func UsersEdge() driver.Collection {
	onceUsersEdgeCollection.Do(func() {
		ctx := context.Background()
		collection, err := DB().Collection(ctx, constants.UsersEdge)
		utility.CheckErr(err)
		usersEdgeCollection = collection
	})
	return usersEdgeCollection
}

//Get games_edge edge
func GamesEdge() driver.Collection {
	onceGamesEdgeCollection.Do(func() {
		ctx := context.Background()
		collection, err := DB().Collection(ctx, constants.GamesEdge)
		utility.CheckErr(err)
		gamesEdgeCollection = collection
	})
	return gamesEdgeCollection
}

//Get games_edge edge
func TeamsEdge() driver.Collection {
	onceGamesEdgeCollection.Do(func() {
		ctx := context.Background()
		collection, err := DB().Collection(ctx, constants.TeamsEdge)
		utility.CheckErr(err)
		gamesEdgeCollection = collection
	})
	return gamesEdgeCollection
}

//Get games collection
func Games() driver.Collection {
	onceGamesCollection.Do(func() {
		ctx := context.Background()
		collection, err := DB().Collection(ctx, constants.Games)
		utility.CheckErr(err)
		gamesCollection = collection
	})
	return gamesCollection
}

//Get tournaments collection
func Tournaments() driver.Collection {
	onceTournamentsCollection.Do(func() {
		ctx := context.Background()
		collection, err := DB().Collection(ctx, constants.Tournaments)
		utility.CheckErr(err)
		tournamentsCollection = collection
	})
	return tournamentsCollection
}

//Get rules collection
func Rules() driver.Collection {
	onceRulesCollection.Do(func() {
		ctx := context.Background()
		collection, err := DB().Collection(ctx, constants.Rules)
		utility.CheckErr(err)
		rulesCollection = collection
	})
	return rulesCollection
}

//Get teams collection
func Teams() driver.Collection {
	onceRulesCollection.Do(func() {
		ctx := context.Background()
		collection, err := DB().Collection(ctx, constants.Teams)
		utility.CheckErr(err)
		rulesCollection = collection
	})
	return rulesCollection
}
