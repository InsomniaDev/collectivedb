package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/insomniadev/collective-db/internal/data"
)

const collectiveDatabase = "collective"

func Start() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/get/{id}/{database}", getWithDatabase)
	myRouter.HandleFunc("/get/{id}", getByKey)
	myRouter.HandleFunc("/update", update)
	myRouter.HandleFunc("/delete", delete)
	log.Fatal(http.ListenAndServe(":10000", myRouter))
}

func getWithDatabase(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["id"]
	database := vars["database"]

	w.Write(GetByDatabase(key, database))
}

// GetByDatabase
//
// Will return the data from the database provided with the key provided
//
//	Returns a byte array
func GetByDatabase(key, database string) []byte {
	if exists, returnedData := data.RetrieveDataFromDatabase(&key, &database); exists {
		return *returnedData
	}
	return nil
}

func getByKey(w http.ResponseWriter, r *http.Request) {
	key := mux.Vars(r)["id"]
	// fmt.Println(key)

	w.Write(Get(key))
}

func Get(key string) []byte {
	return GetByDatabase(key, collectiveDatabase)
}

func update(w http.ResponseWriter, r *http.Request) {
	var body UpdateStruct

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		log.Println(err)
	}

	updated := false
	switch {
	case body.Key == "":
		w.Write([]byte("key must exist"))
		w.WriteHeader(400)
		return
	case body.Database != "":
		updated = UpdateByDatabase(body.Key, body.Database, []byte(body.Data))
	case body.Database == "":
		updated = Update(body.Key, []byte(body.Data))
	}
	if updated {
		w.Write([]byte(body.Key))
		w.WriteHeader(200)
		return
	}
	w.Write([]byte(fmt.Sprintf("failed to update %s", body.Key)))
	w.WriteHeader(400)
}

// UpdateByDatabase
//
// Will store the new data or will update the existing data in the provided location
//
//	Returns a boolean on if it was updated successfully
func UpdateByDatabase(key, database string, newData []byte) bool {
	if updated, _ := data.StoreDataInDatabase(&key, &database, &newData, false, 0); updated {
		return updated
	}
	return false
}

func Update(key string, data []byte) bool {
	return UpdateByDatabase(key, collectiveDatabase, data)
}

func delete(w http.ResponseWriter, r *http.Request) {
	var body DeleteStruct

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		log.Println(err)
	}

	switch {
	case body.Key == "":
		w.Write([]byte("key must exist"))
		w.WriteHeader(400)
		return
	case body.Database != "":
		err = DeleteByDatabase(body.Key, body.Database)
	case body.Database == "":
		err = Delete(body.Key)
	}
	if err == nil {
		w.Write([]byte(body.Key))
		w.WriteHeader(200)
		return
	}
	w.Write([]byte(fmt.Sprintf("failed to delete %s", body.Key)))
	w.WriteHeader(400)
}

// DeleteByDatabase
// 
// Deletes from the provided database with the provided key
// 		Returns either a nil or an error
func DeleteByDatabase(key, database string) error {
	if deleted, err := data.DeleteDataFromDatabase(&key, &database); err != nil {
		if !deleted {
			log.Println("Not deleted correctely ", key, " ", database)
		}
		return err
	}
	return nil
}

func Delete(key string) error {
	return DeleteByDatabase(key, collectiveDatabase)
}
