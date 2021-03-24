package store

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	_ "github.com/lib/pq"
)

//Instance of store
type Store struct {
	config         *Config
	db             *gorm.DB
	userRepository *UserRepository
	autoRepository *AutoRepository
}

// Constructor for store
func New(config *Config) *Store {
	return &Store{
		config: config,
	}
}

//Open store method
func (s *Store) Open() error {
	db, err := gorm.Open(postgres.Open(s.config.DatabaseURL), &gorm.Config{})

	if err != nil {
		return err
	}

	s.db = db
	log.Println("Connection to db successfully")
	return nil
}

//Close store method
func (s *Store) Close() {
	//s.db.Close()
}

//Public for UserRepo
func (s *Store) User() *UserRepository {
	if s.userRepository != nil {
		return s.userRepository
	}

	s.userRepository = NewUserRepository(s)

	return s.userRepository
}

//Public for UserRepo
func (s *Store) Auto() *AutoRepository {
	if s.autoRepository != nil {
		return s.autoRepository
	}

	s.autoRepository = NewAutoRepository(s)

	return s.autoRepository
}
