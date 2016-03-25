package entry

import (
  "fmt"

  "github.com/marcboudreau/godbpf/util"
)

// DBPFEntryTGI is combines three identifiers: the Type, Group, and Instance.
// Together, these identifiers uniquely identify a resource as well as indicating
// the type of the resource.
type DBPFEntryTGI struct {
  // The TypeId indicates the broad type of the resource.
  TypeId uint32

  // The GroupId indicates the sub-type (or category) of the resource, within
  // the declared type.
  GroupId uint32

  // The InstanceId uniquely identifies the resource among all others with the
  // same TypeId and GroupId.
  InstanceId uint32
}

// String returns a string representation of the receiver DBPFEntryTGI.
func (tgi *DBPFEntryTGI) String() string {
  return fmt.Sprintf("T: 0x%08X, G: 0x%08X, I: 0x%08X", tgi.TypeId, tgi.GroupId, tgi.InstanceId)
}

// Equals tests the provided DBPFEntryTGI instance against the receiver to see
// if they match.
func (tgi *DBPFEntryTGI) Equals(other *DBPFEntryTGI) bool {
  if other != nil {
    result := tgi.TypeId == other.TypeId && tgi.GroupId == other.GroupId && tgi.InstanceId == other.InstanceId
    return result
  } else {
    return false
  }
}

// Bytes encodes the receiver as 4 consecutive unsigned 32-bit integers into the
// provided byte slice.
func (tgi *DBPFEntryTGI) Bytes(bytes []byte) {
  copy(bytes[0:4], util.WriteUint32(tgi.TypeId))
  copy(bytes[4:8], util.WriteUint32(tgi.GroupId))
  copy(bytes[8:12], util.WriteUint32(tgi.InstanceId))
}
