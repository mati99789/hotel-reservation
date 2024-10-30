package types

import "go.mongodb.org/mongo-driver/bson/primitive"

//ID    string `bson:"id"`    // For MongoDB
//Email string `yaml:"email"`  // For YAML
//Name  string `json:"name"`   // For JSON

type User struct {
	ID        primitive.ObjectID `bson:"_id, omitempty" json:"id,omitempty"`
	FirstName string             `bson:"firstName" json:"firstName"`
	LastName  string             `bson:"lastName" json:"lastName"`
}
