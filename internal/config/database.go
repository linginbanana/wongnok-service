package config

type Database struct {
	URL string `env:"DATABASE_URL" envDefault:"postgres://postgres:212224@localhost:5432/wongnok"`
}
