// Copyright 2019 The Terraformer Authors.
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

package aws

import (
	"context"

	"github.com/GoogleCloudPlatform/terraformer/terraform_utils"
	"github.com/aws/aws-sdk-go-v2/service/elasticbeanstalk"
)

var beanstalkAllowEmptyValues = []string{"tags."}

type BeanstalkGenerator struct {
	AWSService
}

func (g *BeanstalkGenerator) InitResources() error {
	config, e := g.generateConfig()
	if e != nil {
		return e
	}
	client := elasticbeanstalk.New(config)

	err := g.addApplications(client)
	if err != nil {
		return err
	}
	err = g.addEnvironments(client)
	return err
}

func (g *BeanstalkGenerator) addApplications(client *elasticbeanstalk.Client) error {
	request := client.DescribeApplicationsRequest(&elasticbeanstalk.DescribeApplicationsInput{})
	response, err := request.Send(context.Background())
	if err != nil {
		return err
	}
	for _, application := range response.Applications {
		g.Resources = append(g.Resources, terraform_utils.NewSimpleResource(
			*application.ApplicationName,
			*application.ApplicationName,
			"aws_elastic_beanstalk_application",
			"aws",
			beanstalkAllowEmptyValues,
		))
	}
	return nil
}

func (g *BeanstalkGenerator) addEnvironments(client *elasticbeanstalk.Client) error {
	request := client.DescribeEnvironmentsRequest(&elasticbeanstalk.DescribeEnvironmentsInput{})
	response, err := request.Send(context.Background())
	if err != nil {
		return err
	}
	for _, environment := range response.Environments {
		g.Resources = append(g.Resources, terraform_utils.NewSimpleResource(
			*environment.EnvironmentId,
			*environment.EnvironmentName,
			"aws_elastic_beanstalk_environment",
			"aws",
			beanstalkAllowEmptyValues,
		))
	}
	return nil
}
