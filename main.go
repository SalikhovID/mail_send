package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/smtp"
	"os"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	email_smpt := os.Getenv("EMAIL")
	pass := os.Getenv("PASSWORD")

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
	)
	// Create a database connection
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	for {
		// Execute the query
		rows, err := db.Query("SELECT id, company_id, email, mess_title, mess_body, from_col FROM mail_send_message WHERE status = 0 LIMIT 10")
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		// Iterate over the rows and print the results
		var ids []int
		for rows.Next() {
			var id int
			var company_id int
			var email string
			var mess_title string
			var mess_body string
			var from_col string
			err := rows.Scan(&id, &company_id, &email, &mess_title, &mess_body, &from_col)
			if err != nil {
				log.Fatal(err)
			}
			ids = append(ids, id)

			sendMail(email_smpt, pass, mess_body, email)

			_, err = db.Exec("UPDATE mail_send_message SET status = 1, send_time = NOW() WHERE id = $1", id)
			if err != nil {
				log.Fatal(err)
			}

		}

		// Check for errors during iteration
		if err = rows.Err(); err != nil {
			log.Fatal(err)
		}

		fmt.Println(ids)
		time.Sleep(1 * time.Second)
	}

}

func sendMail(email string, password string, message string, client string) {
	// Set up authentication information.
	auth := smtp.PlainAuth(password, email, password, "smtp.gmail.com")

	// Set up the message headers.
	to := []string{client}
	msg := []byte("To: " + to[0] + "\r\n" +
		"Subject: Test email\r\n" +
		"\r\n" +
		message)

	// Connect to the SMTP server.
	err := smtp.SendMail("smtp.gmail.com:587", auth, email, to, msg)
	if err != nil {
		fmt.Println("Error sending email:", err)
	} else {
		fmt.Println("Email sent successfully.")
	}
}
