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

package aks

import (
	"strconv"
	"strings"

	pkgCommon "github.com/banzaicloud/pipeline/pkg/common"
	pkgErrors "github.com/banzaicloud/pipeline/pkg/errors"
	"github.com/pkg/errors"
)

// ### [ Constants to Azure cluster default values ] ### //
const (
	DefaultAgentName                      = "agentpool1"
	DefaultKubernetesVersion              = "1.9.2"
	MinKubernetesVersionWithAutoscalerStr = "1.9.6"
)

// CreateClusterAKS describes Azure fields of a CreateCluster request
type CreateClusterAKS struct {
	ResourceGroup     string                     `json:"resourceGroup" yaml:"resourceGroup"`
	KubernetesVersion string                     `json:"kubernetesVersion" yaml:"kubernetesVersion"`
	NodePools         map[string]*NodePoolCreate `json:"nodePools,omitempty" yaml:"nodePools,omitempty"`
}

// NodePoolCreate describes Azure's node fields of a CreateCluster request
type NodePoolCreate struct {
	Autoscaling      bool   `json:"autoscaling" yaml:"autoscaling"`
	MinCount         int    `json:"minCount" yaml:"minCount"`
	MaxCount         int    `json:"maxCount" yaml:"maxCount"`
	Count            int    `json:"count" yaml:"count"`
	NodeInstanceType string `json:"instanceType" yaml:"instanceType"`
}

// NodePoolUpdate describes Azure's node count of a UpdateCluster request
type NodePoolUpdate struct {
	Autoscaling bool `json:"autoscaling"`
	MinCount    int  `json:"minCount"`
	MaxCount    int  `json:"maxCount"`
	Count       int  `json:"count"`
}

// UpdateClusterAzure describes Azure's node fields of an UpdateCluster request
type UpdateClusterAzure struct {
	NodePools map[string]*NodePoolUpdate `json:"nodePools,omitempty"`
}

// Validate validates aks cluster create request
func (azure *CreateClusterAKS) Validate() error {

	if azure == nil {
		return pkgErrors.ErrorAzureFieldIsEmpty
	}

	// ---- [ NodePool check ] ---- //
	if azure.NodePools == nil {
		return pkgErrors.ErrorNodePoolEmpty
	}

	if len(azure.ResourceGroup) == 0 {
		return pkgErrors.ErrorResourceGroupRequired
	}

	for _, np := range azure.NodePools {

		// ---- [ Min & Max count fields are required in case of autoscaling ] ---- //
		if np.Autoscaling {
			err := checkVersionsIsNewerThen(azure.KubernetesVersion, MinKubernetesVersionWithAutoscalerStr)
			if err != nil {
				return err
			}
			if np.MinCount == 0 {
				return pkgErrors.ErrorMinFieldRequiredError
			}
			if np.MaxCount == 0 {
				return pkgErrors.ErrorMaxFieldRequiredError
			}
			if np.MaxCount < np.MinCount {
				return pkgErrors.ErrorNodePoolMinMaxFieldError
			}
		}

		if np.Count == 0 {
			np.Count = pkgCommon.DefaultNodeMinCount
		}

		if len(np.NodeInstanceType) == 0 {
			return pkgErrors.ErrorInstancetypeFieldIsEmpty
		}
	}

	if len(azure.KubernetesVersion) == 0 {
		azure.KubernetesVersion = DefaultKubernetesVersion
	}

	return nil
}

func parseVersion(version string) ([]int64, error) {
	iArray := make([]int64, 3)
	vArray := strings.Split(version, ".")
	for idx, n := range vArray {
		v, err := strconv.ParseInt(n, 10, 32)
		if err != nil {
			return nil, err
		}
		iArray[idx] = v
	}
	return iArray, nil
}

// return error if version is not at least minVersionStr
func checkVersionsIsNewerThen(version, minVersionStr string) error {
	minVersion, err := parseVersion(minVersionStr)
	if err != nil {
		return errors.Errorf("min version format is invalid: %s, example of correct format: '1.9.2'", minVersionStr)
	}
	parsedVersion, err := parseVersion(version)
	if err != nil {
		return errors.Errorf("kubernetes version format is invalid: %s, example of correct format: '1.9.2'", version)
	}
	for idx := range parsedVersion {
		if parsedVersion[idx] > minVersion[idx] {
			return nil
		} else if parsedVersion[idx] < minVersion[idx] {
			return errors.Errorf("autoscaler requires at least Kubernetes version: %s", minVersionStr)
		}
	}
	return nil
}

// Validate validates the update request (only aks part). If any of the fields is missing, the method fills
// with stored data.
func (a *UpdateClusterAzure) Validate() error {
	// ---- [ Azure field check ] ---- //
	if a == nil {
		return errors.New("'aks' field is empty") // todo move to errors
	}

	return nil
}

// ClusterProfileAKS describes an Azure profile
type ClusterProfileAKS struct {
	KubernetesVersion string                     `json:"kubernetesVersion"`
	NodePools         map[string]*NodePoolCreate `json:"nodePools,omitempty"`
}
