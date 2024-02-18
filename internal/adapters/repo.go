package adapters

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pkg/errors"
	"virtualization-technologies/internal/entity"
	"virtualization-technologies/internal/entity/user"
)

const (
	pgDuplicateErrCode = "23505"
)

type userRepo struct {
	db *pgx.Conn
}

func New(db *pgx.Conn) userRepo {
	return userRepo{db: db}
}

const (
	getAllUsersQuery = `select * from users order by id, username, email`
	getUserQuery     = `select * from users where id = $1`
	insertUserQuery  = `insert into users (username, email) values ($1, $2) returning id`
	updateUserQuery  = `update users set username = $1, email = $2 where id = $3`
	deleteUserQuery  = `delete from users where id = $1`
)

func (u userRepo) GetAll(ctx context.Context, offset uint64, count uint64) ([]user.User, error) {
	rows, err := u.db.Query(ctx, getAllUsersQuery)
	if err != nil {
		return nil, errors.WithMessage(err, "select all users")
	}
	users := make([]user.User, 0)
	defer rows.Close()
	for i := uint64(0); rows.Next(); i++ {
		if i < offset {
			continue
		}
		if uint64(len(users)) == count {
			break
		}
		var user user.User
		err := rows.Scan(&user.Id, &user.Name, &user.Email)
		if err != nil {
			return nil, errors.WithMessage(err, "scan row")
		}
		users = append(users, user)
	}
	return users, nil
}

func (u userRepo) Get(ctx context.Context, id int) (*user.User, error) {
	user := new(user.User)
	err := u.db.QueryRow(ctx, getUserQuery, id).Scan(&user.Id, &user.Name, &user.Email)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return nil, entity.ErrUserNotFound
	case err != nil:
		return nil, errors.WithMessagef(err, "select user by id %d", id)
	}
	return user, nil
}

func (u userRepo) Create(ctx context.Context, user user.User) (int, error) {
	err := u.db.QueryRow(ctx, insertUserQuery, user.Name, user.Email).Scan(&user.Id)
	pgErr := new(pgconn.PgError)
	switch {
	case errors.As(err, &pgErr) && pgErr.Code == pgDuplicateErrCode:
		return 0, entity.ErrEmailIsAlreadyTaken
	case err != nil:
		return 0, errors.WithMessage(err, "insert new user")
	}
	return user.Id, nil
}

func (u userRepo) Update(ctx context.Context, updatedUser user.User) error {
	user := user.User{}
	err := u.db.QueryRow(ctx, getUserQuery, updatedUser.Id).Scan(&user.Id, &user.Name, &user.Email)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return entity.ErrUserNotFound
	case err != nil:
		return errors.WithMessagef(err, "select user by id %d", updatedUser.Id)
	}
	/* nothing to update */
	if user.Name == updatedUser.Name && user.Email == updatedUser.Email {
		return nil
	}
	_, err = u.db.Exec(ctx, updateUserQuery, updatedUser.Name, updatedUser.Email, updatedUser.Id)
	pgErr := new(pgconn.PgError)
	switch {
	case errors.As(err, &pgErr) && pgErr.Code == pgDuplicateErrCode:
		return entity.ErrEmailIsAlreadyTaken
	case err != nil:
		return errors.WithMessagef(err, "update user with id %d", updatedUser.Id)
	}
	return nil
}

func (u userRepo) Delete(ctx context.Context, id int) (*user.User, error) {
	user := new(user.User)
	err := u.db.QueryRow(ctx, getUserQuery, id).Scan(&user.Id, &user.Name, &user.Email)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return nil, entity.ErrUserNotFound
	case err != nil:
		return nil, errors.WithMessagef(err, "select user by id %d", id)
	}
	if _, err := u.db.Exec(ctx, deleteUserQuery, id); err != nil {
		return nil, errors.WithMessagef(err, "delete user with id %d", id)
	}
	return user, nil
}
