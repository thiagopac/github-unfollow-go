package github

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const baseURL = "https://api.github.com"

type GitHub struct {
	Token    string
	Username string
	Client   *http.Client
}

type User struct {
	Login string `json:"login"`
}

func NewGitHub(token string) *GitHub {
	return &GitHub{
		Token:  token,
		Client: &http.Client{},
	}
}

func (gh *GitHub) GetUserProfile() error {
	req, _ := http.NewRequest("GET", baseURL+"/user", nil)
	req.Header.Add("Authorization", "Bearer "+gh.Token)
	req.Header.Add("Accept", "application/vnd.github.v3+json")

	resp, err := gh.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to fetch user profile, status code: %d", resp.StatusCode)
	}

	body, _ := ioutil.ReadAll(resp.Body)
	var userData User
	json.Unmarshal(body, &userData)
	gh.Username = userData.Login
	return nil
}

func (gh *GitHub) GetFollowingUsers(page int) ([]User, error) {
	req, _ := http.NewRequest("GET", fmt.Sprintf(baseURL+"/user/following?page=%d", page), nil)
	req.Header.Add("Authorization", "Bearer "+gh.Token)
	req.Header.Add("Accept", "application/vnd.github.v3+json")

	resp, err := gh.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch following users, status code: %d", resp.StatusCode)
	}

	body, _ := ioutil.ReadAll(resp.Body)
	var users []User
	json.Unmarshal(body, &users)
	return users, nil
}

func (gh *GitHub) UnfollowUser(target string) error {
	req, _ := http.NewRequest("DELETE", baseURL+"/user/following/"+target, nil)
	req.Header.Add("Authorization", "Bearer "+gh.Token)
	req.Header.Add("Accept", "application/vnd.github.v3+json")

	resp, err := gh.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed to unfollow user %s, status code: %d", target, resp.StatusCode)
	}

	fmt.Printf("Unfollowed %s\n", target)
	return nil
}

func (gh *GitHub) CheckFollowBack(target string) (bool, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf(baseURL+"/users/%s/following/%s", target, gh.Username), nil)
	if err != nil {
		return false, err
	}
	req.Header.Add("Authorization", "Bearer "+gh.Token)
	req.Header.Add("Accept", "application/vnd.github.v3+json")

	resp, err := gh.Client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNoContent {
		return true, nil
	} else if resp.StatusCode == http.StatusNotFound {
		return false, nil
	} else {
		return false, fmt.Errorf("failed to check if user %s is following you, status code: %d", target, resp.StatusCode)
	}
}