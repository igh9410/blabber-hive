package user

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type mockDBTX struct {
	db *sql.DB
	tx *sql.Tx
}

func (m *mockDBTX) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return m.db.ExecContext(ctx, query, args...)
}

func (m *mockDBTX) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	return m.db.PrepareContext(ctx, query)
}

func (m *mockDBTX) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return m.db.QueryContext(ctx, query, args...)
}

func (m *mockDBTX) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return m.db.QueryRowContext(ctx, query, args...)
}

func (m *mockDBTX) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	var err error
	m.tx, err = m.db.BeginTx(ctx, opts)
	return m.tx, err
}

func (m *mockDBTX) Commit() error {
	return m.tx.Commit()
}

func (m *mockDBTX) Rollback() error {
	return m.tx.Rollback()
}

func TestCreateUserRepository(t *testing.T) {
	// Create a new mock database connection
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %s", err)
	}
	defer db.Close()

	mockDB := &mockDBTX{db: db}
	r := NewRepository(mockDB)

	mockUser := &User{
		ID:              uuid.New(), // Generate a new UUID
		Username:        "test1",
		Email:           "test1@gmail.com",
		ProfileImageURL: nil, // This can be omitted
		CreatedAt:       time.Now(),
	}

	rows := sqlmock.NewRows([]string{"id"}).AddRow(uuid.New().String())

	mock.ExpectQuery("INSERT INTO users").WithArgs(sqlmock.AnyArg(), mockUser.Username, mockUser.Email, mockUser.ProfileImageURL, sqlmock.AnyArg()).WillReturnRows(rows)

	user, err := r.CreateUser(context.Background(), mockUser)
	assert.NoError(t, err)
	assert.NotNil(t, user)

}
