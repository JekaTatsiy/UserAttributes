package user

import (
	//"context"
	//"database/sql"
	//"strconv"
	//"errors"
	"database/sql"
	"fmt"
	"math/rand"
	"time"
	"unsafe"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ServiceUserImpl struct {
	db *gorm.DB
}

type UserAttrString struct {
	UserCode  string
	AttrCode  string
	Value     string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt sql.NullTime
}

type ServiceUser interface {
	NewUser() (string, error)
	GetUserAttr(userCode string) ([]*UserAttrString, error)
	GetUsersByAttr(map[string]string) ([]string, error)
	UpdateUserStringAttr(attrId int, userId, value string) error
	DeleteUserStringAttr(attrId int, userId string) error
}

func (s *ServiceUserImpl) NewUser() (string, error) {

	var src = rand.NewSource(time.Now().UnixNano())

	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	const (
		letterIdxBits = 6                    // 6 bits to represent a letter index
		letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
		letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
	)
	b := make([]byte, 16)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := 16-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	code := *(*string)(unsafe.Pointer(&b))
	return code, nil
}

func (s *ServiceUserImpl) GetUserAttr(id string) ([]*UserAttrString, error) {
	attrs := make([]*UserAttrString, 0)
	result := s.db.Table("user_attr_string").
		Where("deleted_at is null").
		Find(&attrs, "user_code=?", id)

	if result.Error != nil {
		return nil, result.Error
	}
	return attrs, result.Error
}
func (s *ServiceUserImpl) GetUsersByAttr(attrs map[string]string) ([]string, error) {
	users := make([]string, 0)
	var attr, value, pair string

	for i, x := range attrs {
		attr += fmt.Sprintf("'%s', ", i)
		value += fmt.Sprintf("'%s', ", x)
		pair += fmt.Sprintf("'%s@%s', ", i, x)
	}
	attr = attr[0 : len(attr)-2]
	value = value[0 : len(value)-2]
	pair = pair[0 : len(pair)-2]

	result := s.db.Debug().Table("user_attr_string").Raw(
		fmt.Sprintf(`
		with
		users_with_args as not materialized(
			select 
				user_code ,
				attr_code ,
				value 
			from user_attr_string
			where deleted_at is null 
				and attr_code = any (array[%s]::varchar[])
				and value = any (array[%s]::varchar[])
		)
		,users as not materialized(
			select 
			user_code ,
			array_agg(concat(attr_code,'@',value))over(partition by user_code)::varchar[] as value
			from users_with_args
		)
		select distinct
			user_code
		from users
			where value @> array[%s]::varchar[]
	`, attr, value, pair)).
		Find(&users)

	return users, result.Error
}

func (s *ServiceUserImpl) UpdateUserStringAttr(userCode, attrCode, value string) error {
	result := s.db.Table("user_attr_string").Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_code"}, {Name: "attr_code"}},
		DoUpdates: clause.AssignmentColumns([]string{"value", "updated_at"}),
	}).
		Create(&UserAttrString{
			UserCode:  userCode,
			AttrCode:  attrCode,
			Value:     value,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			DeletedAt: sql.NullTime{Valid: false},
		})

	return result.Error

}

func (s *ServiceUserImpl) DeleteUserStringAttr(userCode, attrCode string) error {
	result := s.db.Table("user_attr_string").Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_code"}, {Name: "attr_code"}},
		DoNothing: true,
	}).
		Where("user_code = ? and attr_code = ?", userCode, attrCode).
		Update("deleted_at", time.Now())

	return result.Error
}

func New(db *gorm.DB) *ServiceUserImpl {
	serv := ServiceUserImpl{db: db}
	return &serv
}
