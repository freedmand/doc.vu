package docvu

import (
    "fmt"
    "net/http"
    "flag"
    "github.com/golang/glog"
)

var (
  port = flag.Int("port", 8084, "The default server port")
  salt = flag.String("salt", "docvu salt is quite unique", "The salt for hashids")
  dbUrl = flag.String("db_url", "localhost:28015", "The url for the rethink db")
  dbName = flag.String("db_name", "docvu", "The name of the rethink database")
  dbTable = flag.String("db_table", "docs", "The name of the rethink table for docs")
  seed = flag.Int64("seed", 413, "The random seed value")
  hashidsLength = flag.Int("hashids_length", 12, "The minimum length of hashids")
  t *TableState
)

func init() {
    flag.Parse()
    var err error
    t, err = initTableState(*dbUrl, *salt, *hashidsLength, *seed, *dbName, *dbTable)
    if err != nil {
      glog.Fatal(err)
    }
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
  hashid, err := t.createId()
  if err != nil {
    glog.Fatal(err)
  }
  fmt.Fprintf(w, "<body>Your page is <a href=\"#\">doc.vu/%s</a><br>Length: %d</body>", hashid, len(hashid))
  // createDoc(hashid, false, "")
}

// func saveHandler(w http.ResponseWriter, r *http.Request) {
//     decoder := json.NewDecoder(r.Body)
//     var s saveRequest   
//     err := decoder.Decode(&s)
//     if err != nil {
//         glog.Errorf("Error handling save request: %v", err)
//     }
//     err := saveDoc(s.docId, s.previousVersion, s.contents)
//     if err != nil {
//       http.Error(w, "An error occurred during saving", 500)
//     } else {
//       fmt.Fprint(w, "Successful save")
//     }
// }

func main() {
    http.HandleFunc("/", mainHandler)
    // http.HandleFunc("/save", saveHandler)
    glog.Infof("Listening on port %d", *port)
    http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
}