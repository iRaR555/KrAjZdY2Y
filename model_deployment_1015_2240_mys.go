// 代码生成时间: 2025-10-15 22:40:35
package main

import (
    "beego框架的导入路径"
    "encoding/json"
    "errors"
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
    "os"
    "os/exec"
    "path/filepath"
)

// ModelDeploymentController 是模型部署工具的控制器
type ModelDeploymentController struct {
    beego.Controller
}

// DeployModel 部署模型的处理函数
func (c *ModelDeploymentController) DeployModel() {
    // 获取模型文件路径参数
    modelPath := c.GetString("modelPath")
    if modelPath == "" {
        // 参数检查
        c.CustomAbort(http.StatusBadRequest, "Model path is required")
        return
    }

    // 检查模型文件是否存在
    if _, err := os.Stat(modelPath); os.IsNotExist(err) {
        c.CustomAbort(http.StatusNotFound, fmt.Sprintf("Model file not found: %s", modelPath))
        return
    }

    // 部署模型到指定路径
    deployPath := c.GetString("deployPath")
    if deployPath == "" {
        deployPath = "./models" // 默认部署路径
    }

    // 创建部署目录
    if err := os.MkdirAll(deployPath, 0755); err != nil {
        c.CustomAbort(http.StatusInternalServerError, fmt.Sprintf("Failed to create deploy directory: %v", err))
        return
    }

    // 复制模型文件到部署目录
    if err := copyFile(modelPath, deployPath); err != nil {
        c.CustomAbort(http.StatusInternalServerError, fmt.Sprintf("Failed to copy model file: %v", err))
        return
    }

    // 响应成功
    c.Data["json"] = map[string]string{"message": "Model deployed successfully"}
    c.ServeJSON()
}

// copyFile 复制文件的辅助函数
func copyFile(src, dest string) error {
    sourceFile, err := os.Open(src)
    if err != nil {
        return err
