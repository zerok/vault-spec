// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
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

package cmd

import (
	vault "github.com/hashicorp/vault/api"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/zerok/vault-spec/spec"
)

// applyCmd represents the apply command
var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Applies the given specification to Vault",
	Run: func(cmd *cobra.Command, args []string) {
		spc, err := spec.FromPath(specFile)
		if err != nil {
			logger.WithError(err).Fatalf("Failed to load spec from %s.", specFile)
		}
		if err := spc.Validate(); err != nil {
			logger.WithError(err).Fatal("Specification is not valid.")
		}

		for _, path := range spc.SecretPaths() {
			plog := logger.WithField("path", path)
			plog.Info("Applying secret spec")

			secretSpec := spc.Secret(path)
			if secretSpec == nil {
				plog.Warn("No configuration found. Skipping.")
				continue
			}
			if err := applySecretSpec(vaultClient, path, secretSpec, plog); err != nil {
				plog.WithError(err).Fatal("Failed to apply secret spec.")
			}
		}
	},
}

func applySecretSpec(vaultClient *vault.Client, path string, spc spec.Secret, log *logrus.Entry) error {
	var err error
	var secret *vault.Secret
	secret, err = vaultClient.Logical().Read(path)
	if err != nil {
		return err
	}
	var data map[string]interface{}
	if secret != nil && secret.Data != nil {
		data = secret.Data
	} else {
		data = make(map[string]interface{})
	}
	for _, propName := range spc.PropertyNames() {
		propSpec := spc.Property(propName)
		proplog := log.WithField("prop", propName)
		if propSpec == nil {
			proplog.Warn("Property spec empty. Skipping.")
			continue
		}
		err := propSpec.IsValidData(data[propName])
		if err != nil {
			proplog.Debug("Current data not valid. Requesting user input.")
			continue
		}
	}
	return nil
}

func init() {
	rootCmd.AddCommand(applyCmd)
}
