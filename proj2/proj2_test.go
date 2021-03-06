package proj2

// You MUST NOT change what you import.  If you add ANY additional
// imports it will break the autograder, and we will be Very Upset.

import (
	"testing"
	"reflect"
	"github.com/nweaver/cs161-p2/userlib"
	_ "encoding/json"
	_ "encoding/hex"
	_ "github.com/google/uuid"
	_ "strings"
	_ "errors"
	_ "strconv"
)


func TestInit(t *testing.T) {
	t.Log("Initialization test")

	// You may want to turn it off someday
	userlib.SetDebugStatus(true)
	someUsefulThings()
	userlib.SetDebugStatus(false)
	u, err := InitUser("alice", "fubar")
	if err != nil {
		// t.Error says the test fails
		t.Error("Failed to initialize user", err)
	}
	// t.Log() only produces output if you run with "go test -v"
	//t.Log("Got user", u)
	// If you want to comment the line above,
	// write _ = u here to make the compiler happy
	// You probably want many more tests here.
	_ = u
	if u.Username != "alice" {
		t.Error("Not alice")
	}
	_, ok := userlib.KeystoreGet("aliceenc")
	if !ok {
		t.Error("enckey didn't make it to keystore")
	}
	uget, errget := GetUser("alice", "fubar")
	if errget != nil {
		t.Error("Couldn't get alice")
	}
	if !reflect.DeepEqual(u, uget) {
		t.Error("Bad alice")
	}
}


func TestStorage(t *testing.T) {
	// And some more tests, because
	u, err := GetUser("alice", "fubar")
	if err != nil {
		t.Error("Failed to reload user", err)
		return
	}
	t.Log("Loaded user", u)

	v := []byte("This is a test")
	u.StoreFile("file1", v)


	v2, err2 := u.LoadFile("file1")
	if err2 != nil {
		t.Error("Failed to upload and download", err2)
	}
	if !reflect.DeepEqual(v, v2) {
		t.Error("Downloaded file is not the same", v, v2)
	}

	w := []byte(" but this isn't!")
	ww := append(v, w...)
	u.AppendFile("file1", w)

	w2, err3 := u.LoadFile("file1")
	if err3 != nil {
		t.Error("Oh no it didn't work!", err3)
	}
	if !reflect.DeepEqual(w2, ww) {
		t.Error("Not the same!", w2, ww)
	}

	y := []byte(" hehe think again")
	yy := append(ww, y...)
	u.AppendFile("file1", y)

	w3, err4 := u.LoadFile("file1")
	if err4 != nil {
		t.Error("Oh no it didn't work!", err4)
	}
	if !reflect.DeepEqual(w3, yy) {
		t.Error("Not the same!", w3, yy)
	}
}

func TestShare(t *testing.T) {
	u, err := GetUser("alice", "fubar")
	if err != nil {
		t.Error("Failed to reload user", err)
	}
	u2, err2 := InitUser("bob", "foobar")
	if err2 != nil {
		t.Error("Failed to initialize bob", err2)
	}

	var v, v2, v3, v4 []byte
	var magic_string string

	v, err = u.LoadFile("file1")
	if err != nil {
		t.Error("Failed to download the file from alice", err)
	}

	magic_string, err = u.ShareFile("file1", "bob")
	if err != nil {
		t.Error("Failed to share the a file", err)
	}
	err = u2.ReceiveFile("file2", "alice", magic_string)
	if err != nil {
		t.Error("Failed to receive the share message", err)
	}

	v2, err = u2.LoadFile("file2")
	if err != nil {
		t.Error("Failed to download the file after sharing 1", err)
	}
	if !reflect.DeepEqual(v, v2) {
		t.Error("Shared file is not the same", v, v2)
	}

	err = u.RevokeFile("file1")
	if err != nil {
		t.Error("Failed to revoke the a file", err)
	}
	v3, err = u2.LoadFile("file2")
	if err != nil {
		t.Error("Failed to download the file after sharing 2", err)
	}
	if reflect.DeepEqual(v2, v3) {
		t.Error("the file was the same after revoking...")
	}
	u.AppendFile("file1", []byte(" l;kjkl;j kl;j!"))
	v4, err = u2.LoadFile("file2")
	if err != nil {
		t.Error("Failed to download the file after sharing 3", err)
	}
	if !reflect.DeepEqual(v4, v3) {
		t.Error("revoked file changed", v4, v3)
	}

}
