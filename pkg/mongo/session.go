package mongo

import (
	"gopkg.in/mgo.v2"
)

// Session holds mongo session info
type Session struct {
	session *mgo.Session
}

// NewSession creates a new mongo session give host and port
func NewSession(url string) (*Session, error) {
	session, err := mgo.Dial(url)
	if err != nil {
		return nil, err
	}
	return &Session{session}, err
}

// Copy ...
func (s *Session) Copy() *Session {
	return &Session{s.session.Copy()}
}

// GetCollection ...
func (s *Session) GetCollection(db string, col string) *mgo.Collection {
	return s.session.DB(db).C(col)
}

// Close ...
func (s *Session) Close() {
	if s.session != nil {
		s.session.Close()
	}
}

// DropDatabase ...
func (s *Session) DropDatabase(db string) error {
	return s.session.DB(db).DropDatabase()
	/*if s.session != nil {
		return s.session.DB(db).DropDatabase()
	}
	return nil*/
}
