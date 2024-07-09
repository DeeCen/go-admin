package controller

import (
    "encoding/json"

    "github.com/GoAdminGroup/go-admin/modules/logger"

    "github.com/GoAdminGroup/go-admin/context"
    "github.com/GoAdminGroup/go-admin/modules/config"
    "github.com/GoAdminGroup/go-admin/plugins/admin/models"
    "github.com/GoAdminGroup/go-admin/plugins/admin/modules/constant"
    "github.com/GoAdminGroup/go-admin/plugins/admin/modules/response"
)

func (h *Handler) Operation(ctx *context.Context) {
    id := ctx.Query("_key_op_id")
    if !h.OperationHandler(config.Url("/ajax/"+id), ctx) {
        errMsg := "not found"
        if ctx.Headers(constant.PjaxHeader) == "" && ctx.Method() != "GET" {
            response.BadRequest(ctx, errMsg)
        } else {
            response.Alert(ctx, errMsg, errMsg, errMsg, h.conn, h.navButtons)
        }
        return
    }
}

// RecordOperationLog record all operation logs, store into database.
func (h *Handler) RecordOperationLog(ctx *context.Context) {
    if user, ok := ctx.UserValue["user"].(models.UserModel); ok {
        var input []byte
        form := ctx.Request.MultipartForm
        if form != nil {
            input, _ = json.Marshal((*form).Value)
        }
        logger.Infof(`RecordOperationLog [%s] userId=%d path=%s method=%s input=%s`, ctx.LocalIP(), user.ID, ctx.Path(), ctx.Method(), string(input))
    }
}
