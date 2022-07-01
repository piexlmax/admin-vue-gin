package system

import (
	"errors"
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	cp "github.com/otiai10/copy"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AutoCodeApi struct{}

var caser = cases.Title(language.English)

// PreviewTemp
// @Tags AutoCode
// @Summary 预览创建后的代码
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body system.AutoCodeStruct true "预览创建代码"
// @Success 200 {object} response.Response{data=map[string]interface{},msg=string} "预览创建后的代码"
// @Router /autoCode/preview [post]
func (autoApi *AutoCodeApi) PreviewTemp(c *gin.Context) {
	var a system.AutoCodeStruct
	_ = c.ShouldBindJSON(&a)
	if err := utils.Verify(a, utils.AutoCodeVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	a.KeyWord() // 处理go关键字
	a.PackageT = caser.String(a.Package)
	autoCode, err := autoCodeService.PreviewTemp(a)
	if err != nil {
		global.GVA_LOG.Error("预览失败!", zap.Error(err))
		response.FailWithMessage("预览失败", c)
	} else {
		response.OkWithDetailed(gin.H{"autoCode": autoCode}, "预览成功", c)
	}
}

// CreateTemp
// @Tags AutoCode
// @Summary 自动代码模板
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body system.AutoCodeStruct true "创建自动代码"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /autoCode/createTemp [post]
func (autoApi *AutoCodeApi) CreateTemp(c *gin.Context) {
	var a system.AutoCodeStruct
	_ = c.ShouldBindJSON(&a)
	if err := utils.Verify(a, utils.AutoCodeVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	a.KeyWord() // 处理go关键字
	var apiIds []uint
	if a.AutoCreateApiToSql {
		if ids, err := autoCodeService.AutoCreateApi(&a); err != nil {
			global.GVA_LOG.Error("自动化创建失败!请自行清空垃圾数据!", zap.Error(err))
			c.Writer.Header().Add("success", "false")
			c.Writer.Header().Add("msg", url.QueryEscape("自动化创建失败!请自行清空垃圾数据!"))
			return
		} else {
			apiIds = ids
		}
	}
	a.PackageT = caser.String(a.Package)
	err := autoCodeService.CreateTemp(a, apiIds...)
	if err != nil {
		if errors.Is(err, system.AutoMoveErr) {
			c.Writer.Header().Add("success", "true")
			c.Writer.Header().Add("msg", url.QueryEscape(err.Error()))
		} else {
			c.Writer.Header().Add("success", "false")
			c.Writer.Header().Add("msg", url.QueryEscape(err.Error()))
			_ = os.Remove("./ginvueadmin.zip")
		}
	} else {
		c.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s", "ginvueadmin.zip")) // fmt.Sprintf("attachment; filename=%s", filename)对下载的文件重命名
		c.Writer.Header().Add("Content-Type", "application/json")
		c.Writer.Header().Add("success", "true")
		c.File("./ginvueadmin.zip")
		_ = os.Remove("./ginvueadmin.zip")
	}
}

// GetDB
// @Tags AutoCode
// @Summary 获取当前所有数据库
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=map[string]interface{},msg=string} "获取当前所有数据库"
// @Router /autoCode/getDatabase [get]
func (autoApi *AutoCodeApi) GetDB(c *gin.Context) {
	dbs, err := autoCodeService.Database().GetDB()
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
	} else {
		response.OkWithDetailed(gin.H{"dbs": dbs}, "获取成功", c)
	}
}

// GetTables
// @Tags AutoCode
// @Summary 获取当前数据库所有表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=map[string]interface{},msg=string} "获取当前数据库所有表"
// @Router /autoCode/getTables [get]
func (autoApi *AutoCodeApi) GetTables(c *gin.Context) {
	dbName := c.DefaultQuery("dbName", global.GVA_CONFIG.Mysql.Dbname)
	tables, err := autoCodeService.Database().GetTables(dbName)
	if err != nil {
		global.GVA_LOG.Error("查询table失败!", zap.Error(err))
		response.FailWithMessage("查询table失败", c)
	} else {
		response.OkWithDetailed(gin.H{"tables": tables}, "获取成功", c)
	}
}

// GetColumn
// @Tags AutoCode
// @Summary 获取当前表所有字段
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=map[string]interface{},msg=string} "获取当前表所有字段"
// @Router /autoCode/getColumn [get]
func (autoApi *AutoCodeApi) GetColumn(c *gin.Context) {
	dbName := c.DefaultQuery("dbName", global.GVA_CONFIG.Mysql.Dbname)
	tableName := c.Query("tableName")
	columns, err := autoCodeService.Database().GetColumn(tableName, dbName)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
	} else {
		response.OkWithDetailed(gin.H{"columns": columns}, "获取成功", c)
	}
}

// CreatePackage
// @Tags AutoCode
// @Summary 创建package
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body system.SysAutoCode true "创建package"
// @Success 200 {object} response.Response{data=map[string]interface{},msg=string} "创建package成功"
// @Router /autoCode/createPackage [post]
func (autoApi *AutoCodeApi) CreatePackage(c *gin.Context) {
	var a system.SysAutoCode
	_ = c.ShouldBindJSON(&a)
	if err := utils.Verify(a, utils.AutoPackageVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err := autoCodeService.CreateAutoCode(&a)
	if err != nil {
		global.GVA_LOG.Error("创建成功!", zap.Error(err))
		response.FailWithMessage("创建失败", c)
	} else {
		response.OkWithMessage("创建成功", c)
	}
}

// GetPackage
// @Tags AutoCode
// @Summary 获取package
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=map[string]interface{},msg=string} "创建package成功"
// @Router /autoCode/getPackage [post]
func (autoApi *AutoCodeApi) GetPackage(c *gin.Context) {
	pkgs, err := autoCodeService.GetPackage()
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
	} else {
		response.OkWithDetailed(gin.H{"pkgs": pkgs}, "获取成功", c)
	}
}

// DelPackage
// @Tags AutoCode
// @Summary 删除package
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body system.SysAutoCode true "创建package"
// @Success 200 {object} response.Response{data=map[string]interface{},msg=string} "删除package成功"
// @Router /autoCode/delPackage [post]
func (autoApi *AutoCodeApi) DelPackage(c *gin.Context) {
	var a system.SysAutoCode
	_ = c.ShouldBindJSON(&a)
	err := autoCodeService.DelPackage(a)
	if err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
	} else {
		response.OkWithMessage("删除成功", c)
	}
}

// AutoPlug
// @Tags AutoCode
// @Summary 创建插件模板
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body system.SysAutoCode true "创建插件模板"
// @Success 200 {object} response.Response{data=map[string]interface{},msg=string} "创建插件模板成功"
// @Router /autoCode/createPlug [post]
func (autoApi *AutoCodeApi) AutoPlug(c *gin.Context) {
	var a system.AutoPlugReq
	_ = c.ShouldBindJSON(&a)
	a.Snake = strings.ToLower(a.PlugName)
	a.NeedModel = a.HasRequest || a.HasResponse
	err := autoCodeService.CreatePlug(a)
	if err != nil {
		global.GVA_LOG.Error("预览失败!", zap.Error(err))
		response.FailWithMessage("预览失败", c)
	} else {
		response.Ok(c)
	}
}

func (autoApi *AutoCodeApi) InstallPlugin(c *gin.Context) {
	const GVAPLUGPATH = "./gva-plug-temp/"
	defer os.RemoveAll(GVAPLUGPATH)
	header, err := c.FormFile("plug")
	_, err = os.Stat(GVAPLUGPATH)
	if os.IsNotExist(err) {
		os.Mkdir(GVAPLUGPATH, os.ModePerm)
	}
	err = c.SaveUploadedFile(header, GVAPLUGPATH+header.Filename)
	paths, err := utils.Unzip(GVAPLUGPATH+header.Filename, GVAPLUGPATH)
	var webIndex = 0
	var serverIndex = 0
	for i := range paths {
		paths[i] = filepath.ToSlash(paths[i])
		pathArr := strings.Split(paths[i], "/")
		if pathArr[len(pathArr)-2] == "server" && pathArr[len(pathArr)-1] == "plugin" {
			serverIndex = i + 1
		}
		if pathArr[len(pathArr)-2] == "web" && pathArr[len(pathArr)-1] == "plugin" {
			webIndex = i + 1
		}
	}
	if webIndex == 0 && serverIndex == 0 {
		fmt.Println("非标准插件，请按照文档自动迁移使用")
		response.FailWithMessage("非标准插件，请按照文档自动迁移使用", c)
		return
	}

	if webIndex != 0 {
		webNameArr := strings.Split(filepath.ToSlash(paths[webIndex]), "/")
		webName := webNameArr[len(webNameArr)-1]
		var form = filepath.ToSlash(global.GVA_CONFIG.AutoCode.Root + global.GVA_CONFIG.AutoCode.Server + "/" + paths[webIndex])
		var to = filepath.ToSlash(global.GVA_CONFIG.AutoCode.Root + global.GVA_CONFIG.AutoCode.Web + "/plugin/" + webName)
		_, err := os.Stat(to)
		if err == nil {
			fmt.Println("web 已存在同名插件，请自行手动安装")
			response.FailWithMessage("web 已存在同名插件，请自行手动安装", c)
			return
		}
		err = cp.Copy(form, to)
		if err != nil {
			response.FailWithMessage(err.Error(), c)
			return
		}
	}

	if serverIndex != 0 {
		serverNameArr := strings.Split(filepath.ToSlash(paths[serverIndex]), "/")
		serverName := serverNameArr[len(serverNameArr)-1]
		var form = filepath.ToSlash(global.GVA_CONFIG.AutoCode.Root + global.GVA_CONFIG.AutoCode.Server + "/" + paths[serverIndex])
		var to = filepath.ToSlash(global.GVA_CONFIG.AutoCode.Root + global.GVA_CONFIG.AutoCode.Server + "/plugin/" + serverName)
		_, err := os.Stat(to)
		if err == nil {
			fmt.Println("server 已存在同名插件，请自行手动安装")
			response.FailWithMessage("server 已存在同名插件，请自行手动安装", c)
			return
		}
		err = cp.Copy(form, to)
		if err != nil {
			response.FailWithMessage(err.Error(), c)
			return
		}
	}
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	} else {
		response.OkWithMessage("插件安装成功，请按照说明配置使用", c)
	}
}
