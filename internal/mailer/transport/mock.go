package transport

import (
	"sync"
	"time"

	"allaboutapps.dev/aw/go-starter/internal/util"
	"github.com/jordan-wright/email"
)

type MockMailTransport struct {
	sync.RWMutex
	mails      []*email.Email
	OnMailSent func(mail email.Email) // non pointer to prevent concurrent read errors
	wg         sync.WaitGroup
	expected   int
}

func NewMock() *MockMailTransport {
	return &MockMailTransport{
		RWMutex:    sync.RWMutex{},
		mails:      make([]*email.Email, 0),
		OnMailSent: func(mail email.Email) {},
	}
}

func (m *MockMailTransport) Send(mail *email.Email) error {
	m.Lock()
	defer m.Unlock()

	m.mails = append(m.mails, mail)
	m.OnMailSent(*mail)

	if m.expected > 0 {
		m.wg.Done()
	}

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

// Expect adds the mailCnt to a waitgroup. Done() is called by Send
func (m *MockMailTransport) Expect(mailCnt int) {
	m.expected = mailCnt
	m.wg.Add(mailCnt)
}

// Wait until all expected mails have arrived
func (m *MockMailTransport) Wait() {
	_ = util.WaitTimeout(&m.wg, time.Second*10)
}
