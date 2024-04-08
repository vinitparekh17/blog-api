// Code generated by smithy-go-codegen DO NOT EDIT.

package ses

import (
	"context"
	"fmt"
	awsmiddleware "github.com/aws/aws-sdk-go-v2/aws/middleware"
	"github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
	"github.com/aws/smithy-go/middleware"
	smithyhttp "github.com/aws/smithy-go/transport/http"
)

// Returns the current status of Easy DKIM signing for an entity. For domain name
// identities, this operation also returns the DKIM tokens that are required for
// Easy DKIM signing, and whether Amazon SES has successfully verified that these
// tokens have been published. This operation takes a list of identities as input
// and returns the following information for each:
//   - Whether Easy DKIM signing is enabled or disabled.
//   - A set of DKIM tokens that represent the identity. If the identity is an
//     email address, the tokens represent the domain of that address.
//   - Whether Amazon SES has successfully verified the DKIM tokens published in
//     the domain's DNS. This information is only returned for domain name identities,
//     not for email addresses.
//
// This operation is throttled at one request per second and can only get DKIM
// attributes for up to 100 identities at a time. For more information about
// creating DNS records using DKIM tokens, go to the Amazon SES Developer Guide (https://docs.aws.amazon.com/ses/latest/dg/send-email-authentication-dkim-easy-managing.html)
// .
func (c *Client) GetIdentityDkimAttributes(ctx context.Context, params *GetIdentityDkimAttributesInput, optFns ...func(*Options)) (*GetIdentityDkimAttributesOutput, error) {
	if params == nil {
		params = &GetIdentityDkimAttributesInput{}
	}

	result, metadata, err := c.invokeOperation(ctx, "GetIdentityDkimAttributes", params, optFns, c.addOperationGetIdentityDkimAttributesMiddlewares)
	if err != nil {
		return nil, err
	}

	out := result.(*GetIdentityDkimAttributesOutput)
	out.ResultMetadata = metadata
	return out, nil
}

// Represents a request for the status of Amazon SES Easy DKIM signing for an
// identity. For domain identities, this request also returns the DKIM tokens that
// are required for Easy DKIM signing, and whether Amazon SES successfully verified
// that these tokens were published. For more information about Easy DKIM, see the
// Amazon SES Developer Guide (https://docs.aws.amazon.com/ses/latest/dg/send-email-authentication-dkim-easy.html)
// .
type GetIdentityDkimAttributesInput struct {

	// A list of one or more verified identities - email addresses, domains, or both.
	//
	// This member is required.
	Identities []string

	noSmithyDocumentSerde
}

// Represents the status of Amazon SES Easy DKIM signing for an identity. For
// domain identities, this response also contains the DKIM tokens that are required
// for Easy DKIM signing, and whether Amazon SES successfully verified that these
// tokens were published.
type GetIdentityDkimAttributesOutput struct {

	// The DKIM attributes for an email address or a domain.
	//
	// This member is required.
	DkimAttributes map[string]types.IdentityDkimAttributes

	// Metadata pertaining to the operation's result.
	ResultMetadata middleware.Metadata

	noSmithyDocumentSerde
}

func (c *Client) addOperationGetIdentityDkimAttributesMiddlewares(stack *middleware.Stack, options Options) (err error) {
	if err := stack.Serialize.Add(&setOperationInputMiddleware{}, middleware.After); err != nil {
		return err
	}
	err = stack.Serialize.Add(&awsAwsquery_serializeOpGetIdentityDkimAttributes{}, middleware.After)
	if err != nil {
		return err
	}
	err = stack.Deserialize.Add(&awsAwsquery_deserializeOpGetIdentityDkimAttributes{}, middleware.After)
	if err != nil {
		return err
	}
	if err := addProtocolFinalizerMiddlewares(stack, options, "GetIdentityDkimAttributes"); err != nil {
		return fmt.Errorf("add protocol finalizers: %v", err)
	}

	if err = addlegacyEndpointContextSetter(stack, options); err != nil {
		return err
	}
	if err = addSetLoggerMiddleware(stack, options); err != nil {
		return err
	}
	if err = awsmiddleware.AddClientRequestIDMiddleware(stack); err != nil {
		return err
	}
	if err = smithyhttp.AddComputeContentLengthMiddleware(stack); err != nil {
		return err
	}
	if err = addResolveEndpointMiddleware(stack, options); err != nil {
		return err
	}
	if err = v4.AddComputePayloadSHA256Middleware(stack); err != nil {
		return err
	}
	if err = addRetryMiddlewares(stack, options); err != nil {
		return err
	}
	if err = awsmiddleware.AddRawResponseToMetadata(stack); err != nil {
		return err
	}
	if err = awsmiddleware.AddRecordResponseTiming(stack); err != nil {
		return err
	}
	if err = addClientUserAgent(stack, options); err != nil {
		return err
	}
	if err = smithyhttp.AddErrorCloseResponseBodyMiddleware(stack); err != nil {
		return err
	}
	if err = smithyhttp.AddCloseResponseBodyMiddleware(stack); err != nil {
		return err
	}
	if err = addSetLegacyContextSigningOptionsMiddleware(stack); err != nil {
		return err
	}
	if err = addOpGetIdentityDkimAttributesValidationMiddleware(stack); err != nil {
		return err
	}
	if err = stack.Initialize.Add(newServiceMetadataMiddleware_opGetIdentityDkimAttributes(options.Region), middleware.Before); err != nil {
		return err
	}
	if err = awsmiddleware.AddRecursionDetection(stack); err != nil {
		return err
	}
	if err = addRequestIDRetrieverMiddleware(stack); err != nil {
		return err
	}
	if err = addResponseErrorMiddleware(stack); err != nil {
		return err
	}
	if err = addRequestResponseLogging(stack, options); err != nil {
		return err
	}
	if err = addDisableHTTPSMiddleware(stack, options); err != nil {
		return err
	}
	return nil
}

func newServiceMetadataMiddleware_opGetIdentityDkimAttributes(region string) *awsmiddleware.RegisterServiceMetadata {
	return &awsmiddleware.RegisterServiceMetadata{
		Region:        region,
		ServiceID:     ServiceID,
		OperationName: "GetIdentityDkimAttributes",
	}
}
