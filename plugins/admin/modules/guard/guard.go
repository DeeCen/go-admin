// Package guard 表格信息
package guard

import (
    "github.com/GoAdminGroup/go-admin/context"
    "github.com/GoAdminGroup/go-admin/modules/db"
    "github.com/GoAdminGroup/go-admin/modules/errors"
    "github.com/GoAdminGroup/go-admin/modules/service"
    "github.com/GoAdminGroup/go-admin/plugins/admin/modules/constant"
    "github.com/GoAdminGroup/go-admin/plugins/admin/modules/response"
    "github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
    "github.com/GoAdminGroup/go-admin/template"
    "github.com/GoAdminGroup/go-admin/template/types"
)

const (
    editFormParamKey    = "edit_form_param"
    deleteParamKey      = "delete_param"
    exportParamKey      = "export_param"
    serverLoginParamKey = "server_login_param"
    deleteMenuParamKey  = "delete_menu_param"
    editMenuParamKey    = "edit_menu_param"
    newMenuParamKey     = "new_menu_param"
    newFormParamKey     = "new_form_param"
    updateParamKey      = "update_param"
    showFormParamKey    = "show_form_param"
    showNewFormParam    = "show_new_form_param"
)

// Guard 表格
type Guard struct {
    services   service.List
    conn       db.Connection
    tableList  table.GeneratorList
    navButtons *types.Buttons
}

// New 获取Guard实例
func New(s service.List, c db.Connection, t table.GeneratorList, b *types.Buttons) *Guard {
    return &Guard{
        services:   s,
        conn:       c,
        tableList:  t,
        navButtons: b,
    }
}

func (g *Guard) table(ctx *context.Context) (table.Table, string) {
    prefix := ctx.Query(constant.PrefixKey)
    return g.tableList[prefix](ctx), prefix
}

// CheckPrefix 检查前缀是否存在
func (g *Guard) CheckPrefix(ctx *context.Context) {
    prefix := ctx.Query(constant.PrefixKey)

    if _, ok := g.tableList[prefix]; !ok {
        if ctx.Headers(constant.PjaxHeader) == "" && ctx.Method() != "GET" {
            response.BadRequest(ctx, errors.Msg)
        } else {
            response.Alert(ctx, errors.Msg, errors.Msg, "table model not found", g.conn, g.navButtons,
                template.Missing404Page)
        }
        ctx.Abort()
        return
    }

    ctx.Next()
}
