package database

// This file is generated - do not edit.

import (
	"time"

	"github.com/keiwi/utils"
	"github.com/keiwi/utils/models"
	"gopkg.in/mgo.v2/bson"
)

// ClientManager struct
type ClientManager struct {
	db *DB
}

// NewClientManager - Creates a new *ClientManager that can be used for managing clients.
func NewClientManager(db *DB) (*ClientManager, error) {
	manager := ClientManager{}
	manager.db = db
	return &manager, nil
}

func (state ClientManager) Find(options utils.FindOptions) ([]models.Client, error) {
	session := state.db.Session.Copy()
	defer session.Close()
	query := session.DB(state.db.Database).C("clients").Find(options.Filter)

	if len(options.Sort) > 0 {
		query = query.Sort(options.Sort...)
	}

	if options.Limit > 0 {
		query = query.Limit(int(options.Limit))
	}

	var clients []models.Client
	err := query.All(&clients)
	if err != nil {
		return nil, err
	}

	return clients, nil
}

func (state ClientManager) Has(options utils.HasOptions) (bool, error) {
	session := state.db.Session.Copy()
	defer session.Close()
	query := session.DB(state.db.Database).C("clients").Find(options.Filter)

	count, err := query.Count()
	if err != nil {
		return false, err
	}
	return count >= 1, nil
}

func (state ClientManager) Create(c interface{}) (interface{}, error) {
	session := state.db.Session.Copy()
	defer session.Close()

	u := c.(*models.Client)

	if u.ID == "" {
		u.ID = bson.NewObjectId()
	}

	err := session.DB(state.db.Database).C("clients").Insert(u)
	if err != nil {
		return nil, err
	}

	var clients models.Client
	err = session.DB(state.db.Database).C("clients").FindId(u.ID).One(&clients)
	if err != nil {
		return nil, err
	}
	return clients, nil
}

func (state ClientManager) Update(options utils.UpdateOptions) (interface{}, error) {
	session := state.db.Session.Copy()
	defer session.Close()

	err := session.DB(state.db.Database).C("clients").Update(options.Filter, options.Updates)
	if err != nil {
		return nil, err
	}

	return state.Find(utils.FindOptions{Filter: options.Filter})
}

func (state ClientManager) Delete(options utils.DeleteOptions) error {
	session := state.db.Session.Copy()
	defer session.Close()

	_, err := state.Update(utils.UpdateOptions{
		Filter:  options.Filter,
		Updates: utils.Updates{"$set": bson.M{"updated_at": time.Now()}},
	})
	if err != nil {
		return err
	}

	return session.DB(state.db.Database).C("clients").Remove(options.Filter)
}
