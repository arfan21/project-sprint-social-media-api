package migration

func Fresh() error {

	migrations, err := initMigration()
	if err != nil {
		return err
	}

	if err := migrations.Down(); err != nil {
		return err
	}

	if err := migrations.Up(); err != nil {
		return err
	}

	return nil
}
