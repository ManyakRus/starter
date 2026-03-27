Mail 【golang send mail package】
------

##  Mail: send mail support AUTH:

* `LOGIN` : mail.LoginAuth(email.Username, email.Password)
* `CRAM-MD5`: smtp.CRAMMD5Auth(email.Username, email.Password)
* `PLAIN` : smtp.PlainAuth(email.Identity, email.Username, email.Password, email.Host)
* `NTLM`: mail.NTLMAuth(email.Host, email.Username, email.Password, mail.NTLMVersion1) # mail.NTLMVersion2 也支持

## Mail: send mail support Secure:

* SSL
* TLS
* AUTH TYPE ( LOGIN, CRAM-MD5, PLAIN, NTLM )

# Aim of this project 

It was to be able to send email to a NTLM server ( that is not available as default by go lang library )
I have based my library on the https://github.com/cmarkh/smtp that did a great job to package it.
But there were too dependencies and was not up to date ( recent go std library ) ( especially for smtp.go file )

This library has been packaged and is self contained now
Needs no external lib

Can be use directecly in any project by calling `go get gitlab.com/c4sp/go/smtp_ntlm`



## Mail: `go get gitlab.com/c4sp/go/smtp_ntlm`

```
email := NewEMail(`{"port":25}`)
email.From = `farmerx@163.com`
email.Host = `smtp.163.com`
email.Port = int(25) // [587 NTLM AUTH] [465，994]
email.Username = `Farmerx`
email.Secure = `` // SSL，TSL
email.Password = `************`
authType := `LOGIN`
switch authType {
case ``:
	email.Auth = nil
case `LOGIN`:
	email.Auth = LoginAuth(email.Username, email.Password)
case `CRAM-MD5`:
	email.Auth = smtp.CRAMMD5Auth(email.Username, email.Password)
case `PLAIN`:
	email.Auth = smtp.PlainAuth(email.Identity, email.Username, email.Password, email.Host)
case `NTLM`:
	email.Auth = NTLMAuth(email.Host, email.Username, email.Password, NTLMVersion1)
default:
	email.Auth = smtp.PlainAuth(email.Identity, email.Username, email.Password, email.Host)
}

email.To = []string{`farmerx@163.com`}
email.Subject = `send mail success`
email.Text = "Test Email：\r\n   Following is a content send by the lib."
//email.AttachFile(reportFile)
if err := email.Send(); err != nil {
      fmt.Println(err)
}

```




   
