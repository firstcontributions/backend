package mongoqb

import (
	"errors"

	"go.mongodb.org/mongo-driver/bson"
)

type UpdateMap struct {
	set    interface{}
	setMap bson.M
	inc    interface{}
	incMap bson.M
}

func NewUpdateMap() *UpdateMap {
	return &UpdateMap{}
}

func (u *UpdateMap) Set(field string, value interface{}) *UpdateMap {
	u.setMap[field] = value
	return u
}

func (u *UpdateMap) Inc(field string, value int) *UpdateMap {
	u.incMap[field] = value
	return u
}

func (u *UpdateMap) SetFields(updateObj interface{}) *UpdateMap {
	u.set = updateObj
	return u
}

func (u *UpdateMap) IncFields(updateObj interface{}) *UpdateMap {
	u.inc = updateObj
	return u
}

func (u *UpdateMap) BuildUpdate() (bson.M, error) {
	update := bson.M{}
	if u.inc != nil && u.incMap != nil {
		return nil, errors.New("should not use Inc and IncFields together")
	}
	if u.set != nil && u.setMap != nil {
		return nil, errors.New("should not use Set and SetFields together")
	}
	if u.inc != nil {
		update["$inc"] = u.inc
	}
	if u.incMap != nil {
		update["$inc"] = u.incMap
	}
	if u.set != nil {
		update["$set"] = u.set
	}
	if u.setMap != nil {
		update["$set"] = u.setMap
	}
	return update, nil
}
