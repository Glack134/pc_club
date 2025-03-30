package pc

import (
	"encoding/json"
	"net/http"
)

type SteamClient struct {
	APIKey string
}

func (s *SteamClient) GetGameList() ([]Game, error) {
	resp, err := http.Get("https://api.steampowered.com/ISteamApps/GetAppList/v2/?key=" + s.APIKey)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		AppList struct {
			Apps []Game `json:"apps"`
		} `json:"applist"`
	}
	err = json.NewDecoder(resp.Body).Decode(&result)
	return result.AppList.Apps, err
}

type Game struct {
	AppID int    `json:"appid"`
	Name  string `json:"name"`
}
