package main

type (
	Cli struct {
		Api
		*ApiContext
	}
)

func NewCli(db *Database) (*Cli, error) {
	var cli Cli = Cli {}
	//cli.Debug = DEBUG
	//cli.NewRoutes(db)
	return &cli, nil
}

func (c Cli) Run() {
}

func (c Cli) ok(code int, data interface{}) error {
	return nil
}

func (c Cli) ko(code int) error {
	return nil
}
