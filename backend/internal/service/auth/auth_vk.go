package auth

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/vk"
)

const (
	vkGetUsersURL = "https://api.vk.com/method/users.get"
	vkAPIVersion  = "5.199"
)

func (s *service) ConfigVK() *oauth2.Config {
	conf := &oauth2.Config{
		ClientID:     s.authServiceConfig.VKClientID,
		ClientSecret: s.authServiceConfig.VKClientSecret,
		RedirectURL:  s.authServiceConfig.VKCallback,
		Endpoint:     vk.Endpoint,
	}
	return conf
}

func (s *service) GetVKUserID(token string) (int, error) {
	reqURL, err := url.Parse(vkGetUsersURL)
	if err != nil {
		return 0, err
	}
	q := reqURL.Query()
	q.Set("v", vkAPIVersion)
	reqURL.RawQuery = q.Encode()
	tokenHeader := fmt.Sprintf("Bearer %s", token)
	res := &http.Request{
		Method: http.MethodGet,
		URL:    reqURL,
		Header: map[string][]string{
			"Authorization": {tokenHeader},
		},
	}
	req, err := http.DefaultClient.Do(res)
	if err != nil {
		return 0, err
	}
	defer req.Body.Close()
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return 0, err
	}
	var data struct {
		VKResponse []struct {
			VKUserID int `json:"id"`
		} `json:"response"`
	}
	err = json.Unmarshal(body, &data)
	if err != nil || len(data.VKResponse) == 0 {
		return 0, err
	}
	return data.VKResponse[0].VKUserID, nil
}
