package main

import "testing"

var testFancyList = []string{
	"componentspace.saml2.dll",
	"DevExpress.*",
	"System.*.dll",
	"*.aspx",
	"README",
	"!LICense",
	"^!AKEF",
	"!^THOR",
}

func TestFancyList(t *testing.T) {
	if !isFileNameInFancyList("componentspace.saml2.dll", testFancyList) {
		nope(t, "Exact match file name did not match")
	}

	if !isFileNameInFancyList("COMPONENTSPACE.SAML2.DLL", testFancyList) {
		nope(t, "Upper case file name did not match")
	}

	if !isFileNameInFancyList(".abc.aspx", testFancyList) {
		nope(t, "*. wildcard not working")
	}

	if !isFileNameInFancyList("DevExpreSS.abc.DevExpreSS.def", testFancyList) {
		nope(t, ".* wildcard not working")
	}

	if !isFileNameInFancyList("System.fun.system.fun.dll", testFancyList) {
		nope(t, "*.* wildcard not working")
	}

	if !isFileNameInFancyList("LICense", testFancyList) {
		nope(t, "exact case match not working")
	}

	if isFileNameInFancyList("LICENCE", testFancyList) {
		nope(t, "exact case match not working")
	}

	if !isFileNameInFancyList("MAKEFILE", testFancyList) {
		nope(t, "exact case and contains match not working")
	}

	if !isFileNameInFancyList("AUTHORS", testFancyList) {
		nope(t, "exact AA and contains match not working")
	}

}

func nope(t *testing.T, message string) {
	t.Log(message)
	t.Fail()
}
