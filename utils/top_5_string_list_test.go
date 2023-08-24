package utils

import "testing"

func TestTop5StringList(t *testing.T) {
	expectString(t, Top5StringList([]string{"lib1.dll"}), "\"lib1.dll\"")
	expectString(t, Top5StringList([]string{"lib1.dll", "lib1.dll", "lib1.dll"}), "\"lib1.dll\" (x3 instances)")
	expectString(t, Top5StringList([]string{"zlib1.dll", "lib1.dll", "lib1.dll", "lib1.dll"}), "\"lib1.dll\" (x3 instances), \"zlib1.dll\"")
	expectString(t, Top5StringList([]string{"zlib1.dll", "lib1.dll", "lib1.dll", "lib1.dll", "zlib1.dll"}), "\"lib1.dll\" (x3 instances), \"zlib1.dll\" (x2 instances)")
	expectString(t, Top5StringList([]string{"zlib1.dll", "lib1.dll", "lib1.dll", "lib1.dll", "zlib1.dll", "b.dll", "c.dll", "d.dll"}), "\"b.dll\", \"c.dll\", \"d.dll\", \"lib1.dll\" (x3 instances), \"zlib1.dll\" (x2 instances)")
	expectString(t, Top5StringList([]string{"zlib1.dll", "lib1.dll", "lib1.dll", "lib1.dll", "zlib1.dll", "b.dll", "c.dll", "d.dll", "zzlib.dll"}), "\"b.dll\", \"c.dll\", \"d.dll\", \"lib1.dll\" (x3 instances), \"zlib1.dll\" (x2 instances) and 1 other")
	expectString(t, Top5StringList([]string{"zlib1.dll", "lib1.dll", "lib1.dll", "lib1.dll", "zlib1.dll", "b.dll", "c.dll", "d.dll", "zzlib.dll", "zzlib.dll"}), "\"b.dll\", \"c.dll\", \"d.dll\", \"lib1.dll\" (x3 instances), \"zlib1.dll\" (x2 instances) and 2 others")
}
