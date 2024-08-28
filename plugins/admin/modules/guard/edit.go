package guard

import (
    tmpl "html/template"
    "mime/multipart"
    "regexp"
    "strings"

    "github.com/GoAdminGroup/go-admin/context"
    "github.com/GoAdminGroup/go-admin/modules/auth"
    "github.com/GoAdminGroup/go-admin/modules/config"
    "github.com/GoAdminGroup/go-admin/modules/db"
    "github.com/GoAdminGroup/go-admin/modules/errors"
    "github.com/GoAdminGroup/go-admin/plugins/admin/modules/constant"
    "github.com/GoAdminGroup/go-admin/plugins/admin/modules/form"
    "github.com/GoAdminGroup/go-admin/plugins/admin/modules/parameter"
    "github.com/GoAdminGroup/go-admin/plugins/admin/modules/response"
    "github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
    "github.com/GoAdminGroup/go-admin/template"
    "github.com/GoAdminGroup/go-admin/template/types"
)

type ShowFormParam struct {
    Panel  table.Table
    ID     string
    Prefix string
    Param  parameter.Parameters
}

func (g *Guard) ShowForm(ctx *context.Context) {

    panel, prefix := g.table(ctx)

    if !panel.GetEditable() {
        alert(ctx, panel, errors.OperationNotAllow, g.conn, g.navButtons)
        ctx.Abort()
        return
    }

    if panel.GetOnlyInfo() {
        ctx.Redirect(config.Url("/info/" + prefix))
        ctx.Abort()
        return
    }

    if panel.GetOnlyDetail() {
        ctx.Redirect(config.Url("/info/" + prefix + "/detail"))
        ctx.Abort()
        return
    }

    if panel.GetOnlyNewForm() {
        ctx.Redirect(config.Url("/info/" + prefix + "/new"))
        ctx.Abort()
        return
    }

    id := ctx.Query(constant.EditPKKey)

    if id == "" {
        id = "1"
    }

    ctx.SetUserValue(showFormParamKey, &ShowFormParam{
        Panel:  panel,
        ID:     id,
        Prefix: prefix,
        Param: parameter.GetParam(ctx.Request.URL, panel.GetInfo().DefaultPageSize, panel.GetInfo().SortField,
            panel.GetInfo().GetSort()).WithPKs(id),
    })
    ctx.Next()
}

func GetShowFormParam(ctx *context.Context) *ShowFormParam {
    //return ctx.UserValue[showFormParamKey].(*ShowFormParam);
    if v,ok := ctx.UserValue[showFormParamKey];ok{
        return v.(*ShowFormParam);
    }

    return nil
}

type EditFormParam struct {
    Panel        table.Table
    ID           string
    Prefix       string
    Param        parameter.Parameters
    Path         string
    MultiForm    *multipart.Form
    PreviousPath string
    Alert        tmpl.HTML
    FromList     bool
    IsIframe     bool
    IframeID     string
}

func (e EditFormParam) Value() form.Values {
    if e.MultiForm != nil {
        return e.MultiForm.Value
    }

    return nil
}

func (g *Guard) EditForm(ctx *context.Context) {

    panel, prefix := g.table(ctx)

    if !panel.GetEditable() {
        alert(ctx, panel, errors.OperationNotAllow, g.conn, g.navButtons)
        ctx.Abort()
        return
    }
    token := ctx.FormValue(form.TokenKey)

    if !auth.GetTokenService(g.services.Get(auth.TokenServiceKey)).CheckToken(token) {
        alert(ctx, panel, errors.EditFailWrongToken, g.conn, g.navButtons)
        ctx.Abort()
        return
    }

    var (
        previous = ctx.FormValue(form.PreviousKey)
        fromList = isListURL(previous)
        param    = parameter.GetParamFromURL(previous, panel.GetInfo().DefaultPageSize,
            panel.GetInfo().GetSort(), panel.GetPrimaryKey().Name)
    )

    if fromList {
        previous = config.Url("/list/" + prefix + param.GetRouteParamStr())
    }

    var (
        multiForm = ctx.Request.MultipartForm
        id        = ``
    )

    // 修复 multiForm 为nil引发的bug
    if multiForm == nil {
        multiForm = new(multipart.Form)
        multiForm.Value = ctx.Request.Form
    }

    var values map[string][]string
    if multiForm != nil {
        values = multiForm.Value
    }

    if values != nil && len(values[panel.GetPrimaryKey().Name]) > 0 {
        id = values[panel.GetPrimaryKey().Name][0]
    }

    ctx.SetUserValue(editFormParamKey, &EditFormParam{
        Panel:        panel,
        ID:           id,
        Prefix:       prefix,
        Param:        param.WithPKs(id),
        Path:         strings.Split(previous, "?")[0],
        MultiForm:    multiForm,
        IsIframe:     form.Values(values).Get(constant.IframeKey) == "true",
        IframeID:     form.Values(values).Get(constant.IframeIDKey),
        PreviousPath: previous,
        FromList:     fromList,
    })
    ctx.Next()
}

func isListURL(s string) bool {
    reg, _ := regexp.Compile("(.*?)/list/(.*?)$")
    sub := reg.FindStringSubmatch(s)
    return len(sub) > 2 && !strings.Contains(sub[2], "/")
}

func GetEditFormParam(ctx *context.Context) *EditFormParam {
    if v, ok := ctx.UserValue[editFormParamKey]; ok {
        if ret, ok := v.(*EditFormParam); ok {
            return ret
        }
    }
    return nil
}

func alert(ctx *context.Context, panel table.Table, msg string, conn db.Connection, btn *types.Buttons) {
    if ctx.WantJSON() {
        response.BadRequest(ctx, msg)
    } else {
        response.Alert(ctx, panel.GetInfo().Description, panel.GetInfo().Title, msg, conn, btn)
    }
}

func alertWithTitleAndDesc(ctx *context.Context, title, desc, msg string, conn db.Connection, btn *types.Buttons) {
    response.Alert(ctx, desc, title, msg, conn, btn)
}

func getAlert(msg string) tmpl.HTML {
    return template.Get(config.GetTheme()).Alert().Warning(msg)
}
