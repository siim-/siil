package site

import (
	"log"

	"github.com/siim-/siil/entity"
	"github.com/siim-/siil/entity/user"
)

const (
	CLIENT_ID_LENGTH   = 64
	PRIVATE_KEY_LENGTH = 128
)

//The site entity
type Entity struct {
	ClientId    string `db:"client_id"`
	PrivateKey  string `db:"private_id"`
	Owner       uint
	Name        string
	Domain      string
	CallbackURL string `db:"callback_url"`
	CancelURL   string `db:"cancel_url"`
}

func (e *Entity) Load() error {
	loaded := Entity{}
	err := entity.DB.Get(&loaded, "SELECT * FROM site WHERE client_id=?", e.ClientId)
	if err != nil {
		return err
	}
	*e = loaded
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