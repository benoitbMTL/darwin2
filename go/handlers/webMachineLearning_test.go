package handlers

import (
	"net/url"
	"testing"

	"darwin2/utils"
)

func TestPreparePostDataPreservesCompoundAndAccentedNames(t *testing.T) {
	data := utils.FakeData{
		Name:      "Jean Benoît Lévesque-D'Amours",
		FirstName: "Jean Benoît",
		LastName:  "Lévesque-D'Amours",
		EmailU:    "jean.benoit+demo42",
		EmailD:    "team-24.node7.test",
		PhoneH:    "+1 (514) 555-0199",
		Address:   "12 1/2 boulevard René-Lévesque",
		Birthday:  "1980-05-14",
		Username:  "jean-benoit_42",
		Password:  "Maple42!River",
	}

	values, err := url.ParseQuery(preparePostData(data))
	if err != nil {
		t.Fatal(err)
	}
	assertFormValue(t, values, "firstname", data.FirstName)
	assertFormValue(t, values, "lastname", data.LastName)
	assertFormValue(t, values, "email", "jean.benoit+demo42@team-24.node7.test")
	assertFormValue(t, values, "phone", data.PhoneH)
	assertFormValue(t, values, "address", data.Address)
	assertFormValue(t, values, "password", data.Password)
}

func assertFormValue(t *testing.T, values url.Values, key, want string) {
	t.Helper()
	if got := values.Get(key); got != want {
		t.Fatalf("%s = %q, want %q", key, got, want)
	}
}
