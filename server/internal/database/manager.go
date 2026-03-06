package database

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"server/internal/models"
	"sync"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// DBConnManager manages database connections for multitenancy.
// It maintains a pool of connections to tenant databases.
type DBConnManager struct {
	masterDB  *gorm.DB
	tenantDBs map[string]*gorm.DB
	mu        sync.RWMutex
	basePath  string
}

// tenantModels are the models that exist in tenant databases
var tenantModels = []interface{}{
	&models.User{},
	&models.Role{},
	&models.Permission{},
}

// globalModels are models that exist only in the master database
var globalModels = []interface{}{
	&models.Account{},
	&models.Role{},
	&models.User{},
	&models.Subscription{},
}

// Manager is the global database connection manager
var Manager *DBConnManager

// InitManager initializes the database connection manager.
// basePath is the directory where tenant databases are stored (e.g., "data")
func InitManager(masterDSN string, basePath string) error {
	// Connect to master database
	masterDB, err := gorm.Open(sqlite.Open(masterDSN), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect to master database: %w", err)
	}

	// Auto-migrate global models
	if err := masterDB.AutoMigrate(globalModels...); err != nil {
		return fmt.Errorf("failed to migrate global models: %w", err)
	}

	// Ensure base path exists
	if err := os.MkdirAll(basePath, 0755); err != nil {
		return fmt.Errorf("failed to create base path: %w", err)
	}

	Manager = &DBConnManager{
		masterDB:  masterDB,
		tenantDBs: make(map[string]*gorm.DB),
		basePath:  basePath,
	}

	// Initialize default roles in master
	if err := Manager.initDefaultRoles(); err != nil {
		log.Printf("Warning: failed to initialize default roles: %v", err)
	}

	return nil
}

// GetMasterDB returns the master database connection
func (m *DBConnManager) GetMasterDB() *gorm.DB {
	return m.masterDB
}

// GetConnection returns a database connection for the specified account.
// If the connection doesn't exist, it creates a new one.
// The tenant database is stored at {basePath}/{accountID}/data.db
func (m *DBConnManager) GetConnection(accountID string) (*gorm.DB, error) {
	if accountID == "" {
		return nil, fmt.Errorf("accountID cannot be empty")
	}

	// Check cache first
	m.mu.RLock()
	db, exists := m.tenantDBs[accountID]
	m.mu.RUnlock()

	if exists {
		return db, nil
	}

	// Create new connection
	m.mu.Lock()
	defer m.mu.Unlock()

	// Double-check after acquiring write lock
	if db, exists := m.tenantDBs[accountID]; exists {
		return db, nil
	}

	// Create tenant directory
	tenantPath := filepath.Join(m.basePath, accountID)
	if err := os.MkdirAll(tenantPath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create tenant directory: %w", err)
	}

	// Connect to tenant database
	dbPath := filepath.Join(tenantPath, "data.db")
	tenantDB, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to tenant database: %w", err)
	}

	// Auto-migrate tenant models
	if err := tenantDB.AutoMigrate(tenantModels...); err != nil {
		return nil, fmt.Errorf("failed to migrate tenant models: %w", err)
	}

	// Initialize default roles in tenant
	if err := m.initTenantRoles(tenantDB); err != nil {
		log.Printf("Warning: failed to initialize tenant roles: %v", err)
	}

	m.tenantDBs[accountID] = tenantDB
	return tenantDB, nil
}

// CloseConnection closes a specific tenant connection
func (m *DBConnManager) CloseConnection(accountID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	db, exists := m.tenantDBs[accountID]
	if !exists {
		return nil
	}

	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	if err := sqlDB.Close(); err != nil {
		return err
	}

	delete(m.tenantDBs, accountID)
	return nil
}

// CloseAll closes all database connections
func (m *DBConnManager) CloseAll() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	var lastErr error
	for accountID, db := range m.tenantDBs {
		sqlDB, err := db.DB()
		if err != nil {
			lastErr = err
			continue
		}
		if err := sqlDB.Close(); err != nil {
			lastErr = err
		}
		delete(m.tenantDBs, accountID)
	}

	// Close master
	if m.masterDB != nil {
		sqlDB, err := m.masterDB.DB()
		if err != nil {
			return err
		}
		if err := sqlDB.Close(); err != nil {
			lastErr = err
		}
	}

	return lastErr
}

// initDefaultRoles initializes default roles in the master database
func (m *DBConnManager) initDefaultRoles() error {
	defaultRoles := []models.Role{
		{ID: 1, Name: "Admin"},
		{ID: 2, Name: "Editor"},
		{ID: 3, Name: "Viewer"},
	}

	for _, role := range defaultRoles {
		if err := m.masterDB.FirstOrCreate(&role, models.Role{ID: role.ID}).Error; err != nil {
			return err
		}
	}
	return nil
}

// initTenantRoles initializes default roles in a tenant database
func (m *DBConnManager) initTenantRoles(db *gorm.DB) error {
	defaultRoles := []models.Role{
		{ID: 1, Name: "Admin"},
		{ID: 2, Name: "Editor"},
		{ID: 3, Name: "Viewer"},
	}

	for _, role := range defaultRoles {
		if err := db.FirstOrCreate(&role, models.Role{ID: role.ID}).Error; err != nil {
			return err
		}
	}
	return nil
}

// GetAccountByID retrieves an account from the master database
func (m *DBConnManager) GetAccountByID(accountID string) (*models.Account, error) {
	var account models.Account
	if err := m.masterDB.First(&account, "id = ?", accountID).Error; err != nil {
		return nil, err
	}
	return &account, nil
}

// DeleteTenantDB removes a tenant's database file and closes the connection
func (m *DBConnManager) DeleteTenantDB(accountID string) error {
	// Close connection first
	if err := m.CloseConnection(accountID); err != nil {
		return err
	}

	// Remove database file
	tenantPath := filepath.Join(m.basePath, accountID)
	return os.RemoveAll(tenantPath)
}
