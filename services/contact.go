package services

import (
	"time"
	"errors"
	"go_chat/models"
)

type ContactService struct {
}

// 添加好友
func (service *ContactService) AddFriend(
	userid, // 用户id 10086,
	dstid int64) error {
	// 如果加自己
	if userid == dstid {
		return errors.New("不能添加自己为好友啊")
	}
	// 判断是否已经加了好友
	tmp := models.Contact{}
	// 查询是否已经是好友
	// 条件的链式操作
	DBEngine.Where("ownerid = ?", userid).
		And("dstid = ?", dstid).
		And("cate = ?", models.CONCAT_CATE_USER).
		Get(&tmp)
	// 获得1条记录
	// count()
	// 如果存在记录说明已经是好友了不加
	if tmp.Id > 0 {
		return errors.New("该用户已经被添加过啦")
	}
	// 事务,
	session := DBEngine.NewSession()
	session.Begin()
	// 插自己的
	_, e2 := session.InsertOne(models.Contact{
		Ownerid:  userid,
		Dstid:    dstid,
		Cate:     models.CONCAT_CATE_USER,
		Createat: time.Now(),
	})
	// 插对方的
	_, e3 := session.InsertOne(models.Contact{
		Ownerid:  dstid,
		Dstid:    userid,
		Cate:     models.CONCAT_CATE_USER,
		Createat: time.Now(),
	})
	// 没有错误
	if e2 == nil && e3 == nil {
		// 提交
		session.Commit()
		return nil
	} else {
		// 回滚
		session.Rollback()
		if e2 != nil {
			return e2
		} else {
			return e3
		}
	}
}

// 查找好友
func (service *ContactService) SearchFriend(userId int64) ([]models.User) {
	contacts := make([]models.Contact, 0)
	objIds := make([]int64, 0)
	DBEngine.Where("ownerid = ? and cate = ?", userId, models.CONCAT_CATE_USER).Find(&contacts)
	for _, v := range contacts {
		objIds = append(objIds, v.Dstid);
	}
	users := make([]models.User, 0)
	if len(objIds) == 0 {
		return users
	}
	DBEngine.In("id", objIds).Find(&users)
	return users
}
