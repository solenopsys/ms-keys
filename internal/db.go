package internal

import (
	"encoding/json"
	"github.com/dgraph-io/badger/v3"
	"k8s.io/klog/v2"
	"ms-keys/pkg"
)

func NewDb(path string) *DriveDb {
	db := &DriveDb{
		Path: path,
	}
	db.Open()
	return db
}

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

	jsonRegister, err := json.Marshal(register)
	if err != nil {
		klog.Error(err)
	}
	err = d.db.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(register.Hash), jsonRegister)
	})
	if err != nil {
		klog.Error(err)
	}
}
