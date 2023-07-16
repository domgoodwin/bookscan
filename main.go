package main

import (
	"github.com/domgoodwin/bookscan/cmd"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetLevel(logrus.DebugLevel)
	cmd.Execute()
}
