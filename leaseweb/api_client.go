package leaseweb

import (
	"context"

	LSW "github.com/LeaseWeb/leaseweb-go-sdk"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	leasewebAPIToken string
)

func logAPIError(ctx context.Context, err error) {
	fields := map[string]interface{}{}

	if erra, ok := err.(*LSW.ApiError); ok {
		fields["url"] = erra.Url
		fields["method"] = erra.Method
		fields["code"] = erra.Code
		fields["message"] = erra.Message
		fields["correlation_id"] = erra.CorrelationId

		if len(erra.Details) != 0 {
			for field, details := range erra.Details {
				fields["detail_"+field] = details
			}
		}
	} else {
		fields["message"] = err.Error()

		if errd, ok := err.(*LSW.DecodingError); ok {
			fields["url"] = errd.Url
			fields["method"] = errd.Method
		} else if erre, ok := err.(*LSW.EncodingError); ok {
			fields["url"] = erre.Url
			fields["method"] = erre.Method
		}
	}

	tflog.Error(ctx, "API request error", fields)
}

func getAllServers(ctx context.Context, site string) ([]LSW.DedicatedServer, error) {
	var allServers []LSW.DedicatedServer
	offset := 0
	limit := 20

	opts := LSW.DedicatedServerListOptions{
		PaginationOptions: LSW.PaginationOptions{
			Offset: &offset,
			Limit:  &limit,
		},
		Site: &site,
	}

	for {

		result, err := LSW.DedicatedServerApi{}.List(ctx, opts)
		if err != nil {
			return nil, err
		}

		if len(result.Servers) == 0 {
			break
		}

		allServers = append(allServers, result.Servers...)
		*opts.PaginationOptions.Offset += *opts.PaginationOptions.Limit
	}

	return allServers, nil
}
