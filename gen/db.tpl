package database

// This file is generated - do not edit.

import(
    "fmt"
    "math"

    "github.com/keiwi/utils"
    "github.com/keiwi/utils/log"
    "github.com/keiwi/utils/models"
    "github.com/nats-io/go-nats"
    "github.com/pkg/errors"
    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
    "github.com/spf13/viper"
)

// DB abstraction
type DB struct {
    Session  *mgo.Session
    Database string
}

// Database struct
type Database struct {
    db       *DB
    Conn     *nats.Conn

    {{range .}}
        {{.Name}}Manager *{{.Name}}Manager{{end}}
}

// NewDatabase - rethinkdb database
func NewDatabase(user, password, host, port, dbname string) (*Database, error) {
    session, err := mgo.Dial(fmt.Sprintf("%s:%s", host, port))
    if err != nil {
        return nil, err
    }

    session.SetMode(mgo.Monotonic, true)

    db := &DB{session, dbname}

    database := &Database{db: db}

    {{range .}}
        {{.ID}}mgr, err := New{{.Name}}Manager(db)
        if err != nil {
            return nil, err
        }
        database.{{.Name}}Manager = {{.ID}}mgr
    {{end}}

    conn, err := nats.Connect(viper.GetString("nats.url"))
    if err != nil {
        return nil, err
    }
    database.Conn = conn

    return database, nil
}


func (d *Database) Listen() {
    {{ range . }}
        d.Conn.Subscribe("{{.ID}}.retrieve.find", func(m *nats.Msg) {
            var find utils.FindOptions
            err := bson.UnmarshalJSON(m.Data, &find)
            if err != nil {
                log.WithError(errors.Wrap(err, "error decoding event")).Error("{{.ID}}.retrieve.find")
                return
            }
            log.Info("{{.ID}}.retrieve.find")

            manager := d.{{.Name}}Manager

            inter, err := manager.Find(find)
            if err != nil {
                log.WithError(errors.Wrap(err, "error retrieving data")).Error("{{.ID}}.retrieve.find")
                return
            }

            if find.Max > 0 {
            var n []models.{{.Name}}
                step := int(math.Ceil(float64(len(inter)) / float64(find.Max)))
                for i := len(inter) - 1; i >= 0; i -= step {
                    n = append(n, inter[i])
                }
                inter = n
            }
            log.Info("{{.ID}}.retrieve.find - Retrieved")

            data, err := bson.MarshalJSON(inter)
            if err != nil {
                log.WithError(errors.Wrap(err, "error marshaling event")).Error("{{.ID}}.retrieve.find")
                return
            }
            log.Info("{{.ID}}.retrieve.find - Reply")
            d.Conn.Publish(m.Reply, data)
        })

        d.Conn.Subscribe("{{.ID}}.retrieve.has", func(m *nats.Msg) {
            var has utils.HasOptions
            err := bson.UnmarshalJSON(m.Data, &has)
            if err != nil {
                log.WithError(errors.Wrap(err, "error decoding event")).Error("{{.ID}}.retrieve.has")
                return
            }
            log.Info("{{.ID}}.retrieve.has")

            manager := d.{{.Name}}Manager
            inter, err := manager.Has(has)
            if err != nil {
                log.WithError(err).Error("{{.ID}}.retrieve.has")
                return
            }
            log.Info("{{.ID}}.retrieve.has - Retrieved")

            data, err := bson.MarshalJSON(inter)
            if err != nil {
                log.WithError(errors.Wrap(err, "error marshaling event")).Error("{{.ID}}.retrieve.has")
                return
            }
            log.Info("{{.ID}}.retrieve.has - Reply")
            d.Conn.Publish(m.Reply, data)
        })

        d.Conn.Subscribe("{{.ID}}.delete.send", func(m *nats.Msg) {
            var del utils.DeleteOptions
            err := bson.UnmarshalJSON(m.Data, &del)
            if err != nil {
                log.WithError(errors.Wrap(err, "error decoding event")).Error("{{.ID}}.delete")
                return
            }
            log.Info("{{.ID}}.delete")

            d.Conn.Publish("{{.ID}}.delete.before", m.Data)

            manager := d.{{.Name}}Manager

            inter, err := manager.Find(utils.FindOptions{Filter: del.Filter})
            if err != nil {
                log.WithError(errors.Wrap(err, "error finding data")).Error("{{.ID}}.delete")
                return
            }
            log.Info("{{.ID}}.delete - Found")

            err = manager.Delete(del)
            if err != nil {
                log.WithError(errors.Wrap(err, "error deleting data")).Error("{{.ID}}.delete")
                return
            }
            log.Info("{{.ID}}.delete - Deleted")


            data, err := bson.MarshalJSON(inter)
            if err != nil {
                log.WithError(errors.Wrap(err, "error marshaling event")).Error("{{.ID}}.delete")
                return
            }
            log.Info("{{.ID}}.delete - Reply")

            d.Conn.Publish("{{.ID}}.delete.after", data)
        })

        d.Conn.Subscribe("{{.ID}}.update.send", func(m *nats.Msg) {
            var update utils.UpdateOptions
            err := bson.UnmarshalJSON(m.Data, &update)
            if err != nil {
                log.WithError(errors.Wrap(err, "error decoding event")).Error("{{.ID}}.update")
                return
            }
            log.Info("{{.ID}}.update")

            d.Conn.Publish("{{.ID}}.update.before", m.Data)

            manager := d.{{.Name}}Manager

            inter, err := manager.Update(update)
            if err != nil {
                log.WithError(errors.Wrap(err, "error updating data")).Error("{{.ID}}.update")
                return
            }
            log.Info("{{.ID}}.update - Updated")

            data, err := bson.MarshalJSON(inter)
            if err != nil {
                log.WithError(errors.Wrap(err, "error marshaling event")).Error("{{.ID}}.update")
                return
            }
            log.Info("{{.ID}}.update - Reply")

            d.Conn.Publish("{{.ID}}.update.after", data)
        })

        d.Conn.Subscribe("{{.ID}}.create.send", func(m *nats.Msg) {
            var before *models.{{.Name}}
            err := bson.UnmarshalJSON(m.Data, &before)
            if err != nil {
                log.WithError(errors.Wrap(err, "error decoding event")).Error("{{.ID}}.create")
                return
            }
            log.Info("{{.ID}}.create")

            d.Conn.Publish("{{.ID}}.create.before", m.Data)

            manager := d.{{.Name}}Manager

            inter, err := manager.Create(before)
            if err != nil {
                log.WithError(errors.Wrap(err, "error creating data")).Error("{{.ID}}.create")
                return
            }
            log.Info("{{.ID}}.create - Created")

            data, err := bson.MarshalJSON(inter)
            if err != nil {
                log.WithError(errors.Wrap(err, "error marshaling event")).Error("{{.ID}}.create")
                return
            }
            log.Info("{{.ID}}.create - Reply")

            d.Conn.Publish("{{.ID}}.create.after", data)
        })
    {{ end }}
}

func (d *Database) Close() error {
    if d.db.Session != nil {
        d.db.Session.Close()
    }
    if d.Conn != nil {
        d.Conn.Close()
    }
    return nil
}