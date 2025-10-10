// 代码生成时间: 2025-10-11 01:44:35
package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "github.com/astaxie/beego"
    "github.com/astaxie/beego/logs"
    "github.com/go-oauth2/oauth2/v4"
    "github.com/go-oauth2/oauth2/v4/manage"
    "github.com/go-oauth2/oauth2/v4/models"
    "net/http"
)

// CustomStorage 自定义存储结构
type CustomStorage struct {
    // 存储结构
    Storage oauth2.Storage
}

// NewCustomStorage 创建自定义存储
func NewCustomStorage() *CustomStorage {
    return &CustomStorage{
        Storage: oauth2.NewMemoryStorage(
            map[string]interface{}{
                "client_id":   "your_client_id",
                "client_secret": "your_client_secret",
            },
        ),
    }
}

// Authorization 授权接口
func Authorization(rw http.ResponseWriter, r *http.Request) {
    conf := &oauth2.Config{
        RedirectURL:  "http://localhost:8080/oauth2/callback",
        ClientID:     "your_client_id",
        ClientSecret: "your_client_secret",
        Scopes:       "all",
        AuthURL:      "http://localhost:8080/oauth2/authorize",
        TokenURL:     "http://localhost:8080/oauth2/token",
    }

    // 创建授权码URL
    url := conf.AuthCodeURL("state", oauth2.AccessTypeOnline)
    http.Redirect(rw, r, url, http.StatusTemporaryRedirect)
}

// Callback 回调接口
func Callback(rw http.ResponseWriter, r *http.Request) {
    // 从请求中提取授权码
    code := r.FormValue("code")
    if code == "" {
        rw.WriteHeader(http.StatusBadRequest)
        fmt.Fprintf(rw, "empty code")
        return
    }

    token, err := conf.Exchange(oauth2.GrantTypeAuthorizationCode, code)
    if err != nil {
        logs.Error("Exchange:", err)
        rw.WriteHeader(http.StatusInternalServerError)
        fmt.Fprintf(rw, "%s", err.Error())
        return
    }

    // 使用token访问受保护的资源
    resp, err := http.Get(conf.ResourceURL + "?access_token=" + token.AccessToken)
    if err != nil {
        logs.Error("Get:", err)
        rw.WriteHeader(http.StatusInternalServerError)
        fmt.Fprintf(rw, "%s", err.Error())
        return
    }
    defer resp.Body.Close()

    var buf bytes.Buffer
    if _, err := buf.ReadFrom(resp.Body); err != nil {
        logs.Error("ReadFrom:", err)
        rw.WriteHeader(http.StatusInternalServerError)
        fmt.Fprintf(rw, "%s", err.Error())
        return
    }
    fmt.Fprintf(rw, "%s", buf.String())
}

// Token 获取token接口
func Token(rw http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        rw.WriteHeader(http.StatusMethodNotAllowed)
        return
    }

    storage := NewCustomStorage()
    mgr := manage.NewDefaultManager(storage)
    defer mgr.Close()

    if !mgr.ValidateClientID(r.FormValue("client_id")) {
        rw.WriteHeader(http.StatusBadRequest)
        fmt.Fprintf(rw, "invalid client id")
        return
    }

    // 处理授权码模式
    if r.FormValue("grant_type") == oauth2.GrantTypeAuthorizationCode {
        code := r.FormValue("code")
        if code == "" {
            rw.WriteHeader(http.StatusBadRequest)
            fmt.Fprintf(rw, "empty code")
            return
        }

        token, err := mgr.GetCode(code)
        if err != nil || token.GetAccessCodeExpireIn() <= 0 {
            rw.WriteHeader(http.StatusBadRequest)
            fmt.Fprintf(rw, "invalid code")
            return
        }

        if r.FormValue("redirect_uri") != token.GetRedirectUri() {
            rw.WriteHeader(http.StatusBadRequest)
            fmt.Fprintf(rw, "invalid redirect uri")
            return
        }

        // 检查客户端ID和客户端密钥
        if r.FormValue("client_id") != token.GetClientID() || r.FormValue("client_secret") != storage.Storage.GetClientSecret(token.GetClientID()) {
            rw.WriteHeader(http.StatusUnauthorized)
            fmt.Fprintf(rw, "invalid client id or secret")
            return
        }

        // 删除旧的授权码
        mgr.RemoveCode(code)

        // 创建新的token
        newToken := mgr.CreateInfoToken(
            token.GetClientID(),
            token.GetUserID(),
        )

        data, err := json.MarshalIndent(newToken, "", "  ")
        if err != nil {
            logs.Error("json:", err)
            rw.WriteHeader(http.StatusInternalServerError)
            fmt.Fprintf(rw, "%s", err.Error())
            return
        }

        rw.Header().Set("Content-Type", "application/json; charset=utf-8")
        rw.WriteHeader(http.StatusOK)
        fmt.Fprintf(rw, string(data))
    }
}

func main() {
    beego.Router("/oauth2/authorize", &Authorization{}, "get:Authorization")
    beego.Router("/oauth2/callback", &Callback{}, "get:Callback")
    beego.Router("/oauth2/token", &Token{}, "post:Token")
    beego.Run()
}