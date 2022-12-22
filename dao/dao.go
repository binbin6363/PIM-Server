package dao

import (
	"PIM_Server/model"
	"context"
	"errors"

	"PIM_Server/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Dao struct {
	commDB *gorm.DB
}

func New(dsn string, dataCenterId, workerId int64) *Dao {
	d := &Dao{}

	cli, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("dao: New db gorm client error(%v)", err)
	}
	d.commDB = cli
	return d
}

func (d *Dao) db(ctx context.Context) *gorm.DB {
	return d.commDB
}

func (d *Dao) GetGroupMemberList(ctx context.Context, groupId int64) (error, []*model.GroupMembers) {
	r := d.db(ctx)
	if groupId == 0 {
		log.Error("group id is invalid")
		return errors.New("group id is invalid"), nil
	}

	groupMembers := make([]*model.GroupMembers, 0)
	if err := r.Table(model.GroupMembers{}.TableName()).Debug().Where("group_id=?", groupId).Scan(&groupMembers).Error; err != nil {
		log.Infof("GetGroupMemberList read db error(%v) groupId(%d)", err, groupId)
		return err, nil
	}

	log.Infof("GetGroupInfo read db ok groupId(%d), members size:%d", groupId, len(groupMembers))
	return nil, groupMembers
}
