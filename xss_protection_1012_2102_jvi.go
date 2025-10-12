// 代码生成时间: 2025-10-12 21:02:30
package main

import (
    "encoding/json"
    "github.com/astaxie/beego"
    "github.com/astaxie/beego/context"
    "strings"
)

// XSSFilterMiddleware is a middleware that filters out XSS attacks.
func XSSFilterMiddleware(ctx *context.Context) {
    var xssFilter = func(input string) string {
        // Remove all script tags from the input string
        return strings.ReplaceAll(input, "<script", "&lt;script")
    }

    // Check if the request method is POST
    if ctx.Request.Method == "POST" {
        // Iterate over the form values
        for key, value := range ctx.Request.Form {
            // Clean the value from potential XSS attacks
            cleanedValue := xssFilter(value[0])

            // Store the cleaned value back into the form values
            ctx.Request.Form[key] = []string{cleanedValue}
        }
    }

    // Continue to the next middleware
    beego.Next(ctx)
}

func main() {
    // Initialize the Beego framework
    beego.AddFuncMap("xss_filter", XSSFilterMiddleware)

    // Define a route that uses the XSS filter middleware
    beego.Router("/", &MainController{}, "get:SayHello")

    // Run the application
    beego.Run()
}

// MainController is the main controller for the application.
type MainController struct {
    beego.Controller
}

// SayHello is the action that responds with a hello message.
func (c *MainController) SayHello() {
    // Get the user input from the query string
    input := c.GetString("input")

    // Clean the input from potential XSS attacks
    cleanedInput := XSSFilterMiddleware(c.Ctx)

    // Respond with the cleaned input
    c.Data["json"] = map[string]string{"message": cleanedInput}
    c.ServeJSON()
}
