package main

import (
  "os"
  "fmt"
  "sync"
  "net/http"
  "strings"
  "strconv"
  "html/template"
  "log"
)

var mut sync.Mutex
var addr []int
var base = 256
var VPN_PORT string
var VPN_HOST string

type Conf struct {
  IP      string
  InitIP  string
  Port    string
}

func handleReq(w http.ResponseWriter, r *http.Request) {
  a := nextAddr()
  ipStr := ipString(a)
  c := &Conf{ipStr, VPN_HOST, VPN_PORT}
  renderTemplate(w, "peer", c)
  log.Println("Sent config for ip", ipStr)
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

  port := "9000"
  if len(os.Getenv("PORT")) > 0 {
    port = os.Getenv("PORT")
  }

  VPN_PORT = "7000"
  if len(os.Getenv("VPN_PORT")) > 0 {
    VPN_PORT = os.Getenv("VPN_PORT")
  }

  VPN_HOST = "10.8.8.125"
  if len(os.Getenv("VPN_HOST")) > 0 {
    VPN_HOST = os.Getenv("VPN_HOST")
  }

  addr := fmt.Sprintf(":%v", port)
  fmt.Println("Starting server at ", addr)

  http.HandleFunc("/conf/", handleReq)
  http.ListenAndServe(addr, nil)
}
