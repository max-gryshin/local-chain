package main

import (
	leveldb2 "local-chain/internal/adapters/outbound/leveldb"

	"github.com/syndtr/goleveldb/leveldb"
)

func main() {
	db, err := leveldb.OpenFile("./db", nil)
	if err != nil {
		panic(err)
	}
	defer func(db *leveldb.DB) {
		err = db.Close()
		if err != nil {
			panic(err)
		}
	}(db)

	leveldb2.New(db)
}
