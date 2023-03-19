package internal

import (
	"encoding/json"
	"github.com/dgraph-io/badger/v3"
	"k8s.io/klog/v2"
	"ms-keys/pkg"
)

type DriveDb struct {
	Path string
	db   *badger.DB
}

func (d *DriveDb) Open() {
	var err error
	d.db, err = badger.Open(badger.DefaultOptions(d.Path))
	if err != nil {
		klog.Fatal(err)
	}
}

func (d *DriveDb) Close() {
	err := d.db.Close()
	if err != nil {
		klog.Fatal(err)
	}
}

func (d *DriveDb) LoadRegister(hash string) (pkg.RegisterData, error) {
	var register pkg.RegisterData
	err := d.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(hash))
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

func (d *DriveDb) SaveRegister(register pkg.RegisterData) {
	err := d.db.Update(func(txn *badger.Txn) error {
		jsonRegister, err := json.Marshal(register)

		err = txn.Set([]byte(register.Hash), []byte(jsonRegister))
		if err != nil {
			return err
		}
		return txn.Commit()
	})
	if err != nil {
		klog.Fatal(err)
	}
}
