// internal/database/models.go
package database

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Grammar struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Name      string            `bson:"name"`
	Username  string			`bson:"name"`
	Content   string            `bson:"content"`
	CreatedAt time.Time         `bson:"created_at"`
	UpdatedAt time.Time         `bson:"updated_at"`
}

type User struct {
	ID          primitive.ObjectID	`bson:"_id,omitempty"`
	Username    string				`bson:"username"`
}

type GrammarInput struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}

func (g *Grammar) FromInput(input *GrammarInput) {
	now := time.Now()
	g.Name = input.Name
	g.Content = input.Content
	
	if g.CreatedAt.IsZero() {
		g.CreatedAt = now
	}
	g.UpdatedAt = now
}
