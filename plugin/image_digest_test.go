package plugin

import (
	"os"
	"path"
	"testing"
)

func TestResolveImage(t *testing.T) {
	want := "sha256:eb2623f9d4199424b49c7cc3358f30a1da697d48e1432732ae69835e55683f90"

	got, err := resolveToDigest("docker.io/kameshsampath/drone-gcloud-run:v1.0.0", "", "linux/amd64")

	if err != nil {
		t.Fatalf(err.Error())
	}

	if got != want {
		t.Errorf("Got %s but want %s", got, want)
	}
}

func TestImageFromDigestFile(t *testing.T) {
	cwd, _ := os.Getwd()
	f := path.Join(cwd, "testdata", "digest.txt")
	want := "sha256:eb2623f9d4199424b49c7cc3358f30a1da697d48e1432732ae69835e55683f90"
	got, err := resolveToDigest("docker.io/kameshsampath/drone-gcloud-run:v1.0.0", f, "linux/amd64")

	if err != nil {
		t.Fatalf(err.Error())
	}

	if got != want {
		t.Errorf("Got %s but want %s", got, want)
	}
}
