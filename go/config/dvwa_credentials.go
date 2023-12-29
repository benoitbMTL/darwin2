package config

// DVWA Credentials
var dvwaUserCredentials = map[string]string{
	"admin":                          "password",
	"gordonb":                        "abc123",
	"1337":                           "charley",
	"pablo":                          "letmein",
	"smithy":                         "password",
	"pklangdon4@msn.com":             "ppl11266",
	"muldersstan@gmail.com":          "renzo1205",
	"forsternp2@aol.com":             "freedom1",
	"cragsy@msn.com":                 "Snatch01",
	"bjrehdorf@hotmail.com":          "Apollo25504",
	"baz2709@icloud.com":             "sophie12",
	"amysiura@ymail.com":             "active95",
	"jond714@gmail.com":              "jonloveslax",
	"josefahorenstein87@hotmail.com": "qel737Xf3",
	"bizotic6@gmail.com":             "snaker26"}

// GetDVWAPassword returns the password for a given DVWA user
func GetDVWAPassword(username string) string {
	return dvwaUserCredentials[username]
}
