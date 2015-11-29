package session

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/siim-/siil/entity"
	"github.com/siim-/siil/entity/site"
	"github.com/siim-/siil/entity/user"
)

type Entity struct {
	Token     string
	SiteId    string    `db:"site_id"`
	UserId    int       `db:"user_id"`
	CreatedAt time.Time `db:"created_at"`
	ExpiresAt time.Time `db:"expires_at"`
}

func createRandomToken() (string, error) {
	var buffer []byte = make([]byte, hex.DecodedLen(64))
	if _, err := rand.Read(buffer); err != nil {
		return "", err
	} else {
		return hex.EncodeToString(buffer), nil
	}
}

func NewSession(s *site.Entity, u *user.Entity) (*Entity, error) {
	if token, err := createRandomToken(); err != nil {
		return nil, err
	} else {
		sess := Entity{
			Token:     token,
			SiteId:    s.ClientId,
			UserId:    u.Id,
			CreatedAt: time.Now().UTC(),
			ExpiresAt: time.Now().UTC().Add(time.Hour * 2),
		}

		if _, err := entity.DB.NamedExec("INSERT INTO session (token, site_id, user_id, created_at, expires_at) VALUES (:token, :site_id, :user_id, :created_at, :expires_at)", &sess); err != nil {
			return nil, err
		}
		return &sess, nil
	}
}
