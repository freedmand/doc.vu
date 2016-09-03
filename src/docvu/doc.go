package docvu

import (
    "github.com/golang/glog"
    "time"
    "math/rand"
    "github.com/speps/go-hashids"
    r "github.com/dancannon/gorethink"
)

type TableState struct {
  dbName string
  dbTable string
  session *r.Session
  hashid *hashids.HashID
  random *rand.Rand
}

func initTableState(dbUrl string, salt string, hashidsLength int, seed int64, dbName string, dbTable string) (*TableState, error) {
    session, err := r.Connect(r.ConnectOpts{
        Address: dbUrl,
    })

    if err != nil {
      glog.Errorf("Error loading Rethink table: %v", err)
      return nil, err
    }

    hd := hashids.NewData()
    hd.Salt = salt
    hd.MinLength = hashidsLength
    hashid := hashids.NewWithData(hd)

    random := rand.New(rand.NewSource(seed))

    return &TableState{
      dbName: dbName,
      dbTable: dbTable,
      session: session,
      hashid: hashid,
      random: random,
    }, nil
}

func (t TableState) createId() (string, error) {
  milliseconds := time.Now().UnixNano() / int64(time.Millisecond)
  randomness := int64(t.random.Int31() & 0xFFFF)
  id, err := t.hashid.EncodeInt64([]int64{randomness, milliseconds, 0})
  if err != nil {
    return "", err
  }
  return id, nil
}

func (t TableState) createDoc(readOnly bool, contents string) (id string, err error) {
  id, err = t.createId()
  if err != nil {
    return "", err
  }
  r.DB(t.dbName).Table(t.dbTable).Insert(map[string]interface{}{
      "id": id,
      "baseDoc": id,
      "readOnly": readOnly,
      "contents": contents,
      "created": time.Now(),
      "leaf": true,
  }).Run(t.session)
  return id, nil
  // Assign a base version
  // Have a created timestamp
}

func WordCount(doc string) int {
  return 5
}

func (t TableState) saveDoc(docId string, previousVersion string, contents string) error {
  id, err := t.createId()
  if err != nil {
    return err
  }
  r.DB(t.dbName).Table(t.dbTable).Insert(map[string]interface{}{
      "id": id,
      "baseDoc": docId,
      "readOnly": false,
      "contents": contents,
      "created": time.Now(),
      "previousVersion": previousVersion,
  }).Run(t.session)

  return nil
  // Change previous version to say it's not a leaf
  // Set new version with previous version pointing
}