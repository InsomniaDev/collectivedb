package database

// Init will initialize the database and return a Database with the total depth set
func Init(layers int) (db Database) {
	db = Database{
		TotalDepth: layers,
	}
	return
}
