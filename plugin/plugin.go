// Copyright 2020 the Drone Authors. All rights reserved.
// Use of this source code is governed by the Blue Oak Model License
// that can be found in the LICENSE file.

package plugin

import (
	run "cloud.google.com/go/run/apiv2"
	"context"
	"fmt"
	"github.com/googleapis/gax-go/v2/apierror"
	"github.com/sirupsen/logrus"
	runpb "google.golang.org/genproto/googleapis/cloud/run/v2"
	iampb "google.golang.org/genproto/googleapis/iam/v1"
)

// Args provides plugin execution arguments.
type Args struct {
	Pipeline

	// Level defines the plugin log level.
	Level string `envconfig:"PLUGIN_LOG_LEVEL"`

	ProjectName          string `envconfig:"PLUGIN_PROJECT"`
	Location             string `envconfig:"PLUGIN_LOCATION"`
	ServiceName          string `envconfig:"PLUGIN_SERVICE_NAME"`
	Image                string `envconfig:"PLUGIN_IMAGE"`
	Delete               bool   `envconfig:"PLUGIN_DELETE"`
	AllowUnauthenticated bool   `envconfig:"PLUGIN_ALLOW_UNAUTHENTICATED"`
}

// Exec executes the plugin.
func Exec(ctx context.Context, args Args) error {
	c, err := run.NewServicesClient(ctx)
	defer c.Close()
	if err != nil {
		return err
	}

	svc, err := c.GetService(ctx, &runpb.GetServiceRequest{
		Name: fmt.Sprintf("projects/%s/locations/%s/services/%s", args.ProjectName,
			args.Location, args.ServiceName),
	})

	var errCode string

	if err != nil {
		apiErr, _ := apierror.FromError(err)
		errCode = apiErr.GRPCStatus().Code().String()
		// handle only known errors
		if errCode != "NotFound" {
			return err
		}
	}

	if svc == nil || errCode == "NotFound" {
		if err := createService(ctx, args, c); err != nil {
			return err
		}
	} else if svc != nil {
		if args.Delete {
			dOp, err := c.DeleteService(ctx, &runpb.DeleteServiceRequest{
				Name: svc.Name,
			})
			if err != nil {
				return err
			}
			logrus.Infof("\n%s Service deleting\n", args.ServiceName)
			for {
				_, err := dOp.Poll(ctx)
				if err != nil {
					return err
				}
				if dOp.Done() {
					logrus.Infof("\n%s Service successfully deleted\n", args.ServiceName)
					return nil
				}
				logrus.Info(".")
			}
		} else {
			if err := updateService(ctx, args, svc, c); err != nil {
				return err
			}
		}
	}

	return nil
}

func createService(ctx context.Context, args Args, c *run.ServicesClient) error {
	logrus.Infof("\nService %s does not exists, creating it\n", args.ServiceName)

	req := &runpb.CreateServiceRequest{
		Parent: fmt.Sprintf("projects/%s/locations/%s", args.ProjectName, args.Location),
		Service: &runpb.Service{
			Template: &runpb.RevisionTemplate{
				Containers: []*runpb.Container{
					{
						Image: args.Image,
					},
				},
			},
		},
		ServiceId: args.ServiceName,
	}

	logrus.Infof("\n Creating Service %s", args.ServiceName)

	o, err := c.CreateService(ctx, req)
	if err != nil {
		return err
	}

	var svc *runpb.Service
	for {
		svc, err = o.Poll(ctx)
		logrus.Info(".")
		if err != nil {
			return err
		}

		if o.Done() {
			logrus.Infof("\nService %s created\n", svc.Name)
			break
		}
	}

	if args.AllowUnauthenticated {
		if err := setIamPolicy(ctx, args, svc, c); err != nil {
			return err
		}
	}
	return nil
}

func updateService(ctx context.Context, args Args, svc *runpb.Service, c *run.ServicesClient) error {
	logrus.Infof("\nService %s already exists, will update\n", args.ServiceName)
	//Update values from the arguments
	svc.Template = &runpb.RevisionTemplate{
		Containers: []*runpb.Container{
			{
				Image: args.Image,
			},
		},
	}
	req := &runpb.UpdateServiceRequest{
		Service: svc,
	}
	uOp, err := c.UpdateService(ctx, req)
	if err != nil {
		return err
	}

	logrus.Infof("\n Updating Service %s", svc.Name)
	for {
		svc, err = uOp.Poll(ctx)
		logrus.Info(".")
		if err != nil {
			return err
		}

		if uOp.Done() {
			logrus.Infof("\nService %s updated\n", svc.Name)
			break
		}
	}

	if err := setIamPolicy(ctx, args, svc, c); err != nil {
		return err
	}
	return nil
}

func setIamPolicy(ctx context.Context, args Args, svc *runpb.Service, c *run.ServicesClient) error {
	iamReq := &iampb.SetIamPolicyRequest{
		Resource: svc.Name,
	}
	if args.AllowUnauthenticated {
		iamReq.Policy = &iampb.Policy{
			Bindings: []*iampb.Binding{
				{
					Members: []string{"allUsers"},
					Role:    "roles/run.invoker",
				},
			},
		}
	} else {
		logrus.Infoln("Resetting Service to be unauthenticated")
		policy, err := c.GetIamPolicy(ctx, &iampb.GetIamPolicyRequest{
			Resource: svc.Name,
		})
		if err != nil {
			return err
		}

		for _, b := range policy.Bindings {
			// handle only for run.invoker role
			if b.Role == "roles/run.invoker" {
				hasAllUsers := false
				var newMembers []string
				for _, m := range b.Members {
					if m == "allUsers" {
						hasAllUsers = true
						continue
					} else if m != "" {
						newMembers = append(newMembers, m)
					}
				}
				newMembers = append(newMembers, "allAuthenticatedUsers")
				if hasAllUsers {
					b.Members = newMembers
				}
			}
		}
		iamReq.Policy = policy
	}
	_, err := c.SetIamPolicy(ctx, iamReq)

	if err != nil {
		return err
	}

	return nil
}
