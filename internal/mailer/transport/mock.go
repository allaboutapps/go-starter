package transport

import (
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"allaboutapps.dev/aw/go-starter/internal/util"
	"github.com/jordan-wright/email"
)

const defaultWaitTimeout = time.Second * 10

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
		OnMailSent: func(_ email.Email) {},
	}
}

func (m *MockMailTransport) Send(mail *email.Email) error {
	m.Lock()
	defer m.Unlock()

	// Calling wg.Done might panic leaving a user clueless what was the reason of test failure.
	// We will add more information before exiting.
	defer func() {
		rcp := recover()
		if rcp == nil {
			return
		}

		err, ok := rcp.(error)
		if !ok {
			err = fmt.Errorf("%v", rcp)
		}

		log.Fatalf("Unexpected email sent! MockMailTransport panicked: %s", err)
	}()

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
	if err := util.WaitTimeout(&m.wg, defaultWaitTimeout); errors.Is(err, util.ErrWaitTimeout) {
		panic(fmt.Sprintf("Some emails are missing, sent: %v", len(m.GetSentMails())))
	}
}

func (m *MockMailTransport) WaitWithTimeout(timeout time.Duration) {
	if err := util.WaitTimeout(&m.wg, timeout); errors.Is(err, util.ErrWaitTimeout) {
		log.Fatalf("Some emails are missing, found: %v", len(m.GetSentMails()))
	}
}
