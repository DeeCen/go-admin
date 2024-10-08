package action

import (
    "encoding/json"
    "html/template"
    "strings"

    "github.com/GoAdminGroup/go-admin/context"
    "github.com/GoAdminGroup/go-admin/modules/config"
    "github.com/GoAdminGroup/go-admin/modules/utils"
    "github.com/GoAdminGroup/go-admin/template/types"
)

type AjaxData map[string]interface{}

func NewAjaxData() AjaxData {
    return AjaxData{"ids": "{{.IDs}}"}
}

func (a AjaxData) Add(m map[string]interface{}) AjaxData {
    for k, v := range m {
        a[k] = v
    }
    return a
}

func (a AjaxData) JSON() string {
    b, _ := json.Marshal(a)
    //println(`--------------debug----------,`,string(b))
    /*return utils.ReplaceAll(string(b), `"{%id}"`, "{{.ID}}",
        `"{%ids}"`, "{{.IDs}}",
        `"{{.IDs}}"`, "{{.IDs}}",
        `"{{.ID}}"`, "{{.ID}}")*/
    bStr := strings.ReplaceAll(string(b),`"{%id}"`, `{{.ID}}`)
    bStr = strings.ReplaceAll(bStr,`"{%ids}"`, `{{.IDs}}`)
    bStr = strings.ReplaceAll(bStr,`"{{.IDs}}"`, `{{.IDs}}`)
    bStr = strings.ReplaceAll(bStr,`"{{.ID}}"`, `{{.ID}}`)
    return bStr
}

type BaseAction struct {
    BtnId   string
    BtnData interface{}
    JS      template.JS
}

func (base *BaseAction) SetBtnId(btnId string) {
    if btnId[0] != '.' && btnId[0] != '#' {
        base.BtnId = "." + btnId
    } else {
        base.BtnId = btnId
    }
}
func (base *BaseAction) Js() template.JS              { return base.JS }
func (base *BaseAction) BtnClass() template.HTML      { return "" }
func (base *BaseAction) BtnAttribute() template.HTML  { return "" }
func (base *BaseAction) GetCallbacks() context.Node   { return context.Node{} }
func (base *BaseAction) ExtContent() template.HTML    { return template.HTML(``) }
func (base *BaseAction) FooterContent() template.HTML { return template.HTML(``) }
func (base *BaseAction) SetBtnData(data interface{})  { base.BtnData = data }

var _ types.Action = (*AjaxAction)(nil)
var _ types.Action = (*PopUpAction)(nil)
var _ types.Action = (*JumpAction)(nil)
var _ types.Action = (*JumpSelectBoxAction)(nil)

func URL(id string) string {
    return config.Url("/ajax/" + utils.WrapURL(id))
}
