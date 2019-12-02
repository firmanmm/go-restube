package service

import (
	"crypto/rand"
	"crypto/sha512"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"log"
)

type IAuthStorage interface {
	Insert(auth *Authentication) error
	Update(auth *Authentication) error
	FindBySessionID(sessionID string) (*Authentication, error)
	FindByUsername(username string) (*Authentication, error)
	List(limit uint, offset uint) ([]*Authentication, error)
}

type Authentication struct {
	Username       string
	Password       []byte
	SessionID      string
	ByteDownloaded uint
}

type AuthenticationService struct {
	storage IAuthStorage
}

func (a *AuthenticationService) NewAuthentication(username, password string) error {
	oldUser, err := a.storage.FindByUsername(username)
	if err != nil {
		log.Println(err.Error())
		return errors.New("Failed to perform authentication creation")
	}
	if oldUser != nil {
		return errors.New("Username already taken")
	}
	auth := new(Authentication)
	auth.Username = username
	hashedPass := sha512.Sum512([]byte(password)) //SECURITY: Not Secure, Refactor Later
	auth.Password = hashedPass[:]
	sessionID := make([]byte, 64)
	rand.Read(sessionID)
	auth.SessionID = base64.URLEncoding.EncodeToString(sessionID)
	if err := a.storage.Insert(auth); err != nil {
		return errors.New("Failed to create authentication")
	}
	return nil
}

func (a *AuthenticationService) Authenticate(username, password string) (string, error) {
	auth, err := a.storage.FindByUsername(username)
	if err != nil {
		log.Println(err.Error())
		return "", errors.New("Failed to perform authentication")
	}
	if auth == nil {
		return "", errors.New("Username and password combination not found")
	}
	hashedPass := sha512.Sum512([]byte(password))
	if subtle.ConstantTimeCompare(hashedPass[:], auth.Password) != 1 {
		return "", errors.New("Username and Password combination not found")
	}
	return auth.SessionID, nil
}

func (a *AuthenticationService) AddByteDownloaded(sessionID string, size uint) error {
	auth, err := a.storage.FindBySessionID(sessionID)
	if err != nil {
		return errors.New("Failed to perform byte addition")
	}
	if auth == nil {
		return errors.New("Authentication not found, please login!")
	}
	auth.ByteDownloaded += size
	if err := a.storage.Update(auth); err != nil {
		log.Println(err.Error())
		return errors.New("Failed to perform byte addition")
	}
	return nil
}

func (a *AuthenticationService) GetByteDownloaded(sessionID string) (uint, error) {
	auth, err := a.storage.FindBySessionID(sessionID)
	if err != nil {
		return 0, errors.New("Failed to retreive byte")
	}
	if auth == nil {
		return 0, errors.New("Authentication not found, please login!")
	}
	return auth.ByteDownloaded, nil
}

func (a *AuthenticationService) CheckAuthentication(sessionID string) error {
	auth, err := a.storage.FindBySessionID(sessionID)
	if err != nil {
		return errors.New("Failed to perform authentication")
	}
	if auth == nil {
		return errors.New("Authentication not found, please login!")
	}
	return nil
}

func (a *AuthenticationService) FindBySessionID(sessionID string) (*Authentication, error) {
	auth, err := a.storage.FindBySessionID(sessionID)
	if err != nil {
		return nil, errors.New("Failed to perform authentication")
	}
	if auth == nil {
		return nil, errors.New("Authentication not found, please login!")
	}
	return auth, nil
}

func NewAuthenticationService(storage IAuthStorage) *AuthenticationService {
	instance := new(AuthenticationService)
	instance.storage = storage
	return instance
}
