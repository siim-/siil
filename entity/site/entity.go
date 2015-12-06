package site

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"log"
	"regexp"

	"github.com/siim-/siil/entity"
	"github.com/siim-/siil/entity/user"
)

const (
	CLIENT_ID_LENGTH   = 64
	PRIVATE_KEY_LENGTH = 128
)

var (
	SIIL_SITE_ID string
)

//The site entity
type Entity struct {
	ClientId    string `db:"client_id"`
	PrivateKey  string `db:"private_id"`
	Owner       uint   `db:"owner"`
	Name        string `db:"name"`
	Domain      string `db:"domain"`
	CallbackURL string `db:"callback_url"`
	CancelURL   string `db:"cancel_url"`
}

//The site entry
type Entry struct {
	Owner       uint
	Name        string
	Domain      string
	CallbackURL string
	CancelURL   string
}

func (e *Entity) Load() error {
	loaded := Entity{}
	switch {
	case len(e.ClientId) != 0:
		if err := entity.DB.Get(&loaded, "SELECT * FROM site WHERE client_id=?", e.ClientId); err != nil {
			return err
		}
		*e = loaded
	case len(e.Domain) != 0:
		if err := entity.DB.Get(&loaded, "SELECT * FROM site WHERE domain=?", e.Domain); err != nil {
			return err
		}
		*e = loaded
	}

	return nil
}

//Does the user have any active sessions for the site?
func (e *Entity) HasActiveSessionFor(u *user.Entity) bool {
	_, err := entity.DB.NamedExec(
		"SELECT * FROM session WHERE user_id = :user AND site_id = :site AND expires_at > NOW()",
		map[string]interface{}{
			"site": e.ClientId,
			"user": u.Id,
		},
	)
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

func NewSite(entry *Entry) (*Entity, error) {
	if validEntry(entry) {
		clientID, err := createRandomKey(64)
		if err != nil {
			return nil, err
		}
		privateKey, err := createRandomKey(128)
		if err != nil {
			return nil, err
		}

		site := Entity{
			ClientId:    clientID,
			PrivateKey:  privateKey,
			Owner:       entry.Owner,
			Name:        entry.Name,
			Domain:      entry.Domain,
			CallbackURL: entry.CallbackURL,
			CancelURL:   entry.CancelURL,
		}

		if _, err := entity.DB.NamedExec("INSERT INTO site (client_id, private_id, owner, name, domain, callback_url, cancel_url) VALUES (:client_id, :private_id, :owner, :name, :domain, :callback_url, :cancel_url)", &site); err != nil {
			return nil, err
		}

		return &site, nil
	} else {
		return nil, errors.New("Invalid entry.")
	}
}

func GetUsersSites(userID uint) ([]Entity, error) {
	sites := []Entity{}
	if err := entity.DB.Select(&sites, "SELECT * FROM site WHERE owner=?", userID); err != nil {
		return nil, err
	}
	return sites, nil
}

func validEntry(entry *Entry) bool {
	var (
		nameRegex   string = "[A-Za-z1-9]+"
		domainRegex string = "[a-zA-Z0-9][a-zA-Z0-9-]{1,61}[a-zA-Z0-9]\\.[a-zA-Z]{2,}"
		urlRegex    string = "((([A-Za-z]{3,9}:(?:\\/\\/)?)(?:[\\-;:&=\\+\\$,\\w]+@)?[A-Za-z0-9\\.\\-]+|(?:www\\.|[\\-;:&=\\+\\$,\\w]+@)[A-Za-z0-9\\.\\-]+)((?:\\/[\\+~%\\/\\.\\w\\-_]*)?\\??(?:[\\-\\+=&;%@\\.\\w_]*)#?(?:[\\.\\!\\/\\\\\\w]*))?)"
	)

	matched, err := regexp.MatchString(nameRegex, entry.Name)
	if !matched || err != nil {
		return false
	}

	matched, err = regexp.MatchString(domainRegex, entry.Domain)
	if !matched || err != nil {
		return false
	}

	matched, err = regexp.MatchString(urlRegex, entry.CallbackURL)
	if !matched || err != nil {
		return false
	}

	matched, err = regexp.MatchString(urlRegex, entry.CancelURL)
	if !matched || err != nil {
		return false
	}

	return true
}

func createRandomKey(length int) (string, error) {
	var buffer []byte = make([]byte, hex.DecodedLen(length))
	if _, err := rand.Read(buffer); err != nil {
		return "", err
	} else {
		return hex.EncodeToString(buffer), nil
	}
}
