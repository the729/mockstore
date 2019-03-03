package mockstore

import (
	"bytes"
	"encoding/gob"
	"errors"
	"net/http"

	"github.com/gorilla/sessions"
)

type MockStore struct {
	ses map[string]*mockSession
}

type mockSession struct {
	*sessions.Session
	SavedValues map[interface{}]interface{}
}

func NewMockStore() *MockStore {
	return &MockStore{
		ses: make(map[string]*mockSession),
	}
}

func (s *MockStore) Get(r *http.Request, name string) (*sessions.Session, error) {
	mses, ok := s.ses[name]
	if ok {
		return mses.Session, nil
	}
	s.InitValues(name, nil)
	return s.ses[name].Session, nil
}

func (s *MockStore) New(r *http.Request, name string) (*sessions.Session, error) {
	return s.Get(r, name)
}

func (stor *MockStore) Save(r *http.Request, w http.ResponseWriter, s *sessions.Session) error {
	for _, mses := range stor.ses {
		if mses.Session == s {
			clone(s.Values, &mses.SavedValues)
			return nil
		}
	}
	return errors.New("session not found")
}

func (s *MockStore) InitValues(name string, values map[interface{}]interface{}) {
	mses := &mockSession{
		Session:     sessions.NewSession(s, name),
		SavedValues: make(map[interface{}]interface{}),
	}
	if values != nil {
		mses.SavedValues = values
	}
	clone(mses.SavedValues, &mses.Values)
	s.ses[name] = mses
}

func (s *MockStore) GetValues(name string) map[interface{}]interface{} {
	mses, ok := s.ses[name]
	if !ok {
		return nil
	}
	return mses.SavedValues
}

// clone deep-copies a to b
func clone(a, b interface{}) {

	buff := new(bytes.Buffer)
	enc := gob.NewEncoder(buff)
	dec := gob.NewDecoder(buff)
	enc.Encode(a)
	dec.Decode(b)
}
