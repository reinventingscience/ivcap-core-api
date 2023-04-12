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

package artifact

import (
	artifactviews "github.com/reinventingscience/ivcap-core-api/gen/artifact/views"
	"context"
	"io"

	"goa.design/goa/v3/security"
)

// Manage the life cycle of an artifact stored by this deployment.
type Service interface {
	// artifacts
	List(context.Context, *ListPayload) (res *ArtifactListRT, err error)
	// Upload content and create a artifacts.
	Upload(context.Context, *UploadPayload, io.ReadCloser) (res *ArtifactStatusRT, err error)
	// Show artifacts by ID
	Read(context.Context, *ReadPayload) (res *ArtifactStatusRT, err error)
	// Add artifacts to a collection.
	AddCollection(context.Context, *AddCollectionPayload) (err error)
	// Remove artifacts from a collection.
	RemoveCollection(context.Context, *RemoveCollectionPayload) (err error)
	// Add metadata of a partiular schema to artifacts.
	AddMetadata(context.Context, *AddMetadataPayload) (err error)
}

// Auther defines the authorization functions to be implemented by the service.
type Auther interface {
	// JWTAuth implements the authorization logic for the JWT security scheme.
	JWTAuth(ctx context.Context, token string, schema *security.JWTScheme) (context.Context, error)
}

// ServiceName is the name of the service as defined in the design. This is the
// same value that is set in the endpoint request contexts under the ServiceKey
// key.
const ServiceName = "artifact"

// MethodNames lists the service method names as defined in the design. These
// are the same values that are set in the endpoint request contexts under the
// MethodKey key.
var MethodNames = [6]string{"list", "upload", "read", "addCollection", "removeCollection", "addMetadata"}

// AddCollectionPayload is the payload type of the artifact service
// addCollection method.
type AddCollectionPayload struct {
	// Artifact ID
	ID string
	// Name of collection to add this artifacts.
	Name string
	// JWT used for authentication
	JWT string
}

// AddMetadataPayload is the payload type of the artifact service addMetadata
// method.
type AddMetadataPayload struct {
	// Artifact ID
	ID string
	// Schema of metadata. This will replace any existing metadata with this schema.
	Schema string
	// Metadata of a specific schema.
	Meta interface{}
	// JWT used for authentication
	JWT string
}

type ArtifactListItem struct {
	// Artifact ID
	ID *string
	// Optional name
	Name *string
	// Artifact status
	Status *string
	// Size of aritfact in bytes
	Size *int64
	// Mime (content) type of artifact
	MimeType *string
	Links    *SelfT
}

// ArtifactListRT is the result type of the artifact service list method.
type ArtifactListRT struct {
	// Artifacts
	Artifacts []*ArtifactListItem
	// Navigation links
	Links *NavT
}

// ArtifactStatusRT is the result type of the artifact service upload method.
type ArtifactStatusRT struct {
	// Artifact ID
	ID string
	// Optional name
	Name *string
	// List of collections this artifact is part of
	Collections []string
	// Link to retrieve the artifact data
	Data *SelfT
	// Artifact status
	Status string
	// Mime-type of data
	MimeType *string
	// Size of data
	Size *int64
	// List of metadata records associated with this artifact
	Metadata []*MetadataT
	// Reference to billable account
	Account *RefT
	Links   *SelfT
	// link back to record
	Location *string
	// indicate version of TUS supported
	TusResumable *string
	// TUS offset for partially uploaded content
	TusOffset *int64
}

// Bad arguments supplied.
type BadRequestT struct {
	// Information message
	Message string
}

type DescribedByT struct {
	Href *string
	Type *string
}

// Provided credential is not valid.
type InvalidCredentialsT struct {
}

// Caller not authorized to access required scope.
type InvalidScopesT struct {
	// ID of involved resource
	ID *string
	// Message of error
	Message string
}

// ListPayload is the payload type of the artifact service list method.
type ListPayload struct {
	// The $filter system query option allows clients to filter a collection of
	// resources that are addressed by a request URL. The expression specified with
	// $filter
	// is evaluated for each resource in the collection, and only items where the
	// expression
	// evaluates to true are included in the response.
	Filter string
	// The $orderby query option allows clients to request resources in either
	// ascending order using asc or descending order using desc. If asc or desc not
	// specified,
	// then the resources will be ordered in ascending order. The request below
	// orders Trips on
	// property EndsAt in descending order.
	Orderby string
	// The $top system query option requests the number of items in the queried
	// collection to be included in the result.
	Top int
	// The $skip query option requests the number of items in the queried collection
	// that are to be skipped and not included in the result.
	Skip int
	// The $select system query option allows the clients to requests a limited set
	// of properties for each entity or complex type. The example returns Name and
	// IcaoCode
	// of all Airports.
	Select string
	// DEPRECATED: List offset. Use '$skip' instead
	Offset *int
	// DEPRECATED: Max. number of records to return. Use '$top' instead
	Limit *int
	// DEPRECATED: Page token
	PageToken string
	// JWT used for authentication
	JWT string
}

type MetadataT struct {
	Schema *string
	Data   interface{}
}

type NavT struct {
	Self  *string
	First *string
	Next  *string
}

// Method is not yet implemented.
type NotImplementedT struct {
	// Information message
	Message string
}

// ReadPayload is the payload type of the artifact service read method.
type ReadPayload struct {
	// ID of artifacts to show
	ID string
	// JWT used for authentication
	JWT string
}

type RefT struct {
	ID    *string
	Links *SelfT
}

// RemoveCollectionPayload is the payload type of the artifact service
// removeCollection method.
type RemoveCollectionPayload struct {
	// Artifact ID
	ID string
	// Name of collection to remove this artifacts from.
	Name string
	// JWT used for authentication
	JWT string
}

// NotFound is the type returned when attempting to manage a resource that does
// not exist.
type ResourceNotFoundT struct {
	// ID of missing resource
	ID string
	// Message of error
	Message string
}

type SelfT struct {
	Self        *string
	DescribedBy *DescribedByT
}

// Unauthorized access to resource
type UnauthorizedT struct {
}

// UploadPayload is the payload type of the artifact service upload method.
type UploadPayload struct {
	// Content-Type header, MUST define type of uploaded content.
	ContentType *string
	// Content-Encoding header, MAY define encoding of content.
	ContentEncoding *string
	// Content-Length header, MAY define size of expected upload.
	ContentLength *int
	// X-Name header, MAY define a more human friendly name. Reusing a name will
	// NOT override an existing artifact with the same name
	Name *string
	// X-Collection header, MAY define an collection name as a simple way of
	// grouping artifacts
	Collection *string
	// X-Content-Type header, used for initial, empty content creation requests.
	XContentType *string
	// X-Content-Length header, used for initial, empty content creation requests.
	XContentLength *int
	// Upload-Length header, sets the expected content size part of the TUS
	// protocol.
	UploadLength *int
	// Tus-Resumable header, specifies TUS protocol version.
	TusResumable *string
	// JWT used for authentication
	JWT string
}

// Error returns an error description.
func (e *BadRequestT) Error() string {
	return "Bad arguments supplied."
}

// ErrorName returns "BadRequestT".
//
// Deprecated: Use GoaErrorName - https://github.com/goadesign/goa/issues/3105
func (e *BadRequestT) ErrorName() string {
	return e.GoaErrorName()
}

// GoaErrorName returns "BadRequestT".
func (e *BadRequestT) GoaErrorName() string {
	return "bad-request"
}

// Error returns an error description.
func (e *InvalidCredentialsT) Error() string {
	return "Provided credential is not valid."
}

// ErrorName returns "InvalidCredentialsT".
//
// Deprecated: Use GoaErrorName - https://github.com/goadesign/goa/issues/3105
func (e *InvalidCredentialsT) ErrorName() string {
	return e.GoaErrorName()
}

// GoaErrorName returns "InvalidCredentialsT".
func (e *InvalidCredentialsT) GoaErrorName() string {
	return "invalid-credential"
}

// Error returns an error description.
func (e *InvalidScopesT) Error() string {
	return "Caller not authorized to access required scope."
}

// ErrorName returns "InvalidScopesT".
//
// Deprecated: Use GoaErrorName - https://github.com/goadesign/goa/issues/3105
func (e *InvalidScopesT) ErrorName() string {
	return e.GoaErrorName()
}

// GoaErrorName returns "InvalidScopesT".
func (e *InvalidScopesT) GoaErrorName() string {
	return e.Message
}

// Error returns an error description.
func (e *NotImplementedT) Error() string {
	return "Method is not yet implemented."
}

// ErrorName returns "NotImplementedT".
//
// Deprecated: Use GoaErrorName - https://github.com/goadesign/goa/issues/3105
func (e *NotImplementedT) ErrorName() string {
	return e.GoaErrorName()
}

// GoaErrorName returns "NotImplementedT".
func (e *NotImplementedT) GoaErrorName() string {
	return "not-implemented"
}

// Error returns an error description.
func (e *ResourceNotFoundT) Error() string {
	return "NotFound is the type returned when attempting to manage a resource that does not exist."
}

// ErrorName returns "ResourceNotFoundT".
//
// Deprecated: Use GoaErrorName - https://github.com/goadesign/goa/issues/3105
func (e *ResourceNotFoundT) ErrorName() string {
	return e.GoaErrorName()
}

// GoaErrorName returns "ResourceNotFoundT".
func (e *ResourceNotFoundT) GoaErrorName() string {
	return "not-found"
}

// Error returns an error description.
func (e *UnauthorizedT) Error() string {
	return "Unauthorized access to resource"
}

// ErrorName returns "UnauthorizedT".
//
// Deprecated: Use GoaErrorName - https://github.com/goadesign/goa/issues/3105
func (e *UnauthorizedT) ErrorName() string {
	return e.GoaErrorName()
}

// GoaErrorName returns "UnauthorizedT".
func (e *UnauthorizedT) GoaErrorName() string {
	return "not-authorized"
}

// NewArtifactListRT initializes result type ArtifactListRT from viewed result
// type ArtifactListRT.
func NewArtifactListRT(vres *artifactviews.ArtifactListRT) *ArtifactListRT {
	return newArtifactListRT(vres.Projected)
}

// NewViewedArtifactListRT initializes viewed result type ArtifactListRT from
// result type ArtifactListRT using the given view.
func NewViewedArtifactListRT(res *ArtifactListRT, view string) *artifactviews.ArtifactListRT {
	p := newArtifactListRTView(res)
	return &artifactviews.ArtifactListRT{Projected: p, View: "default"}
}

// NewArtifactStatusRT initializes result type ArtifactStatusRT from viewed
// result type ArtifactStatusRT.
func NewArtifactStatusRT(vres *artifactviews.ArtifactStatusRT) *ArtifactStatusRT {
	return newArtifactStatusRT(vres.Projected)
}

// NewViewedArtifactStatusRT initializes viewed result type ArtifactStatusRT
// from result type ArtifactStatusRT using the given view.
func NewViewedArtifactStatusRT(res *ArtifactStatusRT, view string) *artifactviews.ArtifactStatusRT {
	p := newArtifactStatusRTView(res)
	return &artifactviews.ArtifactStatusRT{Projected: p, View: "default"}
}

// newArtifactListRT converts projected type ArtifactListRT to service type
// ArtifactListRT.
func newArtifactListRT(vres *artifactviews.ArtifactListRTView) *ArtifactListRT {
	res := &ArtifactListRT{}
	if vres.Artifacts != nil {
		res.Artifacts = make([]*ArtifactListItem, len(vres.Artifacts))
		for i, val := range vres.Artifacts {
			res.Artifacts[i] = transformArtifactviewsArtifactListItemViewToArtifactListItem(val)
		}
	}
	if vres.Links != nil {
		res.Links = transformArtifactviewsNavTViewToNavT(vres.Links)
	}
	return res
}

// newArtifactListRTView projects result type ArtifactListRT to projected type
// ArtifactListRTView using the "default" view.
func newArtifactListRTView(res *ArtifactListRT) *artifactviews.ArtifactListRTView {
	vres := &artifactviews.ArtifactListRTView{}
	if res.Artifacts != nil {
		vres.Artifacts = make([]*artifactviews.ArtifactListItemView, len(res.Artifacts))
		for i, val := range res.Artifacts {
			vres.Artifacts[i] = transformArtifactListItemToArtifactviewsArtifactListItemView(val)
		}
	}
	if res.Links != nil {
		vres.Links = transformNavTToArtifactviewsNavTView(res.Links)
	}
	return vres
}

// newArtifactStatusRT converts projected type ArtifactStatusRT to service type
// ArtifactStatusRT.
func newArtifactStatusRT(vres *artifactviews.ArtifactStatusRTView) *ArtifactStatusRT {
	res := &ArtifactStatusRT{
		Name:         vres.Name,
		MimeType:     vres.MimeType,
		Size:         vres.Size,
		Location:     vres.Location,
		TusResumable: vres.TusResumable,
		TusOffset:    vres.TusOffset,
	}
	if vres.ID != nil {
		res.ID = *vres.ID
	}
	if vres.Status != nil {
		res.Status = *vres.Status
	}
	if vres.Collections != nil {
		res.Collections = make([]string, len(vres.Collections))
		for i, val := range vres.Collections {
			res.Collections[i] = val
		}
	}
	if vres.Data != nil {
		res.Data = transformArtifactviewsSelfTViewToSelfT(vres.Data)
	}
	if vres.Metadata != nil {
		res.Metadata = make([]*MetadataT, len(vres.Metadata))
		for i, val := range vres.Metadata {
			res.Metadata[i] = transformArtifactviewsMetadataTViewToMetadataT(val)
		}
	}
	if vres.Account != nil {
		res.Account = transformArtifactviewsRefTViewToRefT(vres.Account)
	}
	if vres.Links != nil {
		res.Links = transformArtifactviewsSelfTViewToSelfT(vres.Links)
	}
	return res
}

// newArtifactStatusRTView projects result type ArtifactStatusRT to projected
// type ArtifactStatusRTView using the "default" view.
func newArtifactStatusRTView(res *ArtifactStatusRT) *artifactviews.ArtifactStatusRTView {
	vres := &artifactviews.ArtifactStatusRTView{
		ID:           &res.ID,
		Name:         res.Name,
		Status:       &res.Status,
		MimeType:     res.MimeType,
		Size:         res.Size,
		Location:     res.Location,
		TusResumable: res.TusResumable,
		TusOffset:    res.TusOffset,
	}
	if res.Collections != nil {
		vres.Collections = make([]string, len(res.Collections))
		for i, val := range res.Collections {
			vres.Collections[i] = val
		}
	}
	if res.Data != nil {
		vres.Data = transformSelfTToArtifactviewsSelfTView(res.Data)
	}
	if res.Metadata != nil {
		vres.Metadata = make([]*artifactviews.MetadataTView, len(res.Metadata))
		for i, val := range res.Metadata {
			vres.Metadata[i] = transformMetadataTToArtifactviewsMetadataTView(val)
		}
	}
	if res.Account != nil {
		vres.Account = transformRefTToArtifactviewsRefTView(res.Account)
	}
	if res.Links != nil {
		vres.Links = transformSelfTToArtifactviewsSelfTView(res.Links)
	}
	return vres
}

// transformArtifactviewsArtifactListItemViewToArtifactListItem builds a value
// of type *ArtifactListItem from a value of type
// *artifactviews.ArtifactListItemView.
func transformArtifactviewsArtifactListItemViewToArtifactListItem(v *artifactviews.ArtifactListItemView) *ArtifactListItem {
	if v == nil {
		return nil
	}
	res := &ArtifactListItem{
		ID:       v.ID,
		Name:     v.Name,
		Status:   v.Status,
		Size:     v.Size,
		MimeType: v.MimeType,
	}
	if v.Links != nil {
		res.Links = transformArtifactviewsSelfTViewToSelfT(v.Links)
	}

	return res
}

// transformArtifactviewsSelfTViewToSelfT builds a value of type *SelfT from a
// value of type *artifactviews.SelfTView.
func transformArtifactviewsSelfTViewToSelfT(v *artifactviews.SelfTView) *SelfT {
	res := &SelfT{
		Self: v.Self,
	}
	if v.DescribedBy != nil {
		res.DescribedBy = transformArtifactviewsDescribedByTViewToDescribedByT(v.DescribedBy)
	}

	return res
}

// transformArtifactviewsDescribedByTViewToDescribedByT builds a value of type
// *DescribedByT from a value of type *artifactviews.DescribedByTView.
func transformArtifactviewsDescribedByTViewToDescribedByT(v *artifactviews.DescribedByTView) *DescribedByT {
	if v == nil {
		return nil
	}
	res := &DescribedByT{
		Href: v.Href,
		Type: v.Type,
	}

	return res
}

// transformArtifactviewsNavTViewToNavT builds a value of type *NavT from a
// value of type *artifactviews.NavTView.
func transformArtifactviewsNavTViewToNavT(v *artifactviews.NavTView) *NavT {
	if v == nil {
		return nil
	}
	res := &NavT{
		Self:  v.Self,
		First: v.First,
		Next:  v.Next,
	}

	return res
}

// transformArtifactListItemToArtifactviewsArtifactListItemView builds a value
// of type *artifactviews.ArtifactListItemView from a value of type
// *ArtifactListItem.
func transformArtifactListItemToArtifactviewsArtifactListItemView(v *ArtifactListItem) *artifactviews.ArtifactListItemView {
	res := &artifactviews.ArtifactListItemView{
		ID:       v.ID,
		Name:     v.Name,
		Status:   v.Status,
		Size:     v.Size,
		MimeType: v.MimeType,
	}
	if v.Links != nil {
		res.Links = transformSelfTToArtifactviewsSelfTView(v.Links)
	}

	return res
}

// transformSelfTToArtifactviewsSelfTView builds a value of type
// *artifactviews.SelfTView from a value of type *SelfT.
func transformSelfTToArtifactviewsSelfTView(v *SelfT) *artifactviews.SelfTView {
	res := &artifactviews.SelfTView{
		Self: v.Self,
	}
	if v.DescribedBy != nil {
		res.DescribedBy = transformDescribedByTToArtifactviewsDescribedByTView(v.DescribedBy)
	}

	return res
}

// transformDescribedByTToArtifactviewsDescribedByTView builds a value of type
// *artifactviews.DescribedByTView from a value of type *DescribedByT.
func transformDescribedByTToArtifactviewsDescribedByTView(v *DescribedByT) *artifactviews.DescribedByTView {
	if v == nil {
		return nil
	}
	res := &artifactviews.DescribedByTView{
		Href: v.Href,
		Type: v.Type,
	}

	return res
}

// transformNavTToArtifactviewsNavTView builds a value of type
// *artifactviews.NavTView from a value of type *NavT.
func transformNavTToArtifactviewsNavTView(v *NavT) *artifactviews.NavTView {
	res := &artifactviews.NavTView{
		Self:  v.Self,
		First: v.First,
		Next:  v.Next,
	}

	return res
}

// transformArtifactviewsMetadataTViewToMetadataT builds a value of type
// *MetadataT from a value of type *artifactviews.MetadataTView.
func transformArtifactviewsMetadataTViewToMetadataT(v *artifactviews.MetadataTView) *MetadataT {
	if v == nil {
		return nil
	}
	res := &MetadataT{
		Schema: v.Schema,
		Data:   v.Data,
	}

	return res
}

// transformArtifactviewsRefTViewToRefT builds a value of type *RefT from a
// value of type *artifactviews.RefTView.
func transformArtifactviewsRefTViewToRefT(v *artifactviews.RefTView) *RefT {
	if v == nil {
		return nil
	}
	res := &RefT{
		ID: v.ID,
	}
	if v.Links != nil {
		res.Links = transformArtifactviewsSelfTViewToSelfT(v.Links)
	}

	return res
}

// transformMetadataTToArtifactviewsMetadataTView builds a value of type
// *artifactviews.MetadataTView from a value of type *MetadataT.
func transformMetadataTToArtifactviewsMetadataTView(v *MetadataT) *artifactviews.MetadataTView {
	if v == nil {
		return nil
	}
	res := &artifactviews.MetadataTView{
		Schema: v.Schema,
		Data:   v.Data,
	}

	return res
}

// transformRefTToArtifactviewsRefTView builds a value of type
// *artifactviews.RefTView from a value of type *RefT.
func transformRefTToArtifactviewsRefTView(v *RefT) *artifactviews.RefTView {
	if v == nil {
		return nil
	}
	res := &artifactviews.RefTView{
		ID: v.ID,
	}
	if v.Links != nil {
		res.Links = transformSelfTToArtifactviewsSelfTView(v.Links)
	}

	return res
}

// UploadRequestData holds both the payload and the HTTP request body reader of
// the "upload" method.
type UploadRequestData struct {
	// Payload is the method payload.
	Payload *UploadPayload
	// Body streams the HTTP request body.
	Body io.ReadCloser
}