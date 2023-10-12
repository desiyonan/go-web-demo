package model

type LeaveType string

const (
	AnnualLeave   = LeaveType("年假")
	SickLeave     = LeaveType("病假")
	LieuLeave     = LeaveType("调休假")
	PersonalLeave = LeaveType("事假")
)
