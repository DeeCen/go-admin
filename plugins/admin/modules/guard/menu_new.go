package guard

import (
    "html/template"
    "strconv"

    "github.com/GoAdminGroup/go-admin/context"
    "github.com/GoAdminGroup/go-admin/modules/auth"
    "github.com/GoAdminGroup/go-admin/modules/errors"
    "github.com/GoAdminGroup/go-admin/plugins/admin/modules/form"
)

type MenuNewParam struct {
    Title      string
    Header     string
    ParentId   int64
    Icon       string
    PluginName string
    Uri        string
    Roles      []string
    Alert      template.HTML
}

func (e MenuNewParam) HasAlert() bool {
    return e.Alert != template.HTML("")
}

func (g *Guard) MenuNew(ctx *context.Context) {

    parentId := ctx.FormValue("parentId")
    if parentId == "" {
        parentId = "0"
    }

    var (
        alertHTML template.HTML
        token     = ctx.FormValue(form.TokenKey)
    )

    if !auth.GetTokenService(g.services.Get(auth.TokenServiceKey)).CheckToken(token) {
        alertHTML = getAlert(errors.EditFailWrongToken)
    }

    if alertHTML == "" {
        alertHTML = checkEmpty(ctx, "title", "icon")
    }

    parentIdInt, _ := strconv.Atoi(parentId)

    ctx.SetUserValue(newMenuParamKey, &MenuNewParam{
        Title:      ctx.FormValue("title"),
        Header:     ctx.FormValue("header"),
        PluginName: ctx.FormValue("pluginName"),
        ParentId:   int64(parentIdInt),
        Icon:       ctx.FormValue("icon"),
        Uri:        ctx.FormValue("uri"),
        Roles:      ctx.Request.Form["roles[]"],
        Alert:      alertHTML,
    })
    ctx.Next()
}

func GetMenuNewParam(ctx *context.Context) *MenuNewParam {
    return ctx.UserValue[newMenuParamKey].(*MenuNewParam)
}
