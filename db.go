package main

import (
	"encoding/json"
	"github.com/dgraph-io/badger/v3"
	"k8s.io/klog/v2"
)

type DriveDb struct {
	path string
	db   *badger.DB
}

func (d *DriveDb) open() {
	var err error
	d.db, err = badger.Open(badger.DefaultOptions(d.path))
	if err != nil {
		klog.Fatal(err)
	}
}

func (d *DriveDb) close() {
	err := d.db.Close()
	if err != nil {
		klog.Fatal(err)
	}
}

func (d *DriveDb) loadRegister(email string) (RegisterData, error) {
	var register RegisterData
	err := d.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(email))
		if err != nil {
			return err
		}

		err = item.Value(func(val []byte) error {
			json.Unmarshal(val, &register)
			return nil
		})

		return err
	})

	if err != nil {
		return register, err
	}

	return register, nil
}

func (d *DriveDb) saveRegister(register RegisterData) {
	err := d.db.Update(func(txn *badger.Txn) error {
		jsonRegister, err := json.Marshal(register)

		err = txn.Set([]byte(register.Email), []byte(jsonRegister))
		if err != nil {
			return err
		}
		return txn.Commit()
	})
	if err != nil {
		klog.Fatal(err)
	}
}
