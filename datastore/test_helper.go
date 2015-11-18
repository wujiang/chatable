package datastore

import "github.com/wujiang/chatable"

const (
	// TODO: add readme to setup test db
	testDB             = "postgres://chatable@localhost:5432/chatable_test?sslmode=disable"
	testSenderFirst    = "Sender"
	testSenderLast     = "Send"
	testSenderUname    = "sender"
	testSenderEmail    = "sender@send.com"
	testSenderPhone    = "1234567890"
	testIPAddr         = "0.0.0.0"
	testPassword       = "password123"
	testRecipientFirst = "Recipient"
	testRecipientLast  = "Receive"
	testRecipientUname = "recipient"
	testRecipientEmail = "recipient@receive.com"
	testRecipientPhone = "0123456789"
)

var (
	testStore = NewDataStore(nil)

	testSender = chatable.NewUser(testSenderFirst, testSenderLast,
		testSenderUname, testPassword, testSenderEmail,
		testSenderPhone, testIPAddr)
	testRecipient = chatable.NewUser(testRecipientFirst, testRecipientLast,
		testRecipientUname, testPassword, testRecipientEmail,
		testRecipientPhone, testIPAddr)
)

func newTestUsers() error {
	if err := testStore.UserStore.Create(testSender); err != nil {
		return err
	}
	if err := testStore.UserStore.Create(testRecipient); err != nil {
		return err
	}
	return nil
}
