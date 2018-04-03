package database

// This file is generated - do not edit.

import (
	"time"

	"github.com/keiwi/utils"
	"github.com/keiwi/utils/models"
	"gopkg.in/mgo.v2/bson"
)

// CheckManager struct
type CheckManager struct {
	db *DB
}

// NewCheckManager - Creates a new *CheckManager that can be used for managing checks.
func NewCheckManager(db *DB) (*CheckManager, error) {
	manager := CheckManager{}
	manager.db = db
	return &manager, nil
}

func (state CheckManager) Find(options utils.FindOptions) ([]models.Check, error) {
	session := state.db.Session.Copy()
	defer session.Close()
	query := session.DB(state.db.Database).C("checks").Find(options.Filter)

	if len(options.Sort) > 0 {
		query = query.Sort(options.Sort...)
	}

	if options.Limit > 0 {
		query = query.Limit(int(options.Limit))
	}

	var checks []models.Check
	err := query.All(&checks)
	if err != nil {
		return nil, err
	}

	return checks, nil
}

func (state CheckManager) Has(options utils.HasOptions) (bool, error) {
	session := state.db.Session.Copy()
	defer session.Close()
	query := session.DB(state.db.Database).C("checks").Find(options.Filter)

	count, err := query.Count()
	if err != nil {
		return false, err
	}
	return count >= 1, nil
}

func (state CheckManager) Create(c interface{}) (interface{}, error) {
	session := state.db.Session.Copy()
	defer session.Close()

	u := c.(*models.Check)

	if u.ID == "" {
		u.ID = bson.NewObjectId()
	}

	err := session.DB(state.db.Database).C("checks").Insert(u)
	if err != nil {
		return nil, err
	}

	var checks models.Check
	err = session.DB(state.db.Database).C("checks").FindId(u.ID).One(&checks)
	if err != nil {
		return nil, err
	}
	return checks, nil
}

func (state CheckManager) Update(options utils.UpdateOptions) (interface{}, error) {
	session := state.db.Session.Copy()
	defer session.Close()

	err := session.DB(state.db.Database).C("checks").Update(options.Filter, options.Updates)
	if err != nil {
		return nil, err
	}

	return state.Find(utils.FindOptions{Filter: options.Filter})
}

func (state CheckManager) Delete(options utils.DeleteOptions) error {
	session := state.db.Session.Copy()
	defer session.Close()

	_, err := state.Update(utils.UpdateOptions{
		Filter:  options.Filter,
		Updates: utils.Updates{"$set": bson.M{"updated_at": time.Now()}},
	})
	if err != nil {
		return err
	}

	return session.DB(state.db.Database).C("checks").Remove(options.Filter)
}
