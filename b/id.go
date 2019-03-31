package xidb

import (
	"database/sql/driver"
	"fmt"

	"github.com/rs/xid"
)

type ID struct {
	xid.ID
}

// Value implements the driver.Valuer interface.
func (id ID) Value() (driver.Value, error) {
	if id.ID.IsNil() {
		return nil, nil
	}
	return id.ID[:], nil
}

// Scan implements the sql.Scanner interface.
func (id *ID) Scan(value interface{}) (err error) {
	switch val := value.(type) {
	case []byte:
		i, err := xid.FromBytes(val)
		if err != nil {
			return err
		}
		*id = ID{ID: i}
		return nil
	case nil:
		*id = ID{ID: xid.NilID()}
		return nil
	default:
		return fmt.Errorf("xid: scanning unsupported type: %T", value)
	}
}
