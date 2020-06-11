package postgres

import (
	"errors"
	"fmt"

	log "github.com/sirupsen/logrus"
)

/* Ref: https://www.alexedwards.net/blog/practical-persistence-sql */
func (rdb *RDB) CreatePoster(p Poster) error {
	statement := `INSERT INTO poster (email, username, password, role) VALUES (?,?,?,?);`
	statement = rdb.Poolx.Rebind(statement)

	_, err := rdb.Poolx.Exec(statement,
		p.Email,
		p.Username,
		p.Password,
		p.Role,
	)
	if err != nil {
		log.Errorf("***** [POSTGRES:%s][FAIL] ***** Cannot execute INSERT operation:: %v", "CreatePoster", err)
		return err
	}

	return nil
}

func (rdb *RDB) SelectPosterByEmail(email string) (Poster, error) {
	p := &Poster{}
	statement := `SELECT email, username, password, role, bio, image FROM poster WHERE email = ?;`
	statement = rdb.Poolx.Rebind(statement)

	if err := rdb.Poolx.Get(p, statement, email); err != nil {
		log.Errorf("***** [POSTGRES:%s][FAIL] ***** Cannot execute SELECT operation:: %v", "SelectPosterByEmail", err)
		return *p, err
	}

	return *p, nil
}

func (rdb *RDB) SelectPosterByUsername(username string) (Poster, error) {
	p := &Poster{}
	statement := `SELECT email, username, password, role, bio, image FROM poster WHERE username = ?;`
	statement = rdb.Poolx.Rebind(statement)

	if err := rdb.Poolx.Get(p, statement, username); err != nil {
		log.Errorf("***** [POSTGRES:%s][FAIL] ***** Cannot execute SELECT operation:: %v", "SelectPosterByUsername", err)
		return *p, err
	}

	return *p, nil
}

func (rdb *RDB) UpdatePoster(email string, r *UpdateReq) (Poster, error) {
	p, err := rdb.SelectPosterByEmail(email)
	if err != nil {
		log.Errorf("***** [POSTGRES:%s][FAIL] ***** Cannot execute SELECT operation:: %v", "UpdatePoster", err)
		return Poster{}, err
	}

	if r.Password != "" {
		p.Password = r.Password
	}
	if r.Username != "" {
		p.Username = r.Username
	}
	if r.Image != "" {
		p.Image = r.Image
	}
	if r.Bio != "" {
		p.Bio = r.Bio
	}

	statement := `UPDATE poster SET email = ?, username = ?, password = ?, image = ? , bio = ? WHERE email = ?;`
	statement = rdb.Poolx.Rebind(statement)
	_, err = rdb.Poolx.Exec(statement, p.Email, p.Username, p.Password, p.Image, p.Bio, p.Email)
	if err != nil {
		log.Errorf("***** [POSTGRES:%s][FAIL] ***** Cannot execute UPDATE operation:: %v", "UpdatePoster", err)
		return Poster{}, err
	}

	return p, nil
}

func (rdb *RDB) FetchFollowersByEmail(email string) ([]string, error) {
	var f []string
	statement := `SELECT follower FROM follower WHERE email = ?;`
	statement = rdb.Poolx.Rebind(statement)

	if err := rdb.Poolx.Select(&f, statement, email); err != nil {
		log.Errorf("***** [POSTGRES:%s][FAIL] ***** Cannot execute SELECT operation:: %v", "SelectPosterByUsername", err)
		return f, err
	}

	return f, nil
}

func (rdb *RDB) FollowPoster(poster string, follower string) error {
	statement := `INSERT INTO follower (email, follower) VALUES (?,?);`
	statement = rdb.Poolx.Rebind(statement)

	_, err := rdb.Poolx.Exec(statement, poster, follower)
	if err != nil {
		log.Errorf("***** [POSTGRES:%s][FAIL] ***** Cannot execute INSERT operation:: %v", "CreatePoster", err)
		return err
	}

	return nil
}

func (rdb *RDB) UnFollowPoster(email string, follower string) error {
	statement := `DELETE FROM follower WHERE email = ? AND follower = ?;`
	statement = rdb.Poolx.Rebind(statement)

	result, err := rdb.Poolx.Exec(statement, email, follower)
	if err != nil {
		log.Errorf("***** [POSTGRES:%s][FAIL] ***** Cannot execute DELETE operation:: %v", "UnFollowPoster", err)
		return err
	}
	if row, _ := result.RowsAffected(); row < 1 {
		log.Errorf("***** [POSTGRES:%s][FAIL] ***** Cannot execute DELETE operation", "UnFollowPoster")
		return errors.New(fmt.Sprintf("row(s) affected: %d", row))
	}

	return nil
}
