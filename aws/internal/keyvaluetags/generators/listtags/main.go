//go:build ignore
// +build ignore

package main

import (
	"bytes"
	"go/format"
	"log"
	"os"
	"sort"
	"strings"
	"text/template"

	"github.com/terraform-providers/terraform-provider-aws/aws/internal/keyvaluetags"
)

const filename = `list_tags_gen.go`

var serviceNames = []string{
	"accessanalyzer",
	"acm",
	"acmpca",
	"amplify",
	"apigatewayv2",
	"appconfig",
	"appmesh",
	"apprunner",
	"appstream",
	"appsync",
	"athena",
	"autoscaling",
	"backup",
	"batch",
	"cloud9",
	"cloudfront",
	"cloudhsmv2",
	"cloudtrail",
	"cloudwatch",
	"cloudwatchevents",
	"cloudwatchlogs",
	"codeartifact",
	"codecommit",
	"codedeploy",
	"codepipeline",
	"codestarconnections",
	"codestarnotifications",
	"cognitoidentity",
	"cognitoidentityprovider",
	"configservice",
	"databasemigrationservice",
	"dataexchange",
	"datasync",
	"dax",
	"devicefarm",
	"directconnect",
	"directoryservice",
	"dlm",
	"docdb",
	"dynamodb",
	"ec2",
	"ecr",
	"ecs",
	"efs",
	"eks",
	"elasticache",
	"elasticbeanstalk",
	"elasticsearchservice",
	"elb",
	"elbv2",
	"firehose",
	"fsx",
	"gamelift",
	"glacier",
	"globalaccelerator",
	"glue",
	"guardduty",
	"greengrass",
	"imagebuilder",
	"inspector",
	"iot",
	"iotanalytics",
	"iotevents",
	"kafka",
	"kinesis",
	"kinesisanalytics",
	"kinesisanalyticsv2",
	"kinesisvideo",
	"kms",
	"lambda",
	"licensemanager",
	"mediaconnect",
	"mediaconvert",
	"medialive",
	"mediapackage",
	"mediastore",
	"mq",
	"neptune",
	"networkfirewall",
	"networkmanager",
	"opsworks",
	"organizations",
	"pinpoint",
	"qldb",
	"quicksight",
	"rds",
	"resourcegroups",
	"route53",
	"route53recoveryreadiness",
	"route53resolver",
	"sagemaker",
	"securityhub",
	"servicediscovery",
	"schemas",
	"sfn",
	"shield",
	"signer",
	"sns",
	"sqs",
	"ssm",
	"ssoadmin",
	"storagegateway",
	"swf",
	"timestreamwrite",
	"transfer",
	"waf",
	"wafregional",
	"wafv2",
	"worklink",
	"workspaces",
	"xray",
}

type TemplateData struct {
	ServiceNames []string
}

func main() {
	// Always sort to reduce any potential generation churn
	sort.Strings(serviceNames)

	templateData := TemplateData{
		ServiceNames: serviceNames,
	}
	templateFuncMap := template.FuncMap{
		"ClientType":                           keyvaluetags.ServiceClientType,
		"ListTagsFunction":                     keyvaluetags.ServiceListTagsFunction,
		"ListTagsInputFilterIdentifierName":    keyvaluetags.ServiceListTagsInputFilterIdentifierName,
		"ListTagsInputIdentifierField":         keyvaluetags.ServiceListTagsInputIdentifierField,
		"ListTagsInputIdentifierRequiresSlice": keyvaluetags.ServiceListTagsInputIdentifierRequiresSlice,
		"ListTagsOutputTagsField":              keyvaluetags.ServiceListTagsOutputTagsField,
		"ParentResourceNotFoundError":          keyvaluetags.ServiceParentResourceNotFoundError,
		"TagPackage":                           keyvaluetags.ServiceTagPackage,
		"TagResourceTypeField":                 keyvaluetags.ServiceTagResourceTypeField,
		"TagTypeIdentifierField":               keyvaluetags.ServiceTagTypeIdentifierField,
		"Title":                                strings.Title,
	}

	tmpl, err := template.New("listtags").Funcs(templateFuncMap).Parse(templateBody)

	if err != nil {
		log.Fatalf("error parsing template: %s", err)
	}

	var buffer bytes.Buffer
	err = tmpl.Execute(&buffer, templateData)

	if err != nil {
		log.Fatalf("error executing template: %s", err)
	}

	generatedFileContents, err := format.Source(buffer.Bytes())

	if err != nil {
		log.Fatalf("error formatting generated file: %s", err)
	}

	f, err := os.Create(filename)

	if err != nil {
		log.Fatalf("error creating file (%s): %s", filename, err)
	}

	defer f.Close()

	_, err = f.Write(generatedFileContents)

	if err != nil {
		log.Fatalf("error writing to file (%s): %s", filename, err)
	}
}

var templateBody = `
// Code generated by generators/listtags/main.go; DO NOT EDIT.

package keyvaluetags

import (
	"github.com/aws/aws-sdk-go/aws"
{{- range .ServiceNames }}
	"github.com/aws/aws-sdk-go/service/{{ . }}"
{{- end }}
    "github.com/hashicorp/aws-sdk-go-base/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)
{{ range .ServiceNames }}

// {{ . | Title }}ListTags lists {{ . }} service tags.
// The identifier is typically the Amazon Resource Name (ARN), although
// it may also be a different identifier depending on the service.
func {{ . | Title }}ListTags(conn {{ . | ClientType }}, identifier string{{ if . | TagResourceTypeField }}, resourceType string{{ end }}) (KeyValueTags, error) {
	input := &{{ . | TagPackage  }}.{{ . | ListTagsFunction }}Input{
		{{- if . | ListTagsInputFilterIdentifierName }}
		Filters: []*{{ . | TagPackage  }}.Filter{
			{
				Name:   aws.String("{{ . | ListTagsInputFilterIdentifierName }}"),
				Values: []*string{aws.String(identifier)},
			},
		},
		{{- else }}
		{{- if . | ListTagsInputIdentifierRequiresSlice }}
		{{ . | ListTagsInputIdentifierField }}: aws.StringSlice([]string{identifier}),
		{{- else }}
		{{ . | ListTagsInputIdentifierField }}: aws.String(identifier),
		{{- end }}
		{{- if . | TagResourceTypeField }}
		{{ . | TagResourceTypeField }}:         aws.String(resourceType),
		{{- end }}
		{{- end }}
	}

	output, err := conn.{{ . | ListTagsFunction }}(input)

	{{ . | ParentResourceNotFoundError }}

	if err != nil {
		return New(nil), err
	}

	return {{ . | Title }}KeyValueTags(output.{{ . | ListTagsOutputTagsField }}{{ if . | TagTypeIdentifierField }}, identifier{{ if . | TagResourceTypeField }}, resourceType{{ end }}{{ end }}), nil
}
{{- end }}
`
