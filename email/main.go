package main

import (
    "log"
    "net/smtp"
)

func main() {
    // Sender's email address
    from := "sureshcse2003@gmail.com"
    // Sender's email password
    password := "eaxe rtao mxfg plms"

    // Receiver's email address
    to := "ramakrishnanvaikundam@gmail.com"

    // SMTP server configuration
    smtpHost := "smtp.gmail.com"
    smtpPort := "587"

    // Message to send
    message := []byte("To: " + to + "\r\n" +
        "Subject: Testing Email from Golang\r\n" +
        "\r\n" +
        "This is the email body. This is just for learning. learning is fun\r\n")

    auth := smtp.PlainAuth("", from, password, smtpHost)

    err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, message)
    if err != nil {
        log.Fatal(err)
    }

    log.Println("Email successfully sent!")
}
