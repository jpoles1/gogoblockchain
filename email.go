package main

import (
	"log"
	"os"

	"github.com/SlyMarbo/gmail"
)

var gmailUname string
var gmailPass string

func initEmail() {
	gmailUname = os.Getenv("GMAIL_UNAME")
	if gmailUname == "" {
		log.Fatal("No Gmail username supplied in .env config file!")
	}
	gmailPass = os.Getenv("GMAIL_PASS")
	if gmailPass == "" {
		log.Fatal("No Gmail password supplied in .env config file!")
	}
}

func sendEmail(recipient string, subject string, body string) {
	email := gmail.Compose(subject, body)
	email.From = gmailUname
	email.Password = gmailPass

	// Defaults to "text/plain; charset=utf-8" if unset.
	email.ContentType = "text/html; charset=utf-8"

	// Normally you'll only need one of these, but I thought I'd show both.
	email.AddRecipient(recipient)
	err := email.Send()
	if err != nil {
		// handle error.
	}
}
