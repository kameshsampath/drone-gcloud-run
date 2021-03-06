// Copyright 2020 the Drone Authors. All rights reserved.
// Use of this source code is governed by the Blue Oak Model License
// that can be found in the LICENSE file.

package plugin

import (
	"context"
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/joho/godotenv"
	orderedmap "github.com/wk8/go-ordered-map/v2"
)

func TestPlugin(t *testing.T) {
	dir, _ := os.Getwd()

	if err := godotenv.Load(path.Join(dir, "..", ".env")); err != nil {
		t.Fatalf("Error running tests %#v", err)
	}

	ctx := context.Background()

	project := os.Getenv("GOOGLE_CLOUD_PROJECT")
	region := os.Getenv("GOOGLE_CLOUD_REGION")
	image := fmt.Sprintf("asia.gcr.io/%s/greeter", project)

	tests := orderedmap.New[string, Args]()

	tests.Set("createWithDefaults", Args{
		ServiceAccountJSON: os.Getenv("SERVICE_ACCOUNT_JSON"),
		ProjectName:        project,
		Region:             region,
		Image:              image,
		ServiceName:        "foo",
	})

	tests.Set("createAsUnAuthenticated", Args{
		ServiceAccountJSON:   os.Getenv("SERVICE_ACCOUNT_JSON"),
		ProjectName:          project,
		Region:               region,
		Image:                image,
		ServiceName:          "foo",
		AllowUnauthenticated: true,
	})

	tests.Set("resetToAuthentication", Args{
		ServiceAccountJSON: os.Getenv("SERVICE_ACCOUNT_JSON"),
		ProjectName:        project,
		Region:             region,
		Image:              image,
		ServiceName:        "foo",
	})

	tests.Set("deleteService", Args{
		ServiceAccountJSON: os.Getenv("SERVICE_ACCOUNT_JSON"),
		ProjectName:        project,
		Region:             region,
		Image:              image,
		ServiceName:        "foo",
		Delete:             true,
	})

	os.Setenv("DIGEST_FILE", path.Join(dir, "testdata", "greeter.digest"))
	tests.Set("UseImageDigestFromFile", Args{
		ServiceAccountJSON:   os.Getenv("SERVICE_ACCOUNT_JSON"),
		ProjectName:          project,
		Region:               region,
		Image:                image,
		DigestFile:           "$DIGEST_FILE",
		ServiceName:          "foo",
		Delete:               false,
		AllowUnauthenticated: false,
	})

	for pair := tests.Oldest(); pair != nil; pair = pair.Next() {
		t.Run(pair.Key, func(t *testing.T) {
			err := Exec(ctx, pair.Value)
			if err != nil {
				t.Errorf("Error running command %v", err)
			}
		})
	}
}
