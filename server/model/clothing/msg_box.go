// 自动生成模板MsgBox
package clothing

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// MsgBox 结构体
type MsgBox struct {
	global.GVA_MODEL
	From      uint  `json:"from" form:"from" gorm:"column:from_user;comment:;"`
	To        uint  `json:"to" form:"to" gorm:"column:to_user;comment:;"`
	MsgType   uint  `json:"msgType" form:"msgType" gorm:"column:msg_type;comment:;"`
	MsgID     uint  `json:"msgID" form:"msgID" gorm:"column:msg_id;comment:;"`
	Status    *bool `json:"status" form:"status" gorm:"column:status;comment:;"`
	CreatedBy uint  `gorm:"column:created_by;comment:创建者"`
	UpdatedBy uint  `gorm:"column:updated_by;comment:更新者"`
	DeletedBy uint  `gorm:"column:deleted_by;comment:删除者"`
}

// TableName MsgBox 表名
func (MsgBox) TableName() string {
	return "msg_box"
}
