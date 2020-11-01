package src

import (
	"errors"
	"fmt"

	"log"
	"net/smtp"
	"time"
)

var errorConfig = errors.New("invalid Config")

type Config struct {
	Mail     string
	Password string
}

func SendMessages(d * DataBaseWork, url *string, newPrice *string) error {
	from := d.Config.Mail
	password := d.Config.Password

	fmt.Println("Mail = ", d.Config.Mail)
	fmt.Println("Password = ", d.Config.Password)

	rows, err := d.DBConnection.Query("SELECT DISTINCT Mail FROM pricetable WHERE link=$1", *url)
	if err != nil {
		log.Println(err)
		return err
	}
	if rows == nil {
		log.Println("no such link in database")
		return err
	}
	defer func() {
		if rows != nil {
			_ = rows.Close()
		}
	}()
	to := make([]string, 0)
	for rows.Next() {
		var curMail string
		err := rows.Scan(&curMail)
		if err != nil {
			log.Println(err)
		}
		to = append(to, curMail)
	}
	fmt.Println(to)

	// smtp server configuration.
	smtpHost := "smtp.yandex.ru"
	smtpPort := "587"

	// Message.
	message := []byte("To:" + "" + "\r\n" +
		"Subject: Изменение \r\n" +
		"\r\n" +
		"Добрый день!" + "\r\n" +
		"У товара по ссылке ниже изменилась цена:\r\n" +
		*url + "\r\n" +
		"Новая цена товара: " + *newPrice + "\r\n" +
		"\r\n" +
		"С уважением," + "\r\n" +
		"Сервис Авито\r\n")

	// Authentication.
	auth := smtp.PlainAuth("", from, password, smtpHost)
	// Sending email.
	err = smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println("Email Sent Successfully!")
	return nil
}

func UpdateNewPrice(d * DataBaseWork, url *string, newPrice *string) {
	d.DBConnection.Exec("UPDATE pricetable SET price = $1 WHERE link = $2", newPrice, url)
}

func UpdatingLinkProcess(d * DataBaseWork, url *string, lastPrice string) {
	//row := d.DBConnection.QueryRow("SELECT DISTINCT price FROM $1 WHERE link=$2", PriceTable, url)
	//var lastPrice string
	//row.Scan(&lastPrice)
	for {
		curPrice, err := getPrice(url)
		if err != nil {
			log.Println(err)
		} else {
			if lastPrice != curPrice {
				err = SendMessages(d, url, &curPrice)
				if err == errorConfig {
					return
				}
				UpdateNewPrice(d, url, &curPrice)
				lastPrice = curPrice
			}
		}
		time.Sleep(5 * time.Second)
	}
}
