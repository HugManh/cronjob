package postgres

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupMockDB(t *testing.T) (*sql.DB, *gorm.DB, sqlmock.Sqlmock) {
	sqlDB, mock, err := sqlmock.New()
	assert.NoError(t, err)

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn:                 sqlDB,
		PreferSimpleProtocol: true,
	}), &gorm.Config{})
	assert.NoError(t, err)

	return sqlDB, gormDB, mock
}

func TestOpen(t *testing.T) {
	cfg := Config{
		Host:     "localhost",
		Port:     5432,
		User:     "user",
		Password: "password",
		DBName:   "testdb",
		SSLMode:  "disable",
	}

	db, err := Open(cfg)
	if err == nil {
		assert.NotNil(t, db)
		assert.NotNil(t, db.DB())
		return
	}

	assert.Error(t, err)
}

func TestOpenWithInvalidConfig(t *testing.T) {
	tests := []struct {
		name string
		cfg  Config
	}{
		{
			name: "zero port",
			cfg: Config{
				Host:     "unknown_host_xyz",
				Port:     0,
				User:     "user",
				Password: "password",
				DBName:   "testdb",
				SSLMode:  "disable",
			},
		},
		{
			name: "empty host",
			cfg: Config{
				Host:     "",
				Port:     5432,
				User:     "user",
				Password: "password",
				DBName:   "testdb",
				SSLMode:  "disable",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, err := Open(tt.cfg)
			assert.Error(t, err)
			assert.Nil(t, db)
		})
	}
}

func TestDatabaseDB(t *testing.T) {
	_, gormDB, _ := setupMockDB(t)

	dbInstance := &database{
		client: gormDB,
	}

	client := dbInstance.DB()
	assert.NotNil(t, client)
	assert.Equal(t, gormDB, client)
}

func TestDatabaseDBWithNilClient(t *testing.T) {
	dbInstance := &database{
		client: nil,
	}

	client := dbInstance.DB()
	assert.Nil(t, client)
}

func TestDatabaseClose(t *testing.T) {
	sqlDB, gormDB, mock := setupMockDB(t)

	dbInstance := &database{
		client: gormDB,
	}

	mock.ExpectClose()

	err := dbInstance.Close()
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())

	_ = sqlDB.Close()
}

func TestDatabaseCloseAlreadyClosed(t *testing.T) {
	sqlDB, gormDB, mock := setupMockDB(t)

	dbInstance := &database{
		client: gormDB,
	}

	mock.ExpectClose()

	err := dbInstance.Close()
	assert.NoError(t, err)

	err = dbInstance.Close()
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())

	_ = sqlDB.Close()
}
