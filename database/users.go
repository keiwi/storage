package database

// This file is generated - do not edit.

import (
	"time"

	"github.com/keiwi/utils"
	"github.com/keiwi/utils/models"
	"gopkg.in/mgo.v2/bson"
)

// UserManager struct
type UserManager struct {
	db *DB
}

// NewUserManager - Creates a new *UserManager that can be used for managing users.
func NewUserManager(db *DB) (*UserManager, error) {
	manager := UserManager{}
	manager.db = db
	return &manager, nil
}

func (state UserManager) Find(options utils.FindOptions) ([]models.User, error) {
	session := state.db.Session.Copy()
	defer session.Close()
	query := session.DB(state.db.Database).C("users").Find(options.Filter)

	if len(options.Sort) > 0 {
		query = query.Sort(options.Sort...)
	}

	if options.Limit > 0 {
		query = query.Limit(int(options.Limit))
	}

	var users []models.User
	err := query.All(&users)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (state UserManager) Has(options utils.HasOptions) (bool, error) {
	session := state.db.Session.Copy()
	defer session.Close()
	query := session.DB(state.db.Database).C("users").Find(options.Filter)

	count, err := query.Count()
	if err != nil {
		return false, err
	}
	return count >= 1, nil
}

func (state UserManager) Create(c interface{}) (interface{}, error) {
	session := state.db.Session.Copy()
	defer session.Close()

	u := c.(*models.User)

	if u.ID == "" {
		u.ID = bson.NewObjectId()
	}

	err := session.DB(state.db.Database).C("users").Insert(u)
	if err != nil {
		return nil, err
	}

	var users models.User
	err = session.DB(state.db.Database).C("users").FindId(u.ID).One(&users)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (state UserManager) Update(options utils.UpdateOptions) (interface{}, error) {
	session := state.db.Session.Copy()
	defer session.Close()

	err := session.DB(state.db.Database).C("users").Update(options.Filter, options.Updates)
	if err != nil {
		return nil, err
	}

	return state.Find(utils.FindOptions{Filter: options.Filter})
}

func (state UserManager) Delete(options utils.DeleteOptions) error {
	session := state.db.Session.Copy()
	defer session.Close()

	_, err := state.Update(utils.UpdateOptions{
		Filter:  options.Filter,
		Updates: utils.Updates{"$set": bson.M{"updated_at": time.Now()}},
	})
	if err != nil {
		return err
	}

	return session.DB(state.db.Database).C("users").Remove(options.Filter)
}
