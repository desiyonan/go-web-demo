package model

import "time"

type LeaveProcess struct {
	Id         int64     `json:"id"`
	UserId     string    `json:"user_id"`
	Approver   string    `json:"approver"` // 审批人
	Status     string    `json:"status"`
	Type       LeaveType `json:"type"`
	Len        int       `json:"len"`
	StartTime  string    `json:"start_time"`
	EndTime    string    `json:"end_time"`
	CreateTime time.Time `json:"create_time"`
	UpdateTime time.Time `json:"update_time"`
}
