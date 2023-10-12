package model

import "time"

type User struct {
	Id            int64     `json:"id"`
	Name          string    `json:"name"`
	Leader        int64     `json:"leader"`         // 上级领导
	AnnualLeave   int       `json:"annual_leave"`   // 年假
	SickLeave     int       `json:"sick_leave"`     // 病假
	LieuLeave     int       `json:"lieu_leave"`     // 调休
	PersonalLeave int       `json:"personal_leave"` // 事假
	CreateTime    time.Time `json:"create_time"`
	UpdateTime    time.Time `json:"update_time"`
}
