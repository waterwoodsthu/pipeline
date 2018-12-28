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

package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/banzaicloud/pipeline/config"
	"github.com/spf13/viper"
)

type CloudInfoResponse struct {

	// Products represents a slice of products for a given provider (VMs with attributes and process)
	Products []*MachineDetails `json:"products"`

	// ScrapingTime represents scraping time for a given provider in milliseconds
	ScrapingTime string `json:"scrapingTime,omitempty"`
}

type MachineDetails struct {

	// cpus
	Cpus float64 `json:"cpusPerVm,omitempty"`

	// gpus
	Gpus float64 `json:"gpusPerVm,omitempty"`

	// mem
	Mem float64 `json:"memPerVm,omitempty"`

	// ntw perf
	NtwPerf string `json:"ntwPerf,omitempty"`

	// ntw perf cat
	NtwPerfCat string `json:"ntwPerfCategory,omitempty"`

	// type
	Type string `json:"type,omitempty"`
}

type VMKey struct {
	cloud        string
	service      string
	region       string
	instanceType string
}

var instanceTypeMap = make(map[VMKey]MachineDetails)

func fetchMachineTypes(cloud string, service string, region string) {
	cloudInfoEndPoint := viper.GetString(config.CloudInfoEndPoint)
	if len(cloudInfoEndPoint) == 0 {
		log.Errorf("Missing config: %v", config.CloudInfoEndPoint)
		return
	}
	cloudInfoUrl := fmt.Sprintf(
		"%s/providers/%s/services/%s/regions/%s/products",
		cloudInfoEndPoint, cloud, service, region)
	ciRequest, err := http.NewRequest(http.MethodGet, cloudInfoUrl, nil)
	if err != nil {
		log.Errorf("Error fetching machine types from CloudInfo: %#v", err)
		return
	}

	ciRequest.Header.Set("Content-Type", "application/json")
	httpClient := &http.Client{}

	ciResponse, err := httpClient.Do(ciRequest)
	if err != nil {
		log.Errorf("Error fetching machine types from CloudInfo: %#v", err)
		return
	}
	respBody, _ := ioutil.ReadAll(ciResponse.Body)
	var vmDetails CloudInfoResponse
	json.Unmarshal(respBody, &vmDetails)

	for _, product := range vmDetails.Products {
		instanceTypeMap[VMKey{
			cloud,
			service,
			region,
			product.Type,
		}] = *product
	}
}

//GetMachineDetails returns machine resource details, like cpu/gpu/memory etc. either from local cache or CloudInfo
func GetMachineDetails(cloud string, service string, region string, instanceType string) *MachineDetails {

	vmKey := VMKey{
		cloud,
		service,
		region,
		instanceType,
	}

	vmDetails, ok := instanceTypeMap[vmKey]
	if !ok {
		fetchMachineTypes(cloud, service, region)
		vmDetails = instanceTypeMap[vmKey]
	}

	return &vmDetails
}
