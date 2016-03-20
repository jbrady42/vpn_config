package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/jbrady42/vpn_config/vpn_conf"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/lib/pq"
)

var mut sync.Mutex
var VPN_PORT string
var VPN_HOST string
var baseIP = "2001:db8:1:2:0:0:0:2"

type Conf struct {
	IP     string
	InitIP string
	Port   string
}

type IPCount struct {
	gorm.Model
	IP string
}

func handleReq(w http.ResponseWriter, r *http.Request) {
	addr := getNextAddr()
	c := &Conf{addr, VPN_HOST, VPN_PORT}
	renderTemplate(w, "peer", c)
	log.Println("Sent config for ip", addr)
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Conf) {
	t, _ := template.ParseFiles(tmpl + ".conf")
	t.Execute(w, p)
}

func getNextAddr() string {
	mut.Lock()
	db := connectDB()

	ipM := currentAddress(db)
	addr := ipM.IP

	ret := vpn_conf.NextAddress6(addr)

	//Update latest address
	db.Model(&ipM).Update("IP", ret)

	db.Close()
	mut.Unlock()

	return ret
}

func currentAddress(db *gorm.DB) IPCount {
	var tmp IPCount
	if err := db.First(&tmp).Error; err != nil {
		log.Println(err)
	}
	return tmp
}

func connectDB() *gorm.DB {
	dbStr := os.Getenv("DATABASE_URL")
	db, err := gorm.Open("postgres", dbStr)
	if err != nil {
		log.Fatal("Can't connect to db")
	}
	return db
}

func withDB(fun func(*gorm.DB)) {
	db := connectDB()
	fun(db)
	db.Close()
}

// Migrate and seed
func setupDB(db *gorm.DB) {

	db.AutoMigrate(&IPCount{})

	var count int
	db.Model(&IPCount{}).Count(&count)
	if count == 0 {
		log.Println("Inserting base record")
		db.Create(&IPCount{IP: baseIP})
	}
}

func main() {
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

	// Db setup
	withDB(func(d *gorm.DB) {
		setupDB(d)
	})

	addr := fmt.Sprintf(":%v", port)
	fmt.Println("Starting server at ", addr)

	http.HandleFunc("/conf/", handleReq)
	http.ListenAndServe(addr, nil)
}
