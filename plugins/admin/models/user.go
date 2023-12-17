package models

import (
    "database/sql"
    "fmt"
    "net/url"
    "regexp"
    "strconv"
    "strings"
    "time"

    "github.com/GoAdminGroup/go-admin/modules/config"
    "github.com/GoAdminGroup/go-admin/modules/db"
    "github.com/GoAdminGroup/go-admin/modules/db/dialect"
    //"github.com/GoAdminGroup/go-admin/modules/logger"
    //"github.com/GoAdminGroup/go-admin/modules/utils"
    //"github.com/GoAdminGroup/go-admin/plugins/admin/modules/constant"
)

// UserModel is user model structure.
type UserModel struct {
    Base `json:"-"`

    ID            int64       `json:"id"`
    Name          string      `json:"name"`
    UserName      string      `json:"username"`
    Password      string      `json:"password"`
    Avatar        string      `json:"avatar"`
    RememberToken string      `json:"rememberToken"`
    MenuIds       []int64     `json:"menuIds"`
    Roles         []RoleModel `json:"role"`
    Level         string      `json:"level"`
    LevelName     string      `json:"levelName"`
    CreateAt      int64       `json:"createAt"`
    UpdateAt      int64       `json:"updateAt"`

    cacheReplacer *strings.Replacer
}

// User return a default user model.
func User() UserModel {
    return UserModel{Base: Base{TableName: config.GetAuthUserTable()}}
}

// UserWithId return a default user model of given id.
func UserWithId(id string) UserModel {
    idInt, _ := strconv.Atoi(id)
    return UserModel{Base: Base{TableName: config.GetAuthUserTable()}, ID: int64(idInt)}
}

func (t UserModel) SetConn(con db.Connection) UserModel {
    t.Conn = con
    return t
}

func (t UserModel) WithTx(tx *sql.Tx) UserModel {
    t.Tx = tx
    return t
}

// Find return a default user model of given id.
func (t UserModel) Find(id interface{}) UserModel {
    item, _ := t.Table(t.TableName).Find(id)
    return t.MapToModel(item)
}

// FindByUserName return a default user model of given name.
func (t UserModel) FindByUserName(username interface{}) UserModel {
    item, _ := t.Table(t.TableName).Where("username", "=", username).First()
    return t.MapToModel(item)
}

// IsEmpty check the user model is empty or not.
func (t UserModel) IsEmpty() bool {
    return t.ID == int64(0)
}

// HasMenu check the user has visitable menu or not.
func (t UserModel) HasMenu() bool {
    return len(t.MenuIds) != 0 || t.IsSuperAdmin()
}

// IsSuperAdmin check the user model is super admin or not.
func (t UserModel) IsSuperAdmin() bool {
    /*for _, per := range t.Permissions {
        if len(per.HttpPath) > 0 && per.HttpPath[0] == "*" && per.HttpMethod[0] == "" {
            return true
        }
    }*/
    if t.ID == 1 {
        return true
    }
    return false
}

func (t UserModel) GetCheckPermissionByUrlMethod(path, method string) string {
    if !t.CheckPermissionByUrlMethod(path, method, url.Values{}) {
        return ""
    }
    return path
}

func (t UserModel) IsVisitor() bool {
    return !t.CheckPermissionByUrlMethod(config.Url("/info/normal_manager"), "GET", url.Values{})
}

func (t UserModel) HideUserCenterEntrance() bool {
    return t.IsVisitor() && config.GetHideVisitorUserCenterEntrance()
}

func (t UserModel) Template(str string) string {
    if t.cacheReplacer == nil {
        t.cacheReplacer = strings.NewReplacer("{{.AuthId}}", strconv.Itoa(int(t.ID)),
            "{{.AuthName}}", t.Name, "{{.AuthUserName}}", t.UserName)
    }
    return t.cacheReplacer.Replace(str)
}

func (t UserModel) CheckPermissionByUrlMethod(path, _ string, _ url.Values) bool {
    path, _ = url.PathUnescape(path)
    if t.IsSuperAdmin() {
        return true
    }

    if path == "" {
        return false
    }

    logoutCheck, _ := regexp.Compile(config.Url("/logout") + "(.*?)")
    if logoutCheck.MatchString(path) {
        return true
    }

    // 当前改为判断是否有菜单即可
    if path != "/" && path[len(path)-1] == '/' {
        path = path[:len(path)-1]
    }

    id := getMenuId(t.Conn, path)
    if id == 0 {
        return true
    }

    for _, v := range t.MenuIds {
        if v == id {
            return true
        }
    }

    return false

    /*
       path = utils.ReplaceAll(path, constant.EditPKKey, "id", constant.DetailPKKey, "id")
       path, params := getParam(path)
       for key, value := range formParams {
           if len(value) > 0 {
               params.Add(key, value[0])
           }
       }

       for _, v := range t.Permissions {
           if v.HttpMethod[0] == "" || inMethodArr(v.HttpMethod, method) {

               if v.HttpPath[0] == "*" {
                   return true
               }

               for i := 0; i < len(v.HttpPath); i++ {

                   matchPath := config.Url(t.Template(strings.TrimSpace(v.HttpPath[i])))
                   matchPath, matchParam := getParam(matchPath)

                   if matchPath == path {
                       if t.checkParam(params, matchParam) {
                           return true
                       }
                   }

                   reg, err := regexp.Compile(matchPath)

                   if err != nil {
                       logger.Error("CheckPermissions error: ", err)
                       continue
                   }

                   if reg.FindString(path) == path {
                       if t.checkParam(params, matchParam) {
                           return true
                       }
                   }
               }
           }
       }

       return false
    */
}

// Read Only After Init
var uriToMenuURI = map[string]string{}

// InitURI2MenuURIData 生成uri与菜单填入的uri对应关系,用于权限验证时找到数据库菜单判断是否有权限
func InitURI2MenuURIData(k string) {
    var key string
    f := config.GetURLFormats()

    val := strings.ReplaceAll(f.Info, `:__prefix`, k)
    uriToMenuURI[val] = val

    key = strings.ReplaceAll(f.Detail, `:__prefix`, k)
    uriToMenuURI[key] = val

    key = strings.ReplaceAll(f.Create, `:__prefix`, k)
    uriToMenuURI[key] = val

    key = strings.ReplaceAll(f.Delete, `:__prefix`, k)
    uriToMenuURI[key] = val

    key = strings.ReplaceAll(f.Export, `:__prefix`, k)
    uriToMenuURI[key] = val

    key = strings.ReplaceAll(f.Edit, `:__prefix`, k)
    uriToMenuURI[key] = val

    key = strings.ReplaceAll(f.ShowEdit, `:__prefix`, k)
    uriToMenuURI[key] = val

    key = strings.ReplaceAll(f.ShowCreate, `:__prefix`, k)
    uriToMenuURI[key] = val

    key = strings.ReplaceAll(f.Update, `:__prefix`, k)
    uriToMenuURI[key] = val

    /*fmt.Println(`-----------InitURI2MenuURIData-----------`, k, len(uriToMenuURI))
      for key, val = range uriToMenuURI {
          fmt.Printf(`%s => %s%s`, key, val, "\n")
      }
      fmt.Println(`---------------------------------------------------------`)*/
}

// getMenuId 获取菜单id,不存在时返回0
func getMenuId(db db.Connection, uri string) (ret int64) {
    prefix := config.Url(``)
    uri, _, _ = strings.Cut(uri, `?`)
    uri = strings.TrimPrefix(uri, prefix)
    uri = strings.TrimLeft(uri, `/`)
    uri = `/` + uri
    menuURI := uriToMenuURI[uri]
    if menuURI == `` {
        return
    }

    menuInfo, err := db.Query(`SELECT id FROM goadmin_menu WHERE uri=? LIMIT 1`, menuURI)
    if err == nil && len(menuInfo) > 0 {
        id, _ := strconv.Atoi(fmt.Sprintf(`%d`, menuInfo[0][`id`]))
        ret = int64(id)
    }

    //fmt.Println(`------------getMenuId---------------`, uri, menuURI, ret)

    return
}

/*func getParam(u string) (string, url.Values) {
    m := make(url.Values)
    urr := strings.Split(u, "?")
    if len(urr) > 1 {
        m, _ = url.ParseQuery(urr[1])
    }
    return urr[0], m
}*/

func (t UserModel) checkParam(src, comp url.Values) bool {
    if len(comp) == 0 {
        return true
    }
    if len(src) == 0 {
        return false
    }
    for key, value := range comp {
        v, find := src[key]
        if !find {
            return false
        }
        if len(value) == 0 {
            continue
        }
        if len(v) == 0 {
            return false
        }
        for i := 0; i < len(v); i++ {
            if v[i] == t.Template(value[i]) {
                continue
            } else {
                return false
            }
        }
    }
    return true
}

/*func inMethodArr(arr []string, str string) bool {
    for i := 0; i < len(arr); i++ {
        if strings.EqualFold(arr[i], str) {
            return true
        }
    }
    return false
}*/

// ReleaseConn update the avatar of user.
func (t UserModel) ReleaseConn() UserModel {
    t.Conn = nil
    return t
}

// UpdateAvatar update the avatar of user.
func (t UserModel) UpdateAvatar(avatar string) {
    t.Avatar = avatar
}

// WithRoles query the role info of the user.
func (t UserModel) WithRoles() UserModel {
    roleModel, _ := t.Table("goadmin_role_user").
        LeftJoin("goadmin_role", "goadmin_role.id", "=", "goadmin_role_user.roleId").
        Where("userId", "=", t.ID).
        Select("goadmin_role.id", "goadmin_role.name", "goadmin_role.slug",
            "goadmin_role.createAt", "goadmin_role.updateAt").
        All()

    for _, role := range roleModel {
        t.Roles = append(t.Roles, Role().MapToModel(role))
    }

    if len(t.Roles) > 0 {
        t.Level = t.Roles[0].Slug
        t.LevelName = t.Roles[0].Name
    }

    return t
}

func (t UserModel) GetAllRoleId() []interface{} {

    var ids = make([]interface{}, len(t.Roles))

    for key, role := range t.Roles {
        ids[key] = role.ID
    }

    return ids
}

// WithPermissions query the permission info of the user.
func (t UserModel) WithPermissions() UserModel {
    // 当前更改了用户权限
    /*var permissions = make([]map[string]interface{}, 0)


      roleIds := t.GetAllRoleId()

      if len(roleIds) > 0 {
          permissions, _ = t.Table("goadmin_role_permissions").
              LeftJoin("goadmin_permissions", "goadmin_permissions.id", "=", "goadmin_role_permissions.permission_id").
              WhereIn("roleId", roleIds).
              Select("goadmin_permissions.http_method", "goadmin_permissions.http_path",
                  "goadmin_permissions.id", "goadmin_permissions.name", "goadmin_permissions.slug",
                  "goadmin_permissions.createAt", "goadmin_permissions.updateAt").
              All()
      }


      userPermissions, _ := t.Table("goadmin_user_permissions").
          LeftJoin("goadmin_permissions", "goadmin_permissions.id", "=", "goadmin_user_permissions.permission_id").
          Where("userId", "=", t.ID).
          Select("goadmin_permissions.http_method", "goadmin_permissions.http_path",
              "goadmin_permissions.id", "goadmin_permissions.name", "goadmin_permissions.slug",
              "goadmin_permissions.createAt", "goadmin_permissions.updateAt").
          All()

      permissions = append(permissions, userPermissions...)

      for i := 0; i < len(permissions); i++ {
          exist := false
          for j := 0; j < len(t.Permissions); j++ {
              if t.Permissions[j].ID == permissions[i]["id"] {
                  exist = true
                  break
              }
          }
          if exist {
              continue
          }
          t.Permissions = append(t.Permissions, Permission().MapToModel(permissions[i]))
      }*/

    return t
}

// WithMenus query the menu info of the user.
func (t UserModel) WithMenus() UserModel {

    var menuIdsModel []map[string]interface{}
    if t.IsSuperAdmin() {
        menuIdsModel, _ = t.Conn.Query(`SELECT id FROM goadmin_menu`)
    } else {
        menuIdsModel, _ = t.Conn.Query(`SELECT id,parentId AS pid FROM goadmin_menu WHERE id IN(SELECT menuId FROM goadmin_role_menu WHERE roleId IN(SELECT roleId FROM goadmin_role_user WHERE userId=?))`, t.ID)
    }

    mapIds := make(map[int64]struct{}, len(menuIdsModel)*2)
    for _, mid := range menuIdsModel {
        mid64, _ := strconv.Atoi(fmt.Sprintf(`%d`, mid[`id`]))
        mapIds[int64(mid64)] = struct{}{}

        // 子菜单的父菜单权限也拥有
        if pid, ok := mid["pid"]; ok {
            pid64, _ := strconv.Atoi(fmt.Sprintf(`%d`, pid))
            mapIds[int64(pid64)] = struct{}{}
        }
    }

    for id := range mapIds {
        t.MenuIds = append(t.MenuIds, id)
    }

    return t
}

// New create a user model.
func (t UserModel) New(username, password, name, avatar string) (UserModel, error) {

    id, err := t.WithTx(t.Tx).Table(t.TableName).Insert(dialect.H{
        "username": username,
        "password": password,
        "name":     name,
        "avatar":   avatar,
        `createAt`: uint32(time.Now().Unix()),
        `updateAt`: uint32(time.Now().Unix()),
    })

    t.ID = id
    t.UserName = username
    t.Password = password
    t.Avatar = avatar
    t.Name = name

    return t, err
}

// Update update the user model.
func (t UserModel) Update(username, password, name, avatar string, isUpdateAvatar bool) (int64, error) {

    fieldValues := dialect.H{
        "username": username,
        "name":     name,
        "updateat": uint32(time.Now().Unix()),
    }

    if avatar == "" || isUpdateAvatar {
        fieldValues["avatar"] = avatar
    }

    if password != "" {
        fieldValues["password"] = password
    }

    return t.WithTx(t.Tx).Table(t.TableName).
        Where("id", "=", t.ID).
        Update(fieldValues)
}

// UpdatePwd update the password of the user model.
func (t UserModel) UpdatePwd(password string) UserModel {

    _, _ = t.Table(t.TableName).
        Where("id", "=", t.ID).
        Update(dialect.H{
            "password": password,
        })

    t.Password = password
    return t
}

// CheckRoleId check the role of the user model.
func (t UserModel) CheckRoleId(roleId string) bool {
    checkRole, _ := t.Table("goadmin_role_user").
        Where("roleId", "=", roleId).
        Where("userId", "=", t.ID).
        First()
    return checkRole != nil
}

// DeleteRoles delete all the roles of the user model.
func (t UserModel) DeleteRoles() error {
    return t.Table("goadmin_role_user").
        Where("userId", "=", t.ID).
        Delete()
}

// AddRole add a role of the user model.
func (t UserModel) AddRole(roleId string) (int64, error) {
    if roleId != "" {
        if !t.CheckRoleId(roleId) {
            return t.WithTx(t.Tx).Table("goadmin_role_user").
                Insert(dialect.H{
                    "roleId": roleId,
                    "userId": t.ID,
                })
        }
    }
    return 0, nil
}

// CheckRole check the role of the user.
func (t UserModel) CheckRole(slug string) bool {
    for _, role := range t.Roles {
        if role.Slug == slug {
            return true
        }
    }

    return false
}

// CheckPermissionById check the permission of the user.
func (t UserModel) CheckPermissionById(permissionId string) bool {
    checkPermission, _ := t.Table("goadmin_user_permissions").
        Where("permission_id", "=", permissionId).
        Where("userId", "=", t.ID).
        First()
    return checkPermission != nil
}

// CheckPermission check the permission of the user.
func (t UserModel) CheckPermission(_ string) bool {
    return true
}

// DeletePermissions delete all the permissions of the user model.
func (t UserModel) DeletePermissions() error {
    return t.WithTx(t.Tx).Table("goadmin_role_user").
        Where("userId", "=", t.ID).
        Delete()
}

// AddPermission add a permission of the user model.
func (t UserModel) AddPermission(_ string) (int64, error) {
    return 0, nil
}

// MapToModel get the user model from given map.
func (t UserModel) MapToModel(m map[string]interface{}) UserModel {
    t.ID, _ = m["id"].(int64)
    t.Name, _ = m["name"].(string)
    t.UserName, _ = m["username"].(string)
    t.Password, _ = m["password"].(string)
    t.Avatar, _ = m["avatar"].(string)
    t.RememberToken, _ = m["rememberToken"].(string)
    t.CreateAt, _ = m["createAt"].(int64)
    t.UpdateAt, _ = m["updateAt"].(int64)

    return t
}
