package ias3

import "testing"

func TestValidBucketKey(t *testing.T) {
	keys := map[string]bool{
		"hello123":  true,
		"hello/sd":  false,
		"hello-sd":  true,
		"hello_sd":  true,
		"hello\\df": false,
		"hello+df":  false,
		"hell0*sd":  false,
		"HELLO":     true,
	}

	for k, v := range keys {
		if validateBucketKey(k) != v {
			t.Logf("failed on %s expected %v", k, v)
			t.Fail()
		}
	}
}

func TestValidFileKey(t *testing.T) {
	keys := map[string]bool{
		"hello123":          true,
		"hello/sd":          true,
		"hello-sd":          true,
		"hello_sd":          true,
		"hello\\df":         false,
		"hello+df":          false,
		"hell0*sd":          false,
		"HELLO":             true,
		"heoo/dfsd/df_d-s":  true,
		"/heoo/dfsd/df_d-s": false,
	}

	for k, v := range keys {
		if validateFileKey(k) != v {
			t.Logf("failed on %s expected %v", k, v)
			t.Fail()
		}
	}
}
