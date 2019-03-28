package xidb

import (
	"reflect"
	"testing"

	"github.com/rs/xid"
)

func TestIDValue(t *testing.T) {
	i, _ := xid.FromString("9m4e2mr0ui3e8a215n4g")

	tests := []struct {
		name        string
		id          ID
		expectedVal interface{}
	}{
		{
			name:        "non nil id",
			id:          ID{ID: i},
			expectedVal: i.Bytes(),
		},
		{
			name:        "nil id",
			id:          ID{ID: xid.NilID()},
			expectedVal: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := tt.id.Value()
			if !reflect.DeepEqual(got, tt.expectedVal) {
				t.Errorf("wanted %v, got %v", tt.expectedVal, got)
			}
		})
	}
}

func TestIDScan(t *testing.T) {
	i, _ := xid.FromString("9m4e2mr0ui3e8a215n4g")

	tests := []struct {
		name        string
		val         interface{}
		expectedID  ID
		expectedErr bool
	}{
		{
			name:       "bytes id",
			val:        i.Bytes(),
			expectedID: ID{ID: i},
		},
		{
			name:       "nil id",
			val:        nil,
			expectedID: ID{ID: xid.NilID()},
		},
		{
			name:        "wrong bytes",
			val:         []byte{0x01},
			expectedErr: true,
		},
		{
			name:        "unknown type",
			val:         1,
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id := &ID{}
			err := id.Scan(tt.val)
			if (err != nil) != tt.expectedErr {
				t.Errorf("error expected: %t, got %t", tt.expectedErr, (err != nil))
			}
			if err == nil {
				if !reflect.DeepEqual(id.ID, tt.expectedID.ID) {
					t.Errorf("wanted %v, got %v", tt.expectedID, id)
				}
			}

		})
	}
}
