package main

import "github.com/pulumi/pulumi/sdk/v3/go/pulumi"

// tags returns a standard tag set for all resources.
func tags(name, appEnv string) pulumi.StringMap {
	return pulumi.StringMap{
		"Name":    pulumi.Sprintf("subject-data-%s-%s", name, appEnv),
		"Service": pulumi.String("subject-data"),
		"Env":     pulumi.String(appEnv),
	}
}
