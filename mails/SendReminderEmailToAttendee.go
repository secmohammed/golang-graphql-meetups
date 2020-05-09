package mails

import (
    "log"
    "net/smtp"
    "os"

    "github.com/secmohammed/meetups/models"
)

//SendReminderEmailToAttendee function is used to send an email activation for the recently registered user.
func SendReminderEmailToAttendee(user *models.User, meetup *models.Meetup) {
    from := os.Getenv("MAIL_USERNAME")
    password := os.Getenv("MAIL_PASSWORD")
    host := os.Getenv("MAIL_HOST")
    port := os.Getenv("MAIL_PORT")
    msg := "From: " + from + "\n" +
        "To: " + user.Email + "\n" +
        "Subject: Reminder: For your " + meetup.Name + "\n\n" +
        "Hello, " + user.FirstName + " \n\n Don't forget that the meetup will be held at " + meetup.StartDate.Format("Mon Jan _2 15:04:05 2006") + " at " + meetup.Location
    err := smtp.SendMail(host+":"+port,
        smtp.PlainAuth("", from, password, host),
        from, []string{user.Email}, []byte(msg))
    if err != nil {
        log.Printf("smtp error: %s", err)
        return
    }
}
