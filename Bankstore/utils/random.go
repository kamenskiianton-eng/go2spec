package utils

import (
	"log"
	"math/rand"
	"time"

	"github.com/go-faker/faker/v4"
)

type RandomAccountParams struct {
	Owner    string `faker:"first_name"`
	Balance  int64
	Currency string `faker:"oneof: USD, EUR"`
}

type RandomUserParams struct {
	Username       string `faker:"username"`
	HashedPassword string `faker:"password"`
	FullName       string `faker:"name"`
	Email          string `faker:"email"`
}

func init() {
	rand.New(rand.NewSource(time.Now().UnixNano()))
}

func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

func RandomAccount() RandomAccountParams {
	rap := RandomAccountParams{}
	err := faker.FakeData(&rap)
	if err != nil {
		log.Fatal(err)
	}
	return rap
}

func RandomUser() RandomUserParams {
	rup := RandomUserParams{}
	err := faker.FakeData(&rup)
	if err != nil {
		log.Fatal(err)
	}
	return rup
}
