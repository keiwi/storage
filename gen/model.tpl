package database

// This file is generated - do not edit.

import (
    "time"

    "github.com/keiwi/utils"
    "github.com/keiwi/utils/models"
    "gopkg.in/mgo.v2/bson"
)

// {{.Name}}Manager struct
type {{.Name}}Manager struct {
    db *DB
}

// New{{.Name}}Manager - Creates a new *{{.Name}}Manager that can be used for managing {{.ID}}.
func New{{.Name}}Manager(db *DB) (*{{.Name}}Manager, error) {
    manager := {{.Name}}Manager{}
    manager.db = db
    return &manager, nil
}

func (state {{.Name}}Manager) Find(options utils.FindOptions) ([]models.{{.Name}}, error) {
    session := state.db.Session.Copy()
    defer session.Close()
    query := session.DB(state.db.Database).C("{{.ID}}").Find(options.Filter)

    if len(options.Sort) > 0 {
        query = query.Sort(options.Sort...)
    }

    if options.Limit > 0 {
        query = query.Limit(int(options.Limit))
    }

    var {{.ID}} []models.{{.Name}}
    err := query.All(&{{.ID}})
    if err != nil {
        return nil, err
    }

    return {{.ID}}, nil
}

func (state {{.Name}}Manager) Has(options utils.HasOptions) (bool, error) {
    session := state.db.Session.Copy()
    defer session.Close()
    query := session.DB(state.db.Database).C("{{.ID}}").Find(options.Filter)

    count, err := query.Count()
    if err != nil {
        return false, err
    }
    return count >= 1, nil
}

func (state {{.Name}}Manager) Create(c interface{}) (interface{}, error) {
    session := state.db.Session.Copy()
    defer session.Close()

    u := c.(*models.{{.Name}})

    if u.ID == "" {
        u.ID = bson.NewObjectId()
    }

    err := session.DB(state.db.Database).C("{{.ID}}").Insert(u)
    if err != nil {
        return nil, err
    }

    var {{.ID}} models.{{.Name}}
    err = session.DB(state.db.Database).C("{{.ID}}").FindId(u.ID).One(&{{.ID}})
    if err != nil {
        return nil, err
    }
    return {{.ID}}, nil
}

func (state {{.Name}}Manager) Update(options utils.UpdateOptions) (interface{}, error) {
    session := state.db.Session.Copy()
    defer session.Close()

    err := session.DB(state.db.Database).C("{{.ID}}").Update(options.Filter, options.Updates)
    if err != nil {
        return nil, err
    }

    return state.Find(utils.FindOptions{Filter: options.Filter})
}

func (state {{.Name}}Manager) Delete(options utils.DeleteOptions) error {
    session := state.db.Session.Copy()
    defer session.Close()

    _, err := state.Update(utils.UpdateOptions{
        Filter:  options.Filter,
        Updates: utils.Updates{"$set": bson.M{"updated_at": time.Now()}},
    })
    if err != nil {
        return err
    }

    return session.DB(state.db.Database).C("{{.ID}}").Remove(options.Filter)
}