// -*- Mode: Go; indent-tabs-mode: t -*-

/*
 * Copyright (C) 2015 Canonical Ltd
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License version 3 as
 * published by the Free Software Foundation.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 */

package tests

import (
	"os"
	"os/exec"

	"launchpad.net/snappy/_integration-tests/testutils/build"
	"launchpad.net/snappy/_integration-tests/testutils/cli"
	"launchpad.net/snappy/_integration-tests/testutils/common"
	"launchpad.net/snappy/_integration-tests/testutils/data"

	"gopkg.in/check.v1"
)

var _ = check.Suite(&installAppSuite{})

type installAppSuite struct {
	common.SnappySuite
}

func (s *installAppSuite) TestInstallAppMustPrintPackageInformation(c *check.C) {
	installOutput := common.InstallSnap(c, "hello-world")
	s.AddCleanup(func() {
		common.RemoveSnap(c, "hello-world")
	})

	expected := "(?ms)" +
		"Installing hello-world\n" +
		"Name +Date +Version +Developer \n" +
		".*" +
		"^hello-world +.* +.* +canonical \n" +
		".*"

	c.Assert(installOutput, check.Matches, expected)
}

func (s *installAppSuite) TestCallSuccessfulBinaryFromInstalledSnap(c *check.C) {
	snapPath, err := build.LocalSnap(c, data.BasicWithBinariesSnapName)
	defer os.Remove(snapPath)
	c.Assert(err, check.IsNil)
	common.InstallSnap(c, snapPath)
	defer common.RemoveSnap(c, data.BasicWithBinariesSnapName)

	// Exec command does not fail.
	cli.ExecCommand(c, "basic-with-binaries.success")
}

func (s *installAppSuite) TestCallFailBinaryFromInstalledSnap(c *check.C) {
	snapPath, err := build.LocalSnap(c, data.BasicWithBinariesSnapName)
	defer os.Remove(snapPath)
	c.Assert(err, check.IsNil)
	common.InstallSnap(c, snapPath)
	defer common.RemoveSnap(c, data.BasicWithBinariesSnapName)

	_, err = cli.ExecCommandErr("basic-with-binaries.fail")
	c.Assert(err, check.NotNil, check.Commentf("The binary did not fail"))
}

func (s *installAppSuite) TestInstallUnexistingAppMustPrintError(c *check.C) {
	cmd := exec.Command("sudo", "snappy", "install", "unexisting.canonical")
	output, err := cmd.CombinedOutput()

	c.Assert(err, check.NotNil)
	c.Assert(string(output), check.Equals,
		"Installing unexisting.canonical\n"+
			"unexisting failed to install: snappy package not found\n")
}
