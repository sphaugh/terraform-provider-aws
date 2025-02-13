---
subcategory: "Transfer"
layout: "aws"
page_title: "AWS: aws_transfer_server"
description: |-
  Provides a AWS Transfer Server resource.
---

# Resource: aws_transfer_server

Provides a AWS Transfer Server resource.

~> **NOTE on AWS IAM permissions:** If the `endpoint_type` is set to `VPC`, the `ec2:DescribeVpcEndpoints` and `ec2:ModifyVpcEndpoint` [actions](https://docs.aws.amazon.com/service-authorization/latest/reference/list_amazonec2.html#amazonec2-actions-as-permissions) are used.

## Example Usage

### Basic

```terraform
resource "aws_transfer_server" "example" {
  tags = {
    Name = "Example"
  }
}
```

### Security Policy Name

```terraform
resource "aws_transfer_server" "example" {
  security_policy_name = "TransferSecurityPolicy-2020-06"
}
```

### VPC Endpoint

```terraform
resource "aws_transfer_server" "example" {
  endpoint_type = "VPC"

  endpoint_details {
    address_allocation_ids = [aws_eip.example.id]
    subnet_ids             = [aws_subnet.example.id]
    vpc_id                 = aws_vpc.example.id
  }
}
```

### Protocols

```terraform
resource "aws_transfer_server" "example" {
  endpoint_type = "VPC"

  endpoint_details {
    subnet_ids = [aws_subnet.example.id]
    vpc_id     = aws_vpc.example.id
  }

  protocols   = ["FTP", "FTPS"]
  certificate = aws_acm_certificate.example.arn

  identity_provider_type = "API_GATEWAY"
  url                    = "${aws_api_gateway_deployment.example.invoke_url}${aws_api_gateway_resource.example.path}"
}
```

## Argument Reference

The following arguments are supported:

* `certificate` - (Optional) The Amazon Resource Name (ARN) of the AWS Certificate Manager (ACM) certificate. This is required when `protocols` is set to `FTPS`
* `domain` - (Optional) The domain of the storage system that is used for file transfers. Valid values are: `S3` and `EFS`. The default value is `S3`.
* `protocols` - (Optional) Specifies the file transfer protocol or protocols over which your file transfer protocol client can connect to your server's endpoint. This defaults to `SFTP` . The available protocols are:
    * `SFTP`: File transfer over SSH
    * `FTPS`: File transfer with TLS encryption
    * `FTP`: Unencrypted file transfer
* `endpoint_details` - (Optional) The virtual private cloud (VPC) endpoint settings that you want to configure for your SFTP server. Fields documented below.
* `endpoint_type` - (Optional) The type of endpoint that you want your SFTP server connect to. If you connect to a `VPC` (or `VPC_ENDPOINT`), your SFTP server isn't accessible over the public internet. If you want to connect your SFTP server via public internet, set `PUBLIC`.  Defaults to `PUBLIC`.
* `invocation_role` - (Optional) Amazon Resource Name (ARN) of the IAM role used to authenticate the user account with an `identity_provider_type` of `API_GATEWAY`.
* `host_key` - (Optional) RSA private key (e.g. as generated by the `ssh-keygen -N "" -m PEM -f my-new-server-key` command).
* `url` - (Optional) - URL of the service endpoint used to authenticate users with an `identity_provider_type` of `API_GATEWAY`.
* `identity_provider_type` - (Optional) The mode of authentication enabled for this service. The default value is `SERVICE_MANAGED`, which allows you to store and access SFTP user credentials within the service. `API_GATEWAY` indicates that user authentication requires a call to an API Gateway endpoint URL provided by you to integrate an identity provider of your choice.
* `directory_id` - (Optional) The directory service id of the directory service you want to connect to.
* `logging_role` - (Optional) Amazon Resource Name (ARN) of an IAM role that allows the service to write your SFTP users’ activity to your Amazon CloudWatch logs for monitoring and auditing purposes.
* `force_destroy` - (Optional) A boolean that indicates all users associated with the server should be deleted so that the Server can be destroyed without error. The default value is `false`. This option only applies to servers configured with a `SERVICE_MANAGED` `identity_provider_type`.
* `security_policy_name` - (Optional) Specifies the name of the security policy that is attached to the server. Possible values are `TransferSecurityPolicy-2018-11`, `TransferSecurityPolicy-2020-06`, and  `TransferSecurityPolicy-FIPS-2020-06`. Default value is: `TransferSecurityPolicy-2018-11`.
* `tags` - (Optional) A map of tags to assign to the resource. If configured with a provider [`default_tags` configuration block](https://www.terraform.io/docs/providers/aws/index.html#default_tags-configuration-block) present, tags with matching keys will overwrite those defined at the provider-level.

**endpoint_details** requires the following:

* `address_allocation_ids` - (Optional) A list of address allocation IDs that are required to attach an Elastic IP address to your SFTP server's endpoint. This property can only be used when `endpoint_type` is set to `VPC`.
* `security_group_ids` - (Optional) A list of security groups IDs that are available to attach to your server's endpoint. If no security groups are specified, the VPC's default security groups are automatically assigned to your endpoint. This property can only be used when `endpoint_type` is set to `VPC`.
* `subnet_ids` - (Optional) A list of subnet IDs that are required to host your SFTP server endpoint in your VPC. This property can only be used when `endpoint_type` is set to `VPC`.
* `vpc_endpoint_id` - (Optional) The ID of the VPC endpoint. This property can only be used when `endpoint_type` is set to `VPC_ENDPOINT`
* `vpc_id` - (Optional) The VPC ID of the virtual private cloud in which the SFTP server's endpoint will be hosted. This property can only be used when `endpoint_type` is set to `VPC`.

## Attributes Reference
In addition to all arguments above, the following attributes are exported:

* `arn` - Amazon Resource Name (ARN) of Transfer Server
* `id`  - The Server ID of the Transfer Server (e.g. `s-12345678`)
* `endpoint` - The endpoint of the Transfer Server (e.g. `s-12345678.server.transfer.REGION.amazonaws.com`)
* `host_key_fingerprint` - This value contains the message-digest algorithm (MD5) hash of the server's host key. This value is equivalent to the output of the `ssh-keygen -l -E md5 -f my-new-server-key` command.
* `tags_all` - A map of tags assigned to the resource, including those inherited from the provider [`default_tags` configuration block](https://www.terraform.io/docs/providers/aws/index.html#default_tags-configuration-block).

## Import

Transfer Servers can be imported using the `server id`, e.g.

```
$ terraform import aws_transfer_server.example s-12345678
```

Certain resource arguments, such as `host_key`, cannot be read via the API and imported into Terraform. Terraform will display a difference for these arguments the first run after import if declared in the Terraform configuration for an imported resource.
