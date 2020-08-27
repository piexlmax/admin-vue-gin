package v1

import (
	"fmt"
	"gin-vue-admin/global"
	"gin-vue-admin/global/response"
	"gin-vue-admin/model"
	"gin-vue-admin/model/request"
	resp "gin-vue-admin/model/response"
	"gin-vue-admin/service"
	"gin-vue-admin/utils"
	"github.com/gin-gonic/gin"
	//"github.com/jinzhu/gorm"
)

// @Tags authority
// @Summary 创建角色
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.SysAuthority true "创建角色"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /authority/createAuthority [post]
func CreateAuthority(c *gin.Context) {
	var auth model.SysAuthority
	_ = c.ShouldBindJSON(&auth)
	AuthorityVerify := utils.Rules{
		"AuthorityId":   {utils.NotEmpty()},
		"AuthorityName": {utils.NotEmpty()},
		"ParentId":      {utils.NotEmpty()},
	}
	AuthorityVerifyErr := utils.Verify(auth, AuthorityVerify)
	if AuthorityVerifyErr != nil {
		response.FailWithMessage(AuthorityVerifyErr.Error(), c)
		return
	}
	err, authBack := service.CreateAuthority(auth)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("创建失败，%v", err), c)
	} else {
		response.OkWithData(resp.SysAuthorityResponse{Authority: authBack}, c)
	}
}

// @Tags authority
// @Summary 拷贝角色
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body response.SysAuthorityCopyResponse true "拷贝角色"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"拷贝成功"}"
// @Router /authority/copyAuthority [post]
func CopyAuthority(c *gin.Context) {
	var copyInfo resp.SysAuthorityCopyResponse
	_ = c.ShouldBindJSON(&copyInfo)
	OldAuthorityVerify := utils.Rules{
		"OldAuthorityId": {utils.NotEmpty()},
	}
	OldAuthorityVerifyErr := utils.Verify(copyInfo, OldAuthorityVerify)
	if OldAuthorityVerifyErr != nil {
		response.FailWithMessage(OldAuthorityVerifyErr.Error(), c)
		return
	}
	AuthorityVerify := utils.Rules{
		"AuthorityId":   {utils.NotEmpty()},
		"AuthorityName": {utils.NotEmpty()},
		"ParentId":      {utils.NotEmpty()},
	}
	AuthorityVerifyErr := utils.Verify(copyInfo.Authority, AuthorityVerify)
	if AuthorityVerifyErr != nil {
		response.FailWithMessage(AuthorityVerifyErr.Error(), c)
		return
	}
	err, authBack := service.CopyAuthority(copyInfo)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("拷贝失败，%v", err), c)
	} else {
		response.OkWithData(resp.SysAuthorityResponse{Authority: authBack}, c)
	}
}

// @Tags authority
// @Summary 删除角色
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.SysAuthority true "删除角色"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /authority/deleteAuthority [post]
func DeleteAuthority(c *gin.Context) {
	var a model.SysAuthority
	_ = c.ShouldBindJSON(&a)
	AuthorityIdVerifyErr := utils.Verify(a, utils.CustomizeMap["AuthorityIdVerify"])
	if AuthorityIdVerifyErr != nil {
		response.FailWithMessage(AuthorityIdVerifyErr.Error(), c)
		return
	}
	// 删除角色之前需要判断是否有用户正在使用此角色
	err := service.DeleteAuthority(&a)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("删除失败，%v", err), c)
	} else {
		response.OkWithMessage("删除成功", c)
	}
}

// @Tags authority
// @Summary 设置角色资源权限
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.SysAuthority true "设置角色资源权限"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"设置成功"}"
// @Router /authority/updateAuthority [post]
func UpdateAuthority(c *gin.Context) {
	var auth model.SysAuthority
	_ = c.ShouldBindJSON(&auth)
	AuthorityVerify := utils.Rules{
		"AuthorityId":   {utils.NotEmpty()},
		"AuthorityName": {utils.NotEmpty()},
		"ParentId":      {utils.NotEmpty()},
	}
	AuthorityVerifyErr := utils.Verify(auth, AuthorityVerify)
	if AuthorityVerifyErr != nil {
		response.FailWithMessage(AuthorityVerifyErr.Error(), c)
		return
	}

	claims, _ := c.Get("claims")
	waitUse := claims.(*request.CustomClaims)
	if waitUse.AuthorityId != "0"{
		var data_aggregate []model.SysAuthority
		global.GVA_DB.Find(&data_aggregate)
		blid := auth.ParentId
		for{
			bool := true
			for _,k := range data_aggregate {
				if blid == waitUse.AuthorityId{//判断提交的父角色ID是否与当前ID相同
					bool = false
					break
				}
				if blid == k.AuthorityId{
					blid = k.ParentId
					break
				}
				if blid == "0"{
					response.FailWithMessage(fmt.Sprintf("更新失败，权限不足"), c)
					return
				}
			}
			if bool == false{
				break
			}
		}
	}

	err, authority := service.UpdateAuthority(auth)
		if err != nil {
			response.FailWithMessage(fmt.Sprintf("更新失败，%v", err), c)
		} else {
			response.OkWithData(resp.SysAuthorityResponse{Authority: authority}, c)
		}
	}

// @Tags authority
// @Summary 分页获取角色列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.PageInfo true "分页获取用户列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /authority/getAuthorityList [post]
func GetAuthorityList(c *gin.Context) {
	var pageInfo request.PageInfo
	_ = c.ShouldBindJSON(&pageInfo)
	PageVerifyErr := utils.Verify(pageInfo, utils.CustomizeMap["PageVerify"])
	if PageVerifyErr != nil {
		response.FailWithMessage(PageVerifyErr.Error(), c)
		return
	}
	err, list, total := service.GetAuthorityInfoList(pageInfo)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("获取数据失败，%v", err), c)
	} else {
		response.OkWithData(resp.PageResult{
			List:     list,
			Total:    total,
			Page:     pageInfo.Page,
			PageSize: pageInfo.PageSize,
		}, c)
	}
}

// @Tags authority
// @Summary 设置角色资源权限
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.SysAuthority true "设置角色资源权限"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"设置成功"}"
// @Router /authority/setDataAuthority [post]
func SetDataAuthority(c *gin.Context) {
	var auth model.SysAuthority
	_ = c.ShouldBindJSON(&auth)
	AuthorityIdVerifyErr := utils.Verify(auth, utils.CustomizeMap["AuthorityIdVerify"])
	if AuthorityIdVerifyErr != nil {
		response.FailWithMessage(AuthorityIdVerifyErr.Error(), c)
		return
	}
	err := service.SetDataAuthority(auth)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("设置关联失败，%v", err), c)
	} else {
		response.Ok(c)
	}
}


func loop(user_Authority_Id string,data_tes *[]string)[]model.SysAuthority{
	data_aggregates := *data_tes
	var authority []model.SysAuthority
	global.GVA_DB.Where("parent_id = ?", user_Authority_Id).Find(&authority)
	//return authority

	var data_aggregate []model.SysAuthority
	//data_aggregates := make([]string,1)
	data_aggregate = loop(user_Authority_Id,&data_aggregates)
	data_aggregates = append(data_aggregates, user_Authority_Id)

	for _, l := range data_aggregate {
		if l.AuthorityId == "" {
			break
		}
		data_aggregates = append(data_aggregates, l.AuthorityId) //添加父级ID到切片中
		data_aggregate = loop(l.AuthorityId,&data_aggregates)
		loop(l.AuthorityId,&data_aggregates)
	}
	return data_aggregate
}
