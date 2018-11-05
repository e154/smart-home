package migrations

type MigrationsConfig struct {
	Source string
	Dir    string
}

func NewMigrationsConfig() *MigrationsConfig {
	return &MigrationsConfig{
		Source: "assets",
		Dir:    "migrations",
	}
}
