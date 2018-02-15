package dao

import (
	"database/sql"
	"errors"
	"log"

	"github.com/devng/ldbn/types"
)

var db *sql.DB

// ErrInvalidInput is used when the the input of a dao method is invalid
var ErrInvalidInput = errors.New("dao: invalid input")

// ErrNoUpdate is used when there was not update/delete performed
var ErrNoUpdate = errors.New("dao: SQL table was not updated")

// InitDB initializes the DB by loading the file from the provided file path
func InitDB(filePath string) {
	log.Println("Sqlite3 file used: " + filePath)
	var err error
	db, err = sql.Open("sqlite3", filePath)
	if err != nil {
		log.Panic(err)
	}

	if err = db.Ping(); err != nil {
		log.Panic(err)
	}
}

// SelectAllUsers select all users from the DB and return them in a slice
func SelectAllUsers() ([]*types.User, error) {
	checkDbNotNilOrPanic()
	return handleSelectUserQuery("SELECT user_id, name, email, is_active, is_admin, is_su FROM user")
}

// SelectActiveUsers select only active users from the DB and return them in a slice
func SelectActiveUsers() ([]*types.User, error) {
	checkDbNotNilOrPanic()
	return handleSelectUserQuery("SELECT user_id, name, email, is_active, is_admin, is_su FROM user WHERE is_active = 1")
}

func handleSelectUserQuery(query string) ([]*types.User, error) {
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]*types.User, 0)
	for rows.Next() {
		u := new(types.User)
		err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Active, &u.Admin, &u.SU)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}

// SelectUserByID selects a single person from the ID based on their language level
func SelectUserByID(id int) (*types.User, error) {
	checkDbNotNilOrPanic()
	stmt, err := db.Prepare("SELECT user_id, name, email, is_active, is_admin, is_su FROM user WHERE user_id = ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	user := new(types.User)
	err = stmt.QueryRow(id).Scan(&user.ID, &user.Name, &user.Email, &user.Active, &user.Admin, &user.SU)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return user, nil
}

// SelectAllAssignments select all users from the DB and return them in a slice
func SelectAllAssignments(includeXML bool) ([]*types.Assignment, error) {
	checkDbNotNilOrPanic()
	query := "SELECT a.id, a.name, a.modified_on, u.user_id, u.name, u.email, u.is_active, u.is_admin, u.is_su"
	if includeXML {
		query += ", a.xml"
	}
	query += " FROM assignment a JOIN user u ON u.user_id = a.user_id"

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	assignments := make([]*types.Assignment, 0)
	for rows.Next() {
		a := new(types.Assignment)
		u := new(types.User)

		var err error
		if includeXML {
			err = rows.Scan(&a.ID, &a.Name, &a.ModifiedOn, &u.ID, &u.Name, &u.Email, &u.Active, &u.Admin, &u.SU, &a.XML)
		} else {
			err = rows.Scan(&a.ID, &a.Name, &a.ModifiedOn, &u.ID, &u.Name, &u.Email, &u.Active, &u.Admin, &u.SU)
		}

		if err != nil {
			return nil, err
		}
		a.Author = u
		assignments = append(assignments, a)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return assignments, nil
}

// SelectAssignmentByID finds an assignment in the DB by a given id, if nothing is found nil is returned
func SelectAssignmentByID(id int) (*types.Assignment, error) {
	checkDbNotNilOrPanic()
	stmt, err := db.Prepare(
		"SELECT a.id, a.name, a.modified_on, u.user_id, u.name, u.email, u.is_active, u.is_admin, u.is_su, a.xml " +
			"FROM assignment a JOIN user u ON u.user_id = a.user_id WHERE a.id = ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	a := new(types.Assignment)
	u := new(types.User)

	err = stmt.QueryRow(id).Scan(&a.ID, &a.Name, &a.ModifiedOn, &u.ID, &u.Name, &u.Email, &u.Active, &u.Admin, &u.SU, &a.XML)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	a.Author = u
	return a, nil
}

// InsertAssignment creates an assignment in the DB, it returns the assignment with an updated id
func InsertAssignment(a *types.Assignment) (*types.Assignment, error) {
	checkDbNotNilOrPanic()
	if a == nil || a.Name == "" || a.XML == "" || a.Author == nil || a.Author.ID <= 0 {
		return nil, ErrInvalidInput
	}

	stmt, err := db.Prepare("INSERT INTO assignment(user_id, name, xml) VALUES (?, ?, ?)")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(a.Author.ID, a.Name, a.XML)
	if err != nil {
		return nil, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	a.ID = int(id)
	return a, nil
}

// UpdateAssignment updates an assignment in the DB, if no assignment was updated an error is returned
func UpdateAssignment(a *types.Assignment) error {
	checkDbNotNilOrPanic()
	if a == nil || a.ID <= 0 || a.Name == "" || a.XML == "" || a.Author == nil || a.Author.ID <= 0 {
		return ErrInvalidInput
	}

	stmt, err := db.Prepare("UPDATE assignment SET user_id = ?, name = ?, xml = ? WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	res, err := stmt.Exec(a.Author.ID, a.Name, a.XML, a.ID)
	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected <= 0 {
		return ErrNoUpdate
	}

	return nil
}

// DeleteAssignment deletes an assignment by a given id
func DeleteAssignment(id int) error {
	checkDbNotNilOrPanic()

	stmt, err := db.Prepare("DELETE FROM assignment WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	return err
}

// SelectAssignmentComments returns all comments for an assignment or an empty slice
func SelectAssignmentComments(assignmentID int) ([]*types.Comment, error) {
	checkDbNotNilOrPanic()
	stmt, err := db.Prepare(
		"SELECT c.id, c.assignment_id, c.comment_val, c.modified_on, u.user_id, u.name, u.email, u.is_active, u.is_admin, u.is_su " +
			"FROM comment c JOIN user u ON u.user_id = c.user_id WHERE c.assignment_id = ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(assignmentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	comments := make([]*types.Comment, 0)
	for rows.Next() {
		c := new(types.Comment)
		u := new(types.User)
		err := rows.Scan(&c.ID, &c.AssignmentID, &c.CommentVal, &c.ModifiedOn, &u.ID, &u.Name, &u.Email, &u.Active, &u.Admin, &u.SU)
		if err != nil {
			return nil, err
		}
		c.Author = u
		comments = append(comments, c)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return comments, nil
}

func checkDbNotNilOrPanic() {
	if db == nil {
		log.Panic("DB is nil")
	}
}