package database

// This file is generated - do not edit.

import (
	"fmt"
	"math"

	"github.com/keiwi/utils"
	"github.com/keiwi/utils/models"
	"github.com/nats-io/go-nats"
	"github.com/pkg/errors"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// DB abstraction
type DB struct {
	Session  *mgo.Session
	Database string
}

// Database struct
type Database struct {
	db   *DB
	Conn *nats.Conn

	AlertManager       *AlertManager
	AlertOptionManager *AlertOptionManager
	CheckManager       *CheckManager
	ClientManager      *ClientManager
	CommandManager     *CommandManager
	GroupManager       *GroupManager
	ServerManager      *ServerManager
	UploadManager      *UploadManager
	UserManager        *UserManager
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

	alertsmgr, err := NewAlertManager(db)
	if err != nil {
		return nil, err
	}
	database.AlertManager = alertsmgr

	alert_optionsmgr, err := NewAlertOptionManager(db)
	if err != nil {
		return nil, err
	}
	database.AlertOptionManager = alert_optionsmgr

	checksmgr, err := NewCheckManager(db)
	if err != nil {
		return nil, err
	}
	database.CheckManager = checksmgr

	clientsmgr, err := NewClientManager(db)
	if err != nil {
		return nil, err
	}
	database.ClientManager = clientsmgr

	commandsmgr, err := NewCommandManager(db)
	if err != nil {
		return nil, err
	}
	database.CommandManager = commandsmgr

	groupsmgr, err := NewGroupManager(db)
	if err != nil {
		return nil, err
	}
	database.GroupManager = groupsmgr

	serversmgr, err := NewServerManager(db)
	if err != nil {
		return nil, err
	}
	database.ServerManager = serversmgr

	uploadsmgr, err := NewUploadManager(db)
	if err != nil {
		return nil, err
	}
	database.UploadManager = uploadsmgr

	usersmgr, err := NewUserManager(db)
	if err != nil {
		return nil, err
	}
	database.UserManager = usersmgr

	conn, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		return nil, err
	}
	database.Conn = conn

	return database, nil
}

func (d *Database) Listen() {

	d.Conn.Subscribe("alerts.retrieve.find", func(m *nats.Msg) {
		var find utils.FindOptions
		err := bson.UnmarshalJSON(m.Data, &find)
		if err != nil {
			utils.Log.WithError(errors.Wrap(err, "error decoding event")).Error("alerts.retrieve.find")
			return
		}
		utils.Log.Info("alerts.retrieve.find")

		manager := d.AlertManager

		inter, err := manager.Find(find)
		if err != nil {
			utils.Log.WithError(errors.Wrap(err, "error retrieving data")).Error("alerts.retrieve.find")
			return
		}

		if find.Max > 0 {
			var n []models.Alert
			step := int(math.Ceil(float64(len(inter)) / float64(find.Max)))
			for i := len(inter) - 1; i >= 0; i -= step {
				n = append(n, inter[i])
			}
			inter = n
		}
		utils.Log.Info("alerts.retrieve.find - Retrieved")

		data, err := bson.MarshalJSON(inter)
		if err != nil {
			utils.Log.WithError(errors.Wrap(err, "error marshaling event")).Error("alerts.retrieve.find")
			return
		}
		utils.Log.Info("alerts.retrieve.find - Reply")
		d.Conn.Publish(m.Reply, data)
	})

	d.Conn.Subscribe("alerts.retrieve.has", func(m *nats.Msg) {
		var has utils.HasOptions
		err := bson.UnmarshalJSON(m.Data, &has)
		if err != nil {
			utils.Log.WithError(errors.Wrap(err, "error decoding event")).Error("alerts.retrieve.has")
			return
		}
		utils.Log.Info("alerts.retrieve.has")

		manager := d.AlertManager
		inter, err := manager.Has(has)
		if err != nil {
			utils.Log.WithError(err).Error("alerts.retrieve.has")
			return
		}
		utils.Log.Info("alerts.retrieve.has - Retrieved")

		data, err := bson.MarshalJSON(inter)
		if err != nil {
			utils.Log.WithError(errors.Wrap(err, "error marshaling event")).Error("alerts.retrieve.has")
			return
		}
		utils.Log.Info("alerts.retrieve.has - Reply")
		d.Conn.Publish(m.Reply, data)
	})

	d.Conn.Subscribe("alerts.delete.send", func(m *nats.Msg) {
		var del utils.DeleteOptions
		err := bson.UnmarshalJSON(m.Data, &del)
		if err != nil {
			utils.Log.WithError(errors.Wrap(err, "error decoding event")).Error("alerts.delete")
			return
		}
		utils.Log.Info("alerts.delete")

		d.Conn.Publish("alerts.delete.before", m.Data)

		manager := d.AlertManager

		inter, err := manager.Find(utils.FindOptions{Filter: del.Filter})
		if err != nil {
			utils.Log.WithError(errors.Wrap(err, "error finding data")).Error("alerts.delete")
			return
		}
		utils.Log.Info("alerts.delete - Found")

		err = manager.Delete(del)
		if err != nil {
			utils.Log.WithError(errors.Wrap(err, "error deleting data")).Error("alerts.delete")
			return
		}
		utils.Log.Info("alerts.delete - Deleted")

		data, err := bson.MarshalJSON(inter)
		if err != nil {
			utils.Log.WithError(errors.Wrap(err, "error marshaling event")).Error("alerts.delete")
			return
		}
		utils.Log.Info("alerts.delete - Reply")

		d.Conn.Publish("alerts.delete.after", data)
	})

	d.Conn.Subscribe("alerts.update.send", func(m *nats.Msg) {
		var update utils.UpdateOptions
		err := bson.UnmarshalJSON(m.Data, &update)
		if err != nil {
			utils.Log.WithError(errors.Wrap(err, "error decoding event")).Error("alerts.update")
			return
		}
		utils.Log.Info("alerts.update")

		d.Conn.Publish("alerts.update.before", m.Data)

		manager := d.AlertManager

		inter, err := manager.Update(update)
		if err != nil {
			utils.Log.WithError(errors.Wrap(err, "error updating data")).Error("alerts.update")
			return
		}
		utils.Log.Info("alerts.update - Updated")

		data, err := bson.MarshalJSON(inter)
		if err != nil {
			utils.Log.WithError(errors.Wrap(err, "error marshaling event")).Error("alerts.update")
			return
		}
		utils.Log.Info("alerts.update - Reply")

		d.Conn.Publish("alerts.update.after", data)
	})

	d.Conn.Subscribe("alerts.create.send", func(m *nats.Msg) {
		var before *models.Alert
		err := bson.UnmarshalJSON(m.Data, &before)
		if err != nil {
			utils.Log.WithError(errors.Wrap(err, "error decoding event")).Error("alerts.create")
			return
		}
		utils.Log.Info("alerts.create")

		d.Conn.Publish("alerts.create.before", m.Data)

		manager := d.AlertManager

		inter, err := manager.Create(before)
		if err != nil {
			utils.Log.WithError(errors.Wrap(err, "error creating data")).Error("alerts.create")
			return
		}
		utils.Log.Info("alerts.create - Created")

		data, err := bson.MarshalJSON(inter)
		if err != nil {
			utils.Log.WithError(errors.Wrap(err, "error marshaling event")).Error("alerts.create")
			return
		}
		utils.Log.Info("alerts.create - Reply")

		d.Conn.Publish("alerts.create.after", data)
	})

	d.Conn.Subscribe("alert_options.retrieve.find", func(m *nats.Msg) {
		var find utils.FindOptions
		err := bson.UnmarshalJSON(m.Data, &find)
		if err != nil {
			utils.Log.WithError(errors.Wrap(err, "error decoding event")).Error("alert_options.retrieve.find")
			return
		}
		utils.Log.Info("alert_options.retrieve.find")

		manager := d.AlertOptionManager

		inter, err := manager.Find(find)
		if err != nil {
			utils.Log.WithError(errors.Wrap(err, "error retrieving data")).Error("alert_options.retrieve.find")
			return
		}

		if find.Max > 0 {
			var n []models.AlertOption
			step := int(math.Ceil(float64(len(inter)) / float64(find.Max)))
			for i := len(inter) - 1; i >= 0; i -= step {
				n = append(n, inter[i])
			}
			inter = n
		}
		utils.Log.Info("alert_options.retrieve.find - Retrieved")

		data, err := bson.MarshalJSON(inter)
		if err != nil {
			utils.Log.WithError(errors.Wrap(err, "error marshaling event")).Error("alert_options.retrieve.find")
			return
		}
		utils.Log.Info("alert_options.retrieve.find - Reply")
		d.Conn.Publish(m.Reply, data)
	})

	d.Conn.Subscribe("alert_options.retrieve.has", func(m *nats.Msg) {
		var has utils.HasOptions
		err := bson.UnmarshalJSON(m.Data, &has)
		if err != nil {
			utils.Log.WithError(errors.Wrap(err, "error decoding event")).Error("alert_options.retrieve.has")
			return
		}
		utils.Log.Info("alert_options.retrieve.has")

		manager := d.AlertOptionManager
		inter, err := manager.Has(has)
		if err != nil {
			utils.Log.WithError(err).Error("alert_options.retrieve.has")
			return
		}
		utils.Log.Info("alert_options.retrieve.has - Retrieved")

		data, err := bson.MarshalJSON(inter)
		if err != nil {
			utils.Log.WithError(errors.Wrap(err, "error marshaling event")).Error("alert_options.retrieve.has")
			return
		}
		utils.Log.Info("alert_options.retrieve.has - Reply")
		d.Conn.Publish(m.Reply, data)
	})

	d.Conn.Subscribe("alert_options.delete.send", func(m *nats.Msg) {
		var del utils.DeleteOptions
		err := bson.UnmarshalJSON(m.Data, &del)
		if err != nil {
			utils.Log.WithError(errors.Wrap(err, "error decoding event")).Error("alert_options.delete")
			return
		}
		utils.Log.Info("alert_options.delete")

		d.Conn.Publish("alert_options.delete.before", m.Data)

		manager := d.AlertOptionManager

		inter, err := manager.Find(utils.FindOptions{Filter: del.Filter})
		if err != nil {
			utils.Log.WithError(errors.Wrap(err, "error finding data")).Error("alert_options.delete")
			return
		}
		utils.Log.Info("alert_options.delete - Found")

		err = manager.Delete(del)
		if err != nil {
			utils.Log.WithError(errors.Wrap(err, "error deleting data")).Error("alert_options.delete")
			return
		}
		utils.Log.Info("alert_options.delete - Deleted")

		data, err := bson.MarshalJSON(inter)
		if err != nil {
			utils.Log.WithError(errors.Wrap(err, "error marshaling event")).Error("alert_options.delete")
			return
		}
		utils.Log.Info("alert_options.delete - Reply")

		d.Conn.Publish("alert_options.delete.after", data)
	})

	d.Conn.Subscribe("alert_options.update.send", func(m *nats.Msg) {
		var update utils.UpdateOptions
		err := bson.UnmarshalJSON(m.Data, &update)
		if err != nil {
			utils.Log.WithError(errors.Wrap(err, "error decoding event")).Error("alert_options.update")
			return
		}
		utils.Log.Info("alert_options.update")

		d.Conn.Publish("alert_options.update.before", m.Data)

		manager := d.AlertOptionManager

		inter, err := manager.Update(update)
		if err != nil {
			utils.Log.WithError(errors.Wrap(err, "error updating data")).Error("alert_options.update")
			return
		}
		utils.Log.Info("alert_options.update - Updated")

		data, err := bson.MarshalJSON(inter)
		if err != nil {
			utils.Log.WithError(errors.Wrap(err, "error marshaling event")).Error("alert_options.update")
			return
		}
		utils.Log.Info("alert_options.update - Reply")

		d.Conn.Publish("alert_options.update.after", data)
	})

	d.Conn.Subscribe("alert_options.create.send", func(m *nats.Msg) {
		var before *models.AlertOption
		err := bson.UnmarshalJSON(m.Data, &before)
		if err != nil {
			utils.Log.WithError(errors.Wrap(err, "error decoding event")).Error("alert_options.create")
			return
		}
		utils.Log.Info("alert_options.create")

		d.Conn.Publish("alert_options.create.before", m.Data)

		manager := d.AlertOptionManager

		inter, err := manager.Create(before)
		if err != nil {
			utils.Log.WithError(errors.Wrap(err, "error creating data")).Error("alert_options.create")
			return
		}
		utils.Log.Info("alert_options.create - Created")

		data, err := bson.MarshalJSON(inter)
		if err != nil {
			utils.Log.WithError(errors.Wrap(err, "error marshaling event")).Error("alert_options.create")
			return
		}
		utils.Log.Info("alert_options.create - Reply")

		d.Conn.Publish("alert_options.create.after", data)
	})

	d.Conn.Subscribe("checks.retrieve.find", func(m *nats.Msg) {
		var find utils.FindOptions
		err := bson.UnmarshalJSON(m.Data, &find)
		if err != nil {
			utils.Log.WithError(errors.Wrap(err, "error decoding event")).Error("checks.retrieve.find")
			return
		}
		utils.Log.Info("checks.retrieve.find")

		manager := d.CheckManager

		inter, err := manager.Find(find)
		if err != nil {
			utils.Log.WithError(errors.Wrap(err, "error retrieving data")).Error("checks.retrieve.find")
			return
		}

		if find.Max > 0 {
			var n []models.Check
			step := int(math.Ceil(float64(len(inter)) / float64(find.Max)))
			for i := len(inter) - 1; i >= 0; i -= step {
				n = append(n, inter[i])
			}
			inter = n
		}
		utils.Log.Info("checks.retrieve.find - Retrieved")

		data, err := bson.MarshalJSON(inter)
		if err != nil {
			utils.Log.WithError(errors.Wrap(err, "error marshaling event")).Error("checks.retrieve.find")
			return
		}
		utils.Log.Info("checks.retrieve.find - Reply")
		d.Conn.Publish(m.Reply, data)
	})

	d.Conn.Subscribe("checks.retrieve.has", func(m *nats.Msg) {
		var has utils.HasOptions
		err := bson.UnmarshalJSON(m.Data, &has)
		if err != nil {
			utils.Log.WithError(errors.Wrap(err, "error decoding event")).Error("checks.retrieve.has")
			return
		}
		utils.Log.Info("checks.retrieve.has")

		manager := d.CheckManager
		inter, err := manager.Has(has)
		if err != nil {
			utils.Log.WithError(err).Error("checks.retrieve.has")
			return
		}
		utils.Log.Info("checks.retrieve.has - Retrieved")

		data, err := bson.MarshalJSON(inter)
		if err != nil {
			utils.Log.WithError(errors.Wrap(err, "error marshaling event")).Error("checks.retrieve.has")
			return
		}
		utils.Log.Info("checks.retrieve.has - Reply")
		d.Conn.Publish(m.Reply, data)
	})

	d.Conn.Subscribe("checks.delete.send", func(m *nats.Msg) {
		var del utils.DeleteOptions
		err := bson.UnmarshalJSON(m.Data, &del)
		if err != nil {
			utils.Log.WithError(errors.Wrap(err, "error decoding event")).Error("checks.delete")
			return
		}
		utils.Log.Info("checks.delete")

		d.Conn.Publish("checks.delete.before", m.Data)

		manager := d.CheckManager

		inter, err := manager.Find(utils.FindOptions{Filter: del.Filter})
		if err != nil {
			utils.Log.WithError(errors.Wrap(err, "error finding data")).Error("checks.delete")
			return
		}
		utils.Log.Info("checks.delete - Found")

		err = manager.Delete(del)
		if err != nil {
			utils.Log.WithError(errors.Wrap(err, "error deleting data")).Error("checks.delete")
			return
		}
		utils.Log.Info("checks.delete - Deleted")

		data, err := bson.MarshalJSON(inter)
		if err != nil {
			utils.Log.WithError(errors.Wrap(err, "error marshaling event")).Error("checks.delete")
			return
		}
		utils.Log.Info("checks.delete - Reply")

		d.Conn.Publish("checks.delete.after", data)
	})

	d.Conn.Subscribe("checks.update.send", func(m *nats.Msg) {
		var update utils.UpdateOptions
		err := bson.UnmarshalJSON(m.Data, &update)
		if err != nil {
			utils.Log.WithError(errors.Wrap(err, "error decoding event")).Error("checks.update")
			return
		}
		utils.Log.Info("checks.update")

		d.Conn.Publish("checks.update.before", m.Data)

		manager := d.CheckManager

		inter, err := manager.Update(update)
		if err != nil {
			utils.Log.WithError(errors.Wrap(err, "error updating data")).Error("checks.update")
			return
		}
		utils.Log.Info("checks.update - Updated")

		data, err := bson.MarshalJSON(inter)
		if err != nil {
			utils.Log.WithError(errors.Wrap(err, "error marshaling event")).Error("checks.update")
			return
		}
		utils.Log.Info("checks.update - Reply")

		d.Conn.Publish("checks.update.after", data)
	})

	d.Conn.Subscribe("checks.create.send", func(m *nats.Msg) {
		var before *models.Check
		err := bson.UnmarshalJSON(m.Data, &before)
		if err != nil {
			utils.Log.WithError(errors.Wrap(err, "error decoding event")).Error("checks.create")
			return
		}
		utils.Log.Info("checks.create")

		d.Conn.Publish("checks.create.before", m.Data)

		manager := d.CheckManager

		inter, err := manager.Create(before)
		if err != nil {
			utils.Log.WithError(errors.Wrap(err, "error creating data")).Error("checks.create")
			return
		}
		utils.Log.Info("checks.create - Created")

		data, err := bson.MarshalJSON(inter)
		if err != nil {
			utils.Log.WithError(errors.Wrap(err, "error marshaling event")).Error("checks.create")
			return
		}
		utils.Log.Info("checks.create - Reply")

		d.Conn.Publish("checks.create.after", data)
	})

	d.Conn.Subscribe("clients.retrieve.find", func(m *nats.Msg) {
		var find utils.FindOptions
		err := bson.UnmarshalJSON(m.Data, &find)
		if err != nil {
			utils.Log.WithError(errors.Wrap(err, "error decoding event")).Error("clients.retrieve.find")
			return
		}
		utils.Log.Info("clients.retrieve.find")

		manager := d.ClientManager

		inter, err := manager.Find(find)
		if err != nil {
			utils.Log.WithError(errors.Wrap(err, "error retrieving data")).Error("clients.retrieve.find")
			return
		}

		if find.Max > 0 {
			var n []models.Client
			step := int(math.Ceil(float64(len(inter)) / float64(find.Max)))
			for i := len(inter) - 1; i >= 0; i -= step {
				n = append(n, inter[i])
			}
			inter = n
		}
		utils.Log.Info("clients.retrieve.find - Retrieved")

		data, err := bson.MarshalJSON(inter)
		if err != nil {
			utils.Log.WithError(errors.Wrap(err, "error marshaling event")).Error("clients.retrieve.find")
			return
		}
		utils.Log.Info("clients.retrieve.find - Reply")
		d.Conn.Publish(m.Reply, data)
	})

	d.Conn.Subscribe("clients.retrieve.has", func(m *nats.Msg) {
		var has utils.HasOptions
		err := bson.UnmarshalJSON(m.Data, &has)
		if err != nil {
			utils.Log.WithError(errors.Wrap(err, "error decoding event")).Error("clients.retrieve.has")
			return
		}
		utils.Log.Info("clients.retrieve.has")

		manager := d.ClientManager
		inter, err := manager.Has(has)
		if err != nil {
			utils.Log.WithError(err).Error("clients.retrieve.has")
			return
		}
		utils.Log.Info("clients.retrieve.has - Retrieved")

		data, err := bson.MarshalJSON(inter)
		if err != nil {
			utils.Log.WithError(errors.Wrap(err, "error marshaling event")).Error("clients.retrieve.has")
			return
		}
		utils.Log.Info("clients.retrieve.has - Reply")
		d.Conn.Publish(m.Reply, data)
	})

	d.Conn.Subscribe("clients.delete.send", func(m *nats.Msg) {
		var del utils.DeleteOptions
		err := bson.UnmarshalJSON(m.Data, &del)
		if err != nil {
			utils.Log.WithError(errors.Wrap(err, "error decoding event")).Error("clients.delete")
			return
		}
		utils.Log.Info("clients.delete")

		d.Conn.Publish("clients.delete.before", m.Data)

		manager := d.ClientManager

		inter, err := manager.Find(utils.FindOptions{Filter: del.Filter})
		if err != nil {
			utils.Log.WithError(errors.Wrap(err, "error finding data")).Error("clients.delete")
			return
		}
		utils.Log.Info("clients.delete - Found")

		err = manager.Delete(del)
		if err != nil {
			utils.Log.WithError(errors.Wrap(err, "error deleting data")).Error("clients.delete")
			return
		}
		utils.Log.Info("clients.delete - Deleted")

		data, err := bson.MarshalJSON(inter)
		if err != nil {
			utils.Log.WithError(errors.Wrap(err, "error marshaling event")).Error("clients.delete")
			return
		}
		utils.Log.Info("clients.delete - Reply")

		d.Conn.Publish("clients.delete.after", data)
	})

	d.Conn.Subscribe("clients.update.send", func(m *nats.Msg) {
		var update utils.UpdateOptions
		err := bson.UnmarshalJSON(m.Data, &update)
		if err != nil {
			utils.Log.WithError(errors.Wrap(err, "error decoding event")).Error("clients.update")
			return
		}
		utils.Log.Info("clients.update")

		d.Conn.Publish("clients.update.before", m.Data)

		manager := d.ClientManager

		inter, err := manager.Update(update)
		if err != nil {
			utils.Log.WithError(errors.Wrap(err, "error updating data")).Error("clients.update")
			return
		}
		utils.Log.Info("clients.update - Updated")

		data, err := bson.MarshalJSON(inter)
		if err != nil {
			utils.Log.WithError(errors.Wrap(err, "error marshaling event")).Error("clients.update")
			return
		}
		utils.Log.Info("clients.update - Reply")

		d.Conn.Publish("clients.update.after", data)
	})

	d.Conn.Subscribe("clients.create.send", func(m *nats.Msg) {
		var before *models.Client
		err := bson.UnmarshalJSON(m.Data, &before)
		if err != nil {
			utils.Log.WithError(errors.Wrap(err, "error decoding event")).Error("clients.create")
			return
		}
		utils.Log.Info("clients.create")

		d.Conn.Publish("clients.create.before", m.Data)

		manager := d.ClientManager

		inter, err := manager.Create(before)
		if err != nil {
			utils.Log.WithError(errors.Wrap(err, "error creating data")).Error("clients.create")
			return
		}
		utils.Log.Info("clients.create - Created")

		data, err := bson.MarshalJSON(inter)
		if err != nil {
			utils.Log.WithError(errors.Wrap(err, "error marshaling event")).Error("clients.create")
			return
		}
		utils.Log.Info("clients.create - Reply")

		d.Conn.Publish("clients.create.after", data)
	})

	d.Conn.Subscribe("commands.retrieve.find", func(m *nats.Msg) {
		var find utils.FindOptions
		err := bson.UnmarshalJSON(m.Data, &find)
		if err != nil {
			utils.Log.WithError(errors.Wrap(err, "error decoding event")).Error("commands.retrieve.find")
			return
		}
		utils.Log.Info("commands.retrieve.find")

		manager := d.CommandManager

		inter, err := manager.Find(find)
		if err != nil {
			utils.Log.WithError(errors.Wrap(err, "error retrieving data")).Error("commands.retrieve.find")
			return
		}

		if find.Max > 0 {
			var n []models.Command
			step := int(math.Ceil(float64(len(inter)) / float64(find.Max)))
			for i := len(inter) - 1; i >= 0; i -= step {
				n = append(n, inter[i])
			}
			inter = n
		}
		utils.Log.Info("commands.retrieve.find - Retrieved")

		data, err := bson.MarshalJSON(inter)
		if err != nil {
			utils.Log.WithError(errors.Wrap(err, "error marshaling event")).Error("commands.retrieve.find")
			return
		}
		utils.Log.Info("commands.retrieve.find - Reply")
		d.Conn.Publish(m.Reply, data)
	})

	d.Conn.Subscribe("commands.retrieve.has", func(m *nats.Msg) {
		var has utils.HasOptions
		err := bson.UnmarshalJSON(m.Data, &has)
		if err != nil {
			utils.Log.WithError(errors.Wrap(err, "error decoding event")).Error("commands.retrieve.has")
			return
		}
		utils.Log.Info("commands.retrieve.has")

		manager := d.CommandManager
		inter, err := manager.Has(has)
		if err != nil {
			utils.Log.WithError(err).Error("commands.retrieve.has")
			return
		}
		utils.Log.Info("commands.retrieve.has - Retrieved")

		data, err := bson.MarshalJSON(inter)
		if err != nil {
			utils.Log.WithError(errors.Wrap(err, "error marshaling event")).Error("commands.retrieve.has")
			return
		}
		utils.Log.Info("commands.retrieve.has - Reply")
		d.Conn.Publish(m.Reply, data)
	})

	d.Conn.Subscribe("commands.delete.send", func(m *nats.Msg) {
		var del utils.DeleteOptions
		err := bson.UnmarshalJSON(m.Data, &del)
		if err != nil {
			utils.Log.WithError(errors.Wrap(err, "error decoding event")).Error("commands.delete")
			return
		}
		utils.Log.Info("commands.delete")

		d.Conn.Publish("commands.delete.before", m.Data)

		manager := d.CommandManager

		inter, err := manager.Find(utils.FindOptions{Filter: del.Filter})
		if err != nil {
			utils.Log.WithError(errors.Wrap(err, "error finding data")).Error("commands.delete")
			return
		}
		utils.Log.Info("commands.delete - Found")

		err = manager.Delete(del)
		if err != nil {
			utils.Log.WithError(errors.Wrap(err, "error deleting data")).Error("commands.delete")
			return
		}
		utils.Log.Info("commands.delete - Deleted")

		data, err := bson.MarshalJSON(inter)
		if err != nil {
			utils.Log.WithError(errors.Wrap(err, "error marshaling event")).Error("commands.delete")
			return
		}
		utils.Log.Info("commands.delete - Reply")

		d.Conn.Publish("commands.delete.after", data)
	})

	d.Conn.Subscribe("commands.update.send", func(m *nats.Msg) {
		var update utils.UpdateOptions
		err := bson.UnmarshalJSON(m.Data, &update)
		if err != nil {
			utils.Log.WithError(errors.Wrap(err, "error decoding event")).Error("commands.update")
			return
		}
		utils.Log.Info("commands.update")

		d.Conn.Publish("commands.update.before", m.Data)

		manager := d.CommandManager

		inter, err := manager.Update(update)
		if err != nil {
			utils.Log.WithError(errors.Wrap(err, "error updating data")).Error("commands.update")
			return
		}
		utils.Log.Info("commands.update - Updated")

		data, err := bson.MarshalJSON(inter)
		if err != nil {
			utils.Log.WithError(errors.Wrap(err, "error marshaling event")).Error("commands.update")
			return
		}
		utils.Log.Info("commands.update - Reply")

		d.Conn.Publish("commands.update.after", data)
	})

	d.Conn.Subscribe("commands.create.send", func(m *nats.Msg) {
		var before *models.Command
		err := bson.UnmarshalJSON(m.Data, &before)
		if err != nil {
			utils.Log.WithError(errors.Wrap(err, "error decoding event")).Error("commands.create")
			return
		}
		utils.Log.Info("commands.create")

		d.Conn.Publish("commands.create.before", m.Data)

		manager := d.CommandManager

		inter, err := manager.Create(before)
		if err != nil {
			utils.Log.WithError(errors.Wrap(err, "error creating data")).Error("commands.create")
			return
		}
		utils.Log.Info("commands.create - Created")

		data, err := bson.MarshalJSON(inter)
		if err != nil {
			utils.Log.WithError(errors.Wrap(err, "error marshaling event")).Error("commands.create")
			return
		}
		utils.Log.Info("commands.create - Reply")

		d.Conn.Publish("commands.create.after", data)
	})

	d.Conn.Subscribe("groups.retrieve.find", func(m *nats.Msg) {
		var find utils.FindOptions
		err := bson.UnmarshalJSON(m.Data, &find)
		if err != nil {
			utils.Log.WithError(errors.Wrap(err, "error decoding event")).Error("groups.retrieve.find")
			return
		}
		utils.Log.Info("groups.retrieve.find")

		manager := d.GroupManager

		inter, err := manager.Find(find)
		if err != nil {
			utils.Log.WithError(errors.Wrap(err, "error retrieving data")).Error("groups.retrieve.find")
			return
		}

		if find.Max > 0 {
			var n []models.Group
			step := int(math.Ceil(float64(len(inter)) / float64(find.Max)))
			for i := len(inter) - 1; i >= 0; i -= step {
				n = append(n, inter[i])
			}
			inter = n
		}
		utils.Log.Info("groups.retrieve.find - Retrieved")

		data, err := bson.MarshalJSON(inter)
		if err != nil {
			utils.Log.WithError(errors.Wrap(err, "error marshaling event")).Error("groups.retrieve.find")
			return
		}
		utils.Log.Info("groups.retrieve.find - Reply")
		d.Conn.Publish(m.Reply, data)
	})

	d.Conn.Subscribe("groups.retrieve.has", func(m *nats.Msg) {
		var has utils.HasOptions
		err := bson.UnmarshalJSON(m.Data, &has)
		if err != nil {
			utils.Log.WithError(errors.Wrap(err, "error decoding event")).Error("groups.retrieve.has")
			return
		}
		utils.Log.Info("groups.retrieve.has")

		manager := d.GroupManager
		inter, err := manager.Has(has)
		if err != nil {
			utils.Log.WithError(err).Error("groups.retrieve.has")
			return
		}
		utils.Log.Info("groups.retrieve.has - Retrieved")

		data, err := bson.MarshalJSON(inter)
		if err != nil {
			utils.Log.WithError(errors.Wrap(err, "error marshaling event")).Error("groups.retrieve.has")
			return
		}
		utils.Log.Info("groups.retrieve.has - Reply")
		d.Conn.Publish(m.Reply, data)
	})

	d.Conn.Subscribe("groups.delete.send", func(m *nats.Msg) {
		var del utils.DeleteOptions
		err := bson.UnmarshalJSON(m.Data, &del)
		if err != nil {
			utils.Log.WithError(errors.Wrap(err, "error decoding event")).Error("groups.delete")
			return
		}
		utils.Log.Info("groups.delete")

		d.Conn.Publish("groups.delete.before", m.Data)

		manager := d.GroupManager

		inter, err := manager.Find(utils.FindOptions{Filter: del.Filter})
		if err != nil {
			utils.Log.WithError(errors.Wrap(err, "error finding data")).Error("groups.delete")
			return
		}
		utils.Log.Info("groups.delete - Found")

		err = manager.Delete(del)
		if err != nil {
			utils.Log.WithError(errors.Wrap(err, "error deleting data")).Error("groups.delete")
			return
		}
		utils.Log.Info("groups.delete - Deleted")

		data, err := bson.MarshalJSON(inter)
		if err != nil {
			utils.Log.WithError(errors.Wrap(err, "error marshaling event")).Error("groups.delete")
			return
		}
		utils.Log.Info("groups.delete - Reply")

		d.Conn.Publish("groups.delete.after", data)
	})

	d.Conn.Subscribe("groups.update.send", func(m *nats.Msg) {
		var update utils.UpdateOptions
		err := bson.UnmarshalJSON(m.Data, &update)
		if err != nil {
			utils.Log.WithError(errors.Wrap(err, "error decoding event")).Error("groups.update")
			return
		}
		utils.Log.Info("groups.update")

		d.Conn.Publish("groups.update.before", m.Data)

		manager := d.GroupManager

		inter, err := manager.Update(update)
		if err != nil {
			utils.Log.WithError(errors.Wrap(err, "error updating data")).Error("groups.update")
			return
		}
		utils.Log.Info("groups.update - Updated")

		data, err := bson.MarshalJSON(inter)
		if err != nil {
			utils.Log.WithError(errors.Wrap(err, "error marshaling event")).Error("groups.update")
			return
		}
		utils.Log.Info("groups.update - Reply")

		d.Conn.Publish("groups.update.after", data)
	})

	d.Conn.Subscribe("groups.create.send", func(m *nats.Msg) {
		var before *models.Group
		err := bson.UnmarshalJSON(m.Data, &before)
		if err != nil {
			utils.Log.WithError(errors.Wrap(err, "error decoding event")).Error("groups.create")
			return
		}
		utils.Log.Info("groups.create")

		d.Conn.Publish("groups.create.before", m.Data)

		manager := d.GroupManager

		inter, err := manager.Create(before)
		if err != nil {
			utils.Log.WithError(errors.Wrap(err, "error creating data")).Error("groups.create")
			return
		}
		utils.Log.Info("groups.create - Created")

		data, err := bson.MarshalJSON(inter)
		if err != nil {
			utils.Log.WithError(errors.Wrap(err, "error marshaling event")).Error("groups.create")
			return
		}
		utils.Log.Info("groups.create - Reply")

		d.Conn.Publish("groups.create.after", data)
	})

	d.Conn.Subscribe("servers.retrieve.find", func(m *nats.Msg) {
		var find utils.FindOptions
		err := bson.UnmarshalJSON(m.Data, &find)
		if err != nil {
			utils.Log.WithError(errors.Wrap(err, "error decoding event")).Error("servers.retrieve.find")
			return
		}
		utils.Log.Info("servers.retrieve.find")

		manager := d.ServerManager

		inter, err := manager.Find(find)
		if err != nil {
			utils.Log.WithError(errors.Wrap(err, "error retrieving data")).Error("servers.retrieve.find")
			return
		}

		if find.Max > 0 {
			var n []models.Server
			step := int(math.Ceil(float64(len(inter)) / float64(find.Max)))
			for i := len(inter) - 1; i >= 0; i -= step {
				n = append(n, inter[i])
			}
			inter = n
		}
		utils.Log.Info("servers.retrieve.find - Retrieved")

		data, err := bson.MarshalJSON(inter)
		if err != nil {
			utils.Log.WithError(errors.Wrap(err, "error marshaling event")).Error("servers.retrieve.find")
			return
		}
		utils.Log.Info("servers.retrieve.find - Reply")
		d.Conn.Publish(m.Reply, data)
	})

	d.Conn.Subscribe("servers.retrieve.has", func(m *nats.Msg) {
		var has utils.HasOptions
		err := bson.UnmarshalJSON(m.Data, &has)
		if err != nil {
			utils.Log.WithError(errors.Wrap(err, "error decoding event")).Error("servers.retrieve.has")
			return
		}
		utils.Log.Info("servers.retrieve.has")

		manager := d.ServerManager
		inter, err := manager.Has(has)
		if err != nil {
			utils.Log.WithError(err).Error("servers.retrieve.has")
			return
		}
		utils.Log.Info("servers.retrieve.has - Retrieved")

		data, err := bson.MarshalJSON(inter)
		if err != nil {
			utils.Log.WithError(errors.Wrap(err, "error marshaling event")).Error("servers.retrieve.has")
			return
		}
		utils.Log.Info("servers.retrieve.has - Reply")
		d.Conn.Publish(m.Reply, data)
	})

	d.Conn.Subscribe("servers.delete.send", func(m *nats.Msg) {
		var del utils.DeleteOptions
		err := bson.UnmarshalJSON(m.Data, &del)
		if err != nil {
			utils.Log.WithError(errors.Wrap(err, "error decoding event")).Error("servers.delete")
			return
		}
		utils.Log.Info("servers.delete")

		d.Conn.Publish("servers.delete.before", m.Data)

		manager := d.ServerManager

		inter, err := manager.Find(utils.FindOptions{Filter: del.Filter})
		if err != nil {
			utils.Log.WithError(errors.Wrap(err, "error finding data")).Error("servers.delete")
			return
		}
		utils.Log.Info("servers.delete - Found")

		err = manager.Delete(del)
		if err != nil {
			utils.Log.WithError(errors.Wrap(err, "error deleting data")).Error("servers.delete")
			return
		}
		utils.Log.Info("servers.delete - Deleted")

		data, err := bson.MarshalJSON(inter)
		if err != nil {
			utils.Log.WithError(errors.Wrap(err, "error marshaling event")).Error("servers.delete")
			return
		}
		utils.Log.Info("servers.delete - Reply")

		d.Conn.Publish("servers.delete.after", data)
	})

	d.Conn.Subscribe("servers.update.send", func(m *nats.Msg) {
		var update utils.UpdateOptions
		err := bson.UnmarshalJSON(m.Data, &update)
		if err != nil {
			utils.Log.WithError(errors.Wrap(err, "error decoding event")).Error("servers.update")
			return
		}
		utils.Log.Info("servers.update")

		d.Conn.Publish("servers.update.before", m.Data)

		manager := d.ServerManager

		inter, err := manager.Update(update)
		if err != nil {
			utils.Log.WithError(errors.Wrap(err, "error updating data")).Error("servers.update")
			return
		}
		utils.Log.Info("servers.update - Updated")

		data, err := bson.MarshalJSON(inter)
		if err != nil {
			utils.Log.WithError(errors.Wrap(err, "error marshaling event")).Error("servers.update")
			return
		}
		utils.Log.Info("servers.update - Reply")

		d.Conn.Publish("servers.update.after", data)
	})

	d.Conn.Subscribe("servers.create.send", func(m *nats.Msg) {
		var before *models.Server
		err := bson.UnmarshalJSON(m.Data, &before)
		if err != nil {
			utils.Log.WithError(errors.Wrap(err, "error decoding event")).Error("servers.create")
			return
		}
		utils.Log.Info("servers.create")

		d.Conn.Publish("servers.create.before", m.Data)

		manager := d.ServerManager

		inter, err := manager.Create(before)
		if err != nil {
			utils.Log.WithError(errors.Wrap(err, "error creating data")).Error("servers.create")
			return
		}
		utils.Log.Info("servers.create - Created")

		data, err := bson.MarshalJSON(inter)
		if err != nil {
			utils.Log.WithError(errors.Wrap(err, "error marshaling event")).Error("servers.create")
			return
		}
		utils.Log.Info("servers.create - Reply")

		d.Conn.Publish("servers.create.after", data)
	})

	d.Conn.Subscribe("uploads.retrieve.find", func(m *nats.Msg) {
		var find utils.FindOptions
		err := bson.UnmarshalJSON(m.Data, &find)
		if err != nil {
			utils.Log.WithError(errors.Wrap(err, "error decoding event")).Error("uploads.retrieve.find")
			return
		}
		utils.Log.Info("uploads.retrieve.find")

		manager := d.UploadManager

		inter, err := manager.Find(find)
		if err != nil {
			utils.Log.WithError(errors.Wrap(err, "error retrieving data")).Error("uploads.retrieve.find")
			return
		}

		if find.Max > 0 {
			var n []models.Upload
			step := int(math.Ceil(float64(len(inter)) / float64(find.Max)))
			for i := len(inter) - 1; i >= 0; i -= step {
				n = append(n, inter[i])
			}
			inter = n
		}
		utils.Log.Info("uploads.retrieve.find - Retrieved")

		data, err := bson.MarshalJSON(inter)
		if err != nil {
			utils.Log.WithError(errors.Wrap(err, "error marshaling event")).Error("uploads.retrieve.find")
			return
		}
		utils.Log.Info("uploads.retrieve.find - Reply")
		d.Conn.Publish(m.Reply, data)
	})

	d.Conn.Subscribe("uploads.retrieve.has", func(m *nats.Msg) {
		var has utils.HasOptions
		err := bson.UnmarshalJSON(m.Data, &has)
		if err != nil {
			utils.Log.WithError(errors.Wrap(err, "error decoding event")).Error("uploads.retrieve.has")
			return
		}
		utils.Log.Info("uploads.retrieve.has")

		manager := d.UploadManager
		inter, err := manager.Has(has)
		if err != nil {
			utils.Log.WithError(err).Error("uploads.retrieve.has")
			return
		}
		utils.Log.Info("uploads.retrieve.has - Retrieved")

		data, err := bson.MarshalJSON(inter)
		if err != nil {
			utils.Log.WithError(errors.Wrap(err, "error marshaling event")).Error("uploads.retrieve.has")
			return
		}
		utils.Log.Info("uploads.retrieve.has - Reply")
		d.Conn.Publish(m.Reply, data)
	})

	d.Conn.Subscribe("uploads.delete.send", func(m *nats.Msg) {
		var del utils.DeleteOptions
		err := bson.UnmarshalJSON(m.Data, &del)
		if err != nil {
			utils.Log.WithError(errors.Wrap(err, "error decoding event")).Error("uploads.delete")
			return
		}
		utils.Log.Info("uploads.delete")

		d.Conn.Publish("uploads.delete.before", m.Data)

		manager := d.UploadManager

		inter, err := manager.Find(utils.FindOptions{Filter: del.Filter})
		if err != nil {
			utils.Log.WithError(errors.Wrap(err, "error finding data")).Error("uploads.delete")
			return
		}
		utils.Log.Info("uploads.delete - Found")

		err = manager.Delete(del)
		if err != nil {
			utils.Log.WithError(errors.Wrap(err, "error deleting data")).Error("uploads.delete")
			return
		}
		utils.Log.Info("uploads.delete - Deleted")

		data, err := bson.MarshalJSON(inter)
		if err != nil {
			utils.Log.WithError(errors.Wrap(err, "error marshaling event")).Error("uploads.delete")
			return
		}
		utils.Log.Info("uploads.delete - Reply")

		d.Conn.Publish("uploads.delete.after", data)
	})

	d.Conn.Subscribe("uploads.update.send", func(m *nats.Msg) {
		var update utils.UpdateOptions
		err := bson.UnmarshalJSON(m.Data, &update)
		if err != nil {
			utils.Log.WithError(errors.Wrap(err, "error decoding event")).Error("uploads.update")
			return
		}
		utils.Log.Info("uploads.update")

		d.Conn.Publish("uploads.update.before", m.Data)

		manager := d.UploadManager

		inter, err := manager.Update(update)
		if err != nil {
			utils.Log.WithError(errors.Wrap(err, "error updating data")).Error("uploads.update")
			return
		}
		utils.Log.Info("uploads.update - Updated")

		data, err := bson.MarshalJSON(inter)
		if err != nil {
			utils.Log.WithError(errors.Wrap(err, "error marshaling event")).Error("uploads.update")
			return
		}
		utils.Log.Info("uploads.update - Reply")

		d.Conn.Publish("uploads.update.after", data)
	})

	d.Conn.Subscribe("uploads.create.send", func(m *nats.Msg) {
		var before *models.Upload
		err := bson.UnmarshalJSON(m.Data, &before)
		if err != nil {
			utils.Log.WithError(errors.Wrap(err, "error decoding event")).Error("uploads.create")
			return
		}
		utils.Log.Info("uploads.create")

		d.Conn.Publish("uploads.create.before", m.Data)

		manager := d.UploadManager

		inter, err := manager.Create(before)
		if err != nil {
			utils.Log.WithError(errors.Wrap(err, "error creating data")).Error("uploads.create")
			return
		}
		utils.Log.Info("uploads.create - Created")

		data, err := bson.MarshalJSON(inter)
		if err != nil {
			utils.Log.WithError(errors.Wrap(err, "error marshaling event")).Error("uploads.create")
			return
		}
		utils.Log.Info("uploads.create - Reply")

		d.Conn.Publish("uploads.create.after", data)
	})

	d.Conn.Subscribe("users.retrieve.find", func(m *nats.Msg) {
		var find utils.FindOptions
		err := bson.UnmarshalJSON(m.Data, &find)
		if err != nil {
			utils.Log.WithError(errors.Wrap(err, "error decoding event")).Error("users.retrieve.find")
			return
		}
		utils.Log.Info("users.retrieve.find")

		manager := d.UserManager

		inter, err := manager.Find(find)
		if err != nil {
			utils.Log.WithError(errors.Wrap(err, "error retrieving data")).Error("users.retrieve.find")
			return
		}

		if find.Max > 0 {
			var n []models.User
			step := int(math.Ceil(float64(len(inter)) / float64(find.Max)))
			for i := len(inter) - 1; i >= 0; i -= step {
				n = append(n, inter[i])
			}
			inter = n
		}
		utils.Log.Info("users.retrieve.find - Retrieved")

		data, err := bson.MarshalJSON(inter)
		if err != nil {
			utils.Log.WithError(errors.Wrap(err, "error marshaling event")).Error("users.retrieve.find")
			return
		}
		utils.Log.Info("users.retrieve.find - Reply")
		d.Conn.Publish(m.Reply, data)
	})

	d.Conn.Subscribe("users.retrieve.has", func(m *nats.Msg) {
		var has utils.HasOptions
		err := bson.UnmarshalJSON(m.Data, &has)
		if err != nil {
			utils.Log.WithError(errors.Wrap(err, "error decoding event")).Error("users.retrieve.has")
			return
		}
		utils.Log.Info("users.retrieve.has")

		manager := d.UserManager
		inter, err := manager.Has(has)
		if err != nil {
			utils.Log.WithError(err).Error("users.retrieve.has")
			return
		}
		utils.Log.Info("users.retrieve.has - Retrieved")

		data, err := bson.MarshalJSON(inter)
		if err != nil {
			utils.Log.WithError(errors.Wrap(err, "error marshaling event")).Error("users.retrieve.has")
			return
		}
		utils.Log.Info("users.retrieve.has - Reply")
		d.Conn.Publish(m.Reply, data)
	})

	d.Conn.Subscribe("users.delete.send", func(m *nats.Msg) {
		var del utils.DeleteOptions
		err := bson.UnmarshalJSON(m.Data, &del)
		if err != nil {
			utils.Log.WithError(errors.Wrap(err, "error decoding event")).Error("users.delete")
			return
		}
		utils.Log.Info("users.delete")

		d.Conn.Publish("users.delete.before", m.Data)

		manager := d.UserManager

		inter, err := manager.Find(utils.FindOptions{Filter: del.Filter})
		if err != nil {
			utils.Log.WithError(errors.Wrap(err, "error finding data")).Error("users.delete")
			return
		}
		utils.Log.Info("users.delete - Found")

		err = manager.Delete(del)
		if err != nil {
			utils.Log.WithError(errors.Wrap(err, "error deleting data")).Error("users.delete")
			return
		}
		utils.Log.Info("users.delete - Deleted")

		data, err := bson.MarshalJSON(inter)
		if err != nil {
			utils.Log.WithError(errors.Wrap(err, "error marshaling event")).Error("users.delete")
			return
		}
		utils.Log.Info("users.delete - Reply")

		d.Conn.Publish("users.delete.after", data)
	})

	d.Conn.Subscribe("users.update.send", func(m *nats.Msg) {
		var update utils.UpdateOptions
		err := bson.UnmarshalJSON(m.Data, &update)
		if err != nil {
			utils.Log.WithError(errors.Wrap(err, "error decoding event")).Error("users.update")
			return
		}
		utils.Log.Info("users.update")

		d.Conn.Publish("users.update.before", m.Data)

		manager := d.UserManager

		inter, err := manager.Update(update)
		if err != nil {
			utils.Log.WithError(errors.Wrap(err, "error updating data")).Error("users.update")
			return
		}
		utils.Log.Info("users.update - Updated")

		data, err := bson.MarshalJSON(inter)
		if err != nil {
			utils.Log.WithError(errors.Wrap(err, "error marshaling event")).Error("users.update")
			return
		}
		utils.Log.Info("users.update - Reply")

		d.Conn.Publish("users.update.after", data)
	})

	d.Conn.Subscribe("users.create.send", func(m *nats.Msg) {
		var before *models.User
		err := bson.UnmarshalJSON(m.Data, &before)
		if err != nil {
			utils.Log.WithError(errors.Wrap(err, "error decoding event")).Error("users.create")
			return
		}
		utils.Log.Info("users.create")

		d.Conn.Publish("users.create.before", m.Data)

		manager := d.UserManager

		inter, err := manager.Create(before)
		if err != nil {
			utils.Log.WithError(errors.Wrap(err, "error creating data")).Error("users.create")
			return
		}
		utils.Log.Info("users.create - Created")

		data, err := bson.MarshalJSON(inter)
		if err != nil {
			utils.Log.WithError(errors.Wrap(err, "error marshaling event")).Error("users.create")
			return
		}
		utils.Log.Info("users.create - Reply")

		d.Conn.Publish("users.create.after", data)
	})

}

func (d *Database) Close() {
	if d.db.Session != nil {
		d.db.Session.Close()
	}
	if d.Conn != nil {
		d.Conn.Close()
	}
}
