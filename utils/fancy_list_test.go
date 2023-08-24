package utils

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
	if !IsFileNameInFancyList("componentspace.saml2.dll", testFancyList) {
		nope(t, "Exact match file name did not match")
	}

	if !IsFileNameInFancyList("COMPONENTSPACE.SAML2.DLL", testFancyList) {
		nope(t, "Upper case file name did not match")
	}

	if !IsFileNameInFancyList(".abc.aspx", testFancyList) {
		nope(t, "*. wildcard not working")
	}

	if !IsFileNameInFancyList("DevExpreSS.abc.DevExpreSS.def", testFancyList) {
		nope(t, ".* wildcard not working")
	}

	if !IsFileNameInFancyList("System.fun.system.fun.dll", testFancyList) {
		nope(t, "*.* wildcard not working")
	}

	if !IsFileNameInFancyList("LICense", testFancyList) {
		nope(t, "exact case match not working")
	}

	if IsFileNameInFancyList("LICENCE", testFancyList) {
		nope(t, "exact case match not working")
	}

	if !IsFileNameInFancyList("MAKEFILE", testFancyList) {
		nope(t, "exact case and contains match not working")
	}

	if !IsFileNameInFancyList("AUTHORS", testFancyList) {
		nope(t, "exact AA and contains match not working")
	}
}
