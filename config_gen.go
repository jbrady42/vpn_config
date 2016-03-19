package main

import (
  "fmt"
  "sync"
  "net/http"
  "strings"
  "strconv"
  "html/template"
)

var mut sync.Mutex
var addr []int
var base = 256
var VPN_PORT = "7000"
var VPN_IP = "10.8.8.125"

type Conf struct {
  IP      string
  InitIP  string
  Port    string
}

func handleReq(w http.ResponseWriter, r *http.Request) {
  a := nextAddr()
  c := &Conf{ipString(a), VPN_IP, VPN_PORT}
  renderTemplate(w, "peer", c)
  // fmt.Fprintf(w, ipString(a))
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Conf) {
    t, _ := template.ParseFiles(tmpl + ".conf")
    t.Execute(w, p)
}

func ipString(ip []int) string {
  strs := make([]string, len(addr))
  for i, a := range addr {
    strs[i] = strconv.Itoa(a)
  }
  return strings.Join(strs, ".")
}

func nextAddr() []int {
  mut.Lock()
  for x := len(addr)-1; x > 0; x-- {
    if addr[x] < base - 2 {
      addr[x] += 1
      break
    }
  }
  a := addr
  mut.Unlock()
  return a
}

func main() {
  addr = []int{10, 8, 0, 1}
  port := 9000
  http.HandleFunc("/conf/", handleReq)
  addr := fmt.Sprintf(":%v", port)
  fmt.Println("Starting server at ", addr)
  http.ListenAndServe(addr, nil)
}
