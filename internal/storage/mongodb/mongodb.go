package mongodb

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type MongoDB struct {
}

func NewMongoDB(host, port, username, password, dbname string) *mongo.Database {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	opt := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s", host, port))
	if username != "" && password != "" {
		opt.SetAuth(options.Credential{
			Username: username,
			Password: password,
		})
	}

	client, err := mongo.Connect(ctx, opt)
	if err != nil {
		panic(err)
	}
	if err := client.Ping(ctx, &readpref.ReadPref{}); err != nil {
		panic(err)
	}

	return client.Database(dbname)
}
