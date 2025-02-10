// core/database/mongodb.go
package database

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDB struct {
	client    *mongo.Client
	db        *mongo.Database
	grammars  *mongo.Collection
}

func NewMongoDB(ctx context.Context, uri string) (*MongoDB, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	opts := options.Client().
		ApplyURI(uri).
		SetMinPoolSize(5).
		SetMaxPoolSize(20).
		SetMaxConnIdleTime(30 * time.Second).
		SetTimeout(5 * time.Second)

	client, connErr := mongo.Connect(ctx, opts)
	if connErr != nil {
		return nil, fmt.Errorf("connection error: %w", connErr)
	}

	db := client.Database("resumes-01")
	grammars := db.Collection("grammars")

	return &MongoDB{
		client:    client,
		db:        db,
		grammars:  grammars,
	}, nil
}

func (m *MongoDB) Close(ctx context.Context) error {
	return m.client.Disconnect(ctx)
}

func (m *MongoDB) StoreGrammar(ctx context.Context, grammarID, name, username, content string) error {

	_, err := m.grammars.UpdateOne(
		ctx,
		bson.M{"grammarID": grammarID},
		bson.M{
			"$set": bson.M{
				"content": content,
				"name": name,
				"username":  username,
				"updatedAt": time.Now(),
			},
			"$setOnInsert": bson.M{
				"createdAt": time.Now(),
			},
		},
		options.Update().SetUpsert(true),
	)
	return err
}

func (m *MongoDB) GetGrammar(ctx context.Context, grammarID string) (string, error) {
	var result struct {
		Content string `bson:"content"`
	}
	err := m.grammars.FindOne(ctx, bson.M{"grammarID": grammarID}).Decode(&result)
	return result.Content, err
}
