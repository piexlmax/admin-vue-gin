package request

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/reservation"
	"time"
)

type EqtInfoSearch struct {
	reservation.EqtInfo
	StartCreatedAt   *time.Time `json:"startCreatedAt" form:"startCreatedAt"`
	EndCreatedAt     *time.Time `json:"endCreatedAt" form:"endCreatedAt"`
	StartServiceTime *time.Time `json:"startServiceTime" form:"startServiceTime"`
	EndServiceTime   *time.Time `json:"endServiceTime" form:"endServiceTime"`
	request.PageInfo
}
