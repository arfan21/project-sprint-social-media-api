package migration

func Up() error {

	migrations, err := initMigration()
	if err != nil {
		return err
	}

	if err := migrations.Up(); err != nil {
		return err
	}

	return nil
}
