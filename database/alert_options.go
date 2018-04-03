package database

// This file is generated - do not edit.

import (
	"time"

	"github.com/keiwi/utils"
	"github.com/keiwi/utils/models"
	"gopkg.in/mgo.v2/bson"
)

// AlertOptionManager struct
type AlertOptionManager struct {
	db *DB
}

// NewAlertOptionManager - Creates a new *AlertOptionManager that can be used for managing alert_options.
func NewAlertOptionManager(db *DB) (*AlertOptionManager, error) {
	manager := AlertOptionManager{}
	manager.db = db
	return &manager, nil
}

func (state AlertOptionManager) Find(options utils.FindOptions) ([]models.AlertOption, error) {
	session := state.db.Session.Copy()
	defer session.Close()
	query := session.DB(state.db.Database).C("alert_options").Find(options.Filter)

	if len(options.Sort) > 0 {
		query = query.Sort(options.Sort...)
	}

	if options.Limit > 0 {
		query = query.Limit(int(options.Limit))
	}

	var alert_options []models.AlertOption
	err := query.All(&alert_options)
	if err != nil {
		return nil, err
	}

	return alert_options, nil
}

func (state AlertOptionManager) Has(options utils.HasOptions) (bool, error) {
	session := state.db.Session.Copy()
	defer session.Close()
	query := session.DB(state.db.Database).C("alert_options").Find(options.Filter)

	count, err := query.Count()
	if err != nil {
		return false, err
	}
	return count >= 1, nil
}

func (state AlertOptionManager) Create(c interface{}) (interface{}, error) {
	session := state.db.Session.Copy()
	defer session.Close()

	u := c.(*models.AlertOption)

	if u.ID == "" {
		u.ID = bson.NewObjectId()
	}

	err := session.DB(state.db.Database).C("alert_options").Insert(u)
	if err != nil {
		return nil, err
	}

	var alert_options models.AlertOption
	err = session.DB(state.db.Database).C("alert_options").FindId(u.ID).One(&alert_options)
	if err != nil {
		return nil, err
	}
	return alert_options, nil
}

func (state AlertOptionManager) Update(options utils.UpdateOptions) (interface{}, error) {
	session := state.db.Session.Copy()
	defer session.Close()

	err := session.DB(state.db.Database).C("alert_options").Update(options.Filter, options.Updates)
	if err != nil {
		return nil, err
	}

	return state.Find(utils.FindOptions{Filter: options.Filter})
}

func (state AlertOptionManager) Delete(options utils.DeleteOptions) error {
	session := state.db.Session.Copy()
	defer session.Close()

	_, err := state.Update(utils.UpdateOptions{
		Filter:  options.Filter,
		Updates: utils.Updates{"$set": bson.M{"updated_at": time.Now()}},
	})
	if err != nil {
		return err
	}

	return session.DB(state.db.Database).C("alert_options").Remove(options.Filter)
}
