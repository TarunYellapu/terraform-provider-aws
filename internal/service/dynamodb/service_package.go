// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package dynamodb

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/retry"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	awstypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/errs"
	"github.com/hashicorp/terraform-provider-aws/names"
)

// NewClient returns a new AWS SDK for Go v2 client for this service package's AWS API.
func (p *servicePackage) NewClient(ctx context.Context, config map[string]any) (*dynamodb.Client, error) {
	cfg := *(config["aws_sdkv2_config"].(*aws.Config))

	return dynamodb.NewFromConfig(cfg, func(o *dynamodb.Options) {
		if endpoint := config[names.AttrEndpoint].(string); endpoint != "" {
			o.BaseEndpoint = aws.String(endpoint)
		}

		o.Retryer = conns.AddIsErrorRetryables(cfg.Retryer().(aws.RetryerV2), retry.IsErrorRetryableFunc(func(err error) aws.Ternary {
			if errs.IsAErrorMessageContains[*awstypes.LimitExceededException](err, "Subscriber limit exceeded:") {
				return aws.TrueTernary
			}
			return aws.UnknownTernary // Delegate to configured Retryer.
		}))
	}), nil
}
