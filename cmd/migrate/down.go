package migration

func Down() error {

	migrations, err := initMigration()
	if err != nil {
		return err
	}

	if err := migrations.Down(); err != nil {
		return err
	}

	return nil
}
