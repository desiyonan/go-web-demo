package http

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"go-web-demo/app/demo-api/internal/model"
	"go-web-demo/library/ecode"
	"go-web-demo/library/render"
	"strconv"
	"time"
)

const user_key = "user:%s"

var userIdIdx = 0
var currentUserId = 0

// @Summary 学生列表
// @Produce json
// @Param studName query string true "学生姓名"
// @Success 200 {object} render.JSON
// @Router /api/v1/student/list [get]
func CurrentUser(c *gin.Context) {
	r := render.New(c)

	user, err := getUser(strconv.Itoa(currentUserId))
	if err != nil {
		r.JSON(nil, ecode.RequestErr)
		return
	}

	r.JSON(user, nil)
}

// @Summary 学生列表
// @Produce json
// @Param studName query string true "学生姓名"
// @Success 200 {object} render.JSON
// @Router /api/v1/student/list [get]
func ListUser(c *gin.Context) {
	r := render.New(c)

	studName := c.Query("studName")

	studList, err := srv.ListStudent(c, studName)
	if err != nil {
		r.JSON(nil, ecode.RequestErr)
		return
	}

	r.JSON(studList, nil)
}

// @Summary 添加用户
// @Produce json
// @Param studName query string true "学生姓名"
// @Param studAge query int true "年龄"
// @Param studSex query string true "性别"
// @Success 200 {object} render.JSON
// @Router /api/v1/student/add [post]
func UpdateUser(c *gin.Context) {
	r := render.New(c)

	v := new(struct {
		Id            int64  `json:"id"`
		Name          string `json:"name"  binding:"required,min=1,max=30"`
		Leader        int64  `json:"leader"`         // 上级领导
		AnnualLeave   int64  `json:"annual_leave"`   // 年假
		SickLeave     int64  `json:"sick_leave"`     // 病假
		LieuLeave     int64  `json:"lieu_leave"`     // 调休
		PersonalLeave int64  `json:"personal_leave"` // 事假
	})

	err := c.Bind(v)
	if err != nil {
		r.JSON(nil, ecode.RequestErr)
		return
	}

	user := &model.User{
		Id:            v.Id,
		Name:          v.Name,
		Leader:        v.Leader,
		AnnualLeave:   v.AnnualLeave,
		SickLeave:     v.SickLeave,
		LieuLeave:     v.LieuLeave,
		PersonalLeave: v.PersonalLeave,
		CreateTime:    time.Now(),
		UpdateTime:    time.Now(),
	}

	userJson, err := toJsonString(user)
	if err != nil {
		r.JSON(nil, ecode.RequestErr)
		return
	}

	id := strconv.FormatInt(user.Id, 10)
	key := fmt.Sprintf(user_key, id)
	//id,;
	err1 := srv.SetRedisKey(c, key, userJson, 60*60*24*7)

	r.JSON(id, err1)
}

func getUser(id string) (*model.User, error) {
	key := fmt.Sprintf(user_key, id)
	json, err := srv.GetRedisKey(context.TODO(), key)
	if err != nil {
		return nil, err
	}
	if "" != json {
		user := &model.User{}
		err := toObject(json, user)
		if err != nil {
			return nil, err
		}
		return user, nil
	}
	return nil, err
}

func toJsonString(obj any) (string, error) {
	bytes, err := json.Marshal(obj)
	if err != nil {
		return "", err
	}

	return string(bytes), err
}

func toObject(str string, obj any) error {
	err := json.Unmarshal([]byte(str), obj)
	if err != nil {
		return err
	}

	return err
}
