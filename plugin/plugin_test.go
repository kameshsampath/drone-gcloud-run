// Copyright 2020 the Drone Authors. All rights reserved.
// Use of this source code is governed by the Blue Oak Model License
// that can be found in the LICENSE file.

package plugin

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	orderedmap "github.com/wk8/go-ordered-map/v2"
	"io/ioutil"
	"os"
	"path"
	"testing"
)

func TestPlugin(t *testing.T) {
	dir, _ := os.Getwd()
	credsFile := path.Join(dir, "testdata/sa.json")

	if err := godotenv.Load(path.Join(dir, "..", ".env")); err != nil {
		t.Fatalf("Error running tests %#v", err)
	}

	if err := ioutil.WriteFile("testdata/sa.json",
		[]byte(os.Getenv("SERVICE_ACCOUNT_JSON")), 0600); err != nil {
		t.Fatalf("Error running tests %#v", err)
	}

	if err := os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", credsFile); err != nil {
		t.Fatalf("Error running tests %#v", err)
	}

	ctx := context.Background()

	project := os.Getenv("GOOGLE_CLOUD_PROJECT")
	region := os.Getenv("GOOGLE_CLOUD_REGION")
	image := fmt.Sprintf("asia.gcr.io/%s/greeter", project)

	tests := orderedmap.New[string, Args]()

	tests.Set("createWithDefaults", Args{
		ProjectName: project,
		Location:    region,
		Image:       image,
		ServiceName: "foo",
	})
	tests.Set("createAsUnAuthenticated", Args{
		ProjectName:          project,
		Location:             region,
		Image:                image,
		ServiceName:          "foo",
		AllowUnauthenticated: true,
	})
	tests.Set("resetToAuthentication", Args{
		ProjectName: project,
		Location:    region,
		Image:       image,
		ServiceName: "foo",
	})
	tests.Set("deleteService", Args{
		ProjectName: project,
		Location:    region,
		Image:       image,
		ServiceName: "foo",
		Delete:      true,
	})

	for pair := tests.Oldest(); pair != nil; pair = pair.Next() {
		t.Run(pair.Key, func(t *testing.T) {
			err := Exec(ctx, pair.Value)
			if err != nil {
				t.Errorf("Error running comand %v", err)
			}
		})
	}

	os.Remove(credsFile)
}
