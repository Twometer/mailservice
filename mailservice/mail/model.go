package mail

const (
	PlainText = "text/plain; charset=UTF-8"
	HTML      = "text/html; charset=UTF-8"
)

type Mail struct {
	Subject     string
	From        string
	To          string
	ContentType string
	Body        string
}
