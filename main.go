package main

import (
	"fmt"
	"os"
	"github-unfollow-go/github"
)

func main() {
	token := os.Args[1]
	whiteList := os.Args[2:]

	gh := github.NewGitHub(token)

	err := gh.GetUserProfile()
	if err != nil {
		fmt.Println("Error getting user profile:", err)
		return
	}

	page := 1
	for {
		users, err := gh.GetFollowingUsers(page)
		if err != nil {
			fmt.Println("Error getting following users:", err)
			return
		}
		if len(users) == 0 {
			break
		}

		for _, user := range users {
			if !isUserInWhitelist(user.Login, whiteList) {
				followsBack, err := gh.CheckFollowBack(user.Login)
				if err != nil {
					fmt.Printf("Error checking if user follows back: %s, error: %s\n", user.Login, err)
					continue
				}
				if !followsBack {
					err := gh.UnfollowUser(user.Login)
					if err != nil {
						fmt.Printf("Error unfollowing user: %s, error: %s\n", user.Login, err)
						continue
					}
				}
			}
		}
		page++
	}
	fmt.Println("Unfollow script completed.")
}

func isUserInWhitelist(username string, whiteList []string) bool {
	for _, whiteUser := range whiteList {
		if username == whiteUser {
			return true
		}
	}
	return false
}
