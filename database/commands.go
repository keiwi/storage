package database

// This file is generated - do not edit.

import (
	"time"

	"github.com/keiwi/utils"
	"github.com/keiwi/utils/models"
	"gopkg.in/mgo.v2/bson"
)

// CommandManager struct
type CommandManager struct {
	db *DB
}

// NewCommandManager - Creates a new *CommandManager that can be used for managing commands.
func NewCommandManager(db *DB) (*CommandManager, error) {
	manager := CommandManager{}
	manager.db = db
	return &manager, nil
}

func (state CommandManager) Find(options utils.FindOptions) ([]models.Command, error) {
	session := state.db.Session.Copy()
	defer session.Close()
	query := session.DB(state.db.Database).C("commands").Find(options.Filter)

	if len(options.Sort) > 0 {
		query = query.Sort(options.Sort...)
	}

	if options.Limit > 0 {
		query = query.Limit(int(options.Limit))
	}

	var commands []models.Command
	err := query.All(&commands)
	if err != nil {
		return nil, err
	}

	return commands, nil
}

func (state CommandManager) Has(options utils.HasOptions) (bool, error) {
	session := state.db.Session.Copy()
	defer session.Close()
	query := session.DB(state.db.Database).C("commands").Find(options.Filter)

	count, err := query.Count()
	if err != nil {
		return false, err
	}
	return count >= 1, nil
}

func (state CommandManager) Create(c interface{}) (interface{}, error) {
	session := state.db.Session.Copy()
	defer session.Close()

	u := c.(*models.Command)

	if u.ID == "" {
		u.ID = bson.NewObjectId()
	}

	err := session.DB(state.db.Database).C("commands").Insert(u)
	if err != nil {
		return nil, err
	}

	var commands models.Command
	err = session.DB(state.db.Database).C("commands").FindId(u.ID).One(&commands)
	if err != nil {
		return nil, err
	}
	return commands, nil
}

func (state CommandManager) Update(options utils.UpdateOptions) (interface{}, error) {
	session := state.db.Session.Copy()
	defer session.Close()

	err := session.DB(state.db.Database).C("commands").Update(options.Filter, options.Updates)
	if err != nil {
		return nil, err
	}

	return state.Find(utils.FindOptions{Filter: options.Filter})
}

func (state CommandManager) Delete(options utils.DeleteOptions) error {
	session := state.db.Session.Copy()
	defer session.Close()

	_, err := state.Update(utils.UpdateOptions{
		Filter:  options.Filter,
		Updates: utils.Updates{"$set": bson.M{"updated_at": time.Now()}},
	})
	if err != nil {
		return err
	}

	return session.DB(state.db.Database).C("commands").Remove(options.Filter)
}
