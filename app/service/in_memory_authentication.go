package service

import "errors"

type InMemoryAuthentication struct {
	sessionIDIdx    map[string]*Authentication
	usernameIdx     map[string]*Authentication
	authentications []*Authentication
}

func (i *InMemoryAuthentication) Insert(auth *Authentication) error {
	i.authentications = append(i.authentications, auth)
	i.usernameIdx[auth.Username] = auth
	i.sessionIDIdx[auth.SessionID] = auth
	return nil
}

func (i *InMemoryAuthentication) Update(auth *Authentication) error {
	//PASS THROUGH DUE TO INMEMORY STORAGE
	return nil
}

func (i *InMemoryAuthentication) FindBySessionID(sessionID string) (*Authentication, error) {
	auth, ok := i.sessionIDIdx[sessionID]
	if !ok {
		return nil, nil
	}
	return auth, nil
}

func (i *InMemoryAuthentication) FindByUsername(username string) (*Authentication, error) {
	auth, ok := i.usernameIdx[username]
	if !ok {
		return nil, nil
	}
	return auth, nil
}

func (i *InMemoryAuthentication) List(limit uint, offset uint) ([]*Authentication, error) {
	if int(limit+offset) > len(i.authentications) {
		return nil, errors.New("Out of bound memory")
	}
	return i.authentications[offset : limit+offset], nil
}

func NewInMemoryAuthentication() *InMemoryAuthentication {
	instance := new(InMemoryAuthentication)
	instance.sessionIDIdx = make(map[string]*Authentication)
	instance.usernameIdx = make(map[string]*Authentication)
	instance.authentications = make([]*Authentication, 0, 8)
	return instance
}
