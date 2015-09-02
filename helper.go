package asapp

import (
	"encoding/base32"
	"net/http"
	"reflect"
	"strings"

	"github.com/golang/glog"
	"github.com/gorilla/securecookie"
	"golang.org/x/crypto/bcrypt"
)

const (
	PerPage = 10
)

// GenerateHash creates a hash from a given string using bcrypt with a
// cost which makes brute force cracking hard. bcrypt.DefaultCost
// uses over 0.1 second. Use MinCost (about 40ms) for now.
func GenerateHash(password string) string {
	encrypted, err := bcrypt.GenerateFromPassword([]byte(password),
		bcrypt.MinCost)
	if err != nil {
		glog.Warning(err)
		return password
	}

	return string(encrypted)
}

// CompareHash compares encrypted hash with the plain string. Returns
// true if the hash is generated from the password.
func CompareHash(encrypted, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(encrypted),
		[]byte(password)) == nil
}

type JSONError struct {
	Code    int          `json:"code"`
	Message string       `json:"message"`
	Errors  ErrorDetails `json:"errors"`
}

type JSONResult struct {
	Status      string      `json:"status"`
	Data        interface{} `json:"data"`
	Error       JSONError   `json:"error"`
	Page        int         `json:"page"`
	CurrentPage int         `json:"current_page"`
	PerPage     int         `json:"per_page"`
}

// NewJSONResult returns a unified JSON response.
func NewJSONResult(v interface{}, page int) *JSONResult {
	val := reflect.ValueOf(v)
	return &JSONResult{
		Status:      "success",
		Error:       JSONError{Code: http.StatusOK},
		Data:        v,
		CurrentPage: val.Len(),
		PerPage:     PerPage,
		Page:        page,
	}
}

func NewErrorJSONResult(err JSONError) *JSONResult {
	return &JSONResult{
		Status:      "fail",
		Error:       err,
		Data:        []struct{}{},
		CurrentPage: 0,
		PerPage:     PerPage,
		Page:        1,
	}
}

// GenerateRandomKey generates random key with only alphabetical letters.
func GenerateRandomKey() string {
	rb := securecookie.GenerateRandomKey(32)
	return strings.TrimRight(base32.StdEncoding.EncodeToString(rb), "=")
}

// PersistEnvelope saves thread and 2 envelopes from a PublicEnvelope.
func PersistEnvelope(p PublicEnvelope, us UserService, es EnvelopeService,
	ts ThreadService) CompoundError {
	sender, err := us.GetByUsername(p.Author)
	if err != nil {
		return NewServerError(err.Error())
	}
	recipient, err := us.GetByUsername(p.Recipient)
	if err != nil {
		return NewServerError(err.Error())
	}

	// persist envelopes
	senderEnv, recipientEnv := NewEnvelope(sender.ID, recipient.ID,
		p.Message, p.MessageType)
	if err = es.Create(senderEnv); err != nil {
		return NewServerError(err.Error())
	}
	if err = es.Create(recipientEnv); err != nil {
		return NewServerError(err.Error())
	}

	// persist threads
	t1, t2 := NewThread(sender.ID, recipient.ID, sender.Username, p.Message)
	if _, err = ts.Upsert(t1); err != nil {
		return NewServerError(err.Error())
	}
	if _, err = ts.Upsert(t2); err != nil {
		return NewServerError(err.Error())
	}
	return nil
}
