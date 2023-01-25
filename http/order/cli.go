// Code generated by goa v3.10.2, DO NOT EDIT.
//
// order HTTP client CLI support package
//
// Command:
// $ goa gen github.com/reinventingscience/ivcap-core-api/design

package client

import (
	order "github.com/reinventingscience/ivcap-core-api/gen/order"
	"encoding/json"
	"fmt"
	"strconv"

	goa "goa.design/goa/v3/pkg"
)

// BuildListPayload builds the payload for the order list endpoint from CLI
// flags.
func BuildListPayload(orderListFilter string, orderListOrderby string, orderListTop string, orderListSkip string, orderListSelect string, orderListOffset string, orderListLimit string, orderListPageToken string, orderListJWT string) (*order.ListPayload, error) {
	var err error
	var filter string
	{
		if orderListFilter != "" {
			filter = orderListFilter
		}
	}
	var orderby string
	{
		if orderListOrderby != "" {
			orderby = orderListOrderby
		}
	}
	var top int
	{
		if orderListTop != "" {
			var v int64
			v, err = strconv.ParseInt(orderListTop, 10, strconv.IntSize)
			top = int(v)
			if err != nil {
				return nil, fmt.Errorf("invalid value for top, must be INT")
			}
			if top < 1 {
				err = goa.MergeErrors(err, goa.InvalidRangeError("top", top, 1, true))
			}
			if top > 50 {
				err = goa.MergeErrors(err, goa.InvalidRangeError("top", top, 50, false))
			}
			if err != nil {
				return nil, err
			}
		}
	}
	var skip int
	{
		if orderListSkip != "" {
			var v int64
			v, err = strconv.ParseInt(orderListSkip, 10, strconv.IntSize)
			skip = int(v)
			if err != nil {
				return nil, fmt.Errorf("invalid value for skip, must be INT")
			}
			if skip < 0 {
				err = goa.MergeErrors(err, goa.InvalidRangeError("skip", skip, 0, true))
			}
			if err != nil {
				return nil, err
			}
		}
	}
	var select_ string
	{
		if orderListSelect != "" {
			select_ = orderListSelect
		}
	}
	var offset *int
	{
		if orderListOffset != "" {
			var v int64
			v, err = strconv.ParseInt(orderListOffset, 10, strconv.IntSize)
			val := int(v)
			offset = &val
			if err != nil {
				return nil, fmt.Errorf("invalid value for offset, must be INT")
			}
			if *offset < 0 {
				err = goa.MergeErrors(err, goa.InvalidRangeError("offset", *offset, 0, true))
			}
			if err != nil {
				return nil, err
			}
		}
	}
	var limit *int
	{
		if orderListLimit != "" {
			var v int64
			v, err = strconv.ParseInt(orderListLimit, 10, strconv.IntSize)
			val := int(v)
			limit = &val
			if err != nil {
				return nil, fmt.Errorf("invalid value for limit, must be INT")
			}
			if *limit < 1 {
				err = goa.MergeErrors(err, goa.InvalidRangeError("limit", *limit, 1, true))
			}
			if *limit > 50 {
				err = goa.MergeErrors(err, goa.InvalidRangeError("limit", *limit, 50, false))
			}
			if err != nil {
				return nil, err
			}
		}
	}
	var pageToken string
	{
		if orderListPageToken != "" {
			pageToken = orderListPageToken
		}
	}
	var jwt string
	{
		jwt = orderListJWT
	}
	v := &order.ListPayload{}
	v.Filter = filter
	v.Orderby = orderby
	v.Top = top
	v.Skip = skip
	v.Select = select_
	v.Offset = offset
	v.Limit = limit
	v.PageToken = pageToken
	v.JWT = jwt

	return v, nil
}

// BuildCreatePayload builds the payload for the order create endpoint from CLI
// flags.
func BuildCreatePayload(orderCreateBody string, orderCreateJWT string) (*order.CreatePayload, error) {
	var err error
	var body CreateRequestBody
	{
		err = json.Unmarshal([]byte(orderCreateBody), &body)
		if err != nil {
			return nil, fmt.Errorf("invalid JSON for body, \nerror: %s, \nexample of valid JSON:\n%s", err, "'{\n      \"accountID\": \"123e4567-e89b-12d3-a456-426614174000\",\n      \"metadata\": {\n         \"refID\": \"33-444\"\n      },\n      \"name\": \"Fire risk for Lot2\",\n      \"parameters\": [\n         {\n            \"name\": \"region\",\n            \"value\": \"Upper Valley\"\n         },\n         {\n            \"name\": \"threshold\",\n            \"value\": 10\n         }\n      ],\n      \"serviceID\": \"123e4567-e89b-12d3-a456-426614174000\"\n   }'")
		}
		if body.Parameters == nil {
			err = goa.MergeErrors(err, goa.MissingFieldError("parameters", "body"))
		}
		err = goa.MergeErrors(err, goa.ValidateFormat("body.serviceID", body.ServiceID, goa.FormatURI))
		err = goa.MergeErrors(err, goa.ValidateFormat("body.accountID", body.AccountID, goa.FormatURI))
		if err != nil {
			return nil, err
		}
	}
	var jwt string
	{
		jwt = orderCreateJWT
	}
	v := &order.OrderRequestT{
		ServiceID: body.ServiceID,
		AccountID: body.AccountID,
		Name:      body.Name,
	}
	if body.Metadata != nil {
		v.Metadata = make(map[string]string, len(body.Metadata))
		for key, val := range body.Metadata {
			tk := key
			tv := val
			v.Metadata[tk] = tv
		}
	}
	if body.Parameters != nil {
		v.Parameters = make([]*order.ParameterT, len(body.Parameters))
		for i, val := range body.Parameters {
			v.Parameters[i] = marshalParameterTToOrderParameterT(val)
		}
	}
	res := &order.CreatePayload{
		Orders: v,
	}
	res.JWT = jwt

	return res, nil
}

// BuildReadPayload builds the payload for the order read endpoint from CLI
// flags.
func BuildReadPayload(orderReadID string, orderReadJWT string) (*order.ReadPayload, error) {
	var id string
	{
		id = orderReadID
	}
	var jwt string
	{
		jwt = orderReadJWT
	}
	v := &order.ReadPayload{}
	v.ID = id
	v.JWT = jwt

	return v, nil
}
