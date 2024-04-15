package database

import (
	"reflect"
	"testing"
)

func setupDB(t *testing.T) *PhoneBook {
	pb, err := CreateTable("file::memory:")
	if err != nil {
		t.Fatalf("Failed to open database:\n%s", err)
	}
	return pb
}

func setupDelete(t *testing.T) (*PhoneBook, []*UserEntry) {
	pb := setupDB(t)
	users := []*UserEntry{
		{
			Name:  "Steven Culwell",
			Phone: "+1(123)123-1234",
		},
		{
			Name:  "John Smith",
			Phone: "54321 12344",
		},
		{
			Name:  "Jane Doe",
			Phone: "011 701 111 1234",
		},
		{
			Name:  "John Doe",
			Phone: "011 701 111 1234",
		},
		{
			Name:  "Steven Culwell",
			Phone: "+1(523)111-7578",
		},
	}
	for _, u := range users {
		if err := pb.Append(u); err != nil {
			t.Fatalf("Failed to append:\n%s", err)
		}
	}

	return pb, users
}

func TestAdd(t *testing.T) {
	pb := setupDB(t)
	defer pb.Close()
	users := []*UserEntry{
		{
			Name:  "Steven Culwell",
			Phone: "+1(123)123-1234",
		},
		{
			Name:  "John Smith",
			Phone: "54321 12344",
		},
		{
			Name:  "Jane Doe",
			Phone: "011 701 111 1234",
		},
	}
	for _, u := range users {
		if err := pb.Append(u); err != nil {
			t.Fatalf("Failed to append:\n%s", err)
		}
	}
	list, err := pb.ListAll()
	if err != nil {
		t.Fatalf("Failed to list db:\n%s", err)
	}
	if len(list) != len(users) {
		t.Fatalf("len(list) != len(users) [%d != %d]", len(list), len(users))
	}
	i := 0
	for i < len(users) {
		if !reflect.DeepEqual(list[i], users[i]) {
			t.Errorf("Idx: %d, Value: %v, Expected: %v", i, list[i], users[i])
		}
		i++
	}
}

func TestDeleteByName(t *testing.T) {
	pb, users := setupDelete(t)
	defer pb.Close()
	// Delete all users with name "Steven Culwell"
	name := "Steven Culwell"
	deleted, err := pb.DeleteByName(name)
	if err != nil {
		t.Fatalf("Failed to delete:\n%s", err)
	}
	if !deleted {
		t.Error("Deletion returned failure (false) when success (true) was expected")
	}
	var expected []*UserEntry
	for _, u := range users {
		if u.Name != name {
			expected = append(expected, u)
		}
	}
	list, err := pb.ListAll()
	if err != nil {
		t.Fatalf("Failed to list db:\n%s", err)
	}
	if len(list) != len(expected) {
		t.Fatalf("len(list) != len(expected) [%d != %d]", len(list), len(expected))
	}
	i := 0
	for i < len(expected) {
		if !reflect.DeepEqual(list[i], expected[i]) {
			t.Errorf("Idx: %d, Value: %v, Expected: %v", i, list[i], expected[i])
		}
		i++
	}
}

func TestDeleteByPhone(t *testing.T) {
	pb, users := setupDelete(t)
	defer pb.Close()
	phone := "011 701 111 1234"
	_, deleted, err := pb.DeleteByPhone(phone)
	if err != nil {
		t.Fatalf("Failed to delete:\n%s", err)
	}
	if !deleted {
		t.Error("Deletion returned failure (false) when success (true) was expected")
	}
	var expected []*UserEntry
	for _, u := range users {
		if u.Phone != phone {
			expected = append(expected, u)
		}
	}
	list, err := pb.ListAll()
	if err != nil {
		t.Fatalf("Failed to list db:\n%s", err)
	}
	if len(list) != len(expected) {
		t.Fatalf("len(list) != len(expected) [%d != %d]", len(list), len(expected))
	}
	i := 0
	for i < len(expected) {
		if !reflect.DeepEqual(list[i], expected[i]) {
			t.Errorf("Idx: %d, Value: %v, Expected: %v", i, list[i], expected[i])
		}
		i++
	}
}

func TestDeleteByNonExistantName(t *testing.T) {
	pb, users := setupDelete(t)
	defer pb.Close()
	deleted, err := pb.DeleteByName("Steven Doe")
	if err != nil {
		t.Fatalf("Failed to delete:\n%s", err)
	}
	if deleted {
		t.Error("Deletion returned true (success) when failure (false) was expected")
	}
	list, err := pb.ListAll()
	if err != nil {
		t.Fatalf("Failed to list db:\n%s", err)
	}
	if len(list) != len(users) {
		t.Fatalf("len(list) != len(users) [%d != %d]", len(list), len(users))
	}
	i := 0
	for i < len(users) {
		if !reflect.DeepEqual(list[i], users[i]) {
			t.Errorf("Idx: %d, Value: %v, Expected: %v", i, list[i], users[i])
		}
		i++
	}
}

func TestDeleteByNonExistantPhone(t *testing.T) {
	pb, users := setupDelete(t)
	defer pb.Close()
	_, deleted, err := pb.DeleteByPhone("(703)111-2121")
	if err != nil {
		t.Fatalf("Failed to delete:\n%s", err)
	}
	if deleted {
		t.Error("Deletion returned success (true) when failure (false) was expected")
	}
	list, err := pb.ListAll()
	if err != nil {
		t.Fatalf("Failed to list db:\n%s", err)
	}
	if len(list) != len(users) {
		t.Fatalf("len(list) != len(users) [%d != %d]", len(list), len(users))
	}
	i := 0
	for i < len(users) {
		if !reflect.DeepEqual(list[i], users[i]) {
			t.Errorf("Idx: %d, Value: %v, Expected: %v", i, list[i], users[i])
		}
		i++
	}
}
