// Copyright 2018-2026 the original author or authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"io"
	"net/http"
	"regexp"

	"github.com/Masterminds/semver/v3"
	"github.com/paketo-buildpacks/libdependency/retrieve"
	"github.com/paketo-buildpacks/libdependency/upstream"
	"github.com/paketo-buildpacks/libdependency/versionology"
	"github.com/paketo-buildpacks/packit/v2/cargo"
)

const (
	id   = "skywalking-java"
	name = "Apache SkyWalking Java Agent"
	purl = "pkg:generic/apache-skywalking"

	archive = "https://archive.apache.org/dist/skywalking/java-agent"
)

var versionPattern = regexp.MustCompile(`href="(\d+\.\d+\.\d+)/"`)

func main() {
	retrieve.NewMetadata(id, getAllVersions, generateMetadata)
}

func getAllVersions() (versionology.VersionFetcherArray, error) {
	resp, err := http.Get(fmt.Sprintf("%s/", archive))
	if err != nil {
		return nil, fmt.Errorf("unable to fetch %s\n%w", archive, err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unable to fetch %s: status code %d", archive, resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to read %s\n%w", archive, err)
	}

	var versions versionology.VersionFetcherArray
	seen := map[string]bool{}
	for _, match := range versionPattern.FindAllStringSubmatch(string(body), -1) {
		versionString := match[1]
		if seen[versionString] {
			continue
		}
		seen[versionString] = true

		version, err := semver.NewVersion(versionString)
		if err != nil {
			fmt.Printf("Skipping %s: unable to parse version\n", versionString)
			continue
		}

		versions = append(versions, versionology.NewSimpleVersionFetcher(version))
	}

	return versions, nil
}

func generateMetadata(versionFetcher versionology.VersionFetcher) ([]versionology.Dependency, error) {
	versionString := versionFetcher.Version().String()

	uri := fmt.Sprintf("%s/%[2]s/apache-skywalking-java-agent-%[2]s.tgz", archive, versionString)
	checksum, err := upstream.GetSHA256OfRemoteFile(uri)
	if err != nil {
		return nil, fmt.Errorf("unable to checksum %s\n%w", uri, err)
	}

	source := fmt.Sprintf("%s/%[2]s/apache-skywalking-java-agent-%[2]s-src.tgz", archive, versionString)
	sourceChecksum, err := upstream.GetSHA256OfRemoteFile(source)
	if err != nil {
		return nil, fmt.Errorf("unable to checksum %s\n%w", source, err)
	}

	dependency := cargo.ConfigMetadataDependency{
		Checksum: fmt.Sprintf("sha256:%s", checksum),
		CPE:      fmt.Sprintf("cpe:2.3:a:apache:skywalking:%s:*:*:*:*:*:*:*", versionString),
		ID:       id,
		Licenses: []interface{}{
			map[string]string{
				"type": "Apache-2.0",
				"uri":  "https://www.apache.org/licenses/",
			},
		},
		Name:           name,
		PURL:           fmt.Sprintf("%s@%s", purl, versionString),
		SHA256:         checksum,
		Source:         source,
		SourceChecksum: fmt.Sprintf("sha256:%s", sourceChecksum),
		Stacks:         []string{"*"},
		URI:            uri,
		Version:        versionString,
	}

	return versionology.NewDependencyArray(dependency, "")
}
