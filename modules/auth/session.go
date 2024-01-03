// Package auth 登录验证
package auth

import (
    "encoding/base64"
    "encoding/json"
    "errors"
    "fmt"
    "net/http"
    "time"

    "github.com/GoAdminGroup/go-admin/context"
    "github.com/GoAdminGroup/go-admin/modules/config"
    "github.com/GoAdminGroup/go-admin/modules/logger"
)

// DefaultCookieKey cookie name
const DefaultCookieKey = "adminToken"

// GetSessionByKey get the session value by key.
func GetSessionByKey(ctx *context.Context, key string) (interface{}, error) {
    ses, err := InitSession(ctx)
    if err != nil {
        return nil, err
    }

    return ses.Get(key), nil
}

// InitSession return the default Session.
func InitSession(ctx *context.Context) (*Session, error) {
    sessions := new(Session)

    sessions.Expires = time.Second * time.Duration(config.GetSessionLifeTime())
    sessions.CookieKey = DefaultCookieKey
    sessions.Values = make(map[string]interface{})

    return sessions.StartCtx(ctx)
}

// Session contains info of session.
type Session struct {
    Context   *context.Context
    Expires   time.Duration
    CookieKey string
    Values    map[string]interface{}
}

// Get  the session value.
func (ses *Session) Get(key string) interface{} {
    return ses.Values[key]
}

// Add the session value of key.
func (ses *Session) Add(key string, value interface{}) error {
    ses.Values[key] = value
    return ses.sendCookie()
}

// sendCookie send cookie header.
func (ses *Session) sendCookie() error {
    maxExp := time.Hour * 24 * 30
    if maxExp > ses.Expires {
        maxExp = ses.Expires
    }

    val := ses.Values
    val[`_ttl`] = uint32(time.Now().Add(maxExp).Unix())
    jsonByte, err := json.Marshal(val)
    if err != nil {
        return err
    }

    cookie := http.Cookie{
        Name:     ses.CookieKey,
        Value:    string(cookieEncode(jsonByte)),
        MaxAge:   config.GetSessionLifeTime(),
        Expires:  time.Now().Add(ses.Expires),
        HttpOnly: true,
        Path:     "/",
    }
    if config.GetDomain() != "" {
        cookie.Domain = config.GetDomain()
    }
    ses.Context.SetCookie(&cookie)
    return nil
}

func (ses *Session) Clear() error {
    ses.Values = map[string]interface{}{}
    return ses.sendCookie()
}

// StartCtx return a Session from the given Context.
func (ses *Session) StartCtx(ctx *context.Context) (ret *Session, err error) {
    ret = ses
    ret.Context = ctx

    cookie, e := ctx.Request.Cookie(ses.CookieKey)
    if e != nil {
        //err = fmt.Errorf(`ctx.Request.Cookie: %w`, e)
        return
    }

    if cookie.Value == "" {
        //err = errors.New(`ctx.Request.Cookie empty`)
        return
    }

    jsonByte, e := CookieDecode([]byte(cookie.Value))
    if e != nil {
        err = fmt.Errorf(`ctx.Request.Cookie decode err: %w`, e)
        return
    }

    var values map[string]interface{}
    err = json.Unmarshal(jsonByte, &values)
    if err != nil {
        err = fmt.Errorf(`ctx.Request.Cookie json err: %w`, err)
        return
    }

    ttl, ok := values[`_ttl`]
    if ok == false {
        err = errors.New(`cookie ttl err`)
        return
    }

    // cookie 有效但过期,ttl类型uint32会被自动转为float64
    if ttl64, ok := ttl.(float64); ok == false || ttl64 < float64(time.Now().Unix()) {
        logger.Warn(`cookie ttl change float64 fail`, values)
        return
    }

    ret.Values = values
    return
}

const cookieKey = `6c6b512c4cea6a7cb54655d75c797608`

// cookieEncode cookie加密
func cookieEncode(s []byte) []byte {
    for k, v := range s {
        s[k] = v ^ cookieKey[k%32]
    }
    return Base64Encode(s)
}

// CookieDecode cookie解密
func CookieDecode(s []byte) ([]byte, error) {
    ret, err := Base64Decode(s)
    if err != nil {
        return nil, err
    }

    for k, v := range ret {
        ret[k] = v ^ cookieKey[k%32]
    }
    return ret, nil
}

// Base64Encode 计算base64,url格式
func Base64Encode(s []byte) []byte {
    var ret = make([]byte, base64.URLEncoding.EncodedLen(len(s)))
    base64.URLEncoding.Encode(ret, s)
    return ret
}

// Base64Decode 解密base64,url格式
func Base64Decode(s []byte) ([]byte, error) {
    ret := make([]byte, base64.URLEncoding.DecodedLen(len(s)))
    l, err := base64.URLEncoding.Decode(ret, s)
    return ret[:l], err
}
