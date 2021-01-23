package model

import (
	"database/sql"
	"majoo-backend/usecase/viewmodel"
	"strings"
	"time"
)

var (
	// DefaultUserBy ...
	DefaultUserBy = "def.updated_at"
	// UserBy ...
	UserBy = []string{
		"def.created_at", "def.updated_at",
	}

	userSelectString = `SELECT def.id, def.name, def.user_name, def.password, def.created_at, def.updated_at, def.deleted_at, def.image_path FROM "user" def`
)

func (model userModel) scanRows(rows *sql.Rows) (d UserEntity, err error) {
	err = rows.Scan(
		&d.ID, &d.Name, &d.UserName, &d.Password, &d.CreatedAt,
		&d.UpdatedAt, &d.DeletedAt, &d.ImagePath,
	)

	return d, err
}

func (model userModel) scanRow(row *sql.Row) (d UserEntity, err error) {
	err = row.Scan(
		&d.ID, &d.Name, &d.UserName, &d.Password, &d.CreatedAt,
		&d.UpdatedAt, &d.DeletedAt, &d.ImagePath,
	)

	return d, err
}

// userModel ...
type userModel struct {
	DB *sql.DB
}

// IUser ...
type IUser interface {
	FindAll(search string, offset, limit int, by, sort string) ([]UserEntity, int, error)
	FindByID(id string) (UserEntity, error)
	FindByUserName(userName string) (UserEntity, error)
	Store(body viewmodel.UserVM, changedAt time.Time) (string, error)
	Update(id string, body viewmodel.UserVM, changedAt time.Time) (string, error)
	UpdateImage(id, imagepath string, changedAt time.Time) (string, error)
	Destroy(id string, changedAt time.Time) (string, error)
}

// UserEntity ....
type UserEntity struct {
	ID        string         `db:"id"`
	Name      sql.NullString `db:"name"`
	UserName  string         `db:"user_name"`
	ImagePath sql.NullString `db:"image_path"`
	Password  string         `db:"password"`
	CreatedAt string         `db:"created_at"`
	UpdatedAt string         `db:"updated_at"`
	DeletedAt sql.NullString `db:"deleted_at"`
}

// NewUserModel ...
func NewUserModel(db *sql.DB) IUser {
	return &userModel{DB: db}
}

// FindAll ...
func (model userModel) FindAll(search string, offset, limit int, by, sort string) (res []UserEntity, count int, err error) {
	query := userSelectString + ` WHERE def."deleted_at" IS NULL AND (
	LOWER (def."user_name") LIKE $1 
	OR LOWER ( def."name") LIKE $1 
	) ORDER BY ` + by + ` ` + sort + ` OFFSET $2 LIMIT $3`
	rows, err := model.DB.Query(query, `%`+strings.ToLower(search)+`%`, offset, limit)
	if err != nil {
		return res, count, err
	}
	defer rows.Close()

	for rows.Next() {
		d, err := model.scanRows(rows)
		if err != nil {
			return res, count, err
		}
		res = append(res, d)
	}
	err = rows.Err()
	if err != nil {
		return res, count, err
	}

	query = `SELECT COUNT(def."id") FROM "user" def WHERE def."deleted_at" IS NULL AND (
			LOWER (def."user_name" ) LIKE $1 
			OR LOWER ( def."name" ) like $1 )`
	err = model.DB.QueryRow(query, `%`+strings.ToLower(search)+`%`).Scan(&count)

	return res, count, err
}

// FindByID ...
func (model userModel) FindByID(id string) (res UserEntity, err error) {
	query := userSelectString + ` WHERE def."deleted_at" IS NULL AND def."id" = $1
		ORDER BY def."created_at" DESC LIMIT 1`
	row := model.DB.QueryRow(query, id)
	res, err = model.scanRow(row)

	return res, err
}

// FindByUserName ...
func (model userModel) FindByUserName(userName string) (res UserEntity, err error) {
	query := userSelectString + ` WHERE def."deleted_at" IS NULL  AND LOWER (def."user_name") = $1 ORDER BY def."created_at" DESC  LIMIT 1`
	row := model.DB.QueryRow(query, strings.ToLower(userName))
	res, err = model.scanRow(row)

	return res, err
}

// Store ...
func (model userModel) Store(body viewmodel.UserVM, changedAt time.Time) (res string, err error) {
	sql := `INSERT INTO "user" ("name", "user_name", "password", "created_at", "updated_at"
		) VALUES($1, $2, $3, $4, $4) RETURNING "id"`
	err = model.DB.QueryRow(sql, body.Name, body.UserName, body.Password, changedAt).Scan(&res)

	return res, err
}

// Update ...
func (model userModel) Update(id string, body viewmodel.UserVM, changedAt time.Time) (res string, err error) {
	sql := `UPDATE "user" SET "name" = $1, "user_name" = $2, "updated_at" = $3 WHERE "deleted_at" IS NULL
		AND "id" = $4 RETURNING "id"`
	err = model.DB.QueryRow(sql, body.Name, body.UserName, changedAt, id).Scan(&res)

	return res, err
}

// UpdateImage ...
func (model userModel) UpdateImage(id, ImagePath string, changedAt time.Time) (res string, err error) {
	sql := `UPDATE "user" SET "image_path" = $1, "updated_at" = $2 WHERE "deleted_at" IS NULL
		AND "id" = $3 RETURNING "id"`
	err = model.DB.QueryRow(sql, ImagePath, changedAt, id).Scan(&res)

	return res, err
}

// Destroy ...
func (model userModel) Destroy(id string, changedAt time.Time) (res string, err error) {
	sql := `UPDATE "user" SET "updated_at" = $1, "deleted_at" = $1
		WHERE "deleted_at" IS NULL AND "id" = $2 RETURNING "id"`
	err = model.DB.QueryRow(sql, changedAt, id).Scan(&res)

	return res, err
}
