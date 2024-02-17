package testutil

import (
	"github.com/sirupsen/logrus"
	"os"
)

func IsIntegration() {
	if os.Getenv("TEST_TYPE") != "integration" {
		logrus.Infoln("integration test skipped")
		os.Exit(0)
	}
}
