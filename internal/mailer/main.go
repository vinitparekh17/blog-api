package mailer

import (
	"bytes"
	"context"
	"html/template"
	"log"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
)

func SendEmail(dest string, subject string, body string, sender string) error {

	cfg, err := awsconfig.LoadDefaultConfig(context.TODO())

	if err != nil {
		log.Fatal(err.Error())
	}

	service := ses.NewFromConfig(cfg)
	input := &ses.SendEmailInput{
		Destination: &types.Destination{
			ToAddresses: []string{
				dest,
			},
		},
		Message: &types.Message{
			Body: &types.Body{
				Html: &types.Content{
					Data: &body,
				},
			},
			Subject: &types.Content{
				Data: &subject,
			},
		},
		Source: &sender,
	}

	_, err = service.SendEmail(context.Background(), input)
	if err != nil {
		return err
	}

	return nil
}

var body = `
<!DOCTYPE html>
<html>
<head>
    <title>Account Verification</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            background-color: #f3f4f6;
            margin: 0;
            padding: 0;
        }

        .container {
            max-width: 600px;
            margin: 20px auto;
            background-color: #ffffff;
            padding: 20px;
            border-radius: 8px;
            box-shadow: 0px 0px 10px rgba(0, 0, 0, 0.1);
        }

        h1 {
            color: #3b82f6;
            text-align: center;
        }

        p {
            color: #4b5563;
            font-size: 16px;
            line-height: 1.5;
            margin-bottom: 20px;
        }

        .btn {
            background-color: #3b82f6;
        
            padding: 10px 20px;
            text-decoration: none;
            border-radius: 4px;
            display: inline-block;
			cursor: pointer;
        }

        .btn:hover {
            background-color: #2563eb;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>Account Verification</h1>
		<p>Hi, {{.Name}}
		Thanks for signing up for the blogging site. We're excited to have you on board.
        Click the button below to verify your account for the blogging site.</p>
        <a href="{{.VerificationLink}}">
        <button class="btn">Verify Account</button>
        </a>
    </div>
</body>
</html>

`

func SetupVerificationTemplate(name, link string) (emailBody string, error error) {
	type EmailData struct {
		Name             string
		VerificationLink string
	}

	emailData := EmailData{
		Name:             name,
		VerificationLink: link,
	}

	var bodyBuffer bytes.Buffer

	tmpl, err := template.New("et").Parse(body)
	if err != nil {
		return "", err
	}

	err = tmpl.Execute(&bodyBuffer, emailData)
	if err != nil {
		return "", err
	}

	return bodyBuffer.String(), nil
}
