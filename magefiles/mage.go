//go:build mage

package main

import (
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

var Aliases = map[string]any{
	"Lint": Lint.Run,
	"Test": Test.Run,
	"Docs": Docs.Start,
}

// Setup prepares the development environment.
func Setup() error {
	return sh.Run("git", "config", "core.hooksPath", ".githooks")
}

// ---- Lint ----
type Lint mg.Namespace

// Run runs the linter.
func (Lint) Run() error {
	return sh.RunV("golangci-lint", "run")
}

// Fix fixes style errors.
func (Lint) Fix() error {
	return sh.RunV("golangci-lint", "run", "--fix")
}

// ---- Docs ----
type Docs mg.Namespace

// Start stars doc server.
func (Docs) Start() error {
	return sh.RunV("pkgsite")
}

// --- Test ---
type Test mg.Namespace

// Run runs the tests.
func (Test) Run() error {
	return sh.RunV("go", "test", "./...")
}

// ---- Web ----
type Web mg.Namespace

// Build builds the app.
func (Web) Build() error {
	env := map[string]string{"GOOS": "js", "GOARCH": "wasm"}
	return sh.RunWithV(env, "go", "build", "-o", "docs/public/prism.wasm", "./cmd/prism")
}

// Start starts dev server.
func (Web) Start() error {
	return sh.RunV("npm", "run", "dev", "--prefix", "docs")
}

// Deploy builds the app and starts the dev server.
func (Web) Deploy() {
	mg.SerialDeps(Web.Build, Web.Start)
}

// ---- Eink ----
type Eink mg.Namespace

// Build builds the app.
func (Eink) Build() error {
	env := map[string]string{"GOOS": "linux", "GOARCH": "arm", "GOARM": "7", "CGO_ENABLED": "0"}
	return sh.RunWithV(env, "go", "build", "-o", "prism", "./cmd/prism")
}

// Ship copies the app on the tablet.
func (Eink) Ship() error {
	if err := sh.RunV("scp", "./prism", "rm2:/home/root/"); err != nil {
		return err
	}
	return sh.RunV("ssh", "rm2", "chmod", "+x", "/home/root/prism")
}

// Run runs the app on the tablet.
func (Eink) Run() error {
	mg.Deps(Server.Start)
	return sh.RunV("ssh", "-t", "rm2", "./prism")
}

// Deploy builds, ships and runs the app on the table.
func (Eink) Deploy() {
	mg.SerialDeps(Eink.Build, Eink.Ship, Eink.Run)
}

// ---- Storybook ----
type Sb mg.Namespace

// Build builds storybook.
func (Sb) Build() error {
	env := map[string]string{"GOOS": "js", "GOARCH": "wasm"}
	return sh.RunWithV(env, "go", "build", "-o", "docs/public/storybook.wasm", "./cmd/storybook")
}

// Start starts storybook.
func (Sb) Start() error {
	return sh.RunV("npm", "run", "storybook", "--prefix", "docs")
}

// ---- Server ----
type Server mg.Namespace

// Start starts the rm2fb server.
func (Server) Start() (err error) {
	if err := sh.RunV("ssh", "rm2", "systemctl", "is-active", "--quiet", "rm2fb"); err == nil {
		return nil
	}
	mg.Deps(Front.Stop)
	return sh.RunV("ssh", "rm2", "systemd-run", "--collect", "--unit=rm2fb", "/home/root/rm2fb_server_swtcon")
}

// Stop stops the rm2fb server.
func (Server) Stop() error {
	return sh.RunV("ssh", "rm2",
		"systemctl stop rm2fb 2>/dev/null;"+
			" killall rm2fb_server_swtcon 2>/dev/null; true")
}

// ---- Front ----
type Front mg.Namespace

// Start starts Xochitl (the default front app)
func (Front) Start() error {
	mg.Deps(Server.Stop)
	return sh.RunV("ssh", "rm2", "systemctl", "start", "xochitl")
}

// Stop stops Xochitl (the default front app)
func (Front) Stop() error {
	return sh.RunV("ssh", "rm2", "systemctl", "stop", "xochitl")
}
