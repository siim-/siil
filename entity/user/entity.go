package user

import (
	"strings"

	"github.com/siim-/siil/cert"
	"github.com/siim-/siil/entity"
)

type Entity struct {
	Id        int    `json:"id"`
	Code      string `json:"code"`
	FirstName string `json:"first_name" db:"first_name"`
	LastName  string `json:"last_name" db:"last_name"`
}

func Find(userCert *cert.Cert) (*Entity, error) {
	usr := Entity{}
	if err := entity.DB.Get(&usr, "SELECT * FROM user WHERE code=?", userCert.SerialNumber); err != nil {
		return nil, err
	}
	return &usr, nil
}

func FindById(id int) (*Entity, error) {
	usr := Entity{}
	if err := entity.DB.Get(&usr, "SELECT * FROM user WHERE id=?", id); err != nil {
		return nil, err
	}
	return &usr, nil
}

func FindOrCreate(userCert *cert.Cert) (*Entity, error) {
	usr := Entity{}
	if err := entity.DB.Get(&usr, "SELECT * FROM user WHERE code=?", userCert.SerialNumber); err != nil {
		if strings.Contains(err.Error(), "no rows") {
			usr = Entity{
				Code:      userCert.SerialNumber,
				FirstName: userCert.FirstName,
				LastName:  userCert.LastName,
			}
			if _, err := entity.DB.NamedExec("INSERT INTO user (code, first_name, last_name) VALUES (:code, :first_name, :last_name)", &usr); err != nil {
				return nil, err
			} else {
				return FindOrCreate(userCert)
			}
		} else {
			return nil, err
		}
	}
	return &usr, nil
}
