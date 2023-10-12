package http

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"go-web-demo/app/demo-api/internal/model"
	"go-web-demo/library/ecode"
	"go-web-demo/library/render"
	"strconv"
	"time"
)

const lp_key = "lp:%s"

var lpIdIdx = 0

// @Summary 学生列表
// @Produce json
// @Param studName query string true "学生姓名"
// @Success 200 {object} render.JSON
// @Router /api/v1/student/list [get]
func GetLeaveProcess(c *gin.Context) {
	r := render.New(c)

	studName := c.Query("studName")

	studList, err := srv.ListStudent(c, studName)
	if err != nil {
		r.JSON(nil, ecode.RequestErr)
		return
	}

	r.JSON(studList, nil)
}

// @Summary 学生列表
// @Produce json
// @Param studName query string true "学生姓名"
// @Success 200 {object} render.JSON
// @Router /api/v1/student/list [get]
func LastLeaveProcess(c *gin.Context) {
	r := render.New(c)

	lp, err := getLp(strconv.Itoa(lpIdIdx))
	if err != nil {
		r.JSON(nil, ecode.RequestErr)
		return
	}

	r.JSON(lp, nil)
}

// @Summary 学生列表
// @Produce json
// @Param studName query string true "学生姓名"
// @Success 200 {object} render.JSON
// @Router /api/v1/student/list [get]
func ListLeaveProcess(c *gin.Context) {
	r := render.New(c)

	studName := c.Query("studName")

	studList, err := srv.ListStudent(c, studName)
	if err != nil {
		r.JSON(nil, ecode.RequestErr)
		return
	}

	r.JSON(studList, nil)
}

// @Summary 添加学生
// @Produce json
// @Param studName query string true "学生姓名"
// @Param studAge query int true "年龄"
// @Param studSex query string true "性别"
// @Success 200 {object} render.JSON
// @Router /api/v1/student/add [post]
func AddLeaveProcess(c *gin.Context) {
	r := render.New(c)

	v := &model.LeaveProcess{}

	err := c.Bind(v)
	if err != nil {
		r.JSON(nil, ecode.RequestErr)
		return
	}

	// todo type

	jsonStr, err := toJsonString(v)
	if err != nil {
		r.JSON(nil, ecode.RequestErr)
		return
	}

	v.CreateTime = time.Now()
	v.UpdateTime = time.Now()

	id := strconv.FormatInt(v.Id, 10)
	key := fmt.Sprintf(lp_key, id)
	//id,;
	err1 := srv.SetRedisKey(c, key, jsonStr, 60*60*24*7)
	r.JSON(v, err1)
}

func getLp(id string) (*model.LeaveProcess, error) {
	key := fmt.Sprintf(lp_key, id)
	json, err := srv.GetRedisKey(context.TODO(), key)
	if err != nil {
		return nil, err
	}
	if "" != json {
		lp := &model.LeaveProcess{}
		err := toObject(json, lp)
		if err != nil {
			return nil, err
		}
		return lp, nil
	}
	return nil, err
}

func TakeOut(lp model.LeaveProcess) error {
	user, err := getUser(lp.UserId)
	if err != nil {
		return err
	}

	leaveRemaining := 0
	leaveType := lp.Type
	switch leaveType {
	case model.AnnualLeave:
		if user.AnnualLeave < lp.Len {
			return errors.New("该类型假期不足")
		}
		user.AnnualLeave -= lp.Len
	case model.SickLeave:
		if user.SickLeave < lp.Len {
			return errors.New("该类型假期不足")
		}
		user.SickLeave -= lp.Len
	case model.LieuLeave:
		if user.LieuLeave < lp.Len {
			return errors.New("该类型假期不足")
		}
		user.LieuLeave -= lp.Len
	case model.PersonalLeave:
		// do nothing
	default:
		return errors.New("无效假期类型")
	}

	if leaveType != model.PersonalLeave {
		if leaveRemaining < lp.Len {
			return errors.New("该类型假期不足")
		}

	}

	return nil
}
