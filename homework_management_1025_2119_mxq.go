// 代码生成时间: 2025-10-25 21:19:27
package main

import (
    "beego/logs"
    "encoding/json"
    "github.com/astaxie/beego"
)

// HomeworkModel represents a homework entry in the management platform
type HomeworkModel struct {
    ID       int    `orm:"auto"` // Unique identifier for homework
    Title    string `orm:"size(100)"` // Title of the homework
    Content  string `orm:"type(text)"` // Content of the homework
    Deadline string `orm:"size(100)"` // Deadline of the homework in format YYYY-MM-DD
}

// HomeworkController handles all requests related to homework
type HomeworkController struct {
    beego.Controller
}

// Prepare is called before any controller action
func (c *HomeworkController) Prepare() {
    // You can add your logic here if needed
}

// AddHomework handles the logic to add a new homework entry
func (c *HomeworkController) AddHomework() {
    var homework HomeworkModel
    err := json.Unmarshal(c.Ctx.Input.RequestBody, &homework)
    if err != nil {
        c.Data["json"] = map[string]interface{}{"error": "Invalid input"}
        c.ServeJSON()
        return
    }
    _, err = beego.ORM().Insert(&homework)
    if err != nil {
        c.Data["json"] = map[string]interface{}{"error": "Failed to add homework"}
        c.ServeJSON()
        return
    }
    c.Data["json"] = map[string]interface{}{"success": true, "message": "Homework added successfully"}
    c.ServeJSON()
}

// GetHomework retrieves homework entries
func (c *HomeworkController) GetHomework() {
    homeworks := []HomeworkModel{}
    _, err := beego.ORM().All(&homeworks)
    if err != nil {
        c.Data["json"] = map[string]interface{}{"error": "Failed to retrieve homeworks"}
        c.ServeJSON()
        return
    }
    c.Data["json"] = homeworks
    c.ServeJSON()
}

// UpdateHomework updates an existing homework entry
func (c *HomeworkController) UpdateHomework() {
    var homework HomeworkModel
    err := json.Unmarshal(c.Ctx.Input.RequestBody, &homework)
    if err != nil {
        c.Data["json"] = map[string]interface{}{"error": "Invalid input"}
        c.ServeJSON()
        return
    }
    _, err = beego.ORM().Update(&homework)
    if err != nil {
        c.Data["json"] = map[string]interface{}{"error": "Failed to update homework"}
        c.ServeJSON()
        return
    }
    c.Data["json"] = map[string]interface{}{"success": true, "message": "Homework updated successfully"}
    c.ServeJSON()
}

// DeleteHomework removes a homework entry
func (c *HomeworkController) DeleteHomework() {
    id, _ := c.GetInt("id")
    if id <= 0 {
        c.Data["json"] = map[string]interface{}{"error": "Invalid homework ID"}
        c.ServeJSON()
        return
    }
    _, err := beego.ORM().Delete(&HomeworkModel{ID: id})
    if err != nil {
        c.Data["json"] = map[string]interface{}{"error": "Failed to delete homework"}
        c.ServeJSON()
        return
    }
    c.Data["json"] = map[string]interface{}{"success": true, "message": "Homework deleted successfully"}
    c.ServeJSON()
}

func main() {
    beego.Router("/homework/add", &HomeworkController{}, "post:AddHomework")
    beego.Router("/homework/get", &HomeworkController{}, "get:GetHomework")
    beego.Router("/homework/update", &HomeworkController{}, "put:UpdateHomework")
    beego.Router("/homework/delete", &HomeworkController{}, "delete:DeleteHomework")
    beego.Run()
}
