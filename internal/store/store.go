package store

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"netobject-status-sync/internal/config"
	"netobject-status-sync/internal/model"
)

const DEVICE_QUERY = "SELECT MAC, DEVICE_ID FROM td.device"
const INSERT_STATUS_QUERY = "INSERT INTO td.externalDeviceInfo (device_id, status) VALUES (?, ?) ON DUPLICATE KEY UPDATE status = ?"

type Store struct {
	Db *sql.DB
}

func NewStore() *Store {
	return &Store{}
}

func (store *Store) ConnectToDb(config *config.Config) error {
	db, err := sql.Open("mysql", fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s",
		config.Store.User,
		config.Store.Pwd,
		config.Store.Dsn,
		config.Store.Port,
		config.Store.Database,
	))
	if err != nil {
		return err
	}

	store.Db = db
	return nil
}

func (store *Store) FindMacs() ([]model.Device, error) {
	stmt, err := store.Db.Prepare(DEVICE_QUERY)

	if err != nil {
		return nil, err
	}

	rows, err := stmt.Query()

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var devices []model.Device

	for rows.Next() {
		var device model.Device
		if err = rows.Scan(&device.Mac, &device.Id); err != nil {
			return devices, err
		}
		devices = append(devices, device)
	}

	if err = rows.Err(); err != nil {
		return devices, err
	}
	return devices, nil
}

func (store *Store) SetInstalledStatus(device model.Device, status bool) error {
	stmt, err := store.Db.Prepare(INSERT_STATUS_QUERY)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(device.Id, status, status)
	if err != nil {
		return err
	}

	return nil
}
