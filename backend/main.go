package main

import (
	"log"

	"github.com/colXavGig/XavierGiguereWilliamDescoteauxAlexandreEmond_pf_8inf851/BLL"
	"github.com/colXavGig/XavierGiguereWilliamDescoteauxAlexandreEmond_pf_8inf851/cmd"
)

func main() {
	var (
		connString string
		websiteAddr string
	)

	cmd.SetHostURI(&websiteAddr)
	cmd.SetConnString(&connString)


	srv := BLL.NewServer(websiteAddr, connString)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Error while serving. Error: %s", err.Error())
	}
}
