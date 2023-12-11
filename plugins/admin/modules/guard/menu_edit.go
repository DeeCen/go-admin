// Package guard 列表
package guard

import (
    "html/template"
    "strconv"

    "github.com/GoAdminGroup/go-admin/context"
    "github.com/GoAdminGroup/go-admin/modules/auth"
    "github.com/GoAdminGroup/go-admin/modules/errors"
    "github.com/GoAdminGroup/go-admin/plugins/admin/modules/form"
)

// MenuEditParam 菜单参数
type MenuEditParam struct {
    ID         string
    Title      string
    Header     string
    PluginName string
    ParentID   int64
    Icon       string
    URI        string
    Roles      []string
    Alert      template.HTML
}

// HasAlert 是否存在alert
func (e MenuEditParam) HasAlert() bool {
    return e.Alert != ``
}

// MenuEdit 菜单编辑
func (g *Guard) MenuEdit(ctx *context.Context) {

    parentID := ctx.FormValue("parentId")
    if parentID == "" {
        parentID = "0"
    }

    var (
        parentIDInt, _ = strconv.Atoi(parentID)
        token          = ctx.FormValue(form.TokenKey)
        alert          template.HTML
    )

    if !auth.GetTokenService(g.services.Get(auth.TokenServiceKey)).CheckToken(token) {
        alert = getAlert(errors.EditFailWrongToken)
    }

    if alert == "" {
        alert = checkEmpty(ctx, "id", "title", "icon")
    }

    ctx.SetUserValue(editMenuParamKey, &MenuEditParam{
        ID:         ctx.FormValue("id"),
        Title:      ctx.FormValue("title"),
        Header:     ctx.FormValue("header"),
        PluginName: ctx.FormValue("pluginName"),
        ParentID:   int64(parentIDInt),
        Icon:       ctx.FormValue("icon"),
        URI:        ctx.FormValue("uri"),
        Roles:      ctx.Request.Form["roles[]"],
        Alert:      alert,
    })
    ctx.Next()
}

// GetMenuEditParam 从ctx获取编辑参数
func GetMenuEditParam(ctx *context.Context) *MenuEditParam {
    return ctx.UserValue[editMenuParamKey].(*MenuEditParam)
}

func checkEmpty(ctx *context.Context, key ...string) template.HTML {
    for _, k := range key {
        if ctx.FormValue(k) == "" {
            return getAlert("wrong " + k)
        }
    }
    return ``
}
