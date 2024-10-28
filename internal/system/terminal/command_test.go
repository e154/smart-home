// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2024, Filippov Alex
//
// This library is free software: you can redistribute it and/or
// modify it under the terms of the GNU Lesser General Public
// License as published by the Free Software Foundation; either
// version 3 of the License, or (at your option) any later version.
//
// This library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
// Library General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public
// License along with this library.  If not, see
// <https://www.gnu.org/licenses/>.

package terminal

import (
	"bytes"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func executeCommandC(root *cobra.Command, args ...string) (c *cobra.Command, output string, err error) {
	buf := new(bytes.Buffer)
	root.SetOut(buf)
	root.SetErr(buf)
	root.SetArgs(args)

	c, err = root.ExecuteC()

	return c, buf.String(), err
}

func executeCommand(root *cobra.Command, args ...string) (output string, err error) {
	_, output, err = executeCommandC(root, args...)
	return output, err
}

const onetwo = "one two"

func TestSingleCommand(t *testing.T) {
	var rootCmdArgs []string
	rootCmd := &cobra.Command{
		Use:  "root",
		Args: cobra.ExactArgs(2),
		Run:  func(_ *cobra.Command, args []string) { rootCmdArgs = args },
	}
	aCmd := &cobra.Command{Use: "a", Args: cobra.NoArgs, Run: emptyRun}
	bCmd := &cobra.Command{Use: "b", Args: cobra.NoArgs, Run: emptyRun}
	rootCmd.AddCommand(aCmd, bCmd)

	output, err := executeCommand(rootCmd, "one", "two")
	if output != "" {
		t.Errorf("Unexpected output: %v", output)
	}
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	got := strings.Join(rootCmdArgs, " ")
	if got != onetwo {
		t.Errorf("rootCmdArgs expected: %q, got: %q", onetwo, got)
	}
}
