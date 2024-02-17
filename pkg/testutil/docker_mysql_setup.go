package testutil

import (
	"fmt"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"math/rand"
	"testing"
	"time"
)

// MySQLTestContainer holds the state of the test MySQL container.
type MySQLTestContainer struct {
	pool           *dockertest.Pool     // Pool represents the connection pool to Docker.
	mysql          *dockertest.Resource // MySQL Docker container resource.
	mainConnection *gorm.DB             // Main connection to the MySQL database.
	running        bool                 // Indicates if the MySQL container is running.
}

// Global instance of MySQLTestContainer.
var mysqlTestContainer MySQLTestContainer

// StartMySQLServer initializes and starts the MySQL Docker container for testing.
func StartMySQLServer() error {
	pool, err := dockertest.NewPool("")
	if err != nil {
		logrus.WithError(err).Error("Could not construct pool")
		return err
	}

	// uses pool to try to connect to Docker
	err = pool.Client.Ping()
	if err != nil {
		logrus.WithError(err).Errorln("Could not connect to Docker")
		return err
	}

	// pulls an image, creates a container based on it and runs it
	resource, err := pool.RunWithOptions(
		&dockertest.RunOptions{
			Repository: "mysql",
			Tag:        "latest",
			Env:        []string{"MYSQL_ROOT_PASSWORD=secret"},
		}, func(config *docker.HostConfig) {
			config.AutoRemove = true // Ensure the container is removed after the test.
			config.RestartPolicy = docker.RestartPolicy{
				Name: "no", // Do not automatically restart the container.
			}
		},
	)

	if err != nil {
		logrus.WithError(err).Errorln("Could not start resource")
		return err
	}

	// Set container to expire after two minutes to avoid dangling resources in case of test interruption.
	if err := resource.Expire(120); err != nil {
		logrus.WithError(err).Errorln("Couldn't set MySQL container expiration")
		return err
	}

	var db *gorm.DB
	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	if err := pool.Retry(func() error {
		db, err = connectGorm(resource.GetPort("3306/tcp"), "mysql")
		return err
	}); err != nil {
		logrus.WithError(err).Errorln("Could not connect to database")
		return err
	}

	mysqlTestContainer.pool = pool
	mysqlTestContainer.mysql = resource
	mysqlTestContainer.running = true
	mysqlTestContainer.mainConnection = db

	return nil
}

// TeardownMySQLServer cleans up the MySQL Docker container.
func TeardownMySQLServer() error {
	if err := mysqlTestContainer.pool.Purge(mysqlTestContainer.mysql); err != nil {
		logrus.WithError(err).Errorln("Could not purge resource")
		return err
	}
	mysqlTestContainer.running = false
	return nil
}

// GetConnection creates a new database within the MySQL container and returns a connection to it.
// This allows each test to operate in isolation.
func GetConnection(t *testing.T) *gorm.DB {
	if !mysqlTestContainer.running {
		t.Fatal("mysql container is not running ,did you forget to call StartMySQL?")
	}

	dbName := generateDBName()
	if err := mysqlTestContainer.mainConnection.Exec("CREATE DATABASE " + dbName).Error; err != nil {
		logrus.WithError(err).WithField("dbname", dbName).Errorln("couldn't create the test database")
		t.Fatalf("couldn't create a test database: %v", err)
	}

	db, err := connectGorm(mysqlTestContainer.mysql.GetPort("3306/tcp"), dbName)
	if err != nil {
		t.Fatalf("couldn't connect to the created test database: %v", err)
	}
	return db
}

// connectGorm establishes a connection to a specified database using GORM.
func connectGorm(port, dbName string) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"root:secret@tcp(127.0.0.1:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		port, dbName)
	return gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
}

// generateRandomString generates a random string of a given length.
func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

// generateDBName generates a database name with the specified prefix and random suffix of given length.
func generateDBName() string {
	return fmt.Sprintf("%s%s", "testdb", generateRandomString(10))
}
