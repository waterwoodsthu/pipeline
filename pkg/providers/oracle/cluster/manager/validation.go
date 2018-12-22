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

package manager

import (
	"github.com/banzaicloud/pipeline/pkg/providers/oracle/model"
	"github.com/pkg/errors"
)

// ValidateModel validates model configuration
func (cm *ClusterManager) ValidateModel(clusterModel *model.Cluster) error {

	m := clusterModel

	vn, err := cm.oci.NewVirtualNetworkClient()
	if err != nil {
		return err
	}

	ce, err := cm.oci.NewContainerEngineClient()
	if err != nil {
		return err
	}

	if m.Delete && m.OCID == "" {
		return errors.New("cannot delete cluster without Cluster OCID specified")
	}

	vcn, err := vn.GetVCN(&m.VCNID)
	if err != nil {
		return errors.Errorf("invalid VCN OCID: %s", m.VCNID)
	}

	subnet, err := vn.GetSubnet(&m.LBSubnetID1)
	if err != nil {
		return errors.Errorf("invalid LB 1 Subnet OCID: %s", m.LBSubnetID1)
	}
	if *subnet.VcnId != *vcn.Id {
		return errors.Errorf("invalid LB 1 Subnet OCID: %s not in VCN[%s]", m.LBSubnetID1, *vcn.Id)
	}

	subnet, err = vn.GetSubnet(&m.LBSubnetID2)
	if err != nil {
		return errors.Errorf("invalid LB 2 Subnet OCID: %s", m.LBSubnetID2)
	}
	if *subnet.VcnId != *vcn.Id {
		return errors.Errorf("invalid LB 2 Subnet OCID: %s not in VCN[%s]", m.LBSubnetID2, *vcn.Id)
	}

	k8sVersions, err := ce.GetAvailableKubernetesVersions()
	if err != nil {
		return err
	}

	if !k8sVersions.Has(m.Version) {
		return errors.Errorf("invalid k8s version: %s", m.Version)
	}

	if len(m.NodePools) < 1 {
		return errors.New("at least 1 node pool must be specified")
	}

	nodeOptions, err := ce.GetDefaultNodePoolOptions()
	if err != nil {
		return err
	}

	for _, np := range m.NodePools {
		if !nodeOptions.Images.Has(np.Image) {
			return errors.Errorf("invalid node image '%s' at '%s'", np.Image, np.Name)
		}
		if !nodeOptions.Shapes.Has(np.Shape) {
			return errors.Errorf("invalid shape '%s' at '%s'", np.Shape, np.Name)
		}
		if len(np.Subnets) < 1 {
			return errors.New("there must be at least 1 subnet specified")
		}

		if m.Version != np.Version {
			return errors.Errorf("nodePool[%s]: Different k8s versions were specified for master[%s] and nodes[%s]", np.Name, m.Version, np.Version)
		}

		for _, subnet := range np.Subnets {
			if _, err := vn.GetSubnet(&subnet.SubnetID); err != nil {
				return errors.Errorf("invalid Subnet OCID: %s", subnet.SubnetID)
			}
		}
	}

	return nil
}
