package sending

import (
	"fmt"
	"github.com/KatsayArtemDev/verification/src/sending/parser"
	"gopkg.in/gomail.v2"
	"os"
)

func PinToUser(email, pin string) error {
	password := os.Getenv("PASSWORD")

	parsedTemplate, err := parser.HtmlParser("./sending/template/verification.html", pin)
	if err != nil {
		return err
	}

	mail := gomail.NewMessage()
	mail.SetHeader("From", "katsayartemdev@gmail.com")
	mail.SetHeader("To", email)
	mail.SetAddressHeader("Cc", "katsayartemdev@gmail.com", "Artem")
	mail.SetHeader("Subject", "KatsayInc.")

	mail.SetBody("text/html", parsedTemplate)

	d := gomail.NewDialer("smtp.gmail.com", 587, "katsayartemdev@gmail.com", password)

	err = d.DialAndSend(mail)

	if err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}

	return nil
}
