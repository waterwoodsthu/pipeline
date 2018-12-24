// Copyright © 2018 Banzai Cloud
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package eks

import "github.com/aws/aws-sdk-go/aws/endpoints"

// ### [ Constants to EKS cluster default values ] ### //
const (
	DefaultInstanceType = "m4.xlarge"
	DefaultSpotPrice    = "0.0" // 0 spot price stands for on-demand instances
	DefaultRegion       = endpoints.UsWest2RegionID
)

// DefaultImages in each supported location in EC2 (from https://docs.aws.amazon.com/eks/latest/userguide/launch-workers.html)
var DefaultImages = map[string]string{
	endpoints.UsEast1RegionID:      "ami-0b4eb1d8782fc3aea",
	endpoints.UsEast2RegionID:      "ami-053cbe66e0033ebcf",
	endpoints.UsWest2RegionID:      "ami-094fa4044a2a3cf52",
	endpoints.EuWest1RegionID:      "ami-0a9006fb385703b54",
	endpoints.EuNorth1RegionID:     "ami-082e6cf1c07e60241",
	endpoints.EuCentral1RegionID:   "ami-0ce0ec06e682ee10e",
	endpoints.ApNortheast1RegionID: "ami-063650732b3e8b38c",
	endpoints.ApSoutheast1RegionID: "ami-0549ac6995b998478",
	endpoints.ApSoutheast2RegionID: "ami-03297c04f71690a76",
}
