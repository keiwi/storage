package database

// This file is generated - do not edit.

import (
	"time"

	"github.com/keiwi/utils"
	"github.com/keiwi/utils/models"
	"gopkg.in/mgo.v2/bson"
)

// ServerManager struct
type ServerManager struct {
	db *DB
}

// NewServerManager - Creates a new *ServerManager that can be used for managing servers.
func NewServerManager(db *DB) (*ServerManager, error) {
	manager := ServerManager{}
	manager.db = db
	return &manager, nil
}

func (state ServerManager) Find(options utils.FindOptions) ([]models.Server, error) {
	session := state.db.Session.Copy()
	defer session.Close()
	query := session.DB(state.db.Database).C("servers").Find(options.Filter)

	if len(options.Sort) > 0 {
		query = query.Sort(options.Sort...)
	}

	if options.Limit > 0 {
		query = query.Limit(int(options.Limit))
	}

	var servers []models.Server
	err := query.All(&servers)
	if err != nil {
		return nil, err
	}

	return servers, nil
}

func (state ServerManager) Has(options utils.HasOptions) (bool, error) {
	session := state.db.Session.Copy()
	defer session.Close()
	query := session.DB(state.db.Database).C("servers").Find(options.Filter)

	count, err := query.Count()
	if err != nil {
		return false, err
	}
	return count >= 1, nil
}

func (state ServerManager) Create(c interface{}) (interface{}, error) {
	session := state.db.Session.Copy()
	defer session.Close()

	u := c.(*models.Server)

	if u.ID == "" {
		u.ID = bson.NewObjectId()
	}

	err := session.DB(state.db.Database).C("servers").Insert(u)
	if err != nil {
		return nil, err
	}

	var servers models.Server
	err = session.DB(state.db.Database).C("servers").FindId(u.ID).One(&servers)
	if err != nil {
		return nil, err
	}
	return servers, nil
}

func (state ServerManager) Update(options utils.UpdateOptions) (interface{}, error) {
	session := state.db.Session.Copy()
	defer session.Close()

	err := session.DB(state.db.Database).C("servers").Update(options.Filter, options.Updates)
	if err != nil {
		return nil, err
	}

	return state.Find(utils.FindOptions{Filter: options.Filter})
}

func (state ServerManager) Delete(options utils.DeleteOptions) error {
	session := state.db.Session.Copy()
	defer session.Close()

	_, err := state.Update(utils.UpdateOptions{
		Filter:  options.Filter,
		Updates: utils.Updates{"$set": bson.M{"updated_at": time.Now()}},
	})
	if err != nil {
		return err
	}

	return session.DB(state.db.Database).C("servers").Remove(options.Filter)
}
