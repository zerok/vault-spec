// Copyright © 2018 Horst Gutmann <zerok@zerokspot.com>
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
	"fmt"
	"strconv"

	"github.com/chzyer/readline"
	vault "github.com/hashicorp/vault/api"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/zerok/vault-spec/spec"
)

var prefix string

// UI is a generic interface for handling user interfaces for allowing users to
// enter data about a property.
type UI interface {
	SetLogger(*logrus.Entry)
	RequestInput(spec.SecretProperty) (interface{}, error)
}

type terminalUI struct {
	logger *logrus.Entry
}

func (ui *terminalUI) SetLogger(log *logrus.Entry) {
	ui.logger = log
}

func (ui *terminalUI) RequestInput(prop spec.SecretProperty) (interface{}, error) {
	var err error
	var data string
	var raw []byte
	var output interface{}
	log := ui.logger
	prompt := fmt.Sprintf("%s [%v]: ", prop.Label(), prop.Default())
	rl, err := readline.New(prompt)
	if err != nil {
		return nil, err
	}
	defer rl.Close()

inputLoop:
	for {
		if err != nil {
			log.WithError(err).Error("Invalid input")
		}
		switch prop.Input() {
		case "default":
			data, err = rl.Readline()
		case "hidden":
			raw, err = rl.ReadPassword(prompt)
			data = string(raw)
		}
		if err != nil {
			if err == readline.ErrInterrupt {
				return nil, err
			}
			continue inputLoop
		}

		if data == "" && prop.Default() != nil {
			output = prop.Default()
			break inputLoop
		}

		switch prop.Type() {
		case "string":
			output = data
			break inputLoop
		case "integer":
			output, err = strconv.ParseInt(data, 10, 64)
			if err != nil {
				continue inputLoop
			}
			break inputLoop
		}
	}
	return output, err
}

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
			finalPath := prefix + path
			plog := logger.WithField("path", finalPath)
			plog.Debug("Applying secret spec")
			secretSpec := spc.Secret(path)
			if secretSpec == nil {
				plog.Warn("No configuration found. Skipping.")
				continue
			}
			ui := terminalUI{}
			ui.SetLogger(plog)
			if err := applySecretSpec(vaultClient, finalPath, secretSpec, &ui, plog); err != nil {
				plog.WithError(err).Fatal("Failed to apply secret spec.")
			}
		}
	},
}

func applySecretSpec(vaultClient *vault.Client, path string, spc spec.Secret, ui UI, log *logrus.Entry) error {
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
		if err != nil || data[propName] == nil {
			proplog.WithError(err).Debug("Current data not valid. Requesting user input.")
			proplog.Info("Property has to be set/updated.")
			inp, err := ui.RequestInput(propSpec)
			if err != nil {
				proplog.WithError(err).Error("Failed to retrieve input.")
				return err
			}
			data[propName] = inp
			continue
		}
	}
	_, err = vaultClient.Logical().Write(path, data)
	return err
}

func init() {
	applyCmd.Flags().StringVarP(&prefix, "prefix", "p", "", "Prefix used for the secrets path")
	rootCmd.AddCommand(applyCmd)
}
