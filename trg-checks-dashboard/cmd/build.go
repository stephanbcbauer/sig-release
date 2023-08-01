/*******************************************************************************
 * Copyright (c) 2023 Contributors to the Eclipse Foundation
 *
 * See the NOTICE file(s) distributed with this work for additional
 * information regarding copyright ownership.
 *
 * This program and the accompanying materials are made available under the
 * terms of the Apache License, Version 2.0 which is available at
 * https://www.apache.org/licenses/LICENSE-2.0.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
 * WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
 * License for the specific language governing permissions and limitations
 * under the License.
 *
 * SPDX-License-Identifier: Apache-2.0
 ******************************************************************************/

package cmd

import (
	"bufio"
	"bytes"
	"log"
	"math/rand"
	"os"
	"path"

	cp "github.com/otiai10/copy"
	"github.com/spf13/cobra"
	"trg-checks-dashboard/internal/templating"
)

const buildOutputDir = "build"
const outputFileName = "index.html"

// buildCmd represents the build command
var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Create a statically compiled dashboard with release check status",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		ensureOutputDirExists()
		copyAssets()

		var outputBuffer bytes.Buffer
		products, unhandledRepos := templating.CheckProducts()
		// unhandledRepos := getDemoUnhandledRepos()

		// templating.RenderHtmlTo(&outputBuffer, &templating.TemplateData{Config: getConfig(), CheckedProducts: getDemoChecks(), UnhandledRepos: unhandledRepos})
		templating.RenderHtmlTo(&outputBuffer, &templating.TemplateData{Config: getConfig(), CheckedProducts: products, UnhandledRepos: unhandledRepos})

		writeToFile(outputBuffer)
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// buildCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// buildCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func writeToFile(buffer bytes.Buffer) {
	f, err := os.Create(path.Join(buildOutputDir, outputFileName))
	if err != nil {
		log.Fatalf("Could not create output file: %v", err)
	}

	w := bufio.NewWriter(f)
	_, err = w.WriteString(string(buffer.Bytes()))
	if err != nil {
		log.Fatalf("Could not write to output file: %v", err)
	}

	if err = w.Flush(); err != nil {
		log.Fatalf("Could not flush output: %v", err)
	}
}

func ensureOutputDirExists() {
	if _, err := os.Stat(buildOutputDir); err != nil {
		if err := os.Mkdir(buildOutputDir, 0777); err != nil {
			log.Fatalf("Could not create build output dir: %v", err)
		}
	}
}

func copyAssets() {
	if err := cp.Copy("web/assets", path.Join(buildOutputDir, "assets")); err != nil {
		log.Fatalf("Could not copy CSS to build output dir: %v", err)
	}
}

func getConfig() templating.Config {
	return templating.Config{AssetsPath: "/assets"}
}

func getDemoUnhandledRepos() []templating.Repository {
	return []templating.Repository{
		{
			Name: "BPDM",
			URL:  "https://github.com/eclipse-tractusx/bpdm",
		},
	}
}

func getDemoChecks() []templating.CheckedProduct {
	return []templating.CheckedProduct{
		{
			Name:          "Portal",
			OverallPassed: randomBool(),
			LeadingRepo:   "https://github.com/eclipse-tractusx/portal-cd",
			CheckedRepositories: []templating.CheckedRepository{
				{
					RepoName:        "portal-cd",
					RepoUrl:         "https://github.com/eclipse-tractusx/portal-cd",
					GuidelineChecks: demoChecks(),
				},
				{
					RepoName:        "portal-frontend",
					RepoUrl:         "https://github.com/eclipse-tractusx/portal-frontend",
					GuidelineChecks: demoChecks(),
				},
				{
					RepoName:        "portal-backend",
					RepoUrl:         "https://github.com/eclipse-tractusx/portal-backend",
					GuidelineChecks: demoChecks(),
				},
			},
		},
		{
			Name:          "EDC",
			OverallPassed: randomBool(),
			LeadingRepo:   "https://github.com/eclipse-tractusx/tractusx-edc",
			CheckedRepositories: []templating.CheckedRepository{
				{
					RepoName:        "tractusx-edc",
					RepoUrl:         "https://github.com/eclipse-tractusx/tractusx-edc",
					GuidelineChecks: demoChecks(),
				},
			},
		},
	}
}

func demoChecks() []templating.GuidelineCheck {
	return []templating.GuidelineCheck{
		{
			GuidelineName: "TRG 1.01",
			GuidelineUrl:  "https://eclipse-tractusx.github.io/docs/release/trg-1/trg-1-1",
			Passed:        randomBool(),
		},
		{
			GuidelineName: "TRG 4.02",
			GuidelineUrl:  "https://eclipse-tractusx.github.io/docs/release/trg-4/trg-4-02",
			Passed:        randomBool(),
		},
		{
			GuidelineName: "TRG 5.02",
			GuidelineUrl:  "https://eclipse-tractusx.github.io/docs/release/trg-5/trg-5-02",
			Passed:        randomBool(),
		},
		{
			GuidelineName: "TRG 7.04",
			GuidelineUrl:  "https://eclipse-tractusx.github.io/docs/release/trg-7/trg-7-04",
			Passed:        randomBool(),
		},
	}
}

func randomBool() bool {
	return rand.Int31n(2) == 0
}
