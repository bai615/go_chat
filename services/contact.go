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
		objIds = append(objIds, v.Dstid)
	}
	users := make([]models.User, 0)
	if len(objIds) == 0 {
		return users
	}
	DBEngine.In("id", objIds).Find(&users)
	return users
}

// 查找群
func (service *ContactService) SearchCommunity(userId int64) ([]models.Community) {
	contacts := make([]models.Contact, 0)
	comIds := make([]int64, 0)

	DBEngine.Where("ownerid = ? and cate = ?", userId, models.CONCAT_CATE_COMUNITY).Find(&contacts)
	for _, v := range contacts {
		comIds = append(comIds, v.Dstid)
	}
	coms := make([]models.Community, 0)
	if len(comIds) == 0 {
		return coms
	}
	DBEngine.In("id", comIds).Find(&coms)
	return coms
}

// 建群
func (service *ContactService) CreateCommunity(comm models.Community) (ret models.Community, err error) {
	if len(comm.Name) == 0 {
		err = errors.New("缺少群名称")
		return ret, err
	}
	if comm.Ownerid == 0 {
		err = errors.New("请先登录")
		return ret, err
	}
	com := models.Community{
		Ownerid: comm.Ownerid,
	}
	num, err := DBEngine.Count(&com)

	if (num > 5) {
		err = errors.New("一个用户最多只能创见5个群")
		return com, err
	} else {
		comm.Createat = time.Now()
		session := DBEngine.NewSession()
		session.Begin()
		_, err = session.InsertOne(&comm)
		if err != nil {
			session.Rollback()
			return com, err
		}
		_, err = session.InsertOne(
			models.Contact{
				Ownerid:  comm.Ownerid,
				Dstid:    comm.Id,
				Cate:     models.CONCAT_CATE_COMUNITY,
				Createat: time.Now(),
			})
		if err != nil {
			session.Rollback()
		} else {
			session.Commit()
		}
		return com, err
	}
}

// 加群
func (service *ContactService) JoinCommunity(userId, comId int64) error {
	cot := models.Contact{
		Ownerid: userId,
		Dstid:   comId,
		Cate:    models.CONCAT_CATE_COMUNITY,
	}
	DBEngine.Get(&cot)
	if (cot.Id == 0) {
		cot.Createat = time.Now()
		_, err := DBEngine.InsertOne(cot)
		return err
	} else {
		return nil
	}

}

func (service *ContactService) SearchCommunityIds(userId int64) (comIds []int64) {
	// todo 获取用户全部群 ID
	contacts := make([]models.Contact, 0)
	comIds = make([]int64, 0)

	DBEngine.Where("ownerid = ? and cate = ?", userId, models.CONCAT_CATE_COMUNITY).Find(&contacts)
	for _, v := range contacts {
		comIds = append(comIds, v.Dstid)
	}
	return comIds
}
