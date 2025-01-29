package core

import (
	"os"
	"barbecue/data"
)

type context struct {
	Debug		bool
	Database	*data.Database
	Logger		*Logger
}

var Context *context = &context {
	Debug: false,
}

func (c *context) Ok(data interface{}) error {
	if data != nil {
		c.Logger.Info(data)
	} else {
		c.Logger.Info("OK")
	}
	os.Exit(0)
	return nil
}

func (c *context) Ko(err interface{}) error {
	if err != nil {
		c.Logger.Error(err)
	} else {
		c.Logger.Error("KO")
	}
	os.Exit(1)
	return nil
}
