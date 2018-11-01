package main

import (
	"context"
	"fmt"
	"testing"

	protoNotification "github.com/asciiu/gomo/notification-service/proto/notification"
	repoUser "github.com/asciiu/gomo/user-service/db/sql"
	"github.com/stretchr/testify/assert"
)

func TestEmail(t *testing.T) {
	service, user := setupService()

	defer service.db.Close()

	// The HTML body for the email.
	hBody := "<h1>Fomo SES Test Email</h1><p>This email was sent with " +
		"<a href='https://aws.amazon.com/ses/'>Amazon SES</a> using the " +
		"<a href='https://aws.amazon.com/sdk-for-go/'>AWS SDK for Go</a>.</p>"

	//The email body for recipients with non-HTML email clients.
	txtBody := "This email was sent with Amazon SES using the AWS SDK for Go."

	email := protoNotification.EmailRequest{
		Subject:        "Fomo Test",
		HtmlBody:       hBody,
		TextBody:       txtBody,
		EmailRecipient: "ellyssin.gimhae@gmail.com",
		EmailSender:    "support@projectfomo.com",
	}

	res := protoNotification.EmailResponse{}
	service.SendEmail(context.Background(), &email, &res)

	assert.Equal(t, "success", res.Status, fmt.Sprintf("%s", res.Message))

	repoUser.DeleteUserHard(service.db, user.ID)
}

func TestTemplateCreation(t *testing.T) {
	service, user := setupService()

	defer service.db.Close()

	// The HTML body for the email.
	html := "<h1>Fomo SES Test Email</h1><p>This email was sent with " +
		"<a href='https://aws.amazon.com/ses/'>Amazon SES</a> using the " +
		"<a href='https://aws.amazon.com/sdk-for-go/'>AWS SDK for Go</a>.</p>"

	//The email body for recipients with non-HTML email clients.
	txt := "This email was sent with Amazon SES using the AWS SDK for Go."

	template := protoNotification.CreateTemplateRequest{
		Subject:      "Fomo Test",
		Html:         html,
		Text:         txt,
		TemplateName: "test",
	}

	res := protoNotification.EmailResponse{}
	err := service.CreateTemplate(context.Background(), &template, &res)
	if err != nil {
		t.Fatalf(err.Error())
	}

	list := protoNotification.ListTemplatesRequest{}
	resTemps := protoNotification.TemplatesResponse{}
	err = service.ListTemplates(context.Background(), &list, &resTemps)
	if err != nil {
		t.Fatalf(err.Error())
	}
	fmt.Println(resTemps)

	del := protoNotification.DeleteTemplateRequest{
		TemplateName: "test",
	}

	err = service.DeleteTemplate(context.Background(), &del, &res)
	if err != nil {
		t.Fatalf(err.Error())
	}

	assert.Equal(t, "success", res.Status, fmt.Sprintf("%s", res.Message))

	repoUser.DeleteUserHard(service.db, user.ID)
}
