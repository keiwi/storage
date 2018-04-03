package database

// This file is generated - do not edit.

import (
	"time"

	"github.com/keiwi/utils"
	"github.com/keiwi/utils/models"
	"gopkg.in/mgo.v2/bson"
)

// AlertManager struct
type AlertManager struct {
	db *DB
}

// NewAlertManager - Creates a new *AlertManager that can be used for managing alerts.
func NewAlertManager(db *DB) (*AlertManager, error) {
	manager := AlertManager{}
	manager.db = db
	return &manager, nil
}

func (state AlertManager) Find(options utils.FindOptions) ([]models.Alert, error) {
	session := state.db.Session.Copy()
	defer session.Close()
	query := session.DB(state.db.Database).C("alerts").Find(options.Filter)

	if len(options.Sort) > 0 {
		query = query.Sort(options.Sort...)
	}

	if options.Limit > 0 {
		query = query.Limit(int(options.Limit))
	}

	var alerts []models.Alert
	err := query.All(&alerts)
	if err != nil {
		return nil, err
	}

	return alerts, nil
}

func (state AlertManager) Has(options utils.HasOptions) (bool, error) {
	session := state.db.Session.Copy()
	defer session.Close()
	query := session.DB(state.db.Database).C("alerts").Find(options.Filter)

	count, err := query.Count()
	if err != nil {
		return false, err
	}
	return count >= 1, nil
}

func (state AlertManager) Create(c interface{}) (interface{}, error) {
	session := state.db.Session.Copy()
	defer session.Close()

	u := c.(*models.Alert)

	if u.ID == "" {
		u.ID = bson.NewObjectId()
	}

	err := session.DB(state.db.Database).C("alerts").Insert(u)
	if err != nil {
		return nil, err
	}

	var alerts models.Alert
	err = session.DB(state.db.Database).C("alerts").FindId(u.ID).One(&alerts)
	if err != nil {
		return nil, err
	}
	return alerts, nil
}

func (state AlertManager) Update(options utils.UpdateOptions) (interface{}, error) {
	session := state.db.Session.Copy()
	defer session.Close()

	err := session.DB(state.db.Database).C("alerts").Update(options.Filter, options.Updates)
	if err != nil {
		return nil, err
	}

	return state.Find(utils.FindOptions{Filter: options.Filter})
}

func (state AlertManager) Delete(options utils.DeleteOptions) error {
	session := state.db.Session.Copy()
	defer session.Close()

	_, err := state.Update(utils.UpdateOptions{
		Filter:  options.Filter,
		Updates: utils.Updates{"$set": bson.M{"updated_at": time.Now()}},
	})
	if err != nil {
		return err
	}

	return session.DB(state.db.Database).C("alerts").Remove(options.Filter)
}
