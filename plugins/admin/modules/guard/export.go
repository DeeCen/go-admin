package guard

import (
    "strings"

    "github.com/GoAdminGroup/go-admin/context"
    "github.com/GoAdminGroup/go-admin/modules/errors"
    "github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
)

type ExportParam struct {
    Panel  table.Table
    ID []string
    Prefix string
    IsAll  bool
}

func (g *Guard) Export(ctx *context.Context) {
    panel, prefix := g.table(ctx)
    if !panel.GetExportable() {
        alert(ctx, panel, errors.OperationNotAllow, g.conn, g.navButtons)
        ctx.Abort()
        return
    }

    idStr := make([]string, 0)
    ids := ctx.FormValue("id")
    if ids != "" {
        idStr = strings.Split(ids, ",")
    }

    ctx.SetUserValue(exportParamKey, &ExportParam{
        Panel:  panel,
        ID:     idStr,
        Prefix: prefix,
        IsAll:  ctx.FormValue("is_all") == "true",
    })
    ctx.Next()
}

func GetExportParam(ctx *context.Context) *ExportParam {
    return ctx.UserValue[exportParamKey].(*ExportParam)
}
