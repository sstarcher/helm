/*
Copyright 2016 The Kubernetes Authors All rights reserved.
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

package installer

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"k8s.io/helm/pkg/helm/helmpath"
)

// ErrMissingMetadata indicates that plugin.yaml is missing.
var ErrMissingMetadata = errors.New("plugin metadata (plugin.yaml) missing")

// Debug enables verbose output.
var Debug bool

// Installer provides an interface for installing helm client plugins.
type Installer interface {
	// Install adds a plugin to $HELM_HOME.
	Install() error
	// Path is the directory of the installed plugin.
	Path() string
}

// Install installs a plugin to $HELM_HOME.
func Install(i Installer) error {
	return i.Install()
}

// NewForSource determines the correct Installer for the given source.
func NewForSource(source, version string, home helmpath.Home) (Installer, error) {
	// Check if source is a local directory
	if isLocalReference(source) {
		return NewLocalInstaller(source, home)
	}
	return NewVCSInstaller(source, version, home)
}

// isLocalReference checks if the source exists on the filesystem.
func isLocalReference(source string) bool {
	_, err := os.Stat(source)
	return err == nil
}

// isPlugin checks if the directory contains a plugin.yaml file.
func isPlugin(dirname string) bool {
	_, err := os.Stat(filepath.Join(dirname, "plugin.yaml"))
	return err == nil
}

func debug(format string, args ...interface{}) {
	if Debug {
		format = fmt.Sprintf("[debug] %s\n", format)
		fmt.Printf(format, args...)
	}
}
