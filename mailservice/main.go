package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"html/template"
	"log"
	"net/smtp"
	"os"
	"strconv"
	"strings"
)

type templateCtx struct {
	Site   string
	Fields map[string]string
}

func buildMail(to string, subject string, contentType string, body string) []byte {
	mail := fmt.Sprintf("To: %s\r\n"+
		"Subject: %s\r\n"+
		"Content-Type: %s\r\n"+
		"\r\n"+
		"%s\r\n", to, subject, contentType, body)
	return []byte(mail)
}

func sendMail(config *ServiceConfig, site *SiteConfig, values map[string]string) bool {
	wr := strings.Builder{}
	data, _ := os.ReadFile("config/template.html")
	tpl, err := template.New("mail").Parse(string(data))
	if err != nil {
		log.Fatalln(err)
	}
	err = tpl.Execute(&wr, templateCtx{Site: site.CorsOrigin, Fields: values})
	if err != nil {
		log.Fatalln(err)
	}

	auth := smtp.PlainAuth("", config.SenderAddress, config.SenderPassword, config.ServerName)
	addr := config.ServerName + ":" + strconv.Itoa(int(config.ServerPort))
	msg := buildMail(site.ReceiverAddress, "Twometer Mail Service", "text/html", wr.String())
	err = smtp.SendMail(addr, auth, config.SenderAddress, []string{site.ReceiverAddress}, msg)
	if err != nil {
		fmt.Println("Failed to send message:", err)
		return false
	}
	return true
}

func handleRequest(config *ServiceConfig, site *SiteConfig, context *gin.Context) {
	fieldValues := make(map[string]string)
	for _, key := range site.Fields {
		value := context.PostForm(key)
		if value == "" {
			context.Status(400)
			return
		}
		fieldValues[key] = value
	}

	log.Println("Sending", fieldValues, "to", site.ReceiverAddress)
	if !sendMail(config, site, fieldValues) {
		context.Status(500)
	}
}

func main() {
	log.SetOutput(os.Stdout)

	config, err := LoadConfig("config/config.secret.json")
	if err != nil {
		log.Fatalf("Failed to load configuration: %s", err)
		return
	}
	log.Println("Loaded config:", config)

	r := gin.Default()
	r.POST("/", func(context *gin.Context) {
		siteName := context.Query("site")
		if siteName == "" {
			context.Status(400)
			return
		}

		site, found := config.Sites[siteName]
		if !found {
			context.Status(404)
			return
		}

		handleRequest(&config, &site, context)
	})

	log.Println("Starting REST server...")
	err = r.Run()
	if err != nil {
		log.Fatalf("Failed to start REST server: %s", err)
		return
	}
}
