package command

import (
	"errors"
	"testing"

	"github.com/brozzis/s3fs/handler/command/context"
)

func TestCdCommand_Execute(t *testing.T) {
	{
		var s3 mockS3Client
		var con context.Context
		var args []string
		var cd CdCommand
		var out mockOutputter

		// No args
		cd = NewCd(&s3, &con, args)
		if err := cd.Execute(&out); err != nil {
			t.Fatal(err)
		}

		// Root
		args = []string{context.PathDelimiter}
		cd = NewCd(&s3, &con, args)
		if err := cd.Execute(&out); err != nil {
			t.Fatal(err)
		}
	}

	// Valid bucket
	{
		var s3 mockS3Client
		var con context.Context
		var out mockOutputter
		args := []string{"bucket"}

		s3.bucketExistsCallback = func(b string) (bool, error) {
			if b != args[0] {
				t.Fatalf("Checking for unexpected bucket: %v", b)
			}

			return true, nil
		}

		cd := NewCd(&s3, &con, args)
		if err := cd.Execute(&out); err != nil {
			t.Fatal(err)
		}
	}

	// Invalid bucket
	{
		var s3 mockS3Client
		var con context.Context
		var out mockOutputter
		args := []string{"invalidbucket"}

		s3.bucketExistsCallback = func(b string) (bool, error) {
			if b != args[0] {
				t.Fatalf("Checking for unexpected bucket: %v", b)
			}

			return false, nil
		}

		cd := NewCd(&s3, &con, args)
		if err := cd.Execute(&out); err == nil {
			t.Fatalf("Expected error to be returned for invalid bucket")
		}
	}

	// Error checking bucket
	{
		var s3 mockS3Client
		var con context.Context
		var out mockOutputter
		args := []string{"invalidbucket"}
		fakeErr := errors.New("Error checking bucket")

		s3.bucketExistsCallback = func(b string) (bool, error) {
			if b != args[0] {
				t.Fatalf("Checking for unexpected bucket: %v", b)
			}

			return false, fakeErr
		}

		cd := NewCd(&s3, &con, args)
		if err := cd.Execute(&out); err != fakeErr {
			t.Fatalf("Unexpected error returned for bucket cd")
		}
	}

	// Valid folder
	{
		var s3 mockS3Client
		var con context.Context
		var out mockOutputter
		args := []string{"bucket/folder/"}

		s3.pathExistsCallback = func(b, k string) (bool, error) {
			if b != "bucket" || k != "folder/" {
				t.Fatalf("Checking for unexpected bucket/folder: %v/%v", b, k)
			}

			return true, nil
		}

		cd := NewCd(&s3, &con, args)
		if err := cd.Execute(&out); err != nil {
			t.Fatal(err)
		}
	}

	// Invalid folder
	{
		var s3 mockS3Client
		var con context.Context
		var out mockOutputter
		args := []string{"bucket/invalidfolder/"}

		s3.pathExistsCallback = func(b, k string) (bool, error) {
			if b != "bucket" || k != "invalidfolder/" {
				t.Fatalf("Checking for unexpected bucket/folder: %v/%v", b, k)
			}

			return false, nil
		}

		cd := NewCd(&s3, &con, args)
		if err := cd.Execute(&out); err == nil {
			t.Fatalf("Expected error to be returned for invalid folder")
		}
	}

	// Error checking bucket
	{
		var s3 mockS3Client
		var con context.Context
		var out mockOutputter
		args := []string{"bucket/folder"}
		fakeErr := errors.New("Error checking object")

		s3.pathExistsCallback = func(b, k string) (bool, error) {
			if b != "bucket" || k != "folder/" {
				t.Fatalf("Checking for unexpected bucket/folder: %v/%v", b, k)
			}

			return false, fakeErr
		}

		cd := NewCd(&s3, &con, args)
		if err := cd.Execute(&out); err != fakeErr {
			t.Fatalf("Unexpected error returned for folder cd")
		}
	}
}

func TestCdCommand_IsLongRunning(t *testing.T) {
	var s3 mockS3Client
	var con context.Context
	var args []string
	var cd CdCommand

	// No args
	cd = NewCd(&s3, &con, args)
	if cd.IsLongRunning() {
		t.Fatal("IsLongRunning should be false when there is no target")
	}

	// Root
	args = []string{context.PathDelimiter}
	cd = NewCd(&s3, &con, args)
	if cd.IsLongRunning() {
		t.Fatal("IsLongRunning should be false when the target is the root")
	}

	// Anything else should be long running.
	argSets := [][]string{
		[]string{"bucket"},
		[]string{"bucket/directory"},
		[]string{"bucket/../bucket"},
	}
	for _, argSet := range argSets {
		cd = NewCd(&s3, &con, argSet)

		if !cd.IsLongRunning() {
			t.Fatalf("Expected IsLongRunning to be true for target: %v", argSet)
		}
	}
}

func TestNewCd(t *testing.T) {
	var s3 mockS3Client
	var con context.Context
	args := []string{"directory"}

	cd := NewCd(&s3, &con, args)
	if cd.s3 != &s3 {
		t.Fatalf("Unexpected S3 client stored on cd command: %v", cd.s3)
	} else if cd.con != &con {
		t.Fatalf("Unexpected Context stored on cd command: %v", cd.con)
	} else if cd.args[0] != args[0] {
		t.Fatalf("Unexpected args stored on cd command: %v", cd.args)
	}
}
