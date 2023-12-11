// Package models model
package models

import (
    "encoding/json"
    "strconv"
    "time"

    "github.com/GoAdminGroup/go-admin/modules/db"
    "github.com/GoAdminGroup/go-admin/modules/db/dialect"
)

// MenuModel is menu model structure.
type MenuModel struct {
    Base

    ID       int64
    Title    string
    ParentID int64
    Icon     string
    URI      string
    Header   string
    CreateAt int64
    UpdateAt int64
}

// Menu return a default menu model.
func Menu() MenuModel {
    return MenuModel{Base: Base{TableName: "goadmin_menu"}}
}

// MenuWithID return a default menu model of given id.
func MenuWithID(id string) MenuModel {
    idInt, _ := strconv.Atoi(id)
    return MenuModel{Base: Base{TableName: "goadmin_menu"}, ID: int64(idInt)}
}

// SetConn 设置数据库连接
func (t MenuModel) SetConn(con db.Connection) MenuModel {
    t.Conn = con
    return t
}

// Find return a default menu model of given id.
func (t MenuModel) Find(id interface{}) MenuModel {
    item, _ := t.Table(t.TableName).Find(id)
    return t.MapToModel(item)
}

// New create a new menu model.
func (t MenuModel) New(title, icon, uri, header, pluginName string, parentID, order int64) (MenuModel, error) {

    id, err := t.Table(t.TableName).Insert(dialect.H{
        "title":      title,
        "parentId":   parentID,
        "icon":       icon,
        "uri":        uri,
        "order":      order,
        "header":     header,
        "pluginName": pluginName,
        `createAt`:   uint32(time.Now().Unix()),
    })

    t.ID = id
    t.Title = title
    t.ParentID = parentID
    t.Icon = icon
    t.URI = uri
    t.Header = header

    return t, err
}

// Delete delete the menu model.
func (t MenuModel) Delete() {
    _ = t.Table(t.TableName).Where("id", "=", t.ID).Delete()
    _ = t.Table("goadmin_role_menu").Where("menuId", "=", t.ID).Delete()
    items, _ := t.Table(t.TableName).Where("parentId", "=", t.ID).All()

    if len(items) > 0 {
        ids := make([]interface{}, len(items))
        for i := 0; i < len(ids); i++ {
            ids[i] = items[i]["id"]
        }
        _ = t.Table("goadmin_role_menu").WhereIn("menuId", ids).Delete()
    }

    _ = t.Table(t.TableName).Where("parentId", "=", t.ID).Delete()
}

// Update update the menu model.
func (t MenuModel) Update(title, icon, uri, header, pluginName string, parentID int64) (int64, error) {
    return t.Table(t.TableName).
        Where("id", "=", t.ID).
        Update(dialect.H{
            "title":      title,
            "parentId":   parentID,
            "icon":       icon,
            "pluginName": pluginName,
            "uri":        uri,
            "header":     header,
            "updateat":   uint32(time.Now().Unix()),
        })
}

// OrderItems 排序数组
type OrderItems []OrderItem

// OrderItem 排序
type OrderItem struct {
    ID       uint       `json:"id"`
    Children OrderItems `json:"children"`
}

// ResetOrder update the order of menu models.
func (t MenuModel) ResetOrder(data []byte) {

    var items OrderItems
    _ = json.Unmarshal(data, &items)

    count := 1
    for _, v := range items {
        if len(v.Children) > 0 {
            _, _ = t.Table(t.TableName).
                Where("id", "=", v.ID).
                Update(dialect.H{
                    "order":    count,
                    "parentId": 0,
                })

            for _, v2 := range v.Children {
                if len(v2.Children) > 0 {

                    _, _ = t.Table(t.TableName).
                        Where("id", "=", v2.ID).
                        Update(dialect.H{
                            "order":    count,
                            "parentId": v.ID,
                        })

                    for _, v3 := range v2.Children {
                        _, _ = t.Table(t.TableName).
                            Where("id", "=", v3.ID).
                            Update(dialect.H{
                                "order":    count,
                                "parentId": v2.ID,
                            })
                        count++
                    }
                } else {
                    _, _ = t.Table(t.TableName).
                        Where("id", "=", v2.ID).
                        Update(dialect.H{
                            "order":    count,
                            "parentId": v.ID,
                        })
                    count++
                }
            }
        } else {
            _, _ = t.Table(t.TableName).
                Where("id", "=", v.ID).
                Update(dialect.H{
                    "order":    count,
                    "parentId": 0,
                })
            count++
        }
    }
}

// CheckRole check the role if user has permission to get the menu.
func (t MenuModel) CheckRole(roleID string) bool {
    checkRole, _ := t.Table("goadmin_role_menu").
        Where("roleId", "=", roleID).
        Where("menuId", "=", t.ID).
        First()
    return checkRole != nil
}

// AddRole add a role to the menu.
func (t MenuModel) AddRole(roleID string) (int64, error) {
    if roleID != "" {
        if !t.CheckRole(roleID) {
            return t.Table("goadmin_role_menu").
                Insert(dialect.H{
                    "roleId": roleID,
                    "menuId": t.ID,
                })
        }
    }
    return 0, nil
}

// DeleteRoles delete roles with menu.
func (t MenuModel) DeleteRoles() error {
    return t.Table("goadmin_role_menu").
        Where("menuId", "=", t.ID).
        Delete()
}

// MapToModel get the menu model from given map.
func (t MenuModel) MapToModel(m map[string]interface{}) MenuModel {
    t.ID = m["id"].(int64)
    t.Title, _ = m["title"].(string)
    t.ParentID = m["parentId"].(int64)
    t.Icon, _ = m["icon"].(string)
    t.URI, _ = m["uri"].(string)
    t.Header, _ = m["header"].(string)
    t.CreateAt, _ = m["createAt"].(int64)
    t.UpdateAt, _ = m["updateAt"].(int64)
    return t
}
