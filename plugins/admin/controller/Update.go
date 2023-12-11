package controller

import (
    "github.com/GoAdminGroup/go-admin/context"
    "github.com/GoAdminGroup/go-admin/plugins/admin/modules/guard"
    "github.com/GoAdminGroup/go-admin/plugins/admin/modules/response"
)

// Update the table row of given id.
func (h *Handler) Update(ctx *context.Context) {

    param := guard.GetUpdateParam(ctx)

    err := param.Panel.UpdateDate(param.Value)

    if err != nil {
        response.Error(ctx, err.Error())
        return
    }

    response.Ok(ctx)
}
