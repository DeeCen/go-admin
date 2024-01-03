// Package guard 表格信息
package guard

import (
    "github.com/GoAdminGroup/go-admin/context"
    "github.com/GoAdminGroup/go-admin/modules/errors"
)

// MenuDeleteParam 菜单删除参数
type MenuDeleteParam struct {
    ID string
}

// MenuDelete 删除菜单
func (g *Guard) MenuDelete(ctx *context.Context) {

    id := ctx.Query("id")

    if id == "" {
        alertWithTitleAndDesc(ctx, "Menu", "menu", errors.WrongID, g.conn, g.navButtons)
        ctx.Abort()
        return
    }

    ctx.SetUserValue(deleteMenuParamKey, &MenuDeleteParam{
        ID: id,
    })
    ctx.Next()
}

// GetMenuDeleteParam 获取删除参数
func GetMenuDeleteParam(ctx *context.Context) *MenuDeleteParam {
    return ctx.UserValue[deleteMenuParamKey].(*MenuDeleteParam)
}
