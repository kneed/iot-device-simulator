package form

import "gorm.io/datatypes"

type CreateProtocolForm struct {
	DeviceId int            `form:"device_id" json:"device_id" binding:"required"`
	Name     string         `form:"name" json:"name" binding:"required"`
	Content  datatypes.JSON `form:"content" json:"content" binding:"required"`
	Qos      *int            `form:"qos" json:"qos" binding:"required"`
	Type     *int            `form:"type" json:"type" binding:"required"`
	SubTopic string         `form:"sub_topic" json:"sub_topic"`
	PubTopic string         `form:"pub_topic" json:"pub_topic" binding:"required"`
	Strategy datatypes.JSON `form:"strategy" json:"strategy"`
}
