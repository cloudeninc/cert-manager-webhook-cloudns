/*
 Copyright 2019 IXON B.V.

 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

      http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

package main

import (
	"github.com/ixoncloud/cert-manager-webhook-cloudns/cloudns"
	"github.com/jetstack/cert-manager/pkg/acme/webhook/apis/acme/v1alpha1"
	"github.com/jetstack/cert-manager/pkg/acme/webhook/cmd"
	restclient "k8s.io/client-go/rest"
	"os"
	"log"
)

const ProviderName = "cloudns"

var GroupName = os.Getenv("GROUP_NAME")

func main() {
	if GroupName == "" {
		panic("Please set the GROUP_NAME env variable.")
	}

	log.Printf("Starting group %s webhook", GroupName)

	// Start webhook server
	cmd.RunWebhookServer(GroupName,
		&clouDNSProviderSolver{},
	)
}

// clouDNSProviderSolver implements webhook.Solver
// and will allow cert-manager to create & delete
// DNS TXT records for the DNS01 Challenge
type clouDNSProviderSolver struct {
}

func (c clouDNSProviderSolver) Name() string {
	return ProviderName
}

// Create TXT DNS record for DNS01
func (c clouDNSProviderSolver) Present(ch *v1alpha1.ChallengeRequest) error {
	// Load environment variables and create new ClouDNS provider
	provider, err := cloudns.NewDNSProvider()

	if err != nil {
		return err
	}

	log.Printf("Presenting DNS01 for %s", ch.ResolvedFQDN)

	return provider.Present(ch.ResolvedFQDN, ch.Key)
}

// Delete TXT DNS record for DNS01
func (c clouDNSProviderSolver) CleanUp(ch *v1alpha1.ChallengeRequest) error {
	// Load environment variables and create new ClouDNS provider
	provider, err := cloudns.NewDNSProvider()

	if err != nil {
		return err
	}

	log.Printf("Cleaning up DNS01 for %s", ch.ResolvedFQDN)

	// Remove TXT DNS record
	return provider.CleanUp(ch.ResolvedFQDN, ch.Key)
}

// Could be used to initialise connections or warm up caches, not needed in this case
func (c clouDNSProviderSolver) Initialize(kubeClientConfig *restclient.Config, stopCh <-chan struct{}) error {
	// NOOP
	log.Printf("Initializing ClouDNS DNS01 solver...")
	return nil
}
