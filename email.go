package main

import (
	"context"
	"fmt"
	"log"

	"gopkg.in/mail.v2"

	"github.com/robfig/cron/v3"

	"github.com/jonathanfong098/csci169project/internal/database"
)

type smtpConfig struct {
	server   string
	port     int
	user     string
	password string
}

type smtpServer struct {
	dialer *mail.Dialer
	config *smtpConfig
}

func (es *smtpServer) sendEmail(user database.User, subject, bodyText string) error {
	m := mail.NewMessage()

	m.SetHeader("From", "mailtrap@demomailtrap.com")
	m.SetHeader("To", user.Email)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", bodyText)

	err := es.dialer.DialAndSend(m)
	if err != nil {
		fmt.Println("An error occurred:", err)
		return fmt.Errorf("failed to send email: %w", err)
	}

	fmt.Printf("Email sent successfully to %s\n", user.Email)
	return nil
}

func (es *smtpServer) sendDailyPosts(db *database.Queries) {
	users, err := db.GetSubscribedUsers(context.Background())
	if err != nil {
		log.Printf("Failed to fetch subscribed users: %v", err)
		return
	}

	for _, user := range users {
		posts, err := db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
			UserID: user.ID,
			Limit:  10,
		})
		if err != nil {
			log.Printf("Couldn't get posts for user %v: %v", user.ID, err)
			continue
		}

		body := "Hello " + user.Name + ",\n\nHere are your latest posts:\n\n"
		for _, post := range posts {
			body += fmt.Sprintf("Title: %s\nContent: %s\n\n", post.Title, *nullStringToStringPtr(post.Description))
		}
		body += "\nBest regards,\nYour RSS Aggregator Team"

		if err := es.sendEmail(user, "Your Daily Posts Digest", body); err != nil {
			log.Printf("Failed to send email to %v: %v", user.Email, err)
		}
	}
}

func (es *smtpServer) sendSubscribeEmail(user database.User) error {
	body := "Hello " + user.Name + ",\n\nThank you for subscribing to our service. You will now receive daily posts from your subscribed feeds.\n\nBest regards,\nYour RSS Aggregator Team"
	return es.sendEmail(user, "Welcome to RSS Aggregator", body)
}

func (es *smtpServer) sendUnsubscribeEmail(user database.User) error {
	body := "Hello " + user.Name + ",\n\nYou have successfully unsubscribed from our service. We're sorry to see you go.\n\nBest regards,\nYour RSS Aggregator Team"
	return es.sendEmail(user, "Goodbye from RSS Aggregator", body)
}

func (es *smtpServer) startDailyEmails(db *database.Queries) {
	c := cron.New()
	_, err := c.AddFunc("@daily", func() {
		fmt.Println("Sending daily emails")
		es.sendDailyPosts(db)
	})

	if err != nil {
		log.Fatalf("Failed to add cron job to send daily emails: %v", err)
	} else {
		log.Println("Cron job to send daily emails scheduled successfully")
	}

	fmt.Println("Test daily emails subscription service")
	es.sendDailyPosts(db)

	c.Start()
}
