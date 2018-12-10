package struct2sql

func (db DB) CreateTable(model interface{}) error {
	return createTable(db, model)
}
func (db Tx) CreateTable(model interface{}) error {
	return createTable(db, model)
}
func (db DB) DropTable(model ...interface{}) error {
	return dropTable(db, model...)
}
func (db Tx) DropTable(model ...interface{}) error {
	return dropTable(db, model...)
}

func createTable(db IDB, model interface{}) error {
	table := AnalyzeStruct(model)
	stmt, err := db.Prepare(CreateTableQuery(table))
	if err != nil {
		return err
	}
	_, err = stmt.Exec()
	if err != nil {
		return err
	}
	return nil
}

func dropTable(db IDB, models ...interface{}) error {
	tables := make([]SQLTable, len(models))
	for i, model := range models {
		tables[i] = AnalyzeStruct(model)
	}
	stmt, err := db.Prepare(DropTableQuery(tables...))
	if err != nil {
		return err
	}
	_, err = stmt.Exec()
	if err != nil {
		return err
	}
	return nil
}
