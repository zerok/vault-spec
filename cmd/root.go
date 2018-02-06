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
	"os"

	vault "github.com/hashicorp/vault/api"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var vaultClient *vault.Client
var specFile string
var logger *logrus.Logger
var verbose bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "vault-spec",
	Short: `vault-spec provides various tools around Vault configurations.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&specFile, "spec-file", "f", "", "Path to the spec file")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Verbose logging")
	logger = logrus.New()
	var err error
	cobra.OnInitialize(func() {
		if verbose {
			logger.SetLevel(logrus.DebugLevel)
		}
		vaultCfg := vault.DefaultConfig()
		if err = vaultCfg.ReadEnvironment(); err != nil {
			logger.WithError(err).Fatal("Failed to read Vault environment variables.")
		}
		vaultClient, err = vault.NewClient(vaultCfg)
		if err != nil {
			logger.WithError(err).Fatal("Failed to setup Vault client.")
		}
	})
}
