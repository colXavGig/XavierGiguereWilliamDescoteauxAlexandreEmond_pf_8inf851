package cmd

import (
	"flag"
	"fmt"
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

	flag.StringVar(&info.user, "u", "test", "user for db conn string")
	flag.StringVar(&info.password, "password", "test", "password for db conn string")
	flag.StringVar(&info.host, "db", "", "host URI")
	flag.BoolVar(&info.sysdba, "sysdba", false, "bool for use sysdba")

	*connString = info.getConnString()
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
