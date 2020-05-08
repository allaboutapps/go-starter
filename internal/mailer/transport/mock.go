package transport

import (
	"sync"

	"github.com/jordan-wright/email"
)

type MockMailTransport struct {
	sync.RWMutex
	mails []*email.Email
}

func NewMock() *MockMailTransport {
	return &MockMailTransport{
		RWMutex: sync.RWMutex{},
		mails:   make([]*email.Email, 0),
	}
}

func (m *MockMailTransport) Send(mail *email.Email) error {
	m.Lock()
	defer m.Unlock()

	m.mails = append(m.mails, mail)

	return nil
}

func (m *MockMailTransport) GetLastSentMail() *email.Email {
	m.RLock()
	defer m.RUnlock()

	if len(m.mails) == 0 {
		return nil
	}

	return m.mails[len(m.mails)-1]
}

func (m *MockMailTransport) GetSentMails() []*email.Email {
	m.RLock()
	defer m.RUnlock()

	return m.mails
}
