package hello

import (
  "fmt"
  "net/http"
)

func HelloHandler(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintf(w, "<h1>Hello from Go!</h1>")
}
