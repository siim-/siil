package site

import (
	"github.com/siim-/siil/entity"
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
