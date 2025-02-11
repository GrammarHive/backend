// core/database/models.go
package database

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Grammar struct {
	ID        primitive.ObjectID  `bson:"_id,omitempty"`
	Version   int                 `bson:"version"`
	GrammarID string              `bson:"grammarID"`
	Name      string              `bson:"name"`
	Username  string              `bson:"name"`
	Content   string              `bson:"content"`
	CreatedAt time.Time           `bson:"created_at"`
	UpdatedAt time.Time           `bson:"updated_at"`
}

type User struct {
	ID          primitive.ObjectID  `bson:"_id,omitempty"`
	Username    string              `bson:"username"`
}

