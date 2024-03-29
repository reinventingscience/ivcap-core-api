// Copyright 2023 Commonwealth Scientific and Industrial Research Organisation (CSIRO) ABN 41 687 119 230
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// $ goa gen github.com/reinventingscience/ivcap-core-api/design

package client

import (
	"bytes"
	order "github.com/reinventingscience/ivcap-core-api/gen/order"
	orderviews "github.com/reinventingscience/ivcap-core-api/gen/order/views"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	goahttp "goa.design/goa/v3/http"
)

// BuildReadRequest instantiates a HTTP request object with method and path set
// to call the "order" service "read" endpoint
func (c *Client) BuildReadRequest(ctx context.Context, v interface{}) (*http.Request, error) {
	var (
		id string
	)
	{
		p, ok := v.(*order.ReadPayload)
		if !ok {
			return nil, goahttp.ErrInvalidType("order", "read", "*order.ReadPayload", v)
		}
		id = p.ID
	}
	u := &url.URL{Scheme: c.scheme, Host: c.host, Path: ReadOrderPath(id)}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, goahttp.ErrInvalidURL("order", "read", u.String(), err)
	}
	if ctx != nil {
		req = req.WithContext(ctx)
	}

	return req, nil
}

// EncodeReadRequest returns an encoder for requests sent to the order read
// server.
func EncodeReadRequest(encoder func(*http.Request) goahttp.Encoder) func(*http.Request, interface{}) error {
	return func(req *http.Request, v interface{}) error {
		p, ok := v.(*order.ReadPayload)
		if !ok {
			return goahttp.ErrInvalidType("order", "read", "*order.ReadPayload", v)
		}
		{
			head := p.JWT
			if !strings.Contains(head, " ") {
				req.Header.Set("Authorization", "Bearer "+head)
			} else {
				req.Header.Set("Authorization", head)
			}
		}
		return nil
	}
}

// DecodeReadResponse returns a decoder for responses returned by the order
// read endpoint. restoreBody controls whether the response body should be
// restored after having been read.
// DecodeReadResponse may return the following errors:
//   - "bad-request" (type *order.BadRequestT): http.StatusBadRequest
//   - "invalid-credential" (type *order.InvalidCredentialsT): http.StatusBadRequest
//   - "invalid-scopes" (type *order.InvalidScopesT): http.StatusForbidden
//   - "not-implemented" (type *order.NotImplementedT): http.StatusNotImplemented
//   - "not-found" (type *order.ResourceNotFoundT): http.StatusNotFound
//   - "not-authorized" (type *order.UnauthorizedT): http.StatusUnauthorized
//   - error: internal error
func DecodeReadResponse(decoder func(*http.Response) goahttp.Decoder, restoreBody bool) func(*http.Response) (interface{}, error) {
	return func(resp *http.Response) (interface{}, error) {
		if restoreBody {
			b, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			resp.Body = io.NopCloser(bytes.NewBuffer(b))
			defer func() {
				resp.Body = io.NopCloser(bytes.NewBuffer(b))
			}()
		} else {
			defer resp.Body.Close()
		}
		switch resp.StatusCode {
		case http.StatusOK:
			var (
				body ReadResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("order", "read", err)
			}
			p := NewReadOrderStatusRTOK(&body)
			view := resp.Header.Get("goa-view")
			vres := &orderviews.OrderStatusRT{Projected: p, View: view}
			if err = orderviews.ValidateOrderStatusRT(vres); err != nil {
				return nil, goahttp.ErrValidationError("order", "read", err)
			}
			res := order.NewOrderStatusRT(vres)
			return res, nil
		case http.StatusBadRequest:
			en := resp.Header.Get("goa-error")
			switch en {
			case "bad-request":
				var (
					body ReadBadRequestResponseBody
					err  error
				)
				err = decoder(resp).Decode(&body)
				if err != nil {
					return nil, goahttp.ErrDecodingError("order", "read", err)
				}
				err = ValidateReadBadRequestResponseBody(&body)
				if err != nil {
					return nil, goahttp.ErrValidationError("order", "read", err)
				}
				return nil, NewReadBadRequest(&body)
			case "invalid-credential":
				return nil, NewReadInvalidCredential()
			default:
				body, _ := io.ReadAll(resp.Body)
				return nil, goahttp.ErrInvalidResponse("order", "read", resp.StatusCode, string(body))
			}
		case http.StatusForbidden:
			var (
				body ReadInvalidScopesResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("order", "read", err)
			}
			err = ValidateReadInvalidScopesResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("order", "read", err)
			}
			return nil, NewReadInvalidScopes(&body)
		case http.StatusNotImplemented:
			var (
				body ReadNotImplementedResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("order", "read", err)
			}
			err = ValidateReadNotImplementedResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("order", "read", err)
			}
			return nil, NewReadNotImplemented(&body)
		case http.StatusNotFound:
			var (
				body ReadNotFoundResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("order", "read", err)
			}
			err = ValidateReadNotFoundResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("order", "read", err)
			}
			return nil, NewReadNotFound(&body)
		case http.StatusUnauthorized:
			return nil, NewReadNotAuthorized()
		default:
			body, _ := io.ReadAll(resp.Body)
			return nil, goahttp.ErrInvalidResponse("order", "read", resp.StatusCode, string(body))
		}
	}
}

// BuildListRequest instantiates a HTTP request object with method and path set
// to call the "order" service "list" endpoint
func (c *Client) BuildListRequest(ctx context.Context, v interface{}) (*http.Request, error) {
	u := &url.URL{Scheme: c.scheme, Host: c.host, Path: ListOrderPath()}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, goahttp.ErrInvalidURL("order", "list", u.String(), err)
	}
	if ctx != nil {
		req = req.WithContext(ctx)
	}

	return req, nil
}

// EncodeListRequest returns an encoder for requests sent to the order list
// server.
func EncodeListRequest(encoder func(*http.Request) goahttp.Encoder) func(*http.Request, interface{}) error {
	return func(req *http.Request, v interface{}) error {
		p, ok := v.(*order.ListPayload)
		if !ok {
			return goahttp.ErrInvalidType("order", "list", "*order.ListPayload", v)
		}
		{
			head := p.JWT
			if !strings.Contains(head, " ") {
				req.Header.Set("Authorization", "Bearer "+head)
			} else {
				req.Header.Set("Authorization", head)
			}
		}
		values := req.URL.Query()
		values.Add("limit", fmt.Sprintf("%v", p.Limit))
		if p.Page != nil {
			values.Add("page", *p.Page)
		}
		if p.Filter != nil {
			values.Add("filter", *p.Filter)
		}
		if p.OrderBy != nil {
			values.Add("order-by", *p.OrderBy)
		}
		values.Add("order-desc", fmt.Sprintf("%v", p.OrderDesc))
		if p.AtTime != nil {
			values.Add("at-time", *p.AtTime)
		}
		req.URL.RawQuery = values.Encode()
		return nil
	}
}

// DecodeListResponse returns a decoder for responses returned by the order
// list endpoint. restoreBody controls whether the response body should be
// restored after having been read.
// DecodeListResponse may return the following errors:
//   - "bad-request" (type *order.BadRequestT): http.StatusBadRequest
//   - "invalid-credential" (type *order.InvalidCredentialsT): http.StatusBadRequest
//   - "invalid-parameter" (type *order.InvalidParameterValue): http.StatusUnprocessableEntity
//   - "invalid-scopes" (type *order.InvalidScopesT): http.StatusForbidden
//   - "not-implemented" (type *order.NotImplementedT): http.StatusNotImplemented
//   - "not-authorized" (type *order.UnauthorizedT): http.StatusUnauthorized
//   - error: internal error
func DecodeListResponse(decoder func(*http.Response) goahttp.Decoder, restoreBody bool) func(*http.Response) (interface{}, error) {
	return func(resp *http.Response) (interface{}, error) {
		if restoreBody {
			b, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			resp.Body = io.NopCloser(bytes.NewBuffer(b))
			defer func() {
				resp.Body = io.NopCloser(bytes.NewBuffer(b))
			}()
		} else {
			defer resp.Body.Close()
		}
		switch resp.StatusCode {
		case http.StatusOK:
			var (
				body ListResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("order", "list", err)
			}
			p := NewListOrderListRTOK(&body)
			view := "default"
			vres := &orderviews.OrderListRT{Projected: p, View: view}
			if err = orderviews.ValidateOrderListRT(vres); err != nil {
				return nil, goahttp.ErrValidationError("order", "list", err)
			}
			res := order.NewOrderListRT(vres)
			return res, nil
		case http.StatusBadRequest:
			en := resp.Header.Get("goa-error")
			switch en {
			case "bad-request":
				var (
					body ListBadRequestResponseBody
					err  error
				)
				err = decoder(resp).Decode(&body)
				if err != nil {
					return nil, goahttp.ErrDecodingError("order", "list", err)
				}
				err = ValidateListBadRequestResponseBody(&body)
				if err != nil {
					return nil, goahttp.ErrValidationError("order", "list", err)
				}
				return nil, NewListBadRequest(&body)
			case "invalid-credential":
				return nil, NewListInvalidCredential()
			default:
				body, _ := io.ReadAll(resp.Body)
				return nil, goahttp.ErrInvalidResponse("order", "list", resp.StatusCode, string(body))
			}
		case http.StatusUnprocessableEntity:
			var (
				body ListInvalidParameterResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("order", "list", err)
			}
			err = ValidateListInvalidParameterResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("order", "list", err)
			}
			return nil, NewListInvalidParameter(&body)
		case http.StatusForbidden:
			var (
				body ListInvalidScopesResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("order", "list", err)
			}
			err = ValidateListInvalidScopesResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("order", "list", err)
			}
			return nil, NewListInvalidScopes(&body)
		case http.StatusNotImplemented:
			var (
				body ListNotImplementedResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("order", "list", err)
			}
			err = ValidateListNotImplementedResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("order", "list", err)
			}
			return nil, NewListNotImplemented(&body)
		case http.StatusUnauthorized:
			return nil, NewListNotAuthorized()
		default:
			body, _ := io.ReadAll(resp.Body)
			return nil, goahttp.ErrInvalidResponse("order", "list", resp.StatusCode, string(body))
		}
	}
}

// BuildCreateRequest instantiates a HTTP request object with method and path
// set to call the "order" service "create" endpoint
func (c *Client) BuildCreateRequest(ctx context.Context, v interface{}) (*http.Request, error) {
	u := &url.URL{Scheme: c.scheme, Host: c.host, Path: CreateOrderPath()}
	req, err := http.NewRequest("POST", u.String(), nil)
	if err != nil {
		return nil, goahttp.ErrInvalidURL("order", "create", u.String(), err)
	}
	if ctx != nil {
		req = req.WithContext(ctx)
	}

	return req, nil
}

// EncodeCreateRequest returns an encoder for requests sent to the order create
// server.
func EncodeCreateRequest(encoder func(*http.Request) goahttp.Encoder) func(*http.Request, interface{}) error {
	return func(req *http.Request, v interface{}) error {
		p, ok := v.(*order.CreatePayload)
		if !ok {
			return goahttp.ErrInvalidType("order", "create", "*order.CreatePayload", v)
		}
		{
			head := p.JWT
			if !strings.Contains(head, " ") {
				req.Header.Set("Authorization", "Bearer "+head)
			} else {
				req.Header.Set("Authorization", head)
			}
		}
		body := NewCreateRequestBody(p)
		if err := encoder(req).Encode(&body); err != nil {
			return goahttp.ErrEncodingError("order", "create", err)
		}
		return nil
	}
}

// DecodeCreateResponse returns a decoder for responses returned by the order
// create endpoint. restoreBody controls whether the response body should be
// restored after having been read.
// DecodeCreateResponse may return the following errors:
//   - "bad-request" (type *order.BadRequestT): http.StatusBadRequest
//   - "invalid-credential" (type *order.InvalidCredentialsT): http.StatusBadRequest
//   - "invalid-parameter" (type *order.InvalidParameterValue): http.StatusUnprocessableEntity
//   - "invalid-scopes" (type *order.InvalidScopesT): http.StatusForbidden
//   - "not-implemented" (type *order.NotImplementedT): http.StatusNotImplemented
//   - "not-found" (type *order.ResourceNotFoundT): http.StatusNotFound
//   - "not-available" (type *order.ServiceNotAvailableT): http.StatusServiceUnavailable
//   - "not-authorized" (type *order.UnauthorizedT): http.StatusUnauthorized
//   - error: internal error
func DecodeCreateResponse(decoder func(*http.Response) goahttp.Decoder, restoreBody bool) func(*http.Response) (interface{}, error) {
	return func(resp *http.Response) (interface{}, error) {
		if restoreBody {
			b, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			resp.Body = io.NopCloser(bytes.NewBuffer(b))
			defer func() {
				resp.Body = io.NopCloser(bytes.NewBuffer(b))
			}()
		} else {
			defer resp.Body.Close()
		}
		switch resp.StatusCode {
		case http.StatusOK:
			var (
				body CreateResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("order", "create", err)
			}
			p := NewCreateOrderStatusRTOK(&body)
			view := resp.Header.Get("goa-view")
			vres := &orderviews.OrderStatusRT{Projected: p, View: view}
			if err = orderviews.ValidateOrderStatusRT(vres); err != nil {
				return nil, goahttp.ErrValidationError("order", "create", err)
			}
			res := order.NewOrderStatusRT(vres)
			return res, nil
		case http.StatusBadRequest:
			en := resp.Header.Get("goa-error")
			switch en {
			case "bad-request":
				var (
					body CreateBadRequestResponseBody
					err  error
				)
				err = decoder(resp).Decode(&body)
				if err != nil {
					return nil, goahttp.ErrDecodingError("order", "create", err)
				}
				err = ValidateCreateBadRequestResponseBody(&body)
				if err != nil {
					return nil, goahttp.ErrValidationError("order", "create", err)
				}
				return nil, NewCreateBadRequest(&body)
			case "invalid-credential":
				return nil, NewCreateInvalidCredential()
			default:
				body, _ := io.ReadAll(resp.Body)
				return nil, goahttp.ErrInvalidResponse("order", "create", resp.StatusCode, string(body))
			}
		case http.StatusUnprocessableEntity:
			var (
				body CreateInvalidParameterResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("order", "create", err)
			}
			err = ValidateCreateInvalidParameterResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("order", "create", err)
			}
			return nil, NewCreateInvalidParameter(&body)
		case http.StatusForbidden:
			var (
				body CreateInvalidScopesResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("order", "create", err)
			}
			err = ValidateCreateInvalidScopesResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("order", "create", err)
			}
			return nil, NewCreateInvalidScopes(&body)
		case http.StatusNotImplemented:
			var (
				body CreateNotImplementedResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("order", "create", err)
			}
			err = ValidateCreateNotImplementedResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("order", "create", err)
			}
			return nil, NewCreateNotImplemented(&body)
		case http.StatusNotFound:
			var (
				body CreateNotFoundResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("order", "create", err)
			}
			err = ValidateCreateNotFoundResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("order", "create", err)
			}
			return nil, NewCreateNotFound(&body)
		case http.StatusServiceUnavailable:
			return nil, NewCreateNotAvailable()
		case http.StatusUnauthorized:
			return nil, NewCreateNotAuthorized()
		default:
			body, _ := io.ReadAll(resp.Body)
			return nil, goahttp.ErrInvalidResponse("order", "create", resp.StatusCode, string(body))
		}
	}
}

// BuildLogsRequest instantiates a HTTP request object with method and path set
// to call the "order" service "logs" endpoint
func (c *Client) BuildLogsRequest(ctx context.Context, v interface{}) (*http.Request, error) {
	u := &url.URL{Scheme: c.scheme, Host: c.host, Path: LogsOrderPath()}
	req, err := http.NewRequest("POST", u.String(), nil)
	if err != nil {
		return nil, goahttp.ErrInvalidURL("order", "logs", u.String(), err)
	}
	if ctx != nil {
		req = req.WithContext(ctx)
	}

	return req, nil
}

// EncodeLogsRequest returns an encoder for requests sent to the order logs
// server.
func EncodeLogsRequest(encoder func(*http.Request) goahttp.Encoder) func(*http.Request, interface{}) error {
	return func(req *http.Request, v interface{}) error {
		p, ok := v.(*order.LogsPayload)
		if !ok {
			return goahttp.ErrInvalidType("order", "logs", "*order.LogsPayload", v)
		}
		{
			head := p.JWT
			if !strings.Contains(head, " ") {
				req.Header.Set("Authorization", "Bearer "+head)
			} else {
				req.Header.Set("Authorization", head)
			}
		}
		body := NewLogsRequestBody(p)
		if err := encoder(req).Encode(&body); err != nil {
			return goahttp.ErrEncodingError("order", "logs", err)
		}
		return nil
	}
}

// DecodeLogsResponse returns a decoder for responses returned by the order
// logs endpoint. restoreBody controls whether the response body should be
// restored after having been read.
// DecodeLogsResponse may return the following errors:
//   - "bad-request" (type *order.BadRequestT): http.StatusBadRequest
//   - "invalid-credential" (type *order.InvalidCredentialsT): http.StatusBadRequest
//   - "invalid-parameter" (type *order.InvalidParameterValue): http.StatusUnprocessableEntity
//   - "invalid-scopes" (type *order.InvalidScopesT): http.StatusForbidden
//   - "not-implemented" (type *order.NotImplementedT): http.StatusNotImplemented
//   - "not-found" (type *order.ResourceNotFoundT): http.StatusNotFound
//   - "not-authorized" (type *order.UnauthorizedT): http.StatusUnauthorized
//   - error: internal error
func DecodeLogsResponse(decoder func(*http.Response) goahttp.Decoder, restoreBody bool) func(*http.Response) (interface{}, error) {
	return func(resp *http.Response) (interface{}, error) {
		if restoreBody {
			b, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			resp.Body = io.NopCloser(bytes.NewBuffer(b))
			defer func() {
				resp.Body = io.NopCloser(bytes.NewBuffer(b))
			}()
		}
		switch resp.StatusCode {
		case http.StatusOK:
			return nil, nil
		case http.StatusBadRequest:
			en := resp.Header.Get("goa-error")
			switch en {
			case "bad-request":
				var (
					body LogsBadRequestResponseBody
					err  error
				)
				err = decoder(resp).Decode(&body)
				if err != nil {
					return nil, goahttp.ErrDecodingError("order", "logs", err)
				}
				err = ValidateLogsBadRequestResponseBody(&body)
				if err != nil {
					return nil, goahttp.ErrValidationError("order", "logs", err)
				}
				return nil, NewLogsBadRequest(&body)
			case "invalid-credential":
				return nil, NewLogsInvalidCredential()
			default:
				body, _ := io.ReadAll(resp.Body)
				return nil, goahttp.ErrInvalidResponse("order", "logs", resp.StatusCode, string(body))
			}
		case http.StatusUnprocessableEntity:
			var (
				body LogsInvalidParameterResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("order", "logs", err)
			}
			err = ValidateLogsInvalidParameterResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("order", "logs", err)
			}
			return nil, NewLogsInvalidParameter(&body)
		case http.StatusForbidden:
			var (
				body LogsInvalidScopesResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("order", "logs", err)
			}
			err = ValidateLogsInvalidScopesResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("order", "logs", err)
			}
			return nil, NewLogsInvalidScopes(&body)
		case http.StatusNotImplemented:
			var (
				body LogsNotImplementedResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("order", "logs", err)
			}
			err = ValidateLogsNotImplementedResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("order", "logs", err)
			}
			return nil, NewLogsNotImplemented(&body)
		case http.StatusNotFound:
			var (
				body LogsNotFoundResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("order", "logs", err)
			}
			err = ValidateLogsNotFoundResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("order", "logs", err)
			}
			return nil, NewLogsNotFound(&body)
		case http.StatusUnauthorized:
			return nil, NewLogsNotAuthorized()
		default:
			body, _ := io.ReadAll(resp.Body)
			return nil, goahttp.ErrInvalidResponse("order", "logs", resp.StatusCode, string(body))
		}
	}
}

// BuildTopRequest instantiates a HTTP request object with method and path set
// to call the "order" service "top" endpoint
func (c *Client) BuildTopRequest(ctx context.Context, v interface{}) (*http.Request, error) {
	u := &url.URL{Scheme: c.scheme, Host: c.host, Path: TopOrderPath()}
	req, err := http.NewRequest("POST", u.String(), nil)
	if err != nil {
		return nil, goahttp.ErrInvalidURL("order", "top", u.String(), err)
	}
	if ctx != nil {
		req = req.WithContext(ctx)
	}

	return req, nil
}

// EncodeTopRequest returns an encoder for requests sent to the order top
// server.
func EncodeTopRequest(encoder func(*http.Request) goahttp.Encoder) func(*http.Request, interface{}) error {
	return func(req *http.Request, v interface{}) error {
		p, ok := v.(*order.TopPayload)
		if !ok {
			return goahttp.ErrInvalidType("order", "top", "*order.TopPayload", v)
		}
		{
			head := p.JWT
			if !strings.Contains(head, " ") {
				req.Header.Set("Authorization", "Bearer "+head)
			} else {
				req.Header.Set("Authorization", head)
			}
		}
		body := NewTopRequestBody(p)
		if err := encoder(req).Encode(&body); err != nil {
			return goahttp.ErrEncodingError("order", "top", err)
		}
		return nil
	}
}

// DecodeTopResponse returns a decoder for responses returned by the order top
// endpoint. restoreBody controls whether the response body should be restored
// after having been read.
// DecodeTopResponse may return the following errors:
//   - "bad-request" (type *order.BadRequestT): http.StatusBadRequest
//   - "invalid-credential" (type *order.InvalidCredentialsT): http.StatusBadRequest
//   - "invalid-parameter" (type *order.InvalidParameterValue): http.StatusUnprocessableEntity
//   - "invalid-scopes" (type *order.InvalidScopesT): http.StatusForbidden
//   - "not-implemented" (type *order.NotImplementedT): http.StatusNotImplemented
//   - "not-found" (type *order.ResourceNotFoundT): http.StatusNotFound
//   - "not-authorized" (type *order.UnauthorizedT): http.StatusUnauthorized
//   - error: internal error
func DecodeTopResponse(decoder func(*http.Response) goahttp.Decoder, restoreBody bool) func(*http.Response) (interface{}, error) {
	return func(resp *http.Response) (interface{}, error) {
		if restoreBody {
			b, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			resp.Body = io.NopCloser(bytes.NewBuffer(b))
			defer func() {
				resp.Body = io.NopCloser(bytes.NewBuffer(b))
			}()
		} else {
			defer resp.Body.Close()
		}
		switch resp.StatusCode {
		case http.StatusOK:
			var (
				body TopResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("order", "top", err)
			}
			p := NewTopOrderTopResultItemCollectionOK(body)
			view := "default"
			vres := orderviews.OrderTopResultItemCollection{Projected: p, View: view}
			if err = orderviews.ValidateOrderTopResultItemCollection(vres); err != nil {
				return nil, goahttp.ErrValidationError("order", "top", err)
			}
			res := order.NewOrderTopResultItemCollection(vres)
			return res, nil
		case http.StatusBadRequest:
			en := resp.Header.Get("goa-error")
			switch en {
			case "bad-request":
				var (
					body TopBadRequestResponseBody
					err  error
				)
				err = decoder(resp).Decode(&body)
				if err != nil {
					return nil, goahttp.ErrDecodingError("order", "top", err)
				}
				err = ValidateTopBadRequestResponseBody(&body)
				if err != nil {
					return nil, goahttp.ErrValidationError("order", "top", err)
				}
				return nil, NewTopBadRequest(&body)
			case "invalid-credential":
				return nil, NewTopInvalidCredential()
			default:
				body, _ := io.ReadAll(resp.Body)
				return nil, goahttp.ErrInvalidResponse("order", "top", resp.StatusCode, string(body))
			}
		case http.StatusUnprocessableEntity:
			var (
				body TopInvalidParameterResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("order", "top", err)
			}
			err = ValidateTopInvalidParameterResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("order", "top", err)
			}
			return nil, NewTopInvalidParameter(&body)
		case http.StatusForbidden:
			var (
				body TopInvalidScopesResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("order", "top", err)
			}
			err = ValidateTopInvalidScopesResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("order", "top", err)
			}
			return nil, NewTopInvalidScopes(&body)
		case http.StatusNotImplemented:
			var (
				body TopNotImplementedResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("order", "top", err)
			}
			err = ValidateTopNotImplementedResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("order", "top", err)
			}
			return nil, NewTopNotImplemented(&body)
		case http.StatusNotFound:
			var (
				body TopNotFoundResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("order", "top", err)
			}
			err = ValidateTopNotFoundResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("order", "top", err)
			}
			return nil, NewTopNotFound(&body)
		case http.StatusUnauthorized:
			return nil, NewTopNotAuthorized()
		default:
			body, _ := io.ReadAll(resp.Body)
			return nil, goahttp.ErrInvalidResponse("order", "top", resp.StatusCode, string(body))
		}
	}
}

// unmarshalProductTResponseBodyToOrderviewsProductTView builds a value of type
// *orderviews.ProductTView from a value of type *ProductTResponseBody.
func unmarshalProductTResponseBodyToOrderviewsProductTView(v *ProductTResponseBody) *orderviews.ProductTView {
	if v == nil {
		return nil
	}
	res := &orderviews.ProductTView{
		ID:       v.ID,
		Name:     v.Name,
		Status:   v.Status,
		MimeType: v.MimeType,
		Size:     v.Size,
		Etag:     v.Etag,
	}
	if v.Links != nil {
		res.Links = unmarshalSelfWithDataTResponseBodyToOrderviewsSelfWithDataTView(v.Links)
	}

	return res
}

// unmarshalSelfWithDataTResponseBodyToOrderviewsSelfWithDataTView builds a
// value of type *orderviews.SelfWithDataTView from a value of type
// *SelfWithDataTResponseBody.
func unmarshalSelfWithDataTResponseBodyToOrderviewsSelfWithDataTView(v *SelfWithDataTResponseBody) *orderviews.SelfWithDataTView {
	if v == nil {
		return nil
	}
	res := &orderviews.SelfWithDataTView{
		Self: v.Self,
		Data: v.Data,
	}
	if v.DescribedBy != nil {
		res.DescribedBy = unmarshalDescribedByTResponseBodyToOrderviewsDescribedByTView(v.DescribedBy)
	}

	return res
}

// unmarshalDescribedByTResponseBodyToOrderviewsDescribedByTView builds a value
// of type *orderviews.DescribedByTView from a value of type
// *DescribedByTResponseBody.
func unmarshalDescribedByTResponseBodyToOrderviewsDescribedByTView(v *DescribedByTResponseBody) *orderviews.DescribedByTView {
	if v == nil {
		return nil
	}
	res := &orderviews.DescribedByTView{
		Href: v.Href,
		Type: v.Type,
	}

	return res
}

// unmarshalRefTResponseBodyToOrderviewsRefTView builds a value of type
// *orderviews.RefTView from a value of type *RefTResponseBody.
func unmarshalRefTResponseBodyToOrderviewsRefTView(v *RefTResponseBody) *orderviews.RefTView {
	if v == nil {
		return nil
	}
	res := &orderviews.RefTView{
		ID: v.ID,
	}
	if v.Links != nil {
		res.Links = unmarshalSelfTResponseBodyToOrderviewsSelfTView(v.Links)
	}

	return res
}

// unmarshalSelfTResponseBodyToOrderviewsSelfTView builds a value of type
// *orderviews.SelfTView from a value of type *SelfTResponseBody.
func unmarshalSelfTResponseBodyToOrderviewsSelfTView(v *SelfTResponseBody) *orderviews.SelfTView {
	if v == nil {
		return nil
	}
	res := &orderviews.SelfTView{
		Self: v.Self,
	}
	if v.DescribedBy != nil {
		res.DescribedBy = unmarshalDescribedByTResponseBodyToOrderviewsDescribedByTView(v.DescribedBy)
	}

	return res
}

// unmarshalNavTResponseBodyToOrderviewsNavTView builds a value of type
// *orderviews.NavTView from a value of type *NavTResponseBody.
func unmarshalNavTResponseBodyToOrderviewsNavTView(v *NavTResponseBody) *orderviews.NavTView {
	if v == nil {
		return nil
	}
	res := &orderviews.NavTView{
		Self:  v.Self,
		First: v.First,
		Next:  v.Next,
	}

	return res
}

// unmarshalParameterTResponseBodyToOrderviewsParameterTView builds a value of
// type *orderviews.ParameterTView from a value of type *ParameterTResponseBody.
func unmarshalParameterTResponseBodyToOrderviewsParameterTView(v *ParameterTResponseBody) *orderviews.ParameterTView {
	res := &orderviews.ParameterTView{
		Name:  v.Name,
		Value: v.Value,
	}

	return res
}

// unmarshalOrderListItemResponseBodyToOrderviewsOrderListItemView builds a
// value of type *orderviews.OrderListItemView from a value of type
// *OrderListItemResponseBody.
func unmarshalOrderListItemResponseBodyToOrderviewsOrderListItemView(v *OrderListItemResponseBody) *orderviews.OrderListItemView {
	res := &orderviews.OrderListItemView{
		ID:         v.ID,
		Name:       v.Name,
		Status:     v.Status,
		OrderedAt:  v.OrderedAt,
		StartedAt:  v.StartedAt,
		FinishedAt: v.FinishedAt,
		ServiceID:  v.ServiceID,
		AccountID:  v.AccountID,
	}
	res.Links = unmarshalSelfTResponseBodyToOrderviewsSelfTView(v.Links)

	return res
}

// marshalOrderParameterTToParameterT builds a value of type *ParameterT from a
// value of type *order.ParameterT.
func marshalOrderParameterTToParameterT(v *order.ParameterT) *ParameterT {
	res := &ParameterT{
		Name:  v.Name,
		Value: v.Value,
	}

	return res
}

// marshalParameterTToOrderParameterT builds a value of type *order.ParameterT
// from a value of type *ParameterT.
func marshalParameterTToOrderParameterT(v *ParameterT) *order.ParameterT {
	res := &order.ParameterT{
		Name:  v.Name,
		Value: v.Value,
	}

	return res
}

// unmarshalOrderTopResultItemResponseToOrderviewsOrderTopResultItemView builds
// a value of type *orderviews.OrderTopResultItemView from a value of type
// *OrderTopResultItemResponse.
func unmarshalOrderTopResultItemResponseToOrderviewsOrderTopResultItemView(v *OrderTopResultItemResponse) *orderviews.OrderTopResultItemView {
	res := &orderviews.OrderTopResultItemView{
		Container:        v.Container,
		CPU:              v.CPU,
		Memory:           v.Memory,
		Storage:          v.Storage,
		EphemeralStorage: v.EphemeralStorage,
	}

	return res
}
