package types

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"hotelReservetion/utils"
	"regexp"
)

type UserRole string

const (
	minFirstNameLen          = 2
	minLastNameLen           = 2
	minPasswordLen           = 7
	AdminRole       UserRole = "admin"
	GuestRole       UserRole = "guest"
	StaffRole       UserRole = "staff"
)

type CreateUserParams struct {
	FirstName string
	LastName  string
	Email     string
	Password  string
	Role      UserRole
}

type User struct {
	ID                primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	FirstName         string             `bson:"firstName" json:"firstName"`
	LastName          string             `bson:"lastName" json:"lastName"`
	Email             string             `bson:"email" json:"email"`
	EncryptedPassword string             `bson:"encryptedPassword" json:"-"`
	Role              UserRole           `bson:"role" json:"role"`
}

type UpdateUserParams struct {
	FirstName string
	LastName  string
}

func (p UpdateUserParams) ToBSON() bson.M {
	m := bson.M{}

	if len(p.FirstName) > 0 {
		m["firstName"] = p.FirstName
	}

	if len(p.LastName) > 0 {
		m["lastName"] = p.LastName
	}

	return m
}

func (params CreateUserParams) Validate() []utils.APIError {
	var errors []utils.APIError

	if len(params.FirstName) < minFirstNameLen {
		errors = append(errors, utils.NewErrorResponse("firstName", "FirstName must be at least 2 characters long"))
	}

	if len(params.LastName) < minLastNameLen {
		errors = append(errors, utils.NewErrorResponse("lastName", "LastName must be at least 2 characters long"))
	}

	if len(params.Password) < minPasswordLen {
		errors = append(errors, utils.NewErrorResponse("password", "Password must be at least 8 characters long"))
	}

	if !isEmailValid(params.Email) {
		errors = append(errors, utils.NewErrorResponse("email", "Invalid email"))
	}
	return errors
}

func isEmailValid(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

func IsValidPassword(encpw, pw string) bool {
	return bcrypt.CompareHashAndPassword([]byte([]byte(encpw)), []byte(pw)) == nil

}

func NewUserFromParams(params CreateUserParams) (*User, error) {
	encpw, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)

	if err != nil {
		return nil, err
	}

	return &User{
		FirstName:         params.FirstName,
		LastName:          params.LastName,
		Email:             params.Email,
		EncryptedPassword: string(encpw),
		Role:              params.Role,
	}, nil
}
