package ffapi

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type API struct {
	root        string
	accessToken string
}

type UserInfo struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Type     string `json:"type"`
}

func New(root, accessToken string) *API { return &API{root: root, accessToken: accessToken} }

type WhoAmIResult struct {
	Users UserInfo `json:"users"`
}

func (a *API) WhoAmI() (*WhoAmIResult, error) {
	result := new(WhoAmIResult)
	if err := a.getRequest(a.root+"/v1/users/whoami", result); err != nil {
		return nil, err
	}
	return result, nil
}

type UserInfoResult struct {
	Admins []UserInfo `json:"admins"`
	Users  UserInfo   `json:"users"`
}

func (a *API) UserInfo(username string) (*UserInfoResult, error) {
	result := new(UserInfoResult)
	if err := a.getRequest(a.root+"/v1/users/"+username, result); err != nil {
		return nil, err
	}
	return result, nil
}

// ===============

func (a *API) getRequest(url string, tgt interface{}) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("X-Authentication-Token", a.accessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP error: %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(tgt); err != nil {
		return err
	}

	return nil
}
