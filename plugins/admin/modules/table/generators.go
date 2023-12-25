package table

import (
    "database/sql"
    "errors"
    "fmt"
    tmpl "html/template"
    "net/url"
    "regexp"
    "strconv"
    "strings"
    "time"

    "github.com/GoAdminGroup/go-admin/plugins/admin/modules/constant"

    "github.com/GoAdminGroup/go-admin/modules/auth"

    "github.com/GoAdminGroup/go-admin/context"
    "github.com/GoAdminGroup/go-admin/modules/collection"
    "github.com/GoAdminGroup/go-admin/modules/config"
    "github.com/GoAdminGroup/go-admin/modules/db"
    errs "github.com/GoAdminGroup/go-admin/modules/errors"
    "github.com/GoAdminGroup/go-admin/modules/language"
    "github.com/GoAdminGroup/go-admin/modules/logger"
    "github.com/GoAdminGroup/go-admin/modules/ui"
    "github.com/GoAdminGroup/go-admin/modules/utils"
    "github.com/GoAdminGroup/go-admin/plugins/admin/models"
    form2 "github.com/GoAdminGroup/go-admin/plugins/admin/modules/form"
    "github.com/GoAdminGroup/go-admin/plugins/admin/modules/parameter"
    "github.com/GoAdminGroup/go-admin/plugins/admin/modules/tools"
    "github.com/GoAdminGroup/go-admin/template"
    "github.com/GoAdminGroup/go-admin/template/types"
    "github.com/GoAdminGroup/go-admin/template/types/form"
    selection "github.com/GoAdminGroup/go-admin/template/types/form/select"
    "github.com/GoAdminGroup/html"
    "golang.org/x/crypto/bcrypt"
)

type SystemTable struct {
    conn db.Connection
    c    *config.Config
}

func NewSystemTable(conn db.Connection, c *config.Config) *SystemTable {
    return &SystemTable{conn: conn, c: c}
}

func (s *SystemTable) GetManagerTable(_ *context.Context) (managerTable Table) {
    managerTable = NewDefaultTable(DefaultConfigWithDriver(config.GetDatabases().GetDefault().Driver))

    info := managerTable.GetInfo().AddXssJsFilter().HideFilterArea()

    info.AddField("Id", "id", db.Int).FieldSortable()
    info.AddField(lg("name"), "username", db.Varchar).FieldFilterable()
    info.AddField(lg("nickname"), "name", db.Varchar).FieldFilterable()
    info.AddField(lg("role"), "role", db.Varchar).
        FieldDisplay(func(model types.FieldModel) interface{} {
            return `-`
        }).FieldFilterable()
    info.AddField(lg("createAt"), "createAt", db.Int).FieldDisplay(func(value types.FieldModel) interface{} {
        return toYmdHis(value.Value)
    })
    info.AddField(lg("updateAt"), "updateAt", db.Int).FieldDisplay(func(value types.FieldModel) interface{} {
        return toYmdHis(value.Value)
    })

    info.SetTable("goadmin_user").
        SetTitle(lg("managers")).
        SetDescription(lg("managers")).
        SetDeleteFn(func(idArr []string) error {
            var ids = interfaces(idArr)
            _, txErr := s.connection().WithTransaction(func(tx *sql.Tx) (e error, i map[string]interface{}) {
                deleteUserRoleErr := s.connection().WithTx(tx).
                    Table("goadmin_role_user").
                    WhereIn("userId", ids).
                    Delete()
                if db.CheckError(deleteUserRoleErr, db.DELETE) {
                    return deleteUserRoleErr, nil
                }

                deleteUserErr := s.connection().WithTx(tx).
                    Table("goadmin_user").
                    WhereIn("id", ids).
                    Delete()
                if db.CheckError(deleteUserErr, db.DELETE) {
                    return deleteUserErr, nil
                }

                return nil, nil
            })
            return txErr
        })

    formList := managerTable.GetForm().AddXssJsFilter()

    formList.AddField("Id", "id", db.Int, form.Default).FieldDisplayButCanNotEditWhenUpdate().FieldDisableWhenCreate()
    formList.AddField(lg("name"), "username", db.Varchar, form.Text).
        FieldHelpMsg(template.HTML(lg("use for login"))).FieldMust()
    formList.AddField(lg("nickname"), "name", db.Varchar, form.Text).
        FieldHelpMsg(template.HTML(lg("use to display"))).FieldMust()
    formList.AddField(lg("role"), "roleId", db.Varchar, form.SelectBox).
        FieldMust().
        FieldOptionsFromTable("goadmin_role", "name", "id").
        FieldDisplay(func(model types.FieldModel) interface{} {
            var roles []string

            if model.ID == "" {
                return roles
            }
            roleModel, _ := s.table("goadmin_role_user").Select("roleId").
                Where("userId", "=", model.ID).All()
            for _, v := range roleModel {
                roles = append(roles, strconv.FormatInt(v["roleId"].(int64), 10))
            }
            return roles
        }).FieldHelpMsg(template.HTML(lg("no corresponding options?")) +
        link(config.Url(`/info/roles/new`), "Create here."))
    formList.AddField(lg("avatar"), "avatar", db.Varchar, form.File)
    formList.AddField(lg("password"), "password", db.Varchar, form.Password).
        FieldDisplay(func(value types.FieldModel) interface{} {
            return ""
        })
    formList.AddField(lg("confirm password"), "password_again", db.Varchar, form.Password).
        FieldDisplay(func(value types.FieldModel) interface{} {
            return ""
        })

    formList.SetTable("goadmin_user").SetTitle(lg("managers")).SetDescription(lg("managers"))
    formList.SetUpdateFn(func(values form2.Values) error {
        if values.IsEmpty("name", "username") {
            return errors.New("username and password can not be empty")
        }

        user := models.UserWithId(values.Get("id")).SetConn(s.conn)
        password := values.Get("password")
        if password != "" {
            if password != values.Get("password_again") {
                return errors.New("password does not match")
            }
            password = encodePassword([]byte(values.Get("password")))
        }

        _, txErr := s.connection().WithTransaction(func(tx *sql.Tx) (e error, i map[string]interface{}) {
            avatar := values.Get("avatar")
            if values.Get("avatar__delete_flag") == "1" {
                avatar = ""
            }

            _, updateUserErr := user.WithTx(tx).Update(values.Get("username"),
                password, values.Get("name"), avatar, values.Get("avatar__change_flag") == "1")

            if db.CheckError(updateUserErr, db.UPDATE) {
                return updateUserErr, nil
            }

            delRoleErr := user.WithTx(tx).DeleteRoles()
            if db.CheckError(delRoleErr, db.DELETE) {
                return delRoleErr, nil
            }

            for i := 0; i < len(values["roleId[]"]); i++ {
                _, addRoleErr := user.WithTx(tx).AddRole(values["roleId[]"][i])
                if db.CheckError(addRoleErr, db.INSERT) {
                    return addRoleErr, nil
                }
            }
            return nil, nil
        })

        return txErr
    })
    formList.SetInsertFn(func(values form2.Values) error {
        if values.IsEmpty("name", "username", "password") {
            return errors.New("username and password can not be empty")
        }

        password := values.Get("password")
        if password != values.Get("password_again") {
            return errors.New("password does not match")
        }

        _, txErr := s.connection().WithTransaction(func(tx *sql.Tx) (e error, i map[string]interface{}) {
            user, createUserErr := models.User().WithTx(tx).SetConn(s.conn).New(values.Get("username"),
                encodePassword([]byte(values.Get("password"))),
                values.Get("name"),
                values.Get("avatar"))

            if db.CheckError(createUserErr, db.INSERT) {
                return createUserErr, nil
            }

            for i := 0; i < len(values["roleId[]"]); i++ {
                _, addRoleErr := user.WithTx(tx).AddRole(values["roleId[]"][i])
                if db.CheckError(addRoleErr, db.INSERT) {
                    return addRoleErr, nil
                }
            }

            return nil, nil
        })
        return txErr
    })

    formList.HideContinueEditCheckBox()
    formList.HideContinueNewCheckBox()
    formList.HideResetButton()

    return
}

func (s *SystemTable) GetRoleTable(_ *context.Context) (roleTable Table) {
    roleTable = NewDefaultTable(DefaultConfigWithDriver(config.GetDatabases().GetDefault().Driver))

    info := roleTable.GetInfo().AddXssJsFilter().HideFilterArea()

    info.AddField("Id", "id", db.Int).FieldSortable()
    info.AddField(lg("role"), "name", db.Varchar).FieldFilterable()
    //info.AddField(lg("slug"), "slug", db.Varchar).FieldFilterable()
    info.AddField(lg("createAt"), "createAt", db.Int).FieldDisplay(func(value types.FieldModel) interface{} {
        return toYmdHis(value.Value)
    })
    info.AddField(lg("updateAt"), "updateAt", db.Int).FieldDisplay(func(value types.FieldModel) interface{} {
        return toYmdHis(value.Value)
    })

    info.SetTable("goadmin_role").
        SetTitle(lg("roles manage")).
        SetDescription(lg("roles manage")).
        SetDeleteFn(func(idArr []string) error {

            var ids = interfaces(idArr)

            _, txErr := s.connection().WithTransaction(func(tx *sql.Tx) (e error, i map[string]interface{}) {

                deleteRoleUserErr := s.connection().WithTx(tx).
                    Table("goadmin_role_user").
                    WhereIn("roleId", ids).
                    Delete()

                if db.CheckError(deleteRoleUserErr, db.DELETE) {
                    return deleteRoleUserErr, nil
                }

                deleteRoleMenuErr := s.connection().WithTx(tx).
                    Table("goadmin_role_menu").
                    WhereIn("roleId", ids).
                    Delete()

                if db.CheckError(deleteRoleMenuErr, db.DELETE) {
                    return deleteRoleMenuErr, nil
                }

                deleteRolesErr := s.connection().WithTx(tx).
                    Table("goadmin_role").
                    WhereIn("id", ids).
                    Delete()

                if db.CheckError(deleteRolesErr, db.DELETE) {
                    return deleteRolesErr, nil
                }

                return nil, nil
            })

            return txErr
        })

    formList := roleTable.GetForm().AddXssJsFilter()
    formList.SetTable("goadmin_role").
        SetTitle(lg("roles manage")).
        SetDescription(lg("roles manage"))

    formList.AddField("Id", "id", db.Int, form.Default).FieldDisplayButCanNotEditWhenUpdate().FieldDisableWhenCreate()
    formList.AddField(lg("role"), "name", db.Varchar, form.Text).
        FieldMust().
        FieldHelpMsg(template.HTML(lg("should be unique")))
    //formList.AddField(lg("slug"), "slug", db.Varchar, form.Text).FieldHelpMsg(template.HTML(lg("should be unique"))).FieldMust()
    formList.AddField(lg("permission"), "menu_id", db.Varchar, form.SelectBox).
        FieldOptionsFromTable("goadmin_menu", "title", "id").
        FieldDisplay(func(model types.FieldModel) interface{} {
            var permissions = make([]string, 0)

            if model.ID == "" {
                return permissions
            }
            perModel, _ := s.table("goadmin_role_menu").
                Select("menuId").
                Where("roleId", "=", model.ID).
                All()
            for _, v := range perModel {
                permissions = append(permissions, strconv.FormatInt(v["menuId"].(int64), 10))
            }
            return permissions
        }).FieldHelpMsg(template.HTML(lg("no corresponding options?")) +
        link(config.Url(`/info/menu`), "Create here."))

    formList.AddField(lg("updateAt"), "updateAt", db.Int, form.Default).
        FieldDisableWhenCreate().
        FieldHide().
        FieldDisplay(func(value types.FieldModel) interface{} {
            return toYmdHis(value.Value)
        })
    formList.AddField(lg("createAt"), "createAt", db.Int, form.Default).
        FieldDisableWhenCreate().
        FieldHide().
        FieldDisplay(func(value types.FieldModel) interface{} {
            return toYmdHis(value.Value)
        })

    formList.SetUpdateFn(func(values form2.Values) error {

        if models.Role().SetConn(s.conn).IsNameExist(values.Get("name"), values.Get("id")) {
            return errors.New("slug exists")
        }

        role := models.RoleWithId(values.Get("id")).SetConn(s.conn)

        _, txErr := s.connection().WithTransaction(func(tx *sql.Tx) (e error, i map[string]interface{}) {

            _, updateRoleErr := role.WithTx(tx).Update(values.Get("name"), values.Get("slug"))

            if db.CheckError(updateRoleErr, db.UPDATE) {
                return updateRoleErr, nil
            }

            delPermissionErr := role.WithTx(tx).DeletePermissions()

            if db.CheckError(delPermissionErr, db.DELETE) {
                return delPermissionErr, nil
            }

            for i := 0; i < len(values["menu_id[]"]); i++ {
                _, addPermissionErr := role.WithTx(tx).AddPermission(values["menu_id[]"][i])
                if db.CheckError(addPermissionErr, db.INSERT) {
                    return addPermissionErr, nil
                }
            }

            return nil, nil
        })

        return txErr
    })

    formList.SetInsertFn(func(values form2.Values) error {

        if models.Role().SetConn(s.conn).IsNameExist(values.Get("name"), "") {
            return errors.New("slug exists")
        }

        _, txErr := s.connection().WithTransaction(func(tx *sql.Tx) (e error, i map[string]interface{}) {
            role, createRoleErr := models.Role().WithTx(tx).SetConn(s.conn).New(values.Get("name"), values.Get("slug"))

            if db.CheckError(createRoleErr, db.INSERT) {
                return createRoleErr, nil
            }

            for i := 0; i < len(values["menu_id[]"]); i++ {
                _, addPermissionErr := role.WithTx(tx).AddPermission(values["menu_id[]"][i])
                if db.CheckError(addPermissionErr, db.INSERT) {
                    return addPermissionErr, nil
                }
            }

            return nil, nil
        })

        return txErr
    })

    formList.HideContinueEditCheckBox()
    formList.HideContinueNewCheckBox()
    formList.HideResetButton()

    return
}

func (s *SystemTable) GetMenuTable(ctx *context.Context) (menuTable Table) {
    menuTable = NewDefaultTable(DefaultConfigWithDriver(config.GetDatabases().GetDefault().Driver))

    name := ctx.Query("__pluginName")

    info := menuTable.GetInfo().AddXssJsFilter().HideFilterArea().Where("pluginName", "=", name)

    info.AddField("Id", "id", db.Int).FieldSortable()
    info.AddField(lg("parent"), "parentId", db.Int)
    info.AddField(lg("menu name"), "title", db.Varchar)
    info.AddField(lg("icon"), "icon", db.Varchar)
    info.AddField(lg("uri"), "uri", db.Varchar)
    info.AddField(lg("role"), "roles", db.Varchar)
    info.AddField(lg("header"), "header", db.Varchar)
    info.AddField(lg("createAt"), "createAt", db.Int).FieldDisplay(func(value types.FieldModel) interface{} {
        return toYmdHis(value.Value)
    })
    info.AddField(lg("updateAt"), "updateAt", db.Int).FieldDisplay(func(value types.FieldModel) interface{} {
        return toYmdHis(value.Value)
    })

    info.SetTable("goadmin_menu").
        SetTitle(lg("menus manage")).
        SetDescription(lg("menus manage")).
        SetDeleteFn(func(idArr []string) error {

            var ids = interfaces(idArr)

            _, txErr := s.connection().WithTransaction(func(tx *sql.Tx) (e error, i map[string]interface{}) {

                deleteRoleMenuErr := s.connection().WithTx(tx).
                    Table("goadmin_role_menu").
                    WhereIn("menuId", ids).
                    Delete()

                if db.CheckError(deleteRoleMenuErr, db.DELETE) {
                    return deleteRoleMenuErr, nil
                }

                deleteMenusErr := s.connection().WithTx(tx).
                    Table("goadmin_menu").
                    WhereIn("id", ids).
                    Delete()

                if db.CheckError(deleteMenusErr, db.DELETE) {
                    return deleteMenusErr, nil
                }

                return nil, map[string]interface{}{}
            })

            return txErr
        })

    var parentIDOptions = types.FieldOptions{
        {
            Text:  "ROOT",
            Value: "0",
        },
    }

    allMenus, _ := s.connection().Table("goadmin_menu").
        Where("parentId", "=", 0).
        Where("pluginName", "=", name).
        Select("id", "title").
        OrderBy("order", "asc").
        All()
    allMenuIDs := make([]interface{}, len(allMenus))

    if len(allMenuIDs) > 0 {
        for i := 0; i < len(allMenus); i++ {
            allMenuIDs[i] = allMenus[i]["id"]
        }

        secondLevelMenus, _ := s.connection().Table("goadmin_menu").
            WhereIn("parentId", allMenuIDs).
            Where("pluginName", "=", name).
            Select("id", "title", "parentId").
            All()

        secondLevelMenusCol := collection.Collection(secondLevelMenus)

        for i := 0; i < len(allMenus); i++ {
            parentIDOptions = append(parentIDOptions, types.FieldOption{
                TextHTML: "&nbsp;&nbsp;┝  " + language.GetFromHtml(template.HTML(allMenus[i]["title"].(string))),
                Value:    strconv.Itoa(int(allMenus[i]["id"].(int64))),
            })
            col := secondLevelMenusCol.Where("parentId", "=", allMenus[i]["id"].(int64))
            for j := 0; j < len(col); j++ {
                parentIDOptions = append(parentIDOptions, types.FieldOption{
                    TextHTML: "&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;┝  " +
                        language.GetFromHtml(template.HTML(col[j]["title"].(string))),
                    Value: strconv.Itoa(int(col[j]["id"].(int64))),
                })
            }
        }
    }

    formList := menuTable.GetForm().AddXssJsFilter()
    formList.AddField("Id", "id", db.Int, form.Default).
        FieldDisplayButCanNotEditWhenUpdate().
        FieldDisableWhenCreate()
    formList.AddField(lg("parent"), "parentId", db.Int, form.SelectSingle).
        FieldOptions(parentIDOptions).
        FieldMust().
        FieldDisplay(func(model types.FieldModel) interface{} {
            var menuItem []string

            if model.ID == "" {
                return menuItem
            }

            menuModel, _ := s.table("goadmin_menu").Select("parentId").Find(model.ID)
            parentId := strconv.FormatInt(menuModel["parentId"].(int64), 10)
            menuItem = append(menuItem, parentId)

            return menuItem
        })
    formList.AddField(lg("menu name"), "title", db.Varchar, form.Text).FieldMust().FieldLimit(40)
    formList.AddField(lg("icon"), "icon", db.Varchar, form.IconPicker).FieldMust()
    formList.AddField(lg("uri"), "uri", db.Varchar, form.Text).FieldMust().FieldLimit(40)
    formList.AddField(lg("header"), "header", db.Varchar, form.Text).FieldLimit(40)
    formList.AddField("PluginName", "pluginName", db.Varchar, form.Text).
        FieldDefault(name).
        FieldHide()

    ts := fmt.Sprintf(`%d`, time.Now().Unix())
    formList.AddField("createAt", "createAt", db.Int, form.Default).
        FieldDefault(ts).
        FieldHide().
        FieldDisableWhenUpdate()
    formList.AddField("updateAt", "updateAt", db.Int, form.Default).
        FieldDefault(ts).
        FieldHide()

    formList.SetTable("goadmin_menu").
        SetTitle(lg("menus manage")).
        SetDescription(lg("menus manage"))

    formList.HideContinueEditCheckBox()
    formList.HideContinueNewCheckBox()
    formList.HideResetButton()

    return
}

func (s *SystemTable) GetNormalManagerTable(ctx *context.Context) (managerTable Table) {
    loginUserID := auth.GetUserID(ctx)
    editID, _ := strconv.Atoi(ctx.Query(constant.EditPKKey))
    if editID > 0 && int64(editID) != loginUserID {
        ctx.Write(404, nil, `no permission`)
        ctx.Abort()
    }

    managerTable = NewDefaultTable(DefaultConfigWithDriver(config.GetDatabases().GetDefault().Driver))

    info := managerTable.GetInfo().AddXssJsFilter().HideFilterArea()
    info.HideNewButton()
    info.HideExportButton()
    info.HideDeleteButton()
    info.HideCheckBoxColumn()
    info.HideEditButton()
    info.HideDetailButton()
    info.Where(`id`, `=`, loginUserID)

    info.AddField("Id", "id", db.Int).FieldSortable()
    info.AddField(lg("name"), "username", db.Varchar).FieldFilterable()
    info.AddField(lg("nickname"), "name", db.Varchar).FieldFilterable()
    info.AddField(lg("createAt"), "createAt", db.Int).FieldDisplay(func(value types.FieldModel) interface{} {
        ts, _ := strconv.Atoi(value.Value)
        return time.Unix(int64(ts), 0).Format(`2006-01-02 15:04:05`)
    })

    info.SetTable("goadmin_user").
        SetTitle(lg("managers")).
        SetDescription(``).
        SetDeleteFn(func(_ []string) error {
            return errors.New(`禁止删除操作`)
        })

    formList := managerTable.GetForm().AddXssJsFilter()
    formList.HideResetButton()
    formList.HideContinueEditCheckBox()
    formList.HideContinueNewCheckBox()
    formList.HideBackButton()

    formList.AddField("Id", "id", db.Int, form.Default).FieldDisplayButCanNotEditWhenUpdate().FieldDisableWhenCreate()
    formList.AddField(lg("name"), "username", db.Varchar, form.Text).FieldHelpMsg(template.HTML(lg("use for login"))).FieldMust()
    formList.AddField(lg("nickname"), "name", db.Varchar, form.Text).FieldHelpMsg(template.HTML(lg("use to display"))).FieldMust()
    formList.AddField(lg("avatar"), "avatar", db.Varchar, form.File)
    formList.AddField(lg("password"), "password", db.Varchar, form.Password).
        FieldDisplay(func(value types.FieldModel) interface{} {
            return ""
        })
    formList.AddField(lg("confirm password"), "password_again", db.Varchar, form.Password).
        FieldDisplay(func(value types.FieldModel) interface{} {
            return ""
        })

    formList.SetTable("goadmin_user").SetTitle(lg("managers")).SetDescription(``)
    formList.SetUpdateFn(func(values form2.Values) error {
        if values.IsEmpty("name", "username") {
            return errors.New("username and password can not be empty")
        }

        //user := models.UserWithId(values.Get("id")).SetConn(s.conn)
        user := models.UserWithId(fmt.Sprintf(`%d`, loginUserID)).SetConn(s.conn)
        if values.Has("permission", "role") {
            return errors.New(errs.NoPermission)
        }

        password := values.Get("password")
        if password != "" {
            if password != values.Get("password_again") {
                return errors.New("password does not match")
            }
            password = encodePassword([]byte(values.Get("password")))
        }

        avatar := values.Get("avatar")
        if values.Get("avatar__delete_flag") == "1" {
            avatar = ""
        }

        _, updateUserErr := user.Update(values.Get("username"),
            password, values.Get("name"), avatar, values.Get("avatar__change_flag") == "1")
        if db.CheckError(updateUserErr, db.UPDATE) {
            return updateUserErr
        }

        return nil
    })
    formList.SetInsertFn(func(values form2.Values) error {
        return errors.New(`禁止添加操作`)

        /*if values.IsEmpty("name", "username", "password") {
              return errors.New("username and password can not be empty")
          }

          password := values.Get("password")

          if password != values.Get("password_again") {
              return errors.New("password does not match")
          }

          if values.Has("permission", "role") {
              return errors.New(errs.NoPermission)
          }

          _, createUserErr := models.User().SetConn(s.conn).New(values.Get("username"),
              encodePassword([]byte(values.Get("password"))),
              values.Get("name"),
              values.Get("avatar"))

          if db.CheckError(createUserErr, db.INSERT) {
              return createUserErr
          }

          return nil*/
    })

    return
}

func (s *SystemTable) GetSiteTable(_ *context.Context) (siteTable Table) {
    siteTable = NewDefaultTable(DefaultConfigWithDriver(config.GetDatabases().GetDefault().Driver).
        SetOnlyUpdateForm().
        SetGetDataFun(func(_ parameter.Parameters) ([]map[string]interface{}, int) {
            return []map[string]interface{}{models.Site().SetConn(s.conn).AllToMapInterface()}, 1
        }))

    trueStr := lgWithConfigScore("true")
    falseStr := lgWithConfigScore("false")

    formList := siteTable.GetForm().AddXssJsFilter()
    formList.AddField("Id", "id", db.Varchar, form.Default).FieldDefault("1").FieldHide()
    formList.AddField(lgWithConfigScore("site off"), "site_off", db.Varchar, form.Switch).
        FieldOptions(types.FieldOptions{
            {Text: trueStr, Value: "true"},
            {Text: falseStr, Value: "false"},
        })
    formList.AddField(lgWithConfigScore("debug"), "debug", db.Varchar, form.Switch).
        FieldOptions(types.FieldOptions{
            {Text: trueStr, Value: "true"},
            {Text: falseStr, Value: "false"},
        })
    formList.AddField(lgWithConfigScore("env"), "env", db.Varchar, form.Default).
        FieldDisplay(func(value types.FieldModel) interface{} {
            return s.c.Env
        })

    langOps := make(types.FieldOptions, len(language.Langs))
    for k, t := range language.Langs {
        langOps[k] = types.FieldOption{Text: lgWithConfigScore(t, "language"), Value: t}
    }
    formList.AddField(lgWithConfigScore("language"), "language", db.Varchar, form.SelectSingle).
        FieldDisplay(func(value types.FieldModel) interface{} {
            return language.FixedLanguageKey(value.Value)
        }).
        FieldOptions(langOps)
    themes := template.Themes()
    themesOps := make(types.FieldOptions, len(themes))
    for k, t := range themes {
        themesOps[k] = types.FieldOption{Text: t, Value: t}
    }

    formList.AddField(lgWithConfigScore("theme"), "theme", db.Varchar, form.SelectSingle).
        FieldOptions(themesOps).
        FieldOnChooseShow("adminlte",
            "color_scheme")
    formList.AddField(lgWithConfigScore("title"), "title", db.Varchar, form.Text).FieldMust()
    formList.AddField(lgWithConfigScore("color scheme"), "color_scheme", db.Varchar, form.SelectSingle).
        FieldOptions(types.FieldOptions{
            {Text: "skin-black", Value: "skin-black"},
            {Text: "skin-black-light", Value: "skin-black-light"},
            {Text: "skin-blue", Value: "skin-blue"},
            {Text: "skin-blue-light", Value: "skin-blue-light"},
            {Text: "skin-green", Value: "skin-green"},
            {Text: "skin-green-light", Value: "skin-green-light"},
            {Text: "skin-purple", Value: "skin-purple"},
            {Text: "skin-purple-light", Value: "skin-purple-light"},
            {Text: "skin-red", Value: "skin-red"},
            {Text: "skin-red-light", Value: "skin-red-light"},
            {Text: "skin-yellow", Value: "skin-yellow"},
            {Text: "skin-yellow-light", Value: "skin-yellow-light"},
        }).FieldHelpMsg(template.HTML(lgWithConfigScore("It will work when theme is adminlte")))
    formList.AddField(lgWithConfigScore("login title"), "login_title", db.Varchar, form.Text).FieldMust()
    formList.AddField(lgWithConfigScore("extra"), "extra", db.Varchar, form.TextArea)
    formList.AddField(lgWithConfigScore("logo"), "logo", db.Varchar, form.Code).FieldMust()
    formList.AddField(lgWithConfigScore("mini logo"), "mini_logo", db.Varchar, form.Code).FieldMust()
    if s.c.IsNotProductionEnvironment() {
        formList.AddField(lgWithConfigScore("bootstrap file path"), "bootstrap_file_path", db.Varchar, form.Text)
        formList.AddField(lgWithConfigScore("go mod file path"), "go_mod_file_path", db.Varchar, form.Text)
    }
    formList.AddField(lgWithConfigScore("session life time"), "session_life_time", db.Varchar, form.Number).
        FieldMust().
        FieldHelpMsg(template.HTML(lgWithConfigScore("must bigger than 900 seconds")))
    formList.AddField(lgWithConfigScore("custom head html"), "custom_head_html", db.Varchar, form.Code)
    formList.AddField(lgWithConfigScore("custom foot Html"), "custom_foot_html", db.Varchar, form.Code)
    formList.AddField(lgWithConfigScore("custom 404 html"), "custom_404_html", db.Varchar, form.Code)
    formList.AddField(lgWithConfigScore("custom 403 html"), "custom_403_html", db.Varchar, form.Code)
    formList.AddField(lgWithConfigScore("custom 500 Html"), "custom_500_html", db.Varchar, form.Code)
    formList.AddField(lgWithConfigScore("footer info"), "footer_info", db.Varchar, form.Code)
    formList.AddField(lgWithConfigScore("login logo"), "login_logo", db.Varchar, form.Code)
    formList.AddField(lgWithConfigScore("no limit login ip"), "no_limit_login_ip", db.Varchar, form.Switch).
        FieldOptions(types.FieldOptions{
            {Text: trueStr, Value: "true"},
            {Text: falseStr, Value: "false"},
        })
    formList.AddField(lgWithConfigScore("operation log off"), "operation_log_off", db.Varchar, form.Switch).
        FieldOptions(types.FieldOptions{
            {Text: trueStr, Value: "true"},
            {Text: falseStr, Value: "false"},
        })
    formList.AddField(lgWithConfigScore("allow delete operation log"), "allow_del_operation_log", db.Varchar, form.Switch).
        FieldOptions(types.FieldOptions{
            {Text: trueStr, Value: "true"},
            {Text: falseStr, Value: "false"},
        })
    formList.AddField(lgWithConfigScore("hide config center entrance"), "hide_config_center_entrance", db.Varchar, form.Switch).
        FieldOptions(types.FieldOptions{
            {Text: trueStr, Value: "true"},
            {Text: falseStr, Value: "false"},
        })
    formList.AddField(lgWithConfigScore("hide app info entrance"), "hide_app_info_entrance", db.Varchar, form.Switch).
        FieldOptions(types.FieldOptions{
            {Text: trueStr, Value: "true"},
            {Text: falseStr, Value: "false"},
        })
    formList.AddField(lgWithConfigScore("hide tool entrance"), "hide_tool_entrance", db.Varchar, form.Switch).
        FieldOptions(types.FieldOptions{
            {Text: trueStr, Value: "true"},
            {Text: falseStr, Value: "false"},
        })
    formList.AddField(lgWithConfigScore("hide plugin entrance"), "hide_plugin_entrance", db.Varchar, form.Switch).
        FieldOptions(types.FieldOptions{
            {Text: trueStr, Value: "true"},
            {Text: falseStr, Value: "false"},
        })
    formList.AddField(lgWithConfigScore("animation type"), "animation_type", db.Varchar, form.SelectSingle).
        FieldOptions(types.FieldOptions{
            {Text: "", Value: ""},
            {Text: "bounce", Value: "bounce"}, {Text: "flash", Value: "flash"}, {Text: "pulse", Value: "pulse"},
            {Text: "rubberBand", Value: "rubberBand"}, {Text: "shake", Value: "shake"}, {Text: "swing", Value: "swing"},
            {Text: "tada", Value: "tada"}, {Text: "wobble", Value: "wobble"}, {Text: "jello", Value: "jello"},
            {Text: "heartBeat", Value: "heartBeat"}, {Text: "bounceIn", Value: "bounceIn"}, {Text: "bounceInDown", Value: "bounceInDown"},
            {Text: "bounceInLeft", Value: "bounceInLeft"}, {Text: "bounceInRight", Value: "bounceInRight"}, {Text: "bounceInUp", Value: "bounceInUp"},
            {Text: "fadeIn", Value: "fadeIn"}, {Text: "fadeInDown", Value: "fadeInDown"}, {Text: "fadeInDownBig", Value: "fadeInDownBig"},
            {Text: "fadeInLeft", Value: "fadeInLeft"}, {Text: "fadeInLeftBig", Value: "fadeInLeftBig"}, {Text: "fadeInRight", Value: "fadeInRight"},
            {Text: "fadeInRightBig", Value: "fadeInRightBig"}, {Text: "fadeInUp", Value: "fadeInUp"}, {Text: "fadeInUpBig", Value: "fadeInUpBig"},
            {Text: "flip", Value: "flip"}, {Text: "flipInX", Value: "flipInX"}, {Text: "flipInY", Value: "flipInY"},
            {Text: "lightSpeedIn", Value: "lightSpeedIn"}, {Text: "rotateIn", Value: "rotateIn"}, {Text: "rotateInDownLeft", Value: "rotateInDownLeft"},
            {Text: "rotateInDownRight", Value: "rotateInDownRight"}, {Text: "rotateInUpLeft", Value: "rotateInUpLeft"}, {Text: "rotateInUpRight", Value: "rotateInUpRight"},
            {Text: "slideInUp", Value: "slideInUp"}, {Text: "slideInDown", Value: "slideInDown"}, {Text: "slideInLeft", Value: "slideInLeft"}, {Text: "slideInRight", Value: "slideInRight"},
            {Text: "slideOutRight", Value: "slideOutRight"}, {Text: "zoomIn", Value: "zoomIn"}, {Text: "zoomInDown", Value: "zoomInDown"},
            {Text: "zoomInLeft", Value: "zoomInLeft"}, {Text: "zoomInRight", Value: "zoomInRight"}, {Text: "zoomInUp", Value: "zoomInUp"},
            {Text: "hinge", Value: "hinge"}, {Text: "jackInTheBox", Value: "jackInTheBox"}, {Text: "rollIn", Value: "rollIn"},
        }).FieldOnChooseHide("", "animation_duration", "animation_delay").
        FieldOptionExt(map[string]interface{}{"allowClear": true}).
        FieldHelpMsg(`see more: <a href="https://daneden.github.io/animate.css/">https://daneden.github.io/animate.css/</a>`)

    formList.AddField(lgWithConfigScore("animation duration"), "animation_duration", db.Varchar, form.Number)
    formList.AddField(lgWithConfigScore("animation delay"), "animation_delay", db.Varchar, form.Number)

    formList.AddField(lgWithConfigScore("file upload engine"), "file_upload_engine", db.Varchar, form.Text)

    formList.AddField(lgWithConfigScore("cdn url"), "asset_url", db.Varchar, form.Text).
        FieldHelpMsg(template.HTML(lgWithConfigScore("Do not modify when you have not set up all assets")))

    formList.AddField(lgWithConfigScore("info log path"), "info_log_path", db.Varchar, form.Text)
    formList.AddField(lgWithConfigScore("error log path"), "error_log_path", db.Varchar, form.Text)
    formList.AddField(lgWithConfigScore("access log path"), "access_log_path", db.Varchar, form.Text)
    formList.AddField(lgWithConfigScore("info log off"), "info_log_off", db.Varchar, form.Switch).
        FieldOptions(types.FieldOptions{
            {Text: trueStr, Value: "true"},
            {Text: falseStr, Value: "false"},
        })
    formList.AddField(lgWithConfigScore("error log off"), "error_log_off", db.Varchar, form.Switch).
        FieldOptions(types.FieldOptions{
            {Text: trueStr, Value: "true"},
            {Text: falseStr, Value: "false"},
        })
    formList.AddField(lgWithConfigScore("access log off"), "access_log_off", db.Varchar, form.Switch).
        FieldOptions(types.FieldOptions{
            {Text: trueStr, Value: "true"},
            {Text: falseStr, Value: "false"},
        })
    formList.AddField(lgWithConfigScore("access assets log off"), "access_assets_log_off", db.Varchar, form.Switch).
        FieldOptions(types.FieldOptions{
            {Text: trueStr, Value: "true"},
            {Text: falseStr, Value: "false"},
        })
    formList.AddField(lgWithConfigScore("sql log on"), "sql_log", db.Varchar, form.Switch).
        FieldOptions(types.FieldOptions{
            {Text: trueStr, Value: "true"},
            {Text: falseStr, Value: "false"},
        })
    formList.AddField(lgWithConfigScore("log level"), "logger_level", db.Varchar, form.SelectSingle).
        FieldOptions(types.FieldOptions{
            {Text: "Debug", Value: "-1"},
            {Text: "Info", Value: "0"},
            {Text: "Warn", Value: "1"},
            {Text: "Error", Value: "2"},
        }).FieldDisplay(defaultFilterFn("0"))

    formList.AddField(lgWithConfigScore("logger rotate max size"), "logger_rotate_max_size", db.Varchar, form.Number).
        FieldDivider(lgWithConfigScore("logger rotate")).FieldDisplay(defaultFilterFn("10", "0"))
    formList.AddField(lgWithConfigScore("logger rotate max backups"), "logger_rotate_max_backups", db.Varchar, form.Number).
        FieldDisplay(defaultFilterFn("5", "0"))
    formList.AddField(lgWithConfigScore("logger rotate max age"), "logger_rotate_max_age", db.Varchar, form.Number).
        FieldDisplay(defaultFilterFn("30", "0"))
    formList.AddField(lgWithConfigScore("logger rotate compress"), "logger_rotate_compress", db.Varchar, form.Switch).
        FieldOptions(types.FieldOptions{
            {Text: trueStr, Value: "true"},
            {Text: falseStr, Value: "false"},
        }).FieldDisplay(defaultFilterFn("false"))

    formList.AddField(lgWithConfigScore("logger rotate encoder encoding"), "logger_encoder_encoding", db.Varchar,
        form.SelectSingle).
        FieldDivider(lgWithConfigScore("logger rotate encoder")).
        FieldOptions(types.FieldOptions{
            {Text: "JSON", Value: "json"},
            {Text: "Console", Value: "console"},
        }).FieldDisplay(defaultFilterFn("console")).
        FieldOnChooseHide("Console",
            "logger_encoder_time_key", "logger_encoder_level_key", "logger_encoder_caller_key",
            "logger_encoder_message_key", "logger_encoder_stacktrace_key", "logger_encoder_name_key")

    formList.AddField(lgWithConfigScore("logger rotate encoder time key"), "logger_encoder_time_key", db.Varchar, form.Text).
        FieldDisplay(defaultFilterFn("ts"))
    formList.AddField(lgWithConfigScore("logger rotate encoder level key"), "logger_encoder_level_key", db.Varchar, form.Text).
        FieldDisplay(defaultFilterFn("level"))
    formList.AddField(lgWithConfigScore("logger rotate encoder name key"), "logger_encoder_name_key", db.Varchar, form.Text).
        FieldDisplay(defaultFilterFn("logger"))
    formList.AddField(lgWithConfigScore("logger rotate encoder caller key"), "logger_encoder_caller_key", db.Varchar, form.Text).
        FieldDisplay(defaultFilterFn("caller"))
    formList.AddField(lgWithConfigScore("logger rotate encoder message key"), "logger_encoder_message_key", db.Varchar, form.Text).
        FieldDisplay(defaultFilterFn("msg"))
    formList.AddField(lgWithConfigScore("logger rotate encoder stacktrace key"), "logger_encoder_stacktrace_key", db.Varchar, form.Text).
        FieldDisplay(defaultFilterFn("stacktrace"))

    formList.AddField(lgWithConfigScore("logger rotate encoder level"), "logger_encoder_level", db.Varchar,
        form.SelectSingle).
        FieldOptions(types.FieldOptions{
            {Text: lgWithConfigScore("capital"), Value: "capital"},
            {Text: lgWithConfigScore("capitalColor"), Value: "capitalColor"},
            {Text: lgWithConfigScore("lowercase"), Value: "lowercase"},
            {Text: lgWithConfigScore("lowercaseColor"), Value: "color"},
        }).FieldDisplay(defaultFilterFn("capitalColor"))
    formList.AddField(lgWithConfigScore("logger rotate encoder time"), "logger_encoder_time", db.Varchar,
        form.SelectSingle).
        FieldOptions(types.FieldOptions{
            {Text: "ISO8601(2006-01-02T15:04:05.000Z0700)", Value: "iso8601"},
            {Text: lgWithConfigScore("millisecond"), Value: "millis"},
            {Text: lgWithConfigScore("nanosecond"), Value: "nanos"},
            {Text: "RFC3339(2006-01-02T15:04:05Z07:00)", Value: "rfc3339"},
            {Text: "RFC3339 Nano(2006-01-02T15:04:05.999999999Z07:00)", Value: "rfc3339nano"},
        }).FieldDisplay(defaultFilterFn("iso8601"))
    formList.AddField(lgWithConfigScore("logger rotate encoder duration"), "logger_encoder_duration", db.Varchar,
        form.SelectSingle).
        FieldOptions(types.FieldOptions{
            {Text: lgWithConfigScore("seconds"), Value: "string"},
            {Text: lgWithConfigScore("nanosecond"), Value: "nanos"},
            {Text: lgWithConfigScore("microsecond"), Value: "ms"},
        }).FieldDisplay(defaultFilterFn("string"))
    formList.AddField(lgWithConfigScore("logger rotate encoder caller"), "logger_encoder_caller", db.Varchar,
        form.SelectSingle).
        FieldOptions(types.FieldOptions{
            {Text: lgWithConfigScore("full path"), Value: "full"},
            {Text: lgWithConfigScore("short path"), Value: "short"},
        }).FieldDisplay(defaultFilterFn("full"))

    formList.HideBackButton().HideContinueEditCheckBox().HideContinueNewCheckBox()
    formList.SetTabGroups(types.NewTabGroups("id", "debug", "env", "language", "theme", "color_scheme",
        "asset_url", "title", "login_title", "session_life_time", "bootstrap_file_path", "go_mod_file_path", "no_limit_login_ip",
        "operation_log_off", "allow_del_operation_log", "hide_config_center_entrance", "hide_app_info_entrance", "hide_tool_entrance",
        "hide_plugin_entrance", "animation_type",
        "animation_duration", "animation_delay", "file_upload_engine", "extra").
        AddGroup("access_log_off", "access_assets_log_off", "info_log_off", "error_log_off", "sql_log", "logger_level",
            "info_log_path", "error_log_path",
            "access_log_path", "logger_rotate_max_size", "logger_rotate_max_backups",
            "logger_rotate_max_age", "logger_rotate_compress",
            "logger_encoder_encoding", "logger_encoder_time_key", "logger_encoder_level_key", "logger_encoder_name_key",
            "logger_encoder_caller_key", "logger_encoder_message_key", "logger_encoder_stacktrace_key", "logger_encoder_level",
            "logger_encoder_time", "logger_encoder_duration", "logger_encoder_caller").
        AddGroup("logo", "mini_logo", "custom_head_html", "custom_foot_html", "footer_info", "login_logo",
            "custom_404_html", "custom_403_html", "custom_500_html")).
        SetTabHeaders(lgWithConfigScore("general"), lgWithConfigScore("log"), lgWithConfigScore("custom"))

    formList.SetTable("goadmin_site").
        SetTitle(lgWithConfigScore("site setting")).
        SetDescription(lgWithConfigScore("site setting"))

    formList.SetUpdateFn(func(values form2.Values) error {
        ses := values.Get("session_life_time")
        sesInt, _ := strconv.Atoi(ses)
        if sesInt < 900 {
            return errors.New("wrong session life time, must bigger than 900 seconds")
        }
        if err := checkJSON(values, "file_upload_engine"); err != nil {
            return err
        }

        values["logo"][0] = escape(values.Get("logo"))
        values["mini_logo"][0] = escape(values.Get("mini_logo"))
        values["custom_head_html"][0] = escape(values.Get("custom_head_html"))
        values["custom_foot_html"][0] = escape(values.Get("custom_foot_html"))
        values["custom_404_html"][0] = escape(values.Get("custom_404_html"))
        values["custom_403_html"][0] = escape(values.Get("custom_403_html"))
        values["custom_500_html"][0] = escape(values.Get("custom_500_html"))
        values["footer_info"][0] = escape(values.Get("footer_info"))
        values["login_logo"][0] = escape(values.Get("login_logo"))

        var err error
        if s.c.UpdateProcessFn != nil {
            values, err = s.c.UpdateProcessFn(values)
            if err != nil {
                return err
            }
        }

        ui.GetService(services).RemoveOrShowSiteNavButton(values["hide_config_center_entrance"][0] == "true")
        ui.GetService(services).RemoveOrShowInfoNavButton(values["hide_app_info_entrance"][0] == "true")
        ui.GetService(services).RemoveOrShowToolNavButton(values["hide_tool_entrance"][0] == "true")
        ui.GetService(services).RemoveOrShowPlugNavButton(values["hide_plugin_entrance"][0] == "true")

        // TODO: add transaction
        err = models.Site().SetConn(s.conn).Update(values.RemoveSysRemark())
        if err != nil {
            return err
        }
        return s.c.Update(values.ToMap())
    })

    formList.EnableAjax(lgWithConfigScore("modify site config"),
        lgWithConfigScore("modify site config"),
        "",
        lgWithConfigScore("modify site config success"),
        lgWithConfigScore("modify site config fail"))

    return
}

func (s *SystemTable) GetGenerateForm(_ *context.Context) (generateTool Table) {
    generateTool = NewDefaultTable(DefaultConfigWithDriver(config.GetDatabases().GetDefault().Driver).
        SetOnlyNewForm())

    formList := generateTool.GetForm().AddXssJsFilter().
        SetHeadWidth(1).
        SetInputWidth(4).
        HideBackButton().
        HideContinueNewCheckBox()

    formList.AddField("Id", "id", db.Varchar, form.Default).FieldDefault("1").FieldHide()

    connNames := config.GetDatabases().Connections()
    ops := make(types.FieldOptions, len(connNames))
    for i, name := range connNames {
        ops[i] = types.FieldOption{Text: name, Value: name}
    }

    // General options
    // ================================

    formList.AddField(lgWithScore("connection", "tool"), "conn", db.Varchar, form.SelectSingle).
        FieldOptions(ops).
        FieldOnChooseAjax("table", "/tool/choose/conn",
            func(ctx *context.Context) (success bool, msg string, data interface{}) {
                connName := ctx.FormValue("value")
                if connName == "" {
                    return false, "wrong parameter", nil
                }
                cfg := s.c.Databases[connName]
                conn := db.GetConnectionFromService(services.Get(cfg.Driver))
                tables, err := db.WithDriverAndConnection(connName, conn).Table(cfg.Name).ShowTables()
                if err != nil {
                    return false, err.Error(), nil
                }
                ops := make(selection.Options, len(tables))
                for i, table := range tables {
                    ops[i] = selection.Option{Text: table, ID: table}
                }
                return true, "ok", ops
            })
    formList.AddField(lgWithScore("table", "tool"), "table", db.Varchar, form.SelectSingle).
        FieldOnChooseAjax("ajaxChooseTable", "/tool/choose/table",
            func(ctx *context.Context) (success bool, msg string, data interface{}) {

                var (
                    tableName       = ctx.FormValue("value")
                    connName        = ctx.FormValue("conn")
                    driver          = s.c.Databases[connName].Driver
                    conn            = db.GetConnectionFromService(services.Get(driver))
                    columnsModel, _ = db.WithDriverAndConnection(connName, conn).Table(tableName).ShowColumns()

                    fieldField = "Field"
                    typeField  = "Type"
                )

                if driver == "postgresql" {
                    fieldField = "column_name"
                    typeField = "udt_name"
                }
                if driver == "sqlite" {
                    fieldField = "name"
                    typeField = "type"
                }
                if driver == "mssql" {
                    fieldField = "column_name"
                    typeField = "data_type"
                }

                headName := make([]string, len(columnsModel))
                fieldName := make([]string, len(columnsModel))
                dbTypeList := make([]string, len(columnsModel))
                formTypeList := make([]string, len(columnsModel))

                for i, model := range columnsModel {
                    typeName := getType(model[typeField].(string))

                    // 默认使用注释
                    if comment, ok := model[`Comment`]; ok && comment != `` {
                        headName[i] = comment.(string)
                    } else {
                        headName[i] = strings.ToTitle(model[fieldField].(string))
                    }

                    fieldName[i] = model[fieldField].(string)
                    dbTypeList[i] = typeName
                    formTypeList[i] = form.GetFormTypeFromFieldType(db.DT(strings.ToUpper(typeName)),
                        model[fieldField].(string))
                }

                return true, "ok", [][]string{headName, fieldName, dbTypeList, formTypeList}
            }, template.HTML(utils.ParseText("choose_table_ajax", tmpls["choose_table_ajax"], nil)),
            `"conn":$('.conn').val(),`)
    formList.AddField(lgWithScore("package", "tool"), "package", db.Varchar, form.Text).FieldDefault("tables")
    formList.AddField(lgWithScore("primaryKey", "tool"), "pk", db.Varchar, form.Text).FieldDefault("id")

    formList.AddField(lgWithScore("table permission", "tool"), "permission", db.Varchar, form.Switch).
        FieldOptions(types.FieldOptions{
            {Text: lgWithScore("yes", "tool"), Value: "y"},
            {Text: lgWithScore("no", "tool"), Value: "n"},
        }).FieldDefault("n")

    formList.AddField(lgWithScore("extra import package", "tool"), "extra_import_package", db.Varchar, form.Select).
        FieldOptions(types.FieldOptions{
            {Text: "time", Value: "time"},
            {Text: "log", Value: "log"},
            {Text: "fmt", Value: "fmt"},
            {Text: "github.com/GoAdminGroup/go-admin/modules/db/dialect", Value: "github.com/GoAdminGroup/go-admin/modules/db/dialect"},
            {Text: "github.com/GoAdminGroup/go-admin/modules/db", Value: "github.com/GoAdminGroup/go-admin/modules/db"},
            {Text: "github.com/GoAdminGroup/go-admin/modules/language", Value: "github.com/GoAdminGroup/go-admin/modules/language"},
            {Text: "github.com/GoAdminGroup/go-admin/modules/logger", Value: "github.com/GoAdminGroup/go-admin/modules/logger"},
        }).
        FieldDefault("").
        FieldOptionExt(map[string]interface{}{
            "tags": true,
        })

    formList.AddField(lgWithScore("output", "tool"), "path", db.Varchar, form.Text).
        FieldDefault("").FieldMust().FieldHelpMsg(template.HTML(lgWithScore("use absolute path", "tool")))

    formList.AddField(lgWithScore("extra code", "tool"), "extra_code", db.Varchar, form.Code).
        FieldDefault("").FieldInputWidth(11)

    // Info table generate options
    // ================================

    formList.AddField(lgWithScore("title", "tool"), "table_title", db.Varchar, form.Text)
    formList.AddField(lgWithScore("description", "tool"), "table_description", db.Varchar, form.Text)

    formList.AddRow(func(panel *types.FormPanel) {
        addSwitchForTool(panel, "filter area", "hide_filter_area", "n", 2)
        panel.AddField(lgWithScore("filter form layout", "tool"), "filter_form_layout", db.Varchar, form.SelectSingle).
            FieldOptions(types.FieldOptions{
                {Text: form.LayoutDefault.String(), Value: form.LayoutDefault.String()},
                {Text: form.LayoutTwoCol.String(), Value: form.LayoutTwoCol.String()},
                {Text: form.LayoutThreeCol.String(), Value: form.LayoutThreeCol.String()},
                {Text: form.LayoutFourCol.String(), Value: form.LayoutFourCol.String()},
                {Text: form.LayoutFlow.String(), Value: form.LayoutFlow.String()},
            }).FieldDefault(form.LayoutDefault.String()).
            FieldRowWidth(4).FieldHeadWidth(3)
    })

    formList.AddRow(func(panel *types.FormPanel) {
        addSwitchForTool(panel, "new button", "hide_new_button", "n", 2)
        addSwitchForTool(panel, "export button", "hide_export_button", "n", 4, 3)
        addSwitchForTool(panel, "edit button", "hide_edit_button", "n", 4, 2)
    })

    formList.AddRow(func(panel *types.FormPanel) {
        addSwitchForTool(panel, "pagination", "hide_pagination", "n", 2)
        addSwitchForTool(panel, "delete button", "hide_delete_button", "n", 4, 3)
        addSwitchForTool(panel, "detail button", "hide_detail_button", "n", 4, 2)
    })

    formList.AddRow(func(panel *types.FormPanel) {
        addSwitchForTool(panel, "filter button", "hide_filter_button", "n", 2)
        addSwitchForTool(panel, "row selector", "hide_row_selector", "n", 4, 3)
        addSwitchForTool(panel, "query info", "hide_query_info", "n", 4, 2)
    })

    formList.AddTable(lgWithScore("field", "tool"), "fields", func(pa *types.FormPanel) {
        pa.AddField(lgWithScore("title", "tool"), "field_head", db.Varchar, form.Text).FieldHideLabel().
            FieldDisplay(func(value types.FieldModel) interface{} {
                return []string{""}
            })
        pa.AddField(lgWithScore("field name", "tool"), "field_name", db.Varchar, form.Text).FieldHideLabel().
            FieldDisplay(func(value types.FieldModel) interface{} {
                return []string{""}
            })
        pa.AddField(lgWithScore("field filterable", "tool"), "field_filterable", db.Varchar, form.CheckboxSingle).
            FieldOptions(types.FieldOptions{
                {Text: "", Value: "y"},
                {Text: "", Value: "n"},
            }).
            FieldDefault("n").
            FieldDisplay(func(value types.FieldModel) interface{} {
                return []string{"n"}
            })
        pa.AddField(lgWithScore("field sortable", "tool"), "field_sortable", db.Varchar, form.CheckboxSingle).
            FieldOptions(types.FieldOptions{
                {Text: "", Value: "y"},
                {Text: "", Value: "n"},
            }).
            FieldDefault("n").
            FieldDisplay(func(value types.FieldModel) interface{} {
                return []string{"n"}
            })
        pa.AddField(lgWithScore("field hide", "tool"), "field_hide", db.Varchar, form.CheckboxSingle).
            FieldOptions(types.FieldOptions{
                {Text: "", Value: "y"},
                {Text: "", Value: "n"},
            }).
            FieldDefault("n").
            FieldDisplay(func(value types.FieldModel) interface{} {
                return []string{"n"}
            })
        pa.AddField(lgWithScore("info field editable", "tool"), "info_field_editable", db.Varchar, form.CheckboxSingle).
            FieldOptions(types.FieldOptions{
                {Text: "", Value: "y"},
                {Text: "", Value: "n"},
            }).
            FieldDefault("n").
            FieldDisplay(func(value types.FieldModel) interface{} {
                return []string{"n"}
            })
        //pa.AddField(lgWithScore("db display type", "tool"), "field_display_type", db.Varchar, form.SelectSingle).
        //    FieldOptions(infoFieldDisplayTypeOptions()).
        //    FieldDisplay(func(value types.FieldModel) interface{} {
        //        return []string{""}
        //    })
        pa.AddField(lgWithScore("db type", "tool"), "field_db_type", db.Varchar, form.SelectSingle).
            FieldOptions(databaseTypeOptions()).
            FieldDisplay(func(value types.FieldModel) interface{} {
                return []string{"Int"}
            })
    }).FieldInputWidth(11)

    // Form generate options
    // ================================

    formList.AddField(lgWithScore("title", "tool"), "form_title", db.Varchar, form.Text)
    formList.AddField(lgWithScore("description", "tool"), "form_description", db.Varchar, form.Text)

    formList.AddRow(func(panel *types.FormPanel) {
        addSwitchForTool(panel, "continue edit checkbox", "hide_continue_edit_check_box", "n", 2)
        addSwitchForTool(panel, "reset button", "hide_reset_button", "n", 5, 3)
    })

    formList.AddRow(func(panel *types.FormPanel) {
        addSwitchForTool(panel, "continue new checkbox", "hide_continue_new_check_box", "n", 2)
        addSwitchForTool(panel, "back button", "hide_back_button", "n", 5, 3)
    })

    formList.AddTable(lgWithScore("field", "tool"), "fields_form", func(pa *types.FormPanel) {
        pa.AddField(lgWithScore("title", "tool"), "field_head_form", db.Varchar, form.Text).FieldHideLabel().
            FieldDisplay(func(value types.FieldModel) interface{} {
                return []string{""}
            })
        pa.AddField(lgWithScore("field name", "tool"), "field_name_form", db.Varchar, form.Text).FieldHideLabel().
            FieldDisplay(func(value types.FieldModel) interface{} {
                return []string{""}
            })
        pa.AddField(lgWithScore("field editable", "tool"), "field_canedit", db.Varchar, form.CheckboxSingle).
            FieldOptions(types.FieldOptions{
                {Text: "", Value: "y"},
                {Text: "", Value: "n"},
            }).
            FieldDefault("y").
            FieldDisplay(func(value types.FieldModel) interface{} {
                return []string{"y"}
            })
        pa.AddField(lgWithScore("field can add", "tool"), "field_canadd", db.Varchar, form.CheckboxSingle).
            FieldOptions(types.FieldOptions{
                {Text: "", Value: "y"},
                {Text: "", Value: "n"},
            }).
            FieldDefault("y").
            FieldDisplay(func(value types.FieldModel) interface{} {
                return []string{"y"}
            })
        pa.AddField(lgWithScore("field default", "tool"), "field_default", db.Varchar, form.Text).FieldHideLabel().
            FieldDisplay(func(value types.FieldModel) interface{} {
                return []string{""}
            })
        pa.AddField(lgWithScore("field display", "tool"), "field_display", db.Varchar, form.SelectSingle).
            FieldOptions(types.FieldOptions{
                {Text: lgWithScore("field display normal", "tool"), Value: "0"},
                {Text: lgWithScore("field display hide", "tool"), Value: "1"},
                {Text: lgWithScore("field display edit hide", "tool"), Value: "2"},
                {Text: lgWithScore("field display create hide", "tool"), Value: "3"},
            }).
            FieldDisplay(func(value types.FieldModel) interface{} {
                return []string{"0"}
            })
        pa.AddField(lgWithScore("db type", "tool"), "field_db_type_form", db.Varchar, form.SelectSingle).
            FieldOptions(databaseTypeOptions()).
            FieldDisplay(func(value types.FieldModel) interface{} {
                return []string{"Int"}
            })
        pa.AddField(lgWithScore("form type", "tool"), "field_form_type_form", db.Varchar, form.SelectSingle).
            FieldOptions(formTypeOptions()).FieldDisplay(func(value types.FieldModel) interface{} {
            return []string{"Text"}
        })
    }).FieldInputWidth(11)

    // Detail page generate options
    // ================================

    formList.AddField(lgWithScore("title", "tool"), "detail_title", db.Varchar, form.Text)
    formList.AddField(lgWithScore("description", "tool"), "detail_description", db.Varchar, form.Text)

    formList.AddField(lgWithScore("detail display", "tool"), "detail_display", db.Varchar, form.SelectSingle).
        FieldOptions(types.FieldOptions{
            {Text: lgWithScore("follow list page", "tool"), Value: "0"},
            {Text: lgWithScore("inherit from list page", "tool"), Value: "1"},
            {Text: lgWithScore("independent from list page", "tool"), Value: "2"},
        }).
        FieldDefault("0").
        FieldOnChooseHide("0", "detail_title", "detail_description", "fields_detail")

    formList.AddTable(lgWithScore("field", "tool"), "fields_detail", func(pa *types.FormPanel) {
        pa.AddField(lgWithScore("title", "tool"), "detail_field_head", db.Varchar, form.Text).FieldHideLabel().
            FieldDisplay(func(value types.FieldModel) interface{} {
                return []string{""}
            })
        pa.AddField(lgWithScore("field name", "tool"), "detail_field_name", db.Varchar, form.Text).FieldHideLabel().
            FieldDisplay(func(value types.FieldModel) interface{} {
                return []string{""}
            })
        pa.AddField(lgWithScore("db type", "tool"), "detail_field_db_type", db.Varchar, form.SelectSingle).
            FieldOptions(databaseTypeOptions()).
            FieldDisplay(func(value types.FieldModel) interface{} {
                return []string{"Int"}
            })
    }).FieldInputWidth(11)

    formList.SetTabGroups(types.
        NewTabGroups("conn", "table", "package", "pk", "permission", "extra_import_package", "path", "extra_code").
        AddGroup("table_title", "table_description", "hide_filter_area", "filter_form_layout",
            "hide_new_button", "hide_export_button", "hide_edit_button",
            "hide_pagination", "hide_delete_button", "hide_detail_button",
            "hide_filter_button", "hide_row_selector", "hide_query_info",
            "fields").
        AddGroup("form_title", "form_description", "hide_continue_edit_check_box", "hide_reset_button",
            "hide_continue_new_check_box", "hide_back_button",
            "fields_form").
        AddGroup("detail_display", "detail_title", "detail_description", "fields_detail")).
        SetTabHeaders(lgWithScore("basic info", "tool"), lgWithScore("table info", "tool"),
            lgWithScore("form info", "tool"), lgWithScore("detail info", "tool"))

    formList.SetTable("tool").
        SetTitle(lgWithScore("tool", "tool")).
        SetDescription(lgWithScore("tool", "tool")).
        SetHeader(template.HTML(`<h3 class="box-title">` +
            lgWithScore("generate table model", "tool") + `</h3>`))

    formList.SetInsertFn(func(values form2.Values) error {

        table := values.Get("table")

        if table == "" {
            return errors.New("table is empty")
        }

        if values.Get("permission") == "y" {
            tools.InsertPermissionOfTable(s.conn, table)
        }

        output := values.Get("path")

        if output == "" {
            return errors.New("output path is empty")
        }

        connName := values.Get("conn")

        fields := make(tools.Fields, len(values["field_head"]))

        for i := 0; i < len(values["field_head"]); i++ {
            fields[i] = tools.Field{
                Head:         values["field_head"][i],
                Name:         values["field_name"][i],
                DBType:       values["field_db_type"][i],
                Filterable:   values["field_filterable"][i] == "y",
                Sortable:     values["field_sortable"][i] == "y",
                Hide:         values["field_hide"][i] == "y",
                InfoEditable: values["info_field_editable"][i] == "y",
            }
        }

        extraImport := ""
        for _, pack := range values["extra_import_package[]"] {
            if extraImport != "" {
                extraImport += `
`
            }
            extraImport += `    "` + pack + `"`
        }

        formFields := make(tools.Fields, len(values["field_head_form"]))

        for i := 0; i < len(values["field_head_form"]); i++ {
            extraFun := ""
            if values["field_name_form"][i] == `createAt` {
                extraFun += `.FieldNowWhenInsert()`
            } else if values["field_name_form"][i] == `updateAt` {
                extraFun += `.FieldNowWhenUpdate()`
            } else if values["field_default"][i] != "" && !strings.Contains(values["field_default"][i], `"`) {
                values["field_default"][i] = `"` + values["field_default"][i] + `"`
            }
            formFields[i] = tools.Field{
                Head:       values["field_head_form"][i],
                Name:       values["field_name_form"][i],
                Default:    values["field_default"][i],
                FormType:   values["field_form_type_form"][i],
                DBType:     values["field_db_type_form"][i],
                CanAdd:     values["field_canadd"][i] == "y",
                Editable:   values["field_canedit"][i] == "y",
                FormHide:   values["field_display"][i] == "1",
                CreateHide: values["field_display"][i] == "2",
                EditHide:   values["field_display"][i] == "3",
                ExtraFun:   extraFun,
            }
        }

        detailFields := make(tools.Fields, len(values["detail_field_head"]))

        for i := 0; i < len(values["detail_field_head"]); i++ {
            detailFields[i] = tools.Field{
                Head:   values["detail_field_head"][i],
                Name:   values["detail_field_name"][i],
                DBType: values["detail_field_db_type"][i],
            }
        }

        detailDisplay, _ := strconv.ParseUint(values.Get("detail_display"), 10, 64)

        err := tools.Generate(tools.NewParamWithFields(tools.Config{
            Connection:               connName,
            Driver:                   s.c.Databases[connName].Driver,
            Package:                  values.Get("package"),
            Table:                    table,
            HideFilterArea:           values.Get("hide_filter_area") == "y",
            HideNewButton:            values.Get("hide_new_button") == "y",
            HideExportButton:         values.Get("hide_export_button") == "y",
            HideEditButton:           values.Get("hide_edit_button") == "y",
            HideDeleteButton:         values.Get("hide_delete_button") == "y",
            HideDetailButton:         values.Get("hide_detail_button") == "y",
            HideFilterButton:         values.Get("hide_filter_button") == "y",
            HideRowSelector:          values.Get("hide_row_selector") == "y",
            HidePagination:           values.Get("hide_pagination") == "y",
            HideQueryInfo:            values.Get("hide_query_info") == "y",
            HideContinueEditCheckBox: values.Get("hide_continue_edit_check_box") == "y",
            HideContinueNewCheckBox:  values.Get("hide_continue_new_check_box") == "y",
            HideResetButton:          values.Get("hide_reset_button") == "y",
            HideBackButton:           values.Get("hide_back_button") == "y",
            FilterFormLayout:         form.GetLayoutFromString(values.Get("filter_form_layout")),
            Schema:                   values.Get("schema"),
            Output:                   output,
            DetailDisplay:            uint8(detailDisplay),
            FormTitle:                values.Get("form_title"),
            FormDescription:          values.Get("form_description"),
            DetailTitle:              values.Get("detail_title"),
            DetailDescription:        values.Get("detail_description"),
            TableTitle:               values.Get("table_title"),
            TableDescription:         values.Get("table_description"),
            ExtraImport:              extraImport,
            ExtraCode:                escape(values.Get("extra_code")),
        }, fields, formFields, detailFields))

        if err != nil {
            return err
        }

        return tools.GenerateTables(output, values.Get("package"), []string{table}, false)
    })

    formList.EnableAjaxData(types.AjaxData{
        SuccessTitle: lgWithScore("generate table model", "tool"),
        ErrorTitle:   lgWithScore("generate table model", "tool"),
        SuccessText:  lgWithScore("generate table model success", "tool"),
        ErrorText:    lgWithScore("generate table model fail", "tool"),
        DisableJump:  true,
    })

    formList.SetFooterHtml(utils.ParseHTML("generator", tmpls["generator"], map[string]string{
        "prefix": "go_admin_" + config.GetAppID() + "_generator_",
    }))

    formList.SetFormNewBtnWord(template.HTML(lgWithScore("generate", "tool")))
    formList.SetWrapper(func(content tmpl.HTML) tmpl.HTML {
        headLi := html.LiEl().SetClass("list-group-item", "list-head").
            SetContent(template.HTML(lgWithScore("generated tables", "tool"))).Get()
        return html.UlEl().SetClass("save_table_list", "list-group").SetContent(
            headLi).Get() + content
    })

    formList.SetHideSideBar()

    return generateTool
}

// -------------------------
// helper functions
// -------------------------

func encodePassword(pwd []byte) string {
    hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.DefaultCost)
    if err != nil {
        return ""
    }
    return string(hash)
}

func label() types.LabelAttribute {
    return template.Get(config.GetTheme()).Label().SetType("success")
}

func lg(v string) string {
    return language.Get(v)
}

func defaultFilterFn(val string, def ...string) types.FieldFilterFn {
    return func(value types.FieldModel) interface{} {
        if len(def) > 0 {
            if value.Value == def[0] {
                return val
            }
        } else {
            if value.Value == "" {
                return val
            }
        }
        return value.Value
    }
}

func lgWithScore(v string, score ...string) string {
    return language.GetWithScope(v, score...)
}

func lgWithConfigScore(v string, score ...string) string {
    scores := append([]string{"config"}, score...)
    return language.GetWithScope(v, scores...)
}

func link(url, content string) tmpl.HTML {
    return html.AEl().
        SetAttr("href", url).
        SetContent(template.HTML(lg(content))).
        Get()
}

func escape(s string) string {
    if s == "" {
        return ""
    }
    s, err := url.QueryUnescape(s)
    if err != nil {
        logger.Error("escape error", err)
    }
    return s
}

func checkJSON(values form2.Values, key string) error {
    v := values.Get(key)
    if v != "" && !utils.IsJSON(v) {
        return errors.New("wrong " + key)
    }
    return nil
}

func (s *SystemTable) table(table string) *db.SQL {
    return s.connection().Table(table)
}

func (s *SystemTable) connection() *db.SQL {
    return db.WithDriver(s.conn)
}

func interfaces(arr []string) []interface{} {
    var iArr = make([]interface{}, len(arr))

    for key, v := range arr {
        iArr[key] = v
    }

    return iArr
}

func addSwitchForTool(formList *types.FormPanel, head, field, def string, row ...int) {
    formList.AddField(lgWithScore(head, "tool"), field, db.Varchar, form.Switch).
        FieldOptions(types.FieldOptions{
            {Text: lgWithScore("show", "tool"), Value: "n"},
            {Text: lgWithScore("hide", "tool"), Value: "y"},
        }).FieldDefault(def)
    if len(row) > 0 {
        formList.FieldRowWidth(row[0])
    }
    if len(row) > 1 {
        formList.FieldHeadWidth(row[1])
    }
    if len(row) > 2 {
        formList.FieldInputWidth(row[2])
    }
}

func formTypeOptions() types.FieldOptions {
    opts := make(types.FieldOptions, len(form.AllType))
    for i := 0; i < len(form.AllType); i++ {
        v := form.AllType[i].Name()
        opts[i] = types.FieldOption{Text: v, Value: v}
    }
    return opts
}

func databaseTypeOptions() types.FieldOptions {
    opts := make(types.FieldOptions, len(db.IntTypeList)+
        len(db.StringTypeList)+
        len(db.FloatTypeList)+
        len(db.UintTypeList)+
        len(db.BoolTypeList))
    z := 0
    for _, t := range db.IntTypeList {
        text := string(t)
        v := strings.ToTitle(text)
        opts[z] = types.FieldOption{Text: text, Value: v}
        z++
    }
    for _, t := range db.StringTypeList {
        text := string(t)
        v := strings.ToTitle(text)
        opts[z] = types.FieldOption{Text: text, Value: v}
        z++
    }
    for _, t := range db.FloatTypeList {
        text := string(t)
        v := strings.ToTitle(text)
        opts[z] = types.FieldOption{Text: text, Value: v}
        z++
    }
    for _, t := range db.UintTypeList {
        text := string(t)
        v := strings.ToTitle(text)
        opts[z] = types.FieldOption{Text: text, Value: v}
        z++
    }
    for _, t := range db.BoolTypeList {
        text := string(t)
        v := strings.ToTitle(strings.ToLower(text))
        opts[z] = types.FieldOption{Text: text, Value: v}
        z++
    }
    return opts
}

func getType(typeName string) string {
    r, _ := regexp.Compile(`\(.*?\)`)
    typeName = r.ReplaceAllString(typeName, "")
    r2, _ := regexp.Compile(`unsigned(.*)`)
    return strings.TrimSpace(strings.ToTitle(strings.ToLower(r2.ReplaceAllString(typeName, ""))))
}
