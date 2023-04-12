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
	"fmt"
)

// ListArtifactPath returns the URL path to the artifact service list HTTP endpoint.
func ListArtifactPath() string {
	return "/1/artifacts"
}

// UploadArtifactPath returns the URL path to the artifact service upload HTTP endpoint.
func UploadArtifactPath() string {
	return "/1/artifacts"
}

// ReadArtifactPath returns the URL path to the artifact service read HTTP endpoint.
func ReadArtifactPath(id string) string {
	return fmt.Sprintf("/1/artifacts/%v", id)
}

// AddCollectionArtifactPath returns the URL path to the artifact service addCollection HTTP endpoint.
func AddCollectionArtifactPath(id string, name string) string {
	return fmt.Sprintf("/1/artifacts/%v/.collections/%v", id, name)
}

// RemoveCollectionArtifactPath returns the URL path to the artifact service removeCollection HTTP endpoint.
func RemoveCollectionArtifactPath(id string, name string) string {
	return fmt.Sprintf("/1/artifacts/%v/.collections/%v", id, name)
}

// AddMetadataArtifactPath returns the URL path to the artifact service addMetadata HTTP endpoint.
func AddMetadataArtifactPath(id string, schema string) string {
	return fmt.Sprintf("/1/artifacts/%v/.metadata/%v", id, schema)
}