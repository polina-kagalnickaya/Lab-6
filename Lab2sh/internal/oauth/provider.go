package oauth

import (
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "net/url"
    "strings"
)

type Provider interface {
    GetAuthURL(state string) string
    ExchangeCode(code string) (string, error)
    GetUserInfo(accessToken string) (UserInfo, error)
}

type UserInfo struct {
    ID       string
    Email    string
    FullName string
    Provider string
}

type YandexProvider struct {
    clientID     string
    clientSecret string
    redirectURL  string
}

func NewYandexProvider(clientID, clientSecret, redirectURL string) *YandexProvider {
    return &YandexProvider{
        clientID:     clientID,
        clientSecret: clientSecret,
        redirectURL:  redirectURL,
    }
}


func (p *YandexProvider) GetAuthURL(state string) string {
    params := url.Values{}
    params.Add("response_type", "code")
    params.Add("client_id", p.clientID)
    params.Add("redirect_uri", p.redirectURL)
    params.Add("state", state)
    
    return "https://oauth.yandex.ru/authorize?" + params.Encode()
}


func (p *YandexProvider) ExchangeCode(code string) (string, error) {
    // Используем правильный формат данных (application/x-www-form-urlencoded)
    data := url.Values{}
    data.Set("grant_type", "authorization_code")
    data.Set("code", code)
    data.Set("client_id", p.clientID)
    data.Set("client_secret", p.clientSecret)
    data.Set("redirect_uri", p.redirectURL)
    
    req, err := http.NewRequest("POST", "https://oauth.yandex.ru/token", strings.NewReader(data.Encode()))
    if err != nil {
        return "", err
    }
    req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
    
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()
    
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return "", err
    }
    
    // Проверяем статус ответа
    if resp.StatusCode != http.StatusOK {
        return "", fmt.Errorf("yandex token error: %s, body: %s", resp.Status, string(body))
    }
    
    var result struct {
        AccessToken string `json:"access_token"`
        Error       string `json:"error"`
        ErrorDesc   string `json:"error_description"`
    }
    
    if err := json.Unmarshal(body, &result); err != nil {
        return "", fmt.Errorf("failed to parse response: %v, body: %s", err, string(body))
    }
    
    if result.Error != "" {
        return "", fmt.Errorf("yandex error: %s - %s", result.Error, result.ErrorDesc)
    }
    
    return result.AccessToken, nil
}


func (p *YandexProvider) GetUserInfo(accessToken string) (UserInfo, error) {
    req, err := http.NewRequest("GET", "https://login.yandex.ru/info?format=json", nil)
    if err != nil {
        return UserInfo{}, err
    }
    
    req.Header.Set("Authorization", "OAuth "+accessToken)
    
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return UserInfo{}, err
    }
    defer resp.Body.Close()
    
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return UserInfo{}, err
    }
    
    var yandexUser struct {
        ID       string `json:"id"`
        Login    string `json:"login"`
        RealName string `json:"real_name"`
        Email    string `json:"default_email"`
    }
    
    if err := json.Unmarshal(body, &yandexUser); err != nil {
        return UserInfo{}, fmt.Errorf("failed to parse user info: %v, body: %s", err, string(body))
    }
    
    return UserInfo{
        ID:       yandexUser.ID,
        Email:    yandexUser.Email,
        FullName: yandexUser.RealName,
        Provider: "yandex",
    }, nil
}

