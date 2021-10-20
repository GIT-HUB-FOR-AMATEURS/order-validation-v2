package repository

import (
	"database/sql"

	"clean/internal/entity"
)

type UserPSQL struct {
	db *sql.DB
}

//NewBookMySQL create new repository
func NewUserPSQL(db *sql.DB) *UserPSQL {
	return &UserPSQL{
		db: db,
	}
}

func (r *UserPSQL) Create(u *entity.User) (string, error) {
	stmt, err := r.db.Prepare(`
		INSERT INTO users (id, username, email, password) 
		values($1, $2, md5($3), md5($4))`)
	if err != nil {
		return u.ID, err
	}
	_, err = stmt.Exec(
		u.Username,
		u.Email,
		u.Password,
	)
	if err != nil {
		return u.ID, err
	}
	err = stmt.Close()
	if err != nil {
		return u.ID, err
	}
	return u.ID, nil
}

func (r *UserPSQL) GetbyUsername(username string) (*entity.User, error) {
	stmt, err := r.db.Prepare(`SELECT id, username, email, pswd from users where username = $1`)
	if err != nil {
		return nil, err
	}
	var user entity.User
	row := stmt.QueryRow(username)
	err = row.Scan(&user.ID, &user.Username, &user.Email, &user.Password)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserPSQL) GetbyID(ID string) (*entity.User, error) {
	stmt, err := r.db.Prepare(`SELECT id, username, email from users where ID = $1`)
	if err != nil {
		return nil, err
	}
	var user entity.User
	row := stmt.QueryRow(ID)
	err = row.Scan(&user.ID, &user.Username, &user.Email)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserPSQL) Update(u *entity.User) error {
	_, err := r.db.Exec("UPDATE users SET pswd = md5($1),  username = $2, email = md5($3) where username = $4",
		u.Password, u.Username, u.Email, u.Username)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserPSQL) Search(query string) ([]*entity.User, error) {
	stmt, err := r.db.Prepare(`SELECT id, username, email FROM users WHERE username like $1`)
	if err != nil {
		return nil, err
	}
	var users []*entity.User
	rows, err := stmt.Query("%" + query + "%")
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var u entity.User
		err = rows.Scan(&u.ID, &u.Username, &u.Email)
		if err != nil {
			return nil, err
		}
		users = append(users, &u)
	}

	return users, nil
}

func (r *UserPSQL) List() ([]*entity.User, error) {
	stmt, err := r.db.Prepare(`SELECT ID, username, email FROM users`)
	if err != nil {
		return nil, err
	}
	var users []*entity.User
	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var u entity.User
		err = rows.Scan(&u.ID,
			&u.Email, &u.Username)
		if err != nil {
			return nil, err
		}
		users = append(users, &u)
	}
	return users, nil
}

func (r *UserPSQL) Delete(username string) error {
	_, err := r.db.Exec("DELETE FROM users where username = $1", username)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserPSQL) CustomQuery(query string) (*sql.Rows, error) {
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	return rows, nil
}
