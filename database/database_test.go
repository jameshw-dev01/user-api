package database

import (
	"slices"
	"testing"

	"github.com/jameshw-dev01/user-api/spec"
	"github.com/stretchr/testify/assert"
)

func testData() []spec.User {
	user1 := spec.User{
		Username: "john_doe",
		Hash:     "e1b0c4429f6c8b07838c9ddae756589e472f001d",
		Email:    "johndoe@example.com",
		Name:     "John Doe",
		Age:      30,
	}

	user2 := spec.User{
		Username: "jane_smith",
		Hash:     "acb4c4428f6b9b07953c8dc9e756588f234fc002",
		Email:    "janesmith@example.com",
		Name:     "Jane Smith",
		Age:      28,
	}

	user3 := spec.User{
		Username: "alice_jones",
		Hash:     "bdc5d4439c7b7b08063bdddbe657590f567f003e",
		Email:    "alicejones@example.com",
		Name:     "Alice Jones",
		Age:      32,
	}

	user4 := spec.User{
		Username: "bob_brown",
		Hash:     "ccf5e4440d6d5c09071bddcdf657192f456b004c",
		Email:    "bobbrown@example.com",
		Name:     "Bob Brown",
		Age:      26,
	}

	user5 := spec.User{
		Username: "charlie_garcia",
		Hash:     "ddb6e4550e7e4d01082cdeedf658893e344c005b",
		Email:    "charliegarcia@example.com",
		Name:     "Charlie Garcia",
		Age:      35,
	}
	return []spec.User{user1, user2, user3, user4, user5}
}

func TestCreate1(t *testing.T) {
	db := GetDBConnection(true, "TEST")
	users := testData()
	user1 := users[0]
	db.Create(user1)
	retrieved, err := db.Read("john_doe")
	assert.Equal(t, nil, err)
	assert.Equal(t, user1, retrieved, "users are not equal")
}

func TestReadAll(t *testing.T) {
	db := GetDBConnection(true, "TEST")
	users := testData()
	user1 := users[0]
	user2 := users[1]
	db.Create(user1)
	db.Create(user2)
	retrieved, err := db.ReadAll()
	assert.Equal(t, nil, err)
	idx1 := slices.IndexFunc(retrieved, func(u spec.User) bool { return u.Username == user1.Username })
	idx2 := slices.IndexFunc(retrieved, func(u spec.User) bool { return u.Username == user2.Username })
	assert.NotEqual(t, -1, idx1, "failed to retrieve user1")
	assert.NotEqual(t, -1, idx2, "failed to retrieve user2")
	assert.Equal(t, user1, retrieved[idx1], "user1 not equal")
	assert.Equal(t, user2, retrieved[idx2], "user2 not equal")
}

func TestUpdate(t *testing.T) {
	db := GetDBConnection(true, "TEST")
	users := testData()
	user1 := users[0]
	db.Create(user1)
	user1_updated := user1
	user1_updated.Age = 21
	user1_updated.Email = "john_doe@test.com"
	db.Update(user1_updated)
	retrieved, err := db.Read("john_doe")
	assert.Equal(t, nil, err)
	assert.Equal(t, user1_updated, retrieved)
}

func TestDelete(t *testing.T) {
	db := GetDBConnection(true, "TEST")
	users := testData()
	n := len(users)
	for _, user := range users {
		db.Create(user)
	}
	db.Delete(users[2])
	db.Delete(users[3])

	retrieved, err := db.ReadAll()
	if err != nil {
		t.Fatal(err)
	}
	if len(retrieved) != n-2 {
		t.Fatal("wrong number of retrieved")
	}
	idx1 := slices.IndexFunc(retrieved, func(u spec.User) bool { return u.Username == users[2].Username })
	idx2 := slices.IndexFunc(retrieved, func(u spec.User) bool { return u.Username == users[3].Username })
	assert.Equal(t, -1, idx1, "user2 not deleted")
	assert.Equal(t, -1, idx2, "user3 not deleted")
}
