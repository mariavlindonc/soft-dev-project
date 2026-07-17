package db

import (
	"testing"

	"backend/domain"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	gdb, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	require.NoError(t, err)

	raw, err := gdb.DB()
	require.NoError(t, err)

	raw.Exec(`CREATE TABLE users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		email TEXT NOT NULL UNIQUE,
		password_hash TEXT NOT NULL,
		role TEXT NOT NULL DEFAULT 'client',
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		deleted_at DATETIME
	)`)

	raw.Exec(`CREATE TABLE events (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		description TEXT,
		image_url TEXT,
		category TEXT,
		location TEXT,
		event_date DATETIME NOT NULL,
		duration_minutes INTEGER NOT NULL DEFAULT 0,
		capacity INTEGER NOT NULL DEFAULT 0,
		tickets_sold INTEGER NOT NULL DEFAULT 0,
		price REAL NOT NULL DEFAULT 0.0,
		status TEXT NOT NULL DEFAULT 'active',
		presale_active INTEGER NOT NULL DEFAULT 0,
		presale_code TEXT,
		presale_start_date DATETIME,
		general_sale_date DATETIME,
		created_by_id INTEGER NOT NULL,
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		deleted_at DATETIME,
		FOREIGN KEY (created_by_id) REFERENCES users(id)
	)`)

	raw.Exec(`CREATE TABLE tickets (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		event_id INTEGER NOT NULL,
		status TEXT NOT NULL DEFAULT 'active',
		purchase_price REAL NOT NULL DEFAULT 0.0,
		purchased_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		cancelled_at DATETIME,
		transferred_at DATETIME,
		transferred_to_id INTEGER,
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id),
		FOREIGN KEY (event_id) REFERENCES events(id),
		FOREIGN KEY (transferred_to_id) REFERENCES users(id)
	)`)

	t.Cleanup(func() {
		raw.Close()
	})
	return gdb
}

func TestUserDAO_Create(t *testing.T) {
	gdb := setupTestDB(t)
	dao := NewUserDAO(gdb)

	t.Run("creates user successfully", func(t *testing.T) {
		u := &domain.User{Name: "Alice", Email: "alice@example.com", PasswordHash: "hash123", Role: "client"}
		err := dao.Create(u)
		require.NoError(t, err)
		assert.NotZero(t, u.ID)
	})

	t.Run("duplicate email fails", func(t *testing.T) {
		u1 := &domain.User{Name: "Bob", Email: "bob@example.com", PasswordHash: "h1", Role: "client"}
		require.NoError(t, dao.Create(u1))
		u2 := &domain.User{Name: "Bob2", Email: "bob@example.com", PasswordHash: "h2", Role: "client"}
		err := dao.Create(u2)
		assert.Error(t, err)
	})
}

func TestUserDAO_FindByID(t *testing.T) {
	gdb := setupTestDB(t)
	dao := NewUserDAO(gdb)

	t.Run("finds existing user", func(t *testing.T) {
		u := &domain.User{Name: "Carol", Email: "carol@example.com", PasswordHash: "h", Role: "client"}
		require.NoError(t, dao.Create(u))
		found, err := dao.FindByID(u.ID)
		require.NoError(t, err)
		assert.Equal(t, "Carol", found.Name)
		assert.Equal(t, "carol@example.com", found.Email)
	})

	t.Run("returns error for non-existent user", func(t *testing.T) {
		_, err := dao.FindByID(99999)
		assert.Error(t, err)
	})
}

func TestUserDAO_FindByEmail(t *testing.T) {
	gdb := setupTestDB(t)
	dao := NewUserDAO(gdb)

	t.Run("finds existing user by email", func(t *testing.T) {
		u := &domain.User{Name: "Dave", Email: "dave@example.com", PasswordHash: "h", Role: "admin"}
		require.NoError(t, dao.Create(u))
		found, err := dao.FindByEmail("dave@example.com")
		require.NoError(t, err)
		assert.Equal(t, "Dave", found.Name)
		assert.Equal(t, "admin", found.Role)
	})

	t.Run("returns error for non-existent email", func(t *testing.T) {
		_, err := dao.FindByEmail("nobody@example.com")
		assert.Error(t, err)
	})
}
