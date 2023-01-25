// Code generated by goa v3.10.2, DO NOT EDIT.
//
// metadata HTTP client encoders and decoders
//
// Command:
// $ goa gen cayp/api_gateway/design

package client

import (
	"bytes"
	metadata "cayp/api_gateway/gen/metadata"
	metadataviews "cayp/api_gateway/gen/metadata/views"
	"context"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	goahttp "goa.design/goa/v3/http"
)

// BuildReadRequest instantiates a HTTP request object with method and path set
// to call the "metadata" service "read" endpoint
func (c *Client) BuildReadRequest(ctx context.Context, v interface{}) (*http.Request, error) {
	var (
		entityID string
	)
	{
		p, ok := v.(*metadata.ReadPayload)
		if !ok {
			return nil, goahttp.ErrInvalidType("metadata", "read", "*metadata.ReadPayload", v)
		}
		entityID = p.EntityID
	}
	u := &url.URL{Scheme: c.scheme, Host: c.host, Path: ReadMetadataPath(entityID)}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, goahttp.ErrInvalidURL("metadata", "read", u.String(), err)
	}
	if ctx != nil {
		req = req.WithContext(ctx)
	}

	return req, nil
}

// EncodeReadRequest returns an encoder for requests sent to the metadata read
// server.
func EncodeReadRequest(encoder func(*http.Request) goahttp.Encoder) func(*http.Request, interface{}) error {
	return func(req *http.Request, v interface{}) error {
		p, ok := v.(*metadata.ReadPayload)
		if !ok {
			return goahttp.ErrInvalidType("metadata", "read", "*metadata.ReadPayload", v)
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
		if p.SchemaFilter != nil {
			values.Add("$schema-filter", *p.SchemaFilter)
		}
		if p.AtTime != nil {
			values.Add("$at-time", *p.AtTime)
		}
		req.URL.RawQuery = values.Encode()
		return nil
	}
}

// DecodeReadResponse returns a decoder for responses returned by the metadata
// read endpoint. restoreBody controls whether the response body should be
// restored after having been read.
// DecodeReadResponse may return the following errors:
//   - "bad-request" (type *metadata.BadRequestT): http.StatusBadRequest
//   - "invalid-credential" (type *metadata.InvalidCredentialsT): http.StatusBadRequest
//   - "invalid-parameter" (type *metadata.InvalidParameterValue): http.StatusUnprocessableEntity
//   - "invalid-scopes" (type *metadata.InvalidScopesT): http.StatusForbidden
//   - "not-implemented" (type *metadata.NotImplementedT): http.StatusNotImplemented
//   - "not-authorized" (type *metadata.UnauthorizedT): http.StatusUnauthorized
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
				return nil, goahttp.ErrDecodingError("metadata", "read", err)
			}
			p := NewReadMetaRTViewOK(&body)
			view := "default"
			vres := &metadataviews.ReadMetaRT{Projected: p, View: view}
			if err = metadataviews.ValidateReadMetaRT(vres); err != nil {
				return nil, goahttp.ErrValidationError("metadata", "read", err)
			}
			res := metadata.NewReadMetaRT(vres)
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
					return nil, goahttp.ErrDecodingError("metadata", "read", err)
				}
				err = ValidateReadBadRequestResponseBody(&body)
				if err != nil {
					return nil, goahttp.ErrValidationError("metadata", "read", err)
				}
				return nil, NewReadBadRequest(&body)
			case "invalid-credential":
				return nil, NewReadInvalidCredential()
			default:
				body, _ := io.ReadAll(resp.Body)
				return nil, goahttp.ErrInvalidResponse("metadata", "read", resp.StatusCode, string(body))
			}
		case http.StatusUnprocessableEntity:
			var (
				body ReadInvalidParameterResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("metadata", "read", err)
			}
			err = ValidateReadInvalidParameterResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("metadata", "read", err)
			}
			return nil, NewReadInvalidParameter(&body)
		case http.StatusForbidden:
			var (
				body ReadInvalidScopesResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("metadata", "read", err)
			}
			err = ValidateReadInvalidScopesResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("metadata", "read", err)
			}
			return nil, NewReadInvalidScopes(&body)
		case http.StatusNotImplemented:
			var (
				body ReadNotImplementedResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("metadata", "read", err)
			}
			err = ValidateReadNotImplementedResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("metadata", "read", err)
			}
			return nil, NewReadNotImplemented(&body)
		case http.StatusUnauthorized:
			return nil, NewReadNotAuthorized()
		default:
			body, _ := io.ReadAll(resp.Body)
			return nil, goahttp.ErrInvalidResponse("metadata", "read", resp.StatusCode, string(body))
		}
	}
}

// BuildAddRequest instantiates a HTTP request object with method and path set
// to call the "metadata" service "add" endpoint
func (c *Client) BuildAddRequest(ctx context.Context, v interface{}) (*http.Request, error) {
	var (
		entityID string
		schema   string
		body     io.Reader
	)
	{
		rd, ok := v.(*metadata.AddRequestData)
		if !ok {
			return nil, goahttp.ErrInvalidType("metadata", "add", "metadata.AddRequestData", v)
		}
		p := rd.Payload
		body = rd.Body
		entityID = p.EntityID
		schema = p.Schema
	}
	u := &url.URL{Scheme: c.scheme, Host: c.host, Path: AddMetadataPath(entityID, schema)}
	req, err := http.NewRequest("PUT", u.String(), body)
	if err != nil {
		return nil, goahttp.ErrInvalidURL("metadata", "add", u.String(), err)
	}
	if ctx != nil {
		req = req.WithContext(ctx)
	}

	return req, nil
}

// EncodeAddRequest returns an encoder for requests sent to the metadata add
// server.
func EncodeAddRequest(encoder func(*http.Request) goahttp.Encoder) func(*http.Request, interface{}) error {
	return func(req *http.Request, v interface{}) error {
		data, ok := v.(*metadata.AddRequestData)
		if !ok {
			return goahttp.ErrInvalidType("metadata", "add", "*metadata.AddRequestData", v)
		}
		p := data.Payload
		{
			head := p.JWT
			if !strings.Contains(head, " ") {
				req.Header.Set("Authorization", "Bearer "+head)
			} else {
				req.Header.Set("Authorization", head)
			}
		}
		if p.ContentType != nil {
			head := *p.ContentType
			req.Header.Set("Content-Type", head)
		}
		return nil
	}
}

// DecodeAddResponse returns a decoder for responses returned by the metadata
// add endpoint. restoreBody controls whether the response body should be
// restored after having been read.
// DecodeAddResponse may return the following errors:
//   - "bad-request" (type *metadata.BadRequestT): http.StatusBadRequest
//   - "invalid-credential" (type *metadata.InvalidCredentialsT): http.StatusBadRequest
//   - "invalid-parameter" (type *metadata.InvalidParameterValue): http.StatusUnprocessableEntity
//   - "invalid-scopes" (type *metadata.InvalidScopesT): http.StatusForbidden
//   - "not-implemented" (type *metadata.NotImplementedT): http.StatusNotImplemented
//   - "not-authorized" (type *metadata.UnauthorizedT): http.StatusUnauthorized
//   - error: internal error
func DecodeAddResponse(decoder func(*http.Response) goahttp.Decoder, restoreBody bool) func(*http.Response) (interface{}, error) {
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
				body AddResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("metadata", "add", err)
			}
			p := NewAddMetaRTViewOK(&body)
			view := "default"
			vres := &metadataviews.AddMetaRT{Projected: p, View: view}
			if err = metadataviews.ValidateAddMetaRT(vres); err != nil {
				return nil, goahttp.ErrValidationError("metadata", "add", err)
			}
			res := metadata.NewAddMetaRT(vres)
			return res, nil
		case http.StatusBadRequest:
			en := resp.Header.Get("goa-error")
			switch en {
			case "bad-request":
				var (
					body AddBadRequestResponseBody
					err  error
				)
				err = decoder(resp).Decode(&body)
				if err != nil {
					return nil, goahttp.ErrDecodingError("metadata", "add", err)
				}
				err = ValidateAddBadRequestResponseBody(&body)
				if err != nil {
					return nil, goahttp.ErrValidationError("metadata", "add", err)
				}
				return nil, NewAddBadRequest(&body)
			case "invalid-credential":
				return nil, NewAddInvalidCredential()
			default:
				body, _ := io.ReadAll(resp.Body)
				return nil, goahttp.ErrInvalidResponse("metadata", "add", resp.StatusCode, string(body))
			}
		case http.StatusUnprocessableEntity:
			var (
				body AddInvalidParameterResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("metadata", "add", err)
			}
			err = ValidateAddInvalidParameterResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("metadata", "add", err)
			}
			return nil, NewAddInvalidParameter(&body)
		case http.StatusForbidden:
			var (
				body AddInvalidScopesResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("metadata", "add", err)
			}
			err = ValidateAddInvalidScopesResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("metadata", "add", err)
			}
			return nil, NewAddInvalidScopes(&body)
		case http.StatusNotImplemented:
			var (
				body AddNotImplementedResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("metadata", "add", err)
			}
			err = ValidateAddNotImplementedResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("metadata", "add", err)
			}
			return nil, NewAddNotImplemented(&body)
		case http.StatusUnauthorized:
			return nil, NewAddNotAuthorized()
		default:
			body, _ := io.ReadAll(resp.Body)
			return nil, goahttp.ErrInvalidResponse("metadata", "add", resp.StatusCode, string(body))
		}
	}
}

// // BuildAddStreamPayload creates a streaming endpoint request payload from the
// method payload and the path to the file to be streamed
func BuildAddStreamPayload(payload interface{}, fpath string) (*metadata.AddRequestData, error) {
	f, err := os.Open(fpath)
	if err != nil {
		return nil, err
	}
	return &metadata.AddRequestData{
		Payload: payload.(*metadata.AddPayload),
		Body:    f,
	}, nil
}

// BuildRevokeRequest instantiates a HTTP request object with method and path
// set to call the "metadata" service "revoke" endpoint
func (c *Client) BuildRevokeRequest(ctx context.Context, v interface{}) (*http.Request, error) {
	var (
		entityID string
	)
	{
		p, ok := v.(*metadata.RevokePayload)
		if !ok {
			return nil, goahttp.ErrInvalidType("metadata", "revoke", "*metadata.RevokePayload", v)
		}
		if p.EntityID != nil {
			entityID = *p.EntityID
		}
	}
	u := &url.URL{Scheme: c.scheme, Host: c.host, Path: RevokeMetadataPath(entityID)}
	req, err := http.NewRequest("DELETE", u.String(), nil)
	if err != nil {
		return nil, goahttp.ErrInvalidURL("metadata", "revoke", u.String(), err)
	}
	if ctx != nil {
		req = req.WithContext(ctx)
	}

	return req, nil
}

// EncodeRevokeRequest returns an encoder for requests sent to the metadata
// revoke server.
func EncodeRevokeRequest(encoder func(*http.Request) goahttp.Encoder) func(*http.Request, interface{}) error {
	return func(req *http.Request, v interface{}) error {
		p, ok := v.(*metadata.RevokePayload)
		if !ok {
			return goahttp.ErrInvalidType("metadata", "revoke", "*metadata.RevokePayload", v)
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

// DecodeRevokeResponse returns a decoder for responses returned by the
// metadata revoke endpoint. restoreBody controls whether the response body
// should be restored after having been read.
// DecodeRevokeResponse may return the following errors:
//   - "bad-request" (type *metadata.BadRequestT): http.StatusBadRequest
//   - "invalid-credential" (type *metadata.InvalidCredentialsT): http.StatusBadRequest
//   - "invalid-parameter" (type *metadata.InvalidParameterValue): http.StatusUnprocessableEntity
//   - "invalid-scopes" (type *metadata.InvalidScopesT): http.StatusForbidden
//   - "not-implemented" (type *metadata.NotImplementedT): http.StatusNotImplemented
//   - "not-authorized" (type *metadata.UnauthorizedT): http.StatusUnauthorized
//   - error: internal error
func DecodeRevokeResponse(decoder func(*http.Response) goahttp.Decoder, restoreBody bool) func(*http.Response) (interface{}, error) {
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
		case http.StatusNoContent:
			return nil, nil
		case http.StatusBadRequest:
			en := resp.Header.Get("goa-error")
			switch en {
			case "bad-request":
				var (
					body RevokeBadRequestResponseBody
					err  error
				)
				err = decoder(resp).Decode(&body)
				if err != nil {
					return nil, goahttp.ErrDecodingError("metadata", "revoke", err)
				}
				err = ValidateRevokeBadRequestResponseBody(&body)
				if err != nil {
					return nil, goahttp.ErrValidationError("metadata", "revoke", err)
				}
				return nil, NewRevokeBadRequest(&body)
			case "invalid-credential":
				return nil, NewRevokeInvalidCredential()
			default:
				body, _ := io.ReadAll(resp.Body)
				return nil, goahttp.ErrInvalidResponse("metadata", "revoke", resp.StatusCode, string(body))
			}
		case http.StatusUnprocessableEntity:
			var (
				body RevokeInvalidParameterResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("metadata", "revoke", err)
			}
			err = ValidateRevokeInvalidParameterResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("metadata", "revoke", err)
			}
			return nil, NewRevokeInvalidParameter(&body)
		case http.StatusForbidden:
			var (
				body RevokeInvalidScopesResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("metadata", "revoke", err)
			}
			err = ValidateRevokeInvalidScopesResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("metadata", "revoke", err)
			}
			return nil, NewRevokeInvalidScopes(&body)
		case http.StatusNotImplemented:
			var (
				body RevokeNotImplementedResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("metadata", "revoke", err)
			}
			err = ValidateRevokeNotImplementedResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("metadata", "revoke", err)
			}
			return nil, NewRevokeNotImplemented(&body)
		case http.StatusUnauthorized:
			return nil, NewRevokeNotAuthorized()
		default:
			body, _ := io.ReadAll(resp.Body)
			return nil, goahttp.ErrInvalidResponse("metadata", "revoke", resp.StatusCode, string(body))
		}
	}
}

// unmarshalMetadataListItemResponseBodyToMetadataviewsMetadataListItemView
// builds a value of type *metadataviews.MetadataListItemView from a value of
// type *MetadataListItemResponseBody.
func unmarshalMetadataListItemResponseBodyToMetadataviewsMetadataListItemView(v *MetadataListItemResponseBody) *metadataviews.MetadataListItemView {
	res := &metadataviews.MetadataListItemView{
		RecordID: v.RecordID,
		Entity:   v.Entity,
		Schema:   v.Schema,
		Aspect:   v.Aspect,
	}

	return res
}

// unmarshalNavTResponseBodyToMetadataviewsNavTView builds a value of type
// *metadataviews.NavTView from a value of type *NavTResponseBody.
func unmarshalNavTResponseBodyToMetadataviewsNavTView(v *NavTResponseBody) *metadataviews.NavTView {
	res := &metadataviews.NavTView{
		Self:  v.Self,
		First: v.First,
		Next:  v.Next,
	}

	return res
}
