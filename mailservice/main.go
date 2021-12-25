package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"html/template"
	"log"
	"os"
	"time"
	"twometer.dev/mailservice/config"
)

func loadConfig() config.ServiceConfig {
	path := os.Getenv("CONFIG_FILE")
	if path == "" {
		path = "config/config.json"
	}

	conf, err := config.Read(path)
	if err != nil {
		log.Fatalln("Failed to read config file:", err)
	}
	return conf
}

func loadTemplate() *template.Template {
	path := os.Getenv("TEMPLATE_FILE")
	if path == "" {
		path = "config/template.html"
	}

	content, err := os.ReadFile(path)
	if err != nil {
		log.Fatalln("Failed to read template:", err)
	}

	tmpl, err := template.New("mail").Parse(string(content))
	if err != nil {
		log.Fatalln("Failed to parse template:", err)
	}

	return tmpl
}

func init() {
	log.SetOutput(os.Stdout)
}

func main() {
	log.Println("Loading resources...")
	conf := loadConfig()
	tmpl := loadTemplate()

	log.Println("Initializing web server...")
	r := gin.Default()

	var allowedOrigins []string
	for _, v := range conf.Sites {
		allowedOrigins = append(allowedOrigins, v.CorsOrigin)
	}
	r.Use(cors.New(cors.Config{
		AllowOrigins: allowedOrigins,
		AllowMethods: []string{"POST"},
		AllowHeaders: []string{"Origin", "Content-Length", "Content-Type"},
		MaxAge:       12 * time.Hour,
	}))

	r.POST("/", func(context *gin.Context) {
		handleRequest(&conf, tmpl, context)
	})

	log.Println("Launching web server")
	err := r.Run()
	if err != nil {
		log.Fatalln("Failed to launch web server:", err)
		return
	}
}
