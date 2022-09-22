package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/ra111eo/TestTaskITgo/pb"
)

var ErrAlreadyExists = errors.New("record already exist")
var ErrNotInDB = errors.New("record not in db")
var ErrNotEnoughMoney = errors.New("not enough money")

type EwalletStore interface {
	Save(ewallet *pb.Ewallet) error
	Find(id string) (*pb.Ewallet, error)
	Send(transaction *pb.Transaction) error
	Getlast() []byte
	// UpdateDBStore() error
}

type DBEwalletStore struct {
	mutex  sync.RWMutex
	inDBid map[string]bool
}

func NewDBEwalletStore() *DBEwalletStore {
	log.Printf("time to create new ewallet store")
	var dbStore DBEwalletStore
	dbStore.inDBid = make(map[string]bool)
	UpdateDBStore(&dbStore)
	return &dbStore
}

func UpdateDBStore(store *DBEwalletStore) error {
	store.mutex.Lock()
	defer store.mutex.Unlock()
	dbInfo, err := GetDBinfo()
	if err != nil {
		log.Fatalln(err)
		return err
	}
	for i := range dbInfo.Rows {
		store.inDBid[dbInfo.Rows[i].Id] = true
	}
	return nil
}

func (store *DBEwalletStore) Save(ewallet *pb.Ewallet) error {
	store.mutex.Lock()
	defer store.mutex.Unlock()
	log.Printf("time to save ewallet with id %s ", ewallet.Id)
	if store.inDBid[ewallet.Id] == true {
		log.Printf("already exists")
		return ErrAlreadyExists
	}
	err := UpdateDataInDB(ewallet)
	if err != nil {
		log.Fatalln(err)
		return err
	}
	store.inDBid[ewallet.Id] = true
	log.Printf("saved succesfully with id %s ", ewallet.Id)
	return nil
}

func (store *DBEwalletStore) Find(id string) (*pb.Ewallet, error) {
	store.mutex.RLock()
	defer store.mutex.RUnlock()
	log.Printf("time to find ewallet with id %s ", id)

	if store.inDBid[id] != true {
		log.Printf("there are no ewallet with that id %s ", id)
		return nil, ErrNotInDB
	}
	var ewallet *pb.Ewallet
	dbUrl := makeDBurl()
	jsonData, err := getInfoFromDB(dbUrl, id)
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}
	err = json.Unmarshal(jsonData, &ewallet)
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}
	log.Printf("id %s  have found successfully", id)
	return ewallet, nil
}

func (store *DBEwalletStore) Send(transaction *pb.Transaction) error {
	sender, err := store.Find(transaction.From)
	if err != nil {
		return err
	}
	if sender.Balance < transaction.Amount {
		return ErrNotEnoughMoney
	}
	recipient, err := store.Find(transaction.To)
	if err != nil {
		return err
	}
	store.mutex.Lock()
	defer store.mutex.Unlock()
	sender.Balance -= transaction.Amount
	transaction.Datetime = time.Now().String()
	sender.Send = append(sender.Send, transaction)
	recipient.Balance += transaction.Amount
	recipient.Get = append(recipient.Get, transaction)

	err = UpdateDataInDB(sender)
	if err != nil {
		log.Fatalln("can not update sender", err)
		return err
	}
	err = UpdateDataInDB(recipient)
	if err != nil {
		log.Fatalln("can not update recipient", err)
		return err
	}
	return nil
}

func (store *DBEwalletStore) Getlast() []byte {
	var allFields []*pb.TransactionField

	for id := range store.inDBid {
		var oneField pb.TransactionField
		ewallet, _ := store.Find(id)
		if ewallet.Get != nil {
			fmt.Println(ewallet.Get)
			oneField.Tr = ewallet.Get
			allFields = append(allFields, &oneField)
			ewallet.Get = nil
			UpdateDataInDB(ewallet)
		}
	}
	jsonData, err := json.Marshal(allFields)
	if err != nil {
		log.Fatal(err)
	}
	return jsonData
}
