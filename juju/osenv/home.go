// Copyright 2013 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package osenv

import (
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/juju/utils/v4"
)

// jujuXDGDataHome stores the path to the juju configuration
// folder, which is only meaningful when running the juju
// CLI tool, and is typically defined by $JUJU_DATA or
// $XDG_DATA_HOME/juju or ~/.local/share/juju as default if none
// of the aforementioned variables are defined.
var (
	jujuXDGDataHomeMu sync.Mutex
	jujuXDGDataHome   string
)

const (
	// DirectorySubPathSSH is the sub directory under Juju home that holds ssh
	// related information for the Juju client.
	DirectorySubPathSSH = "ssh"
)

// SetJujuXDGDataHome sets the value of juju home and
// returns the current one.
func SetJujuXDGDataHome(newJujuXDGDataHomeHome string) string {
	jujuXDGDataHomeMu.Lock()
	defer jujuXDGDataHomeMu.Unlock()

	oldJujuXDGDataHomeHome := jujuXDGDataHome
	jujuXDGDataHome = newJujuXDGDataHomeHome
	return oldJujuXDGDataHomeHome
}

// JujuXDGDataHome returns the current juju home.
func JujuXDGDataHome() string {
	jujuXDGDataHomeMu.Lock()
	defer jujuXDGDataHomeMu.Unlock()
	return jujuXDGDataHome
}

// JujuXDGDDataHomeFS returns a file system rooted at home directory for the
// Juju data directory.
func JujuXDGDataHomeFS() fs.FS {
	return os.DirFS(JujuXDGDataHomeDir())
}

// JujuXDGDataSSHFS return a file system rooted at the ssh directory in the Juju
// data directory.
func JujuXDGDataSSHFS() fs.FS {
	return os.DirFS(filepath.Join(JujuXDGDataHomeDir(), DirectorySubPathSSH))
}

// JujuXDGDataHomePath returns the path to a file in the
// current juju home.
func JujuXDGDataHomePath(names ...string) string {
	all := append([]string{JujuXDGDataHomeDir()}, names...)
	return filepath.Join(all...)
}

// JujuXDGDataHomeDir returns the directory where juju should store application-specific files
func JujuXDGDataHomeDir() string {
	homeDir := JujuXDGDataHome()
	if homeDir != "" {
		return homeDir
	}
	homeDir = os.Getenv(JujuXDGDataHomeEnvKey)
	if homeDir == "" {
		if runtime.GOOS == "windows" {
			homeDir = jujuXDGDataHomeWin()
		} else {
			homeDir = jujuXDGDataHomeLinux()
		}
	}
	return homeDir
}

// jujuXDGDataHomeLinux returns the directory where juju should store application-specific files on Linux.
func jujuXDGDataHomeLinux() string {
	xdgConfig := os.Getenv(XDGDataHome)
	if xdgConfig != "" {
		return filepath.Join(xdgConfig, "juju")
	}
	// If xdg config home is not defined, the standard indicates that its default value
	// is $HOME/.local/share
	home := utils.Home()
	return filepath.Join(home, ".local", "share", "juju")
}

// jujuXDGDataHomeWin returns the directory where juju should store application-specific files on Windows.
func jujuXDGDataHomeWin() string {
	appdata := os.Getenv("APPDATA")
	if appdata == "" {
		return ""
	}
	return filepath.Join(appdata, "Juju")
}
