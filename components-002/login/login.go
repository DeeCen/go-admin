package login

import (
    "bytes"
    "encoding/base64"
    "fmt"
    "html/template"
    textTemplate "text/template"

    "github.com/GoAdminGroup/components/login/theme1"
    ctx2 "github.com/GoAdminGroup/go-admin/context"
    "github.com/GoAdminGroup/go-admin/engine"
    "github.com/GoAdminGroup/go-admin/modules/logger"
    captcha2 "github.com/GoAdminGroup/go-admin/plugins/admin/modules/captcha"
    template2 "github.com/GoAdminGroup/go-admin/template"
    "github.com/GoAdminGroup/go-admin/template/login"
    "github.com/GoAdminGroup/go-admin/template/types"
    "github.com/dchest/captcha"
)

var themes = map[string]Theme{
    "theme1": new(theme1.Theme1),
}

func Register(key string, theme Theme) {
    if _, ok := themes[key]; ok {
        panic("duplicate login theme")
    }
    themes[key] = theme
}

type Login struct {
    TencentWaterProofWallData TencentWaterProofWallData `json:"tencent_water_proof_wall_data"`
    CaptchaDigits             int                       `json:"captcha_digits"`
    CaptchaID                 string                    `json:"captcha_id"`
    CaptchaImgSrc             string                    `json:"captcha_img_src"`
    Theme                     string                    `json:"theme"`
}

type TencentWaterProofWallData struct {
    AppID     string `json:"app_id"`
    AppSecret string `json:"app_secret"`
}

type Config struct {
    TencentWaterProofWallData TencentWaterProofWallData `json:"tencent_water_proof_wall_data"`
    CaptchaDigits             int                       `json:"captcha_digits"`
    Theme                     string                    `json:"theme"`
}

func Init(e *engine.Engine, cfg Config) {
    template2.AddLoginComp(Get(e, cfg))
}

func Get(e *engine.Engine, cfg Config) *Login {
    if cfg.CaptchaDigits != 0 && cfg.TencentWaterProofWallData.AppID == "" {
        //captchaData.Clean()
        captcha2.Add(CaptchaDriverKeyDefault, new(DigitsCaptcha))

        // 设置验证码方式 & 注册验证码获取api
        e.SetCaptchaDriver(CaptchaDriverKeyDefault)
        e.Data(`POST`, `/captchaRefresh`, captchaRefreshAPI, true)
    }

    if cfg.TencentWaterProofWallData.AppID != "" {
        captcha2.Add(CaptchaDriverKeyTencent, &TencentCaptcha{
            AppID:     cfg.TencentWaterProofWallData.AppID,
            AppSecret: cfg.TencentWaterProofWallData.AppSecret,
        })
    }

    if cfg.Theme == "" {
        cfg.Theme = "theme1"
    }

    return &Login{
        TencentWaterProofWallData: cfg.TencentWaterProofWallData,
        CaptchaDigits:             cfg.CaptchaDigits,
        Theme:                     cfg.Theme,
    }
}

func captchaRefreshAPI(ctx *ctx2.Context) {
    digitByte := captcha.RandomDigits(4)
    id := MakeCaptchaToken(byteToStr(digitByte))
    img := captcha.NewImage(id, digitByte, 110, 34)
    buf := new(bytes.Buffer)
    _, _ = img.WriteTo(buf)

    ctx.JSON(200, map[string]interface{}{
        "id":  id,
        "img": "data:image/png;base64," + base64.StdEncoding.EncodeToString(buf.Bytes()),
    })
}

func byteToStr(b []byte) string {
    s := ""
    for i := 0; i < len(b); i++ {
        s += fmt.Sprintf("%v", b[i])
    }
    return s
}

func (l *Login) GetTemplate() (*template.Template, string) {
    t := textTemplate.New("login").Delims("{%", "%}")
    t, err := t.Parse(themes[l.Theme].GetHTML())
    if err != nil {
        logger.Error("login component, get template parse error: ", err)
    }
    buf := new(bytes.Buffer)
    err = t.Execute(buf, l)
    if err != nil {
        logger.Error("login component, get template execute error: ", err)
    }

    tmpl, err := template.New("login").
        Funcs(login.DefaultFuncMap).
        Parse(buf.String())

    if err != nil {
        logger.Error("login component, get template error: ", err)
    }

    return tmpl, "login"
}

func (l *Login) GetAssetList() []string               { return themes[l.Theme].GetAssetList() }
func (l *Login) GetAsset(name string) ([]byte, error) { return themes[l.Theme].GetAsset(name[1:]) }
func (l *Login) GetName() string                      { return "login" }
func (l *Login) IsAPage() bool                        { return true }
func (l *Login) GetJS() template.JS                   { return "" }
func (l *Login) GetCSS() template.CSS                 { return "" }
func (l *Login) GetCallbacks() types.Callbacks        { return make(types.Callbacks, 0) }

func (l *Login) GetContent() template.HTML {
    buffer := new(bytes.Buffer)
    tmpl, defineName := l.GetTemplate()
    err := tmpl.ExecuteTemplate(buffer, defineName, l)
    if err != nil {
        logger.Error("login component, compose html error:", err)
    }
    return template.HTML(buffer.String())
}

type Theme interface {
    GetAssetList() []string
    GetAsset(name string) ([]byte, error)
    GetHTML() string
}
