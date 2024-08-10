package cleanrooms

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cleanrooms"
	"github.com/aws/aws-sdk-go-v2/service/cleanrooms/types"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-provider-aws/internal/errs"

	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

func findConfiguredTableAssociationByID(ctx context.Context, conn *cleanrooms.Client, id string) (*cleanrooms.GetConfiguredTableAssociationOutput, error) {
	in := &cleanrooms.GetConfiguredTableAssociationInput{
		ConfiguredTableAssociationIdentifier: aws.String(id),
	}
	out, err := conn.GetConfiguredTableAssociation(ctx, in)

	if errs.IsA[*types.AccessDeniedException](err) {
		return nil, &retry.NotFoundError{
			LastError:   err,
			LastRequest: in,
		}
	}

	if err != nil {
		return nil, err
	}

	if out == nil || out.ConfiguredTableAssociation == nil {
		return nil, tfresource.NewEmptyResultError(in)
	}

	return out, nil
}
