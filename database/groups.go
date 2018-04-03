package database

// This file is generated - do not edit.

import (
	"time"

	"github.com/keiwi/utils"
	"github.com/keiwi/utils/models"
	"gopkg.in/mgo.v2/bson"
)

// GroupManager struct
type GroupManager struct {
	db *DB
}

// NewGroupManager - Creates a new *GroupManager that can be used for managing groups.
func NewGroupManager(db *DB) (*GroupManager, error) {
	manager := GroupManager{}
	manager.db = db
	return &manager, nil
}

func (state GroupManager) Find(options utils.FindOptions) ([]models.Group, error) {
	session := state.db.Session.Copy()
	defer session.Close()
	query := session.DB(state.db.Database).C("groups").Find(options.Filter)

	if len(options.Sort) > 0 {
		query = query.Sort(options.Sort...)
	}

	if options.Limit > 0 {
		query = query.Limit(int(options.Limit))
	}

	var groups []models.Group
	err := query.All(&groups)
	if err != nil {
		return nil, err
	}

	return groups, nil
}

func (state GroupManager) Has(options utils.HasOptions) (bool, error) {
	session := state.db.Session.Copy()
	defer session.Close()
	query := session.DB(state.db.Database).C("groups").Find(options.Filter)

	count, err := query.Count()
	if err != nil {
		return false, err
	}
	return count >= 1, nil
}

func (state GroupManager) Create(c interface{}) (interface{}, error) {
	session := state.db.Session.Copy()
	defer session.Close()

	u := c.(*models.Group)

	if u.ID == "" {
		u.ID = bson.NewObjectId()
	}

	err := session.DB(state.db.Database).C("groups").Insert(u)
	if err != nil {
		return nil, err
	}

	var groups models.Group
	err = session.DB(state.db.Database).C("groups").FindId(u.ID).One(&groups)
	if err != nil {
		return nil, err
	}
	return groups, nil
}

func (state GroupManager) Update(options utils.UpdateOptions) (interface{}, error) {
	session := state.db.Session.Copy()
	defer session.Close()

	err := session.DB(state.db.Database).C("groups").Update(options.Filter, options.Updates)
	if err != nil {
		return nil, err
	}

	return state.Find(utils.FindOptions{Filter: options.Filter})
}

func (state GroupManager) Delete(options utils.DeleteOptions) error {
	session := state.db.Session.Copy()
	defer session.Close()

	_, err := state.Update(utils.UpdateOptions{
		Filter:  options.Filter,
		Updates: utils.Updates{"$set": bson.M{"updated_at": time.Now()}},
	})
	if err != nil {
		return err
	}

	return session.DB(state.db.Database).C("groups").Remove(options.Filter)
}
