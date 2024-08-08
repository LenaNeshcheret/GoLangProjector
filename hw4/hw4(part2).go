package main

import (
	"fmt"
	"sort"
)

func main() {
	users := []User{
		{Id: 10},
		{Id: 1},
		{Id: 3},
		{Id: 7},
		{Id: 1},
		{Id: 5},
		{Id: 3},
	}
	uniqueUsers := uniqueUser(users)
	fmt.Println("Unique and Sorted Users:", uniqueUsers)
}

type User struct {
	Id int
}

func uniqueUser(users []User) []User {
	uniqueMap := make(map[int]bool)

	var uniqueUsers []User

	for _, user := range users {
		if _, exists := uniqueMap[user.Id]; !exists {
			uniqueMap[user.Id] = true
			uniqueUsers = append(uniqueUsers, user)
		} else {
			fmt.Printf("User with Id %d already exists\n", user.Id)
		}
	}
	userSort(uniqueUsers)

	return uniqueUsers
}

func userSort(users []User) []User {
	sort.Slice(users, func(i, j int) bool {
		return users[i].Id < users[j].Id
	})
	return users
}
