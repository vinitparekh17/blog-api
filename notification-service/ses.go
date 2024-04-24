package main

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
	pb "github.com/jay-bhogayata/notifyHub/proto"
)

// func (app *application) sendEmail(params mail_input) error {

// 	var destinations []string

// 	destinations = append(destinations, params.DestinationEmail)

// 	service := ses.NewFromConfig(app.config.awsConfig)
// 	input := &ses.SendEmailInput{
// 		Destination: &types.Destination{
// 			ToAddresses: []string{
// 				destinations[0],
// 			},
// 		},
// 		Message: &types.Message{
// 			Body: &types.Body{
// 				Text: &types.Content{
// 					Data: &params.Body,
// 				},
// 			},
// 			Subject: &types.Content{
// 				Data: &params.Subject,
// 			},
// 		},
// 		Source: &app.config.sender_email,
// 	}

// _, err := service.SendEmail(context.Background(), input)
// if err != nil {
// 	logger.Error("could not send email", "error", err)
// }

// 	return nil
// }

func (app *application) SendEmail(ctx context.Context, req *pb.EmailRequest) (*pb.EmailResponse, error) {
	var destinations []string

	destinations = append(destinations, req.To)

	service := ses.NewFromConfig(app.config.awsConfig)
	input := &ses.SendEmailInput{
		Destination: &types.Destination{
			ToAddresses: []string{
				destinations[0],
			},
		},
		Message: &types.Message{
			Body: &types.Body{
				Html: &types.Content{
					Data: &req.Body,
				},
			},
			Subject: &types.Content{
				Data: &req.Subject,
			},
		},
		Source: &app.config.sender_email,
	}

	_, err := service.SendEmail(ctx, input)
	if err != nil {
		logger.Error("could not send email", "error", err)
		return nil, err
	}

	return &pb.EmailResponse{
		Status:  true,
		Message: "Email sent successfully",
	}, nil
}
