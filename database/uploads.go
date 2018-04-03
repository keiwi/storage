package database

// This file is generated - do not edit.

import (
	"time"

	"github.com/keiwi/utils"
	"github.com/keiwi/utils/models"
	"gopkg.in/mgo.v2/bson"
)

// UploadManager struct
type UploadManager struct {
	db *DB
}

// NewUploadManager - Creates a new *UploadManager that can be used for managing uploads.
func NewUploadManager(db *DB) (*UploadManager, error) {
	manager := UploadManager{}
	manager.db = db
	return &manager, nil
}

func (state UploadManager) Find(options utils.FindOptions) ([]models.Upload, error) {
	session := state.db.Session.Copy()
	defer session.Close()
	query := session.DB(state.db.Database).C("uploads").Find(options.Filter)

	if len(options.Sort) > 0 {
		query = query.Sort(options.Sort...)
	}

	if options.Limit > 0 {
		query = query.Limit(int(options.Limit))
	}

	var uploads []models.Upload
	err := query.All(&uploads)
	if err != nil {
		return nil, err
	}

	return uploads, nil
}

func (state UploadManager) Has(options utils.HasOptions) (bool, error) {
	session := state.db.Session.Copy()
	defer session.Close()
	query := session.DB(state.db.Database).C("uploads").Find(options.Filter)

	count, err := query.Count()
	if err != nil {
		return false, err
	}
	return count >= 1, nil
}

func (state UploadManager) Create(c interface{}) (interface{}, error) {
	session := state.db.Session.Copy()
	defer session.Close()

	u := c.(*models.Upload)

	if u.ID == "" {
		u.ID = bson.NewObjectId()
	}

	err := session.DB(state.db.Database).C("uploads").Insert(u)
	if err != nil {
		return nil, err
	}

	var uploads models.Upload
	err = session.DB(state.db.Database).C("uploads").FindId(u.ID).One(&uploads)
	if err != nil {
		return nil, err
	}
	return uploads, nil
}

func (state UploadManager) Update(options utils.UpdateOptions) (interface{}, error) {
	session := state.db.Session.Copy()
	defer session.Close()

	err := session.DB(state.db.Database).C("uploads").Update(options.Filter, options.Updates)
	if err != nil {
		return nil, err
	}

	return state.Find(utils.FindOptions{Filter: options.Filter})
}

func (state UploadManager) Delete(options utils.DeleteOptions) error {
	session := state.db.Session.Copy()
	defer session.Close()

	_, err := state.Update(utils.UpdateOptions{
		Filter:  options.Filter,
		Updates: utils.Updates{"$set": bson.M{"updated_at": time.Now()}},
	})
	if err != nil {
		return err
	}

	return session.DB(state.db.Database).C("uploads").Remove(options.Filter)
}
