package cmd

import (
	"flag"
	"fmt"
	"log"
)

// TODO: code cli
type dbConnInfo struct {
	user string
	password string
	host string 
	sysdba bool
}

type hostInfo struct {
	host string
	port int
}

func (info *dbConnInfo) getConnString() string {
	return fmt.Sprintf("user=\"%s\" password=\"%s\" connectString=\"%s\" sysdba=%t", info.user, info.password, info.host, info.sysdba)
}

func SetConnString(connString *string) {
	info := dbConnInfo{}

	flag.StringVar(&info.user, "u", "db2_fall2024", "user for db conn string")
	flag.StringVar(&info.password, "password", "Baddemon665!", "password for db conn string")
	flag.StringVar(&info.host, "db", "aedev.pro:1521/XEPDB1", "host URI")
	flag.BoolVar(&info.sysdba, "sysdba", false, "bool for use sysdba")

	*connString = info.getConnString()
	log.Printf("Connectring: %s",info.getConnString())
}

func SetHostURI(host *string) {
	info := hostInfo{}

	flag.StringVar(&info.host, "h", "", "website host url")
	flag.IntVar(&info.port, "port", 8000, "website port")

	*host = info.GetHostURI()
}

func (info *hostInfo) GetHostURI() string {
	return fmt.Sprintf("%s:%d", info.host, info.port)
}
