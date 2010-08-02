/*
Copyright (c) 2009-2010 Jim Teeuwen.

This software is provided 'as-is', without any express or implied
warranty. In no event will the authors be held liable for any damages
arising from the use of this software.

Permission is granted to anyone to use this software for any purpose,
including commercial applications, and to alter it and redistribute it
freely, subject to the following restrictions:

    1. The origin of this software must not be misrepresented; you must not
    claim that you wrote the original software. If you use this software
    in a product, an acknowledgment in the product documentation would be
    appreciated but is not required.

    2. Altered source versions must be plainly marked as such, and must not be
    misrepresented as being the original software.

    3. This notice may not be removed or altered from any source distribution.

*/
package optarg

import "os"
import "testing"

func Test(t *testing.T) {
	os.Args = []string{ // manually rebuild os.Args for testing purposes.
		os.Args[0],
		"--bin", "/a/b/foo/bin",
		"--arch", os.Getenv("GOARCH"),
		"-nps", "/a/b/foo/src",
		"foo.go", "bar.go",
	}

	// Add some flags
	Add("s", "source", "Path to the source folder. Here is some added description information which is completely useless, but it makes sure we can pimp our sexy Usage() output when dealing with lenghty, multi-line description texts.", "")
	Add("b", "bin", "Path to the binary folder.", "")
	Add("a", "arch", "Target architecture.", os.Getenv("GOARCH"))
	Add("n", "noproc", "Skip pre/post processing.", false)
	Add("p", "purge", "Clean compiled packages after linking is complete.", false)

	// These will hold the option values.
	var src, bin, arch string
	var noproc, purge bool

	// Parse os.Args
	for opt := range Parse() {
		switch opt.ShortName {
		case "s":
			src = opt.String()
		case "b":
			bin = opt.String()
		case "a":
			arch = opt.String()
		case "p":
			purge = opt.Bool()
		case "n":
			noproc = opt.Bool()
		}
	}

	// Make sure everything went ok.

	if arch != os.Getenv("GOARCH") {
		t.Errorf("Parse(): incorrect value for arch: %s", arch)
	}

	if bin != "/a/b/foo/bin" {
		t.Errorf("Parse(): incorrect value for bin: %s", bin)
	}

	if src != "/a/b/foo/src" {
		t.Errorf("Parse(): incorrect value for src: %s", src)
	}

	if !purge {
		t.Errorf("Parse(): purge is not set")
	}

	if !noproc {
		t.Errorf("Parse(): noproc is not set")
	}

	if len(Remainder) != 2 { // should contain: foo.go, bar.go
		t.Errorf("Parse(): incorrect number of remaining arguments. Expected 2. got %d", len(Remainder))
	}

	// This outputs the usage information. No need to do this in a test case.
	//Usage();
}
