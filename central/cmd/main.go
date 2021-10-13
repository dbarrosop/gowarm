package main

import (
	"fmt"

	"tinygo.org/x/bluetooth"

	"github.com/dbarrosop/gowarm/central/core"
	"github.com/sirupsen/logrus"
)

func main() {
	l := logrus.New()
	l.SetLevel(logrus.DebugLevel)

	logger := logrus.NewEntry(l)

	adapter := bluetooth.DefaultAdapter
	if err := adapter.Enable(); err != nil {
		panic(fmt.Sprintf("problem enabling adapter: %s", err))
	}

	c := core.New(adapter, logger, "F5:AC:32:45:C4:AD")
	if err := c.Init(); err != nil {
		panic(fmt.Sprintf("problem initializing central: %s", err))
	}
}
