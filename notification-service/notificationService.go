package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	protoGorush "github.com/appleboy/gorush/rpc/proto"
	constRes "github.com/asciiu/gomo/common/constants/response"
	protoDevice "github.com/asciiu/gomo/device-service/proto/device"
	repoNotification "github.com/asciiu/gomo/notification-service/db/sql"
	protoNotification "github.com/asciiu/gomo/notification-service/proto/notification"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/google/uuid"
	micro "github.com/micro/go-micro"
	"google.golang.org/grpc"
)

const (
	// The character encoding for the email.
	AwsRegion = "us-east-1"
	CharSet   = "UTF-8"
)

type NotificationService struct {
	db      *sql.DB
	devices protoDevice.DeviceServiceClient
	client  protoGorush.GorushClient
	topic   string
}

func NewNotificationService(db *sql.DB, service micro.Service) *NotificationService {
	address := fmt.Sprintf("%s", os.Getenv("GORUSH_ADDRESS"))
	topic := fmt.Sprintf("%s", os.Getenv("APNS_TOPIC"))

	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	client := protoGorush.NewGorushClient(conn)

	hs := NotificationService{
		db:      db,
		client:  client,
		topic:   topic,
		devices: protoDevice.NewDeviceServiceClient("devices", service.Client()),
	}

	return &hs
}

func (service *NotificationService) LogActivity(ctx context.Context, history *protoNotification.Activity) error {

	_, error := repoNotification.InsertActivity(service.db, history)
	if error != nil {
		log.Println("could not insert new notification ", error)
	}

	log.Println("notification ", history.Description)

	getRequest := protoDevice.GetUserDevicesRequest{
		UserID: history.UserID,
	}

	// get device tokens from DB for user ID
	ds, _ := service.devices.GetUserDevices(context.Background(), &getRequest)

	if ds.Status != "success" {
		log.Println("error from GetUserDevices ", ds.Message)
		return errors.New(ds.Message)
	}

	iosTokens := make([]string, 0)

	for _, thing := range ds.Data.Devices {
		deviceType := thing.DeviceType
		deviceToken := thing.DeviceToken
		if deviceType == "ios" {
			iosTokens = append(iosTokens, deviceToken)
		}
	}

	// loop over device tokens and send
	// TODO fill in the rest
	r, err := service.client.Send(context.Background(), &protoGorush.NotificationRequest{
		Platform: 1,
		Tokens:   iosTokens,
		Message:  "fomo",
		Badge:    1,
		Category: "test",
		Sound:    "1.caf",
		Topic:    service.topic,
		Alert: &protoGorush.Alert{
			Title:    history.Title,
			Body:     history.Description,
			Subtitle: history.Subtitle,
			LocKey:   "Test loc key",
			LocArgs:  []string{"test", "test"},
		},
	})

	if err != nil {
		log.Printf("could not send: %v\n", err)
	} else {
		log.Printf("Success: %t\n", r.Success)
	}

	return nil
}

func (service *NotificationService) FindUserActivity(ctx context.Context, req *protoNotification.ActivityRequest, res *protoNotification.ActivityPagedResponse) error {
	var pagedResult *protoNotification.UserActivityPage
	var err error

	if req.ObjectID != "" {
		if _, err := uuid.Parse(req.ObjectID); err != nil {
			res.Status = constRes.Fail
			res.Message = fmt.Sprintf("object %s not found", req.ObjectID)
			return nil
		}
	}

	if req.ObjectID != "" {
		// history associated with object ID only
		pagedResult, err = repoNotification.FindObjectActivity(service.db, req)
	} else {
		// all user history
		//pagedResult, err = repoNotification.FindUserActivity(service.db, req.UserID, req.Page, req.PageSize)
		pagedResult, err = repoNotification.FindUserPlansActivity(service.db, req.UserID, req.Page, req.PageSize)
	}

	if err == nil {
		res.Status = constRes.Success
		res.Data = pagedResult
	} else {
		res.Status = constRes.Error
		res.Message = err.Error()
	}
	return nil
}

func (service *NotificationService) FindMostRecentActivity(ctx context.Context, req *protoNotification.RecentActivityRequest, res *protoNotification.ActivityListResponse) error {
	history, err := repoNotification.FindRecentObjectActivity(service.db, req)
	if err == nil {
		res.Status = constRes.Success
		res.Data = &protoNotification.ActivityList{
			Activity: history,
		}
	} else {
		res.Status = constRes.Error
		res.Message = err.Error()
	}
	return nil
}

func (service *NotificationService) FindActivityCount(ctx context.Context, req *protoNotification.ActivityCountRequest, res *protoNotification.ActivityCountResponse) error {
	count := repoNotification.FindObjectActivityCount(service.db, req.ObjectID)
	res.Status = constRes.Success
	res.Data = &protoNotification.ActivityCount{
		Count: count,
	}
	return nil
}

func (service *NotificationService) UpdateActivity(ctx context.Context, req *protoNotification.UpdateActivityRequest, res *protoNotification.ActivityResponse) error {

	history, err := repoNotification.FindActivity(service.db, req.ActivityID)

	if err == sql.ErrNoRows {
		res.Status = constRes.Nonentity
		res.Message = "activity not found"
		return nil
	}

	if req.ClickedAt != "" {
		// time.Parse(time.RFC3339, "2006-01-02T15:04:05Z")
		_, err = time.Parse(time.RFC3339, req.ClickedAt)
		if err != nil {
			res.Status = constRes.Fail
			res.Message = "clickedAt must be RFC3339 format: e.g. 2006-01-02T15:04:05Z"
			return nil
		}

		history, err = repoNotification.UpdateActivityClickedAt(service.db, req.ActivityID, req.ClickedAt)
		if err != nil {
			res.Status = constRes.Error
			res.Message = err.Error()
			return nil
		}
	}

	if req.SeenAt != "" {
		_, err = time.Parse(time.RFC3339, req.SeenAt)
		if err != nil {
			res.Status = constRes.Fail
			res.Message = "seenAt must be RFC3339 format: e.g. 2006-01-02T15:04:05Z"
			return nil
		}

		history, err = repoNotification.UpdateActivitySeenAt(service.db, req.ActivityID, req.SeenAt)
		if err != nil {
			res.Status = constRes.Error
			res.Message = err.Error()
			return nil
		}
	}

	res.Status = constRes.Success
	res.Data = &protoNotification.ActivityData{
		Activity: history,
	}

	return nil
}

func (service *NotificationService) SendEmail(ctx context.Context, req *protoNotification.EmailRequest, res *protoNotification.EmailResponse) error {
	// Create a new session in the us-west-2 region.
	// Replace us-west-2 with the AWS Region you're using for Amazon SES.
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(AwsRegion)},
	)

	// Create an SES session.
	svc := ses.New(sess)

	// Assemble the email.
	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			CcAddresses: []*string{},
			ToAddresses: []*string{
				aws.String(req.EmailRecipient),
			},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String(CharSet),
					Data:    aws.String(req.HtmlBody),
				},
				Text: &ses.Content{
					Charset: aws.String(CharSet),
					Data:    aws.String(req.TextBody),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String(CharSet),
				Data:    aws.String(req.Subject),
			},
		},
		Source: aws.String(req.EmailSender),
		// Uncomment to use a configuration set
		//ConfigurationSetName: aws.String(ConfigurationSet),
	}

	// Attempt to send the email.
	_, err = svc.SendEmail(input)

	// Display error messages if they occur.
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case ses.ErrCodeMessageRejected:
				return errors.New(fmt.Sprintf("%s: %s", ses.ErrCodeMessageRejected, aerr.Error()))
			case ses.ErrCodeMailFromDomainNotVerifiedException:
				return errors.New(fmt.Sprintf("%s: %s", ses.ErrCodeMailFromDomainNotVerifiedException, aerr.Error()))
			case ses.ErrCodeConfigurationSetDoesNotExistException:
				return errors.New(fmt.Sprintf("%s: %s", ses.ErrCodeConfigurationSetDoesNotExistException, aerr.Error()))
			default:
				return errors.New(aerr.Error())
			}
		}
		return err
	}

	res.Status = constRes.Success

	return nil
}

func (service *NotificationService) CreateTemplate(ctx context.Context, req *protoNotification.CreateTemplateRequest, res *protoNotification.EmailResponse) error {
	// Create a new session in the us-west-2 region.
	// Replace us-west-2 with the AWS Region you're using for Amazon SES.
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(AwsRegion)},
	)

	// Create an SES session.
	svc := ses.New(sess)

	// Assemble the email.
	input := &ses.CreateTemplateInput{
		Template: &ses.Template{
			TemplateName: aws.String(req.TemplateName),
			HtmlPart:     aws.String(req.Html),
			TextPart:     aws.String(req.Text),
			SubjectPart:  aws.String(req.Subject),
		},
	}

	// Attempt to send the email.
	_, err = svc.CreateTemplate(input)

	// Display error messages if they occur.
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			return errors.New(aerr.Code())
		}
		return err
	}

	res.Status = constRes.Success

	return nil
}

func (service *NotificationService) DeleteTemplate(ctx context.Context, req *protoNotification.DeleteTemplateRequest, res *protoNotification.EmailResponse) error {
	// Create a new session in the us-west-2 region.
	// Replace us-west-2 with the AWS Region you're using for Amazon SES.
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(AwsRegion)},
	)

	// Create an SES session.
	svc := ses.New(sess)

	// Assemble the email.
	input := &ses.DeleteTemplateInput{
		TemplateName: aws.String(req.TemplateName),
	}

	// Attempt to send the email.
	_, err = svc.DeleteTemplate(input)

	// Display error messages if they occur.
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			return errors.New(aerr.Code())
		}
		return err
	}

	res.Status = constRes.Success

	return nil
}
