package postgres

import (
	"database/sql"

	// This is imported for migrations
	_ "github.com/Kount/pq-timeouts"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"strings"

	"github.com/stretchr/testify/suite"

	migrate "github.com/golang-migrate/migrate/v4"
	_postgres "github.com/golang-migrate/migrate/v4/database/postgres"
)

const (
	migrationDbName = "postgres"
	postgresDriver  = "pq-timeouts"
	// DB_HOST=localhost DB_USER=postgres DB_PASS=postgres DB_NAME=traingolang DB_PORT=5432
	// DefaultTestDsn is the default url for testing postgresql in the postgres test suites
	DefaultTestDsn = "user=postgres password=postgres dbname=traingolang_integration host=localhost port=5432 sslmode=disable read_timeout=300000 write_timeout=300000"
)

type migration struct {
	Migrate *migrate.Migrate
}

func (m *migration) Up() (bool, error) {
	err := m.Migrate.Up()
	if err != nil {
		if err == migrate.ErrNoChange {
			return true, nil
		}
		return false, err
	}
	return true, nil
}

func (m *migration) Down() (bool, error) {
	err := m.Migrate.Down()
	if err != nil {
		if err == migrate.ErrNoChange {
			return true, nil
		}
		return false, err
	}
	return true, err
}

// Suite struct for MySQL Suite
type Suite struct {
	suite.Suite
	DSN                     string
	DBConn                  *sql.DB
	Migration               *migration
	MigrationLocationFolder string
	DBName                  string
}

// SetupSuite setup at the beginning of test
func (s *Suite) SetupSuite() {
	var err error
	s.DBConn, err = sql.Open(postgresDriver, s.DSN)
	s.Require().NoError(err)
	err = s.DBConn.Ping()
	s.Require().NoError(err)
	s.Migration, err = runMigration(s.DBConn, s.MigrationLocationFolder)
	s.Require().NoError(err)
}

// TearDownSuite teardown at the end of test
func (s *Suite) TearDownSuite() {
	err := s.DBConn.Close()
	s.Require().NoError(err)
}

func runMigration(dbConn *sql.DB, migrationsFolderLocation string) (*migration, error) {
	dataPath := []string{}
	dataPath = append(dataPath, "file://")
	dataPath = append(dataPath, migrationsFolderLocation)

	pathToMigrate := strings.Join(dataPath, "")

	driver, err := _postgres.WithInstance(dbConn, &_postgres.Config{})
	if err != nil {
		return nil, err
	}

	m, err := migrate.NewWithDatabaseInstance(pathToMigrate, migrationDbName, driver)
	if err != nil {
		return nil, err
	}
	return &migration{Migrate: m}, nil
}
