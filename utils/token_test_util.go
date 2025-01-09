package utils

import "hotelReservetion/types"

func CreateTestToken(user *types.User) string {
	token, err := types.GenerateToken(user, types.AccessToken)
	if err != nil {
		panic(err) // for tests, we can panic
	}
	return token
}
