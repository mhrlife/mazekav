package integration

import (
	"github.com/sirupsen/logrus"
	"mazekav/pkg/testutil"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	testutil.IsIntegration()

	// Start the MySQL test server
	if err := testutil.StartMySQLServer(); err != nil {
		logrus.WithError(err).Fatal("Failed to start MySQL test server")
		os.Exit(1)
	}

	// Run the tests
	code := m.Run()

	// Teardown the MySQL test server
	testutil.TeardownMySQLServer()
	os.Exit(code)
}
