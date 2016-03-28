package entry

import (
  "fmt"

  "github.com/marcboudreau/godbpf/util"
)

var DIR_ENTRY_TGI *DBPFEntryTGI = &DBPFEntryTGI{TypeId: 0xE86B1EEF, GroupId: 0xE86B1EEF, InstanceId: 0x286B1F03}

// AddEntry modifies the receiver DBPFEntry to include the provided DBPFEntryTGI
// and size.  The method will make sure that the receiver DBPFEntry instance has
// the correct DBPFEntryTGI value, otherwise it will panic.
func (e *DBPFEntry) AddEntry(tgi *DBPFEntryTGI, size uint32) {
  if !e.TGI.Equals(DIR_ENTRY_TGI) {
    panic(fmt.Sprintf("dbpfdirentry.AddEntry() can only be called with a DBPFEntry that has this TGI: {%s}\n", DIR_ENTRY_TGI))
  }

  entryData := make([]byte, 16)
  tgi.Bytes(entryData[0:12])
  copy(entryData[12:], util.WriteUint32(size))

  data := e.GetData()
  len := len(data)
  temp := make([]byte, len + 16)
  copy(temp[:len], data)
  copy(temp[len:], entryData)

  e.SetData(temp)
}

// CreateDirEntry creates a DBPFEntry with the DBPFEntryTGI reserved for the DIR
// entry.
func CreateDirEntry() *DBPFEntry {
  return &DBPFEntry{TGI: DIR_ENTRY_TGI}
}
