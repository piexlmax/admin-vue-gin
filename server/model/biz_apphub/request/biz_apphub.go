package request

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"time"
)

type BizAppHubSearch struct{
    StartCreatedAt *time.Time `json:"startCreatedAt" form:"startCreatedAt"`
    EndCreatedAt   *time.Time `json:"endCreatedAt" form:"endCreatedAt"`
    AppName  string `json:"appName" form:"appName" `
    AppCode  string `json:"appCode" form:"appCode" `
    Title  string `json:"title" form:"title" `
    Desc  string `json:"desc" form:"desc" `
    Classify  string `json:"classify" form:"classify" `
    Version  string `json:"version" form:"version" `
    Mode  string `json:"mode" form:"mode" `
    Tags  string `json:"tags" form:"tags" `
    Video  string `json:"video" form:"video" `
    request.PageInfo
    Sort  string `json:"sort" form:"sort"`
    Order string `json:"order" form:"order"`
}
