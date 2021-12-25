package main

import (
	"github.com/gin-gonic/gin"
	"html/template"
	"log"
	"strings"
	"twometer.dev/mailservice/config"
	"twometer.dev/mailservice/mail"
)

type mailContext struct {
	Site string
	Body map[string]string
}

func buildMailContext(siteName string, siteConf config.SiteConfig, ginCtx *gin.Context) (ctx mailContext, ok bool) {
	formBodyValues := make(map[string]string)
	for _, key := range siteConf.Fields {
		value := ginCtx.PostForm(key)
		if value == "" {
			return mailContext{}, false
		}
		formBodyValues[key] = value
	}
	return mailContext{
		Site: siteName,
		Body: formBodyValues,
	}, true
}

func buildMailBody(ctx mailContext, tmpl *template.Template) (string, error) {
	builder := strings.Builder{}
	err := tmpl.Execute(&builder, ctx)
	return builder.String(), err
}

func handleRequest(conf *config.ServiceConfig, tmpl *template.Template, ginCtx *gin.Context) {
	siteName := ginCtx.Query("site")
	if siteName == "" {
		ginCtx.Status(400)
		return
	}

	site, found := conf.Sites[siteName]
	if !found {
		ginCtx.Status(404)
		return
	}

	mailCtx, ok := buildMailContext(siteName, site, ginCtx)
	if !ok {
		ginCtx.Status(400)
		return
	}

	mailBody, err := buildMailBody(mailCtx, tmpl)
	if err != nil {
		log.Println("Failed to build mail:", err)
		ginCtx.Status(500)
		return
	}

	m := mail.Mail{
		Subject:     "Twometer Mail Service",
		From:        conf.SenderAddress,
		To:          site.ReceiverAddress,
		ContentType: mail.HTML,
		Body:        mailBody,
	}

	err = m.Send(conf.ServerName, conf.ServerPort, conf.SenderPassword)
	if err != nil {
		log.Println("Failed to send mail:", err)
		ginCtx.Status(500)
	}
}
