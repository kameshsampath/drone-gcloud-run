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

type validationTestCase struct {
	args Args
	err  string
}

func TestValidatePluginParameters(t *testing.T) {
	dir, _ := os.Getwd()

	if err := godotenv.Load(path.Join(dir, "..", ".env")); err != nil {
		t.Fatalf("Error running tests %#v", err)
	}

	ctx := context.Background()

	project := os.Getenv("GOOGLE_CLOUD_PROJECT")
	region := os.Getenv("GOOGLE_CLOUD_REGION")
	image := fmt.Sprintf("asia.gcr.io/%s/greeter", project)

	tests := orderedmap.New[string, validationTestCase]()

	tests.Set("errWhenNoServiceAccountJson", validationTestCase{
		err: "no service account json specified",
		args: Args{
			ProjectName: project,
			Region:      region,
			Image:       image,
			ServiceName: "foo",
		},
	})

	tests.Set("errWhenNoProjectName", validationTestCase{
		err: "no google cloud project specified",
		args: Args{
			ServiceAccountJSON: "Some Content",
			Region:             region,
			Image:              image,
			ServiceName:        "foo",
		},
	})

	tests.Set("errWhenNoRegion", validationTestCase{
		err: "no google cloud region specified",
		args: Args{
			ServiceAccountJSON: "Some Content",
			ProjectName:        project,
			Image:              image,
			ServiceName:        "foo",
		},
	})

	tests.Set("errWhenNoServiceName", validationTestCase{
		err: "no google cloud run service name specified",
		args: Args{
			ServiceAccountJSON: "Some Content",
			ProjectName:        project,
			Region:             region,
			Image:              image,
		},
	})

	tests.Set("errWhenNoServiceImageWhenNotDelete", validationTestCase{
		err: "no google cloud run service image specified",
		args: Args{
			ServiceAccountJSON: "Some Content",
			ProjectName:        project,
			Region:             region,
			ServiceName:        "foo",
		},
	})

	for pair := tests.Oldest(); pair != nil; pair = pair.Next() {
		t.Run(pair.Key, func(t *testing.T) {
			err := Exec(ctx, pair.Value.args)
			if err != nil {
				if err.Error() != pair.Value.err {
					t.Errorf("Expecting error %s but got %s", pair.Value.err, err.Error())
				}
			} else {
				t.Error("Expecting error but got nil")
			}
		})
	}
}
