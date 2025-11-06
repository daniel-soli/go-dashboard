package auth

import (
	"errors"
	"sync"

	"golang.org/x/crypto/bcrypt"
)

// User represents a user in the system
type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"-"` // Never serialize password
}

// UserStore manages user storage (in-memory for now)
type UserStore struct {
	users  map[string]*User // key: username or email
	mu     sync.RWMutex
	nextID int
}

var (
	store *UserStore
	once  sync.Once
)

// GetUserStore returns the singleton user store
func GetUserStore() *UserStore {
	once.Do(func() {
		store = &UserStore{
			users:  make(map[string]*User),
			nextID: 1,
		}
		// Create default admin user
		store.createDefaultUser()
	})
	return store
}

// createDefaultUser creates a default admin user for development
func (s *UserStore) createDefaultUser() {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	admin := &User{
		ID:       1,
		Username: "admin",
		Email:    "admin@dashboard.com",
		Password: string(hashedPassword),
	}
	s.users["admin"] = admin
	s.users["admin@dashboard.com"] = admin
}

// CreateUser creates a new user
func (s *UserStore) CreateUser(username, email, password string) (*User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Check if user already exists
	if _, exists := s.users[username]; exists {
		return nil, errors.New("username already exists")
	}
	if _, exists := s.users[email]; exists {
		return nil, errors.New("email already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &User{
		ID:       s.nextID,
		Username: username,
		Email:    email,
		Password: string(hashedPassword),
	}

	s.users[username] = user
	s.users[email] = user
	s.nextID++

	return user, nil
}

// GetUserByUsernameOrEmail retrieves a user by username or email
func (s *UserStore) GetUserByUsernameOrEmail(identifier string) (*User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	user, exists := s.users[identifier]
	if !exists {
		return nil, errors.New("user not found")
	}

	return user, nil
}

// VerifyPassword checks if the provided password matches the user's password
func (s *UserStore) VerifyPassword(user *User, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	return err == nil
}
