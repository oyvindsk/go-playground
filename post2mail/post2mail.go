package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"html/template"
	"log"
	"net/http"
	"net/smtp"
	"strconv"
	txttemplate "text/template"
)

// Status returned to the HTTP client
type StatusMSG struct {
	Status string
	Data   Email
}

// Information regarding the SMTP server
type SMTPUser struct {
	Name   string
	Pass   string
	Server string
	Port   int
}

// Data received in the form and goes into the email
type Email struct {
	Name    string
	From    string
	Subject string
	Text    string
	To      string
}

// html template (for testing) and Email template
// guess they are here mostly to avoid compiling them for each request
var indexTemplate = template.Must(template.New("index").Parse(templateStr))
var emailTemplate = txttemplate.Must(txttemplate.New("emailTemplate").Parse(templateEmail))

// Some structs to hold parameters given on the command-line
var smtpUser SMTPUser
var email Email

func init() {

    // Handle command line parameters
	flag.StringVar(&smtpUser.Server, "smtp-server", "smtp.gmail.com", "SMTP server to send mail through")
	flag.IntVar(&smtpUser.Port, "smtp-port", 587, "Port to connect to on the smtp server, usually 25 or 587")
	flag.StringVar(&smtpUser.Name, "smtp-user", "", "Username for the SMTP server")
	flag.StringVar(&smtpUser.Pass, "smtp-pass", "", "Password for the SMTP server")

	flag.StringVar(&email.To, "email-to", "", "Email address to send the emails to")

	flag.Parse()
}

func main() {
	http.Handle("/", http.HandlerFunc(handleGet))       // Get, for testing
	http.Handle("/post", http.HandlerFunc(handlePost))  // Receive the POST from the form
	err := http.ListenAndServe(":1234", nil)           // Remeber to use the same port in the <form> :)
	checkErr("ListenAndServe", err)
}

func handleGet(w http.ResponseWriter, req *http.Request) {
    // Just return a simple html form - for testing
	indexTemplate.Execute(w, req.FormValue("s"))
}

func handlePost(w http.ResponseWriter, req *http.Request) {
	log.Println("handlePost")
	log.Println(req.Form)

	// so this is kind of stupid.. ?
	email.Name = req.FormValue("name")
	email.From = req.FormValue("from")
	email.Subject = "FASTPRISAPPS Henvendelse: " + req.FormValue("subject")
	email.Text = req.FormValue("text")

	// FIXME: Error checking ..

	// Send the info on email
	go formatAndSendEmail(email, smtpUser) // don't rely on variables beeing global

	// Return a JSON status hash
	m := StatusMSG{"OK", email}
	b, err := json.Marshal(m)
	checkErr("Json Marshal", err)
	log.Printf("Email sent:\n\t%s", b)

	w.Write(b)

}

func checkErr(msg string, err error) {
	if err != nil {
		log.Fatal(msg, err)
	}
}

func formatAndSendEmail(data Email, smtpInfo SMTPUser) {

	var buffer bytes.Buffer
	err := emailTemplate.Execute(&buffer, &data)

	auth := smtp.PlainAuth(
		"",
		smtpInfo.Name,
		smtpInfo.Pass,
		smtpInfo.Server,
	)

	err = smtp.SendMail(
		smtpInfo.Server+":"+strconv.Itoa(smtpInfo.Port),
		auth,
		smtpInfo.Name,
		[]string{"os@odots.org"}, //[]string{data.To},
		buffer.Bytes(),
	)

	checkErr("SendMail: ", err)
}

//
// Templates for email and html pages

const templateEmail = `From: {{.From}}
To: {{.To}}
Subject: {{.Subject}}
MIME-version: 1.0
Content-Type: text/html; charset="UTF-8"

	Name    {{.Name}}
    <br>
	From    {{.From}}
    <br>
	Subject {{.Subject}}
    <br>
	Text    {{.Text}}
    <br>
`

const templateStr = `
<html>
<head>
<title>test</title>
</head>
<body>
{{if .}}
<h1>?</h1>
{{.}}
<br>
<br>
{{end}}
<form action="/post" name=f method="POST">
    <input maxLength=1024 size=70   name=name       value="" placeholder="Name"    title="Name">
    <input maxLength=1024 size=70   name=from       value="" placeholder="From"    title="From Email">
    <input maxLength=1024 size=70   name=subject    value="" placeholder="Subject" title="Subject">
    <input maxLength=1024 size=500  name=text       value="" placeholder="Message" title="Message">
    <input type=submit value="Send" name=send>
</form>

</body>
</html>
`
