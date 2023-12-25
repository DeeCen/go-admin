// Package controller 控制器
package controller

import (
    "fmt"
    template2 "html/template"
    "net/http"
    "net/url"

    "github.com/GoAdminGroup/go-admin/modules/logger"

    "github.com/GoAdminGroup/go-admin/template"

    "github.com/GoAdminGroup/go-admin/plugins/admin/modules/response"

    "github.com/GoAdminGroup/go-admin/context"
    "github.com/GoAdminGroup/go-admin/modules/auth"
    "github.com/GoAdminGroup/go-admin/modules/file"
    "github.com/GoAdminGroup/go-admin/modules/language"
    "github.com/GoAdminGroup/go-admin/plugins/admin/modules"
    "github.com/GoAdminGroup/go-admin/plugins/admin/modules/constant"
    form2 "github.com/GoAdminGroup/go-admin/plugins/admin/modules/form"
    "github.com/GoAdminGroup/go-admin/plugins/admin/modules/guard"
    "github.com/GoAdminGroup/go-admin/plugins/admin/modules/parameter"
    "github.com/GoAdminGroup/go-admin/template/types"
    "github.com/GoAdminGroup/go-admin/template/types/form"
)

// ShowForm show form page.
func (h *Handler) ShowForm(ctx *context.Context) {
    param := guard.GetShowFormParam(ctx)
    h.showForm(ctx, "", param.Prefix, param.Param, false)
}

func (h *Handler) showForm(ctx *context.Context, alert template2.HTML, prefix string, param parameter.Parameters, isEdit bool, animation ...bool) {

    panel := h.table(prefix, ctx)

    if panel.GetForm().HasError() {
        if panel.GetForm().PageErrorHTML != template2.HTML("") {
            h.HTML(ctx, auth.Auth(ctx),
                types.Panel{Content: panel.GetForm().PageErrorHTML}, template.ExecuteOptions{Animation: param.Animation})
            return
        }
        h.HTML(ctx, auth.Auth(ctx),
            template.WarningPanel(panel.GetForm().PageError.Error(),
                template.GetPageTypeFromPageError(panel.GetForm().PageError)), template.ExecuteOptions{Animation: param.Animation})
        return
    }

    var (
        user       = auth.Auth(ctx)
        paramStr   = param.GetRouteParamStr()
        newURL     = modules.AorEmpty(panel.GetCanAdd(), h.routePathWithPrefix("show_new", prefix)+paramStr)
        footerKind = "edit"
    )

    if newURL == "" || !user.CheckPermissionByUrlMethod(newURL, h.route("show_new").Method(), url.Values{}) {
        footerKind = "edit_only"
    }

    formInfo, err := panel.GetDataWithId(param)

    if err != nil {
        logger.Error("receive data error: ", err)
        h.HTML(ctx, user, template.
            WarningPanelWithDescAndTitle(err.Error(), panel.GetForm().Description, panel.GetForm().Title),
            template.ExecuteOptions{Animation: alert == "" || ((len(animation) > 0) && animation[0])})

        if isEdit {
            ctx.AddHeader(constant.PjaxURLHeader, h.routePathWithPrefix("show_edit", prefix)+
                param.DeletePK().GetRouteParamStr())
        }
        return
    }

    showEditURL := h.routePathWithPrefix("show_edit", prefix) + param.DeletePK().GetRouteParamStr()
    infoURL := h.routePathWithPrefix("info", prefix) + param.DeleteField(constant.EditPKKey).GetRouteParamStr()
    editURL := h.routePathWithPrefix("edit", prefix)
    referer := ctx.Referer()

    if referer != "" && !isInfoUrl(referer) && !isEditUrl(referer, ctx.Query(constant.PrefixKey)) {
        infoURL = referer
    }

    f := panel.GetForm()

    isNotIframe := ctx.Query(constant.IframeKey) != "true"

    hiddenFields := map[string]string{
        form2.TokenKey:    h.authSrv().AddToken(),
        form2.PreviousKey: infoURL,
    }

    if ctx.Query(constant.IframeKey) != "" {
        hiddenFields[constant.IframeKey] = ctx.Query(constant.IframeKey)
    }

    if ctx.Query(constant.IframeIDKey) != "" {
        hiddenFields[constant.IframeIDKey] = ctx.Query(constant.IframeIDKey)
    }

    content := formContent(aForm().
        SetContent(formInfo.FieldList).
        SetFieldsHTML(f.HTMLContent).
        SetTabContents(formInfo.GroupFieldList).
        SetTabHeaders(formInfo.GroupFieldHeaders).
        SetPrefix(h.config.PrefixFixSlash()).
        SetInputWidth(f.InputWidth).
        SetHeadWidth(f.HeadWidth).
        SetPrimaryKey(panel.GetPrimaryKey().Name).
        SetUrl(editURL).
        SetTitle(f.FormEditTitle).
        SetAjax(f.AjaxSuccessJS, f.AjaxErrorJS).
        SetLayout(f.Layout).
        SetHiddenFields(hiddenFields).
        SetOperationFooter(formFooter(footerKind,
            f.IsHideContinueEditCheckBox,
            f.IsHideContinueNewCheckBox,
            f.IsHideResetButton, f.FormEditBtnWord)).
        SetHeader(f.HeaderHtml).
        SetFooter(f.FooterHtml), len(formInfo.GroupFieldHeaders) > 0, !isNotIframe, f.IsHideBackButton, f.Header)

    if f.Wrapper != nil {
        content = f.Wrapper(content)
    }

    h.HTML(ctx, user, types.Panel{
        Content:     alert + content,
        Description: template2.HTML(formInfo.Description),
        Title:       modules.AorBHTML(isNotIframe, template2.HTML(formInfo.Title), ""),
        MiniSidebar: f.HideSideBar,
    }, template.ExecuteOptions{Animation: alert == "" || ((len(animation) > 0) && animation[0]), NoCompress: f.NoCompress})

    if isEdit {
        ctx.AddHeader(constant.PjaxURLHeader, showEditURL)
    }
}

// EditForm 编辑表单
func (h *Handler) EditForm(ctx *context.Context) {
    param := guard.GetEditFormParam(ctx)
    if param==nil{
        return
    }

    if param.MultiForm!=nil && len(param.MultiForm.File) > 0 {
        err := file.GetFileEngine(h.config.FileUploadEngine.Name).Upload(param.MultiForm)
        if err != nil {
            logger.Error("get file engine error: ", err)
            if ctx.WantJSON() {
                response.Error(ctx, err.Error())
            } else {
                h.showForm(ctx, aAlert().Warning(err.Error()), param.Prefix, param.Param, true)
            }
            return
        }
    }

    formPanel := param.Panel.GetForm()

    for i := 0; i < len(formPanel.FieldList); i++ {
        key := formPanel.FieldList[i].Field
        delFlag := key + `__delete_flag`
        changeFlag := key + `__change_flag`
        if formPanel.FieldList[i].FormType == form.File &&
            len(param.MultiForm.File[key]) == 0 &&
            len(param.MultiForm.Value[delFlag]) > 0 &&
            param.MultiForm.Value[delFlag][0] != "1" {
            param.MultiForm.Value[key] = []string{""}
        }
        if formPanel.FieldList[i].FormType == form.File &&
            len(param.MultiForm.Value[changeFlag]) > 0 &&
            param.MultiForm.Value[changeFlag][0] != "1" {
            delete(param.MultiForm.Value, key)
        }
    }

    err := param.Panel.UpdateDate(param.Value())
    if err != nil {
        logger.Error("update data error: ", err)
        if ctx.WantJSON() {
            response.Error(ctx, err.Error(), map[string]interface{}{
                "token": h.authSrv().AddToken(),
            })
        } else {
            /*
             *    前端提交表单后,报错会导致所有表单数据清空,这里返回404,不再刷新清空原表单数据
             */
            ctx.Write(404, nil, `pjaxError:`+err.Error())
            //h.showForm(ctx, aAlert().Warning(err.Error()), param.Prefix, param.Param, true)
        }
        return
    }

    if formPanel.Responder != nil {
        formPanel.Responder(ctx)
        return
    }

    if ctx.WantJSON() && !param.IsIframe {
        response.OkWithData(ctx, map[string]interface{}{
            "url":   param.PreviousPath,
            "token": h.authSrv().AddToken(),
        })
        return
    }

    if !param.FromList {
        if isNewUrl(param.PreviousPath, param.Prefix) {
            h.showNewForm(ctx, param.Alert, param.Prefix, param.Param.DeleteEditPk().GetRouteParamStr(), true)
            return
        }

        if isEditUrl(param.PreviousPath, param.Prefix) {
            h.showForm(ctx, param.Alert, param.Prefix, param.Param, true, false)
            return
        }

        ctx.HTML(http.StatusOK, fmt.Sprintf(`<script>location.href="%s"</script>`, param.PreviousPath))
        ctx.AddHeader(constant.PjaxURLHeader, param.PreviousPath)
        return
    }

    if param.IsIframe {
        ctx.HTML(http.StatusOK, fmt.Sprintf(`<script>
        swal('%s', '', 'success');
        setTimeout(function(){
            $("#%s", window.parent.document).hide();
            $('.modal-backdrop.fade.in', window.parent.document).hide();
        }, 1000)
</script>`, language.Get("success"), param.IframeID))
        return
    }

    buf := h.showTable(ctx, param.Prefix, param.Param.DeletePK().DeleteEditPk(), nil)
    ctx.HTML(http.StatusOK, buf.String())
    ctx.AddHeader(constant.PjaxURLHeader, param.PreviousPath)
}
