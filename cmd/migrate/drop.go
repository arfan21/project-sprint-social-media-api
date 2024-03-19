package migration

func Drop() error {

	migrations, err := initMigration()
	if err != nil {
		return err
	}

	if err := migrations.Drop(); err != nil {
		return err
	}

	return nil
}
