package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/imsilence/account-help/handlers"
)

const token = "xxxxxxxxxxxxx"

func main() {
	// req.Debug = true
	var (
		addr string
		help bool
	)

	os.Setenv("account.help.github.user", "imsilence")
	os.Setenv("account.help.github.password", token)

	flag.StringVar(&addr, "addr", ":8091", "listen addr")
	flag.BoolVar(&help, "help", false, "help")
	flag.Usage = func() {
		fmt.Println("usage: account-help --addr :8091")
		flag.PrintDefaults()
	}

	flag.Parse()

	if help {
		flag.Usage()
		os.Exit(0)
	}

	http.HandleFunc("/", handlers.Index)
	http.HandleFunc("/member/", handlers.Member)
	http.HandleFunc("/invitations/", handlers.Invitations)
	http.HandleFunc("/members/", handlers.Members)
	http.HandleFunc("/repos/", handlers.Repos)
	log.Fatal(http.ListenAndServe(addr, nil))
}
