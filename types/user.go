package types

//ID    string `bson:"id"`    // For MongoDB
//Email string `yaml:"email"`  // For YAML
//Name  string `json:"name"`   // For JSON

type User struct {
	ID        string `bson:"_id" json:"id,omitempty"`
	FirstName string `bson:"firstName" json:"firstName"`
	LastName  string `bson:"lastName" json:"lastName"`
}
