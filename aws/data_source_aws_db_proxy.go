package aws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/terraform-providers/terraform-provider-aws/aws/internal/service/rds/finder"
)

func dataSourceAwsDbProxy() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAwsDbProxyRead,
		Schema: map[string]*schema.Schema{
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"auth": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"auth_scheme": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"iam_auth": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"secret_arn": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"debug_logging": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"engine_family": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"idle_client_timeout": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"require_tls": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"role_arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vpc_security_group_ids": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"vpc_subnet_ids": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceAwsDbProxyRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).rdsconn

	name := d.Get("name").(string)
	dbProxy, err := finder.DBProxyByName(conn, name)

	if err != nil {
		return fmt.Errorf("error reading RDS DB Proxy (%s): %w", name, err)
	}

	d.SetId(name)
	d.Set("arn", dbProxy.DBProxyArn)
	d.Set("auth", flattenDbProxyAuths(dbProxy.Auth))
	d.Set("debug_logging", dbProxy.DebugLogging)
	d.Set("endpoint", dbProxy.Endpoint)
	d.Set("engine_family", dbProxy.EngineFamily)
	d.Set("idle_client_timeout", dbProxy.IdleClientTimeout)
	d.Set("require_tls", dbProxy.RequireTLS)
	d.Set("role_arn", dbProxy.RoleArn)
	d.Set("vpc_id", dbProxy.VpcId)
	d.Set("vpc_security_group_ids", aws.StringValueSlice(dbProxy.VpcSecurityGroupIds))
	d.Set("vpc_subnet_ids", aws.StringValueSlice(dbProxy.VpcSubnetIds))

	return nil
}
