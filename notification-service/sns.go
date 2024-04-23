package main

// func isInList(numberToCheck string, arrayOfNumbers []string) error {
// 	for _, number := range arrayOfNumbers {
// 		if number == numberToCheck {
// 			return nil
// 		}
// 	}
// 	return fmt.Errorf("number not in allowed list")
// }

// var totalNumberIWant int32 = 50

// func (app *application) sendSms(smsParams sms_input) error {

// 	cfg, err := awsconfig.LoadDefaultConfig(context.TODO())

// 	if err != nil {
// 		log.Fatal(err.Error())

// 	}

// 	service := sns.NewFromConfig(cfg)

// 	input := &sns.ListSMSSandboxPhoneNumbersInput{
// 		MaxResults: &totalNumberIWant,
// 	}

// 	res, err := service.ListSMSSandboxPhoneNumbers(context.TODO(), input)
// 	if err != nil {
// 		logger.Error(err.Error())
// 		return err
// 	}

// 	var numbers []string

// 	for _, number := range res.PhoneNumbers {
// 		numbers = append(numbers, *number.PhoneNumber)
// 	}

// 	err = isInList(smsParams.Recipient, numbers)
// 	if err != nil {

// 		logger.Error("number not in list")
// 		return err
// 	}

// 	publishInput := &sns.PublishInput{
// 		Message:     &smsParams.Message,
// 		PhoneNumber: &smsParams.Recipient,
// 	}

// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()

// 	result, err := service.Publish(ctx, publishInput)
// 	if err != nil {
// 		logger.Error(err.Error())
// 		return err
// 	}

// 	if result.MessageId != nil {
// 		logger.Info("Message sent successfully.", " Message ID:", *result.MessageId)
// 	}

// 	return nil
// }
