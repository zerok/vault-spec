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
	"github.com/spf13/cobra"
	"github.com/zerok/vault-spec/spec"
)

// validateCmd represents the validate command
var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validates the given speciciation file",
	Run: func(cmd *cobra.Command, args []string) {
		if specFile == "" {
			logger.Fatal("No spec file specified. Use --spec-file.")
		}
		s, err := spec.FromPath(specFile)
		if err != nil {
			logger.WithError(err).Fatalf("Failed to load spec file from %s", specFile)
		}
		if err := s.Validate(); err != nil {
			logger.WithError(err).Fatal("Spec is not valid.")
		}
		logger.Info("Spec is valid.")
	},
}

func init() {
	rootCmd.AddCommand(validateCmd)
}
