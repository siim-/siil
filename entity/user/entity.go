package user

import (
	"strings"

	"github.com/siim-/siil/cert"
	"github.com/siim-/siil/entity"
)

type Entity struct {
	Id        int
	Code      string
	FirstName string `db:"first_name"`
	LastName  string `db:"last_name"`
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
			}
		}
	}
	return &usr, nil
}
