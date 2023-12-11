package models

import (
    "database/sql"
    "strconv"
    "time"

    "github.com/GoAdminGroup/go-admin/modules/db"
    "github.com/GoAdminGroup/go-admin/modules/db/dialect"
)

// RoleModel is role model structure.
type RoleModel struct {
    Base

    ID       int64
    Name     string
    Slug     string
    CreateAt int64
    UpdateAt int64
}

// Role return a default role model.
func Role() RoleModel {
    return RoleModel{Base: Base{TableName: "goadmin_role"}}
}

// RoleWithId return a default role model of given id.
func RoleWithId(id string) RoleModel {
    idInt, _ := strconv.Atoi(id)
    return RoleModel{Base: Base{TableName: "goadmin_role"}, ID: int64(idInt)}
}

func (t RoleModel) SetConn(con db.Connection) RoleModel {
    t.Conn = con
    return t
}

func (t RoleModel) WithTx(tx *sql.Tx) RoleModel {
    t.Tx = tx
    return t
}

// Find return a default role model of given id.
func (t RoleModel) Find(id interface{}) RoleModel {
    item, _ := t.Table(t.TableName).Find(id)
    return t.MapToModel(item)
}

// IsNameExist check the row exist with given name and id.
func (t RoleModel) IsNameExist(name string, id string) bool {
    if id == "" {
        check, _ := t.Table(t.TableName).Where("name", "=", name).First()
        return check != nil
    }
    check, _ := t.Table(t.TableName).
        Where("name", "=", name).
        Where("id", "!=", id).
        First()
    return check != nil
}

// New create a role model.
func (t RoleModel) New(name, slug string) (RoleModel, error) {

    id, err := t.WithTx(t.Tx).Table(t.TableName).Insert(dialect.H{
        "name":     name,
        "slug":     slug,
        "createat": time.Now().Unix(),
    })

    t.ID = id
    t.Name = name
    t.Slug = slug

    return t, err
}

// Update update the role model.
func (t RoleModel) Update(name, slug string) (int64, error) {

    return t.WithTx(t.Tx).Table(t.TableName).
        Where("id", "=", t.ID).
        Update(dialect.H{
            "name":     name,
            "slug":     slug,
            "updateat": uint32(time.Now().Unix()),
        })
}

// CheckPermission check the permission of role.
func (t RoleModel) CheckPermission(menuId string) bool {
    checkPermission, _ := t.Table("goadmin_role_menu").
        Where("menuId", "=", menuId).
        Where("roleId", "=", t.ID).
        First()
    return checkPermission != nil
}

// DeletePermissions delete all the permissions of role.
func (t RoleModel) DeletePermissions() error {
    return t.WithTx(t.Tx).Table("goadmin_role_menu").
        Where("roleId", "=", t.ID).
        Delete()
}

// AddPermission add the permissions to the role.
func (t RoleModel) AddPermission(menuId string) (int64, error) {

    if menuId != "" && !t.CheckPermission(menuId) {
        return t.WithTx(t.Tx).Table("goadmin_role_menu").
            Insert(dialect.H{
                "menuId":   menuId,
                "roleId":   t.ID,
                "createat": uint32(time.Now().Unix()),
            })
    }

    return 0, nil
}

// MapToModel get the role model from given map.
func (t RoleModel) MapToModel(m map[string]interface{}) RoleModel {
    t.ID = m["id"].(int64)
    t.Name, _ = m["name"].(string)
    t.Slug, _ = m["slug"].(string)
    t.CreateAt, _ = m["createAt"].(int64)
    t.UpdateAt, _ = m["updateAt"].(int64)
    return t
}
