package mail

import (
	"fmt"
	"io"
	"log"
	"net/smtp"
	"strings"
)

func writeMailHeader(writer io.Writer, key string, value string) {
	_, err := fmt.Fprintf(writer, "%s: %s\r\n", key, value)
	if err != nil {
		log.Fatalln("Failed to write mail header", err)
	}
}

func (mail *Mail) serialize() []byte {
	builder := strings.Builder{}
	writeMailHeader(&builder, "Subject", mail.Subject)
	writeMailHeader(&builder, "From", mail.From)
	writeMailHeader(&builder, "To", mail.To)
	writeMailHeader(&builder, "Content-Type", mail.ContentType)
	builder.WriteString("\r\n")
	builder.WriteString(mail.Body)
	builder.WriteString("\r\n")

	return []byte(builder.String())
}

func (mail *Mail) Send(serverAddress string, serverPort uint16, password string) error {
	auth := smtp.PlainAuth("", mail.From, password, serverAddress)
	server := fmt.Sprintf("%s:%d", serverAddress, serverPort)
	return smtp.SendMail(server, auth, mail.From, []string{mail.To}, mail.serialize())
}
