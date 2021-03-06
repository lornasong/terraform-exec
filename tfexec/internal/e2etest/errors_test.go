// This file contains tests that only compile/work in Go 1.13 and forward
// +build go1.13

package e2etest

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/hashicorp/go-version"

	"github.com/hashicorp/terraform-exec/tfexec"
)

func TestUnparsedError(t *testing.T) {
	// This simulates an unparsed error from the Cmd.Run method (in this case file not found). This
	// is to ensure we don't miss raising unexpected errors in addition to parsed / well known ones.
	runTest(t, "", func(t *testing.T, tfv *version.Version, tf *tfexec.Terraform) {

		// force delete the working dir to cause an os.PathError
		err := os.RemoveAll(tf.WorkingDir())
		if err != nil {
			t.Fatal(err)
		}

		err = tf.Init(context.Background())
		if err == nil {
			t.Fatalf("expected error running Init, none returned")
		}
		var e *os.PathError
		if !errors.As(err, &e) {
			t.Fatalf("expected os.PathError, got %T, %s", err, err)
		}
	})
}

func TestMissingVar(t *testing.T) {
	runTest(t, "var", func(t *testing.T, tfv *version.Version, tf *tfexec.Terraform) {
		err := tf.Init(context.Background())
		if err != nil {
			t.Fatalf("err during init: %s", err)
		}

		_, err = tf.Plan(context.Background())
		if err == nil {
			t.Fatalf("expected error running Plan, none returned")
		}
		var e *tfexec.ErrMissingVar
		if !errors.As(err, &e) {
			t.Fatalf("expected ErrMissingVar, got %T, %s", err, err)
		}

		if e.VariableName != "no_default" {
			t.Fatalf("expected missing no_default, got %q", e.VariableName)
		}

		_, err = tf.Plan(context.Background(), tfexec.Var("no_default=foo"))
		if err != nil {
			t.Fatalf("expected no error, got %s", err)
		}
	})
}
