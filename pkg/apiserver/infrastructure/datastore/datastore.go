package datastore

import "fmt"

var (
	// ErrPrimaryEmpty Error that primary key is empty.
	ErrPrimaryEmpty = NewDBError(fmt.Errorf("entity primary is empty"))

	// ErrTableNameEmpty Error that table name is empty.
	ErrTableNameEmpty = NewDBError(fmt.Errorf("entity table name is empty"))

	// ErrNilEntity Error that entity is nil
	ErrNilEntity = NewDBError(fmt.Errorf("entity is nil"))

	// ErrRecordExist Error that entity primary key is exist
	ErrRecordExist = NewDBError(fmt.Errorf("data record is exist"))

	// ErrRecordNotExist Error that entity primary key is not exist
	ErrRecordNotExist = NewDBError(fmt.Errorf("data record is not exist"))

	// ErrIndexInvalid Error that entity index is invalid
	ErrIndexInvalid = NewDBError(fmt.Errorf("entity index is invalid"))

	// ErrEntityInvalid Error that entity is invalid
	ErrEntityInvalid = NewDBError(fmt.Errorf("entity is invalid"))
)

type DBError struct {
	err error
}

func (d *DBError) Error() string {
	return d.err.Error()
}

func NewDBError(err error) error {
	return &DBError{err: err}
}