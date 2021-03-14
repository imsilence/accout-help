package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/imroc/req"
	"github.com/imsilence/account-help/handlers"
)

const token = "xxx"

func main() {
	var (
		addr  string
		help  bool
		debug bool
	)

	os.Setenv("account.help.github.user", "imsilence")
	os.Setenv("account.help.github.password", token)

	flag.StringVar(&addr, "addr", ":8091", "listen addr")
	flag.BoolVar(&help, "help", false, "help")
	flag.BoolVar(&debug, "debug", false, "debug")
	flag.Usage = func() {
		fmt.Println("usage: account-help --addr :8091")
		flag.PrintDefaults()
	}

	flag.Parse()

	if help {
		flag.Usage()
		os.Exit(0)
	}
	req.Debug = debug

	http.HandleFunc("/", handlers.Index)
	http.HandleFunc("/member/", handlers.Member)
	http.HandleFunc("/invitations/", handlers.Invitations)
	http.HandleFunc("/invitations/cancel/", handlers.CancelInvitation)
	http.HandleFunc("/members/", handlers.Members)
	http.HandleFunc("/repos/", handlers.Repos)
	log.Fatal(http.ListenAndServe(addr, nil))
}
