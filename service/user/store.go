package user

import (
	"database/sql"
	"fmt"

	"github.com/varnit-ta/Ecom-API/types"
)

type Store struct {
	db *sql.DB
}

/*
NewStore creates a new instance of Store with the provided database connection.

@param db - *sql.DB: The database connection to be used by the Store.

@return *Store: A new Store instance.
*/
func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) GetUserByID(id int) (*types.User, error) {
	return nil, nil
}

func (s *Store) CreateUser(user types.User) error {
	return nil
}

/*
GetUserByEmail retrieves a user from the database by their email.

@param email - string: The email address to search for in the database.

@return *types.User: The user found by the email, or nil if not found.
@return error: An error, if any occurs during the database query or data retrieval.
*/
func (s *Store) GetUserByEmail(email string) (*types.User, error) {
	rows, err := s.db.Query("SELECT * FROM users WHERE email = ?", email)

	if err != nil {
		return nil, err
	}

	u := new(types.User)

	for rows.Next() {
		u, err = scanRowsIntoUser(rows)

		if err != nil {
			return nil, err
		}
	}

	if u.ID == 0 {
		return nil, fmt.Errorf("user not found")
	}

	return u, nil
}

/*
scanRowsIntoUser scans the rows from the database query result into a User object.

@param rows - *sql.Rows: The result rows from the database query.

@return *types.User: A User object populated with the data from the query result.
@return error: An error, if any occurs during the row scanning process.
*/
func scanRowsIntoUser(rows *sql.Rows) (*types.User, error) {
	user := new(types.User)

	err := rows.Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Password,
		&user.Password,
		&user.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}
