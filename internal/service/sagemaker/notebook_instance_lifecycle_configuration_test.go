package sagemaker_test

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sagemaker"
	"github.com/hashicorp/aws-sdk-go-base/tfawserr"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tfsagemaker "github.com/hashicorp/terraform-provider-aws/internal/service/sagemaker"
	"github.com/hashicorp/terraform-provider-aws/internal/sweep"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
)

func TestAccSageMakerNotebookInstanceLifecycleConfiguration_basic(t *testing.T) {
	var lifecycleConfig sagemaker.DescribeNotebookInstanceLifecycleConfigOutput
	rName := sdkacctest.RandomWithPrefix(tfsagemaker.NotebookInstanceLifecycleConfigurationResourcePrefix)
	resourceName := "aws_sagemaker_notebook_instance_lifecycle_configuration.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, sagemaker.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckNotebookInstanceLifecycleConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSagemakerNotebookInstanceLifecycleConfigurationConfig_Basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNotebookInstanceLifecycleConfigurationExists(resourceName, &lifecycleConfig),

					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckNoResourceAttr(resourceName, "on_create"),
					resource.TestCheckNoResourceAttr(resourceName, "on_start"),
					acctest.CheckResourceAttrRegionalARN(resourceName, "arn", "sagemaker", fmt.Sprintf("notebook-instance-lifecycle-config/%s", rName)),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccSageMakerNotebookInstanceLifecycleConfiguration_update(t *testing.T) {
	var lifecycleConfig sagemaker.DescribeNotebookInstanceLifecycleConfigOutput
	rName := sdkacctest.RandomWithPrefix(tfsagemaker.NotebookInstanceLifecycleConfigurationResourcePrefix)
	resourceName := "aws_sagemaker_notebook_instance_lifecycle_configuration.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, sagemaker.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckNotebookInstanceLifecycleConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSagemakerNotebookInstanceLifecycleConfigurationConfig_Basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNotebookInstanceLifecycleConfigurationExists(resourceName, &lifecycleConfig),

					resource.TestCheckResourceAttr(resourceName, "name", rName),
				),
			},
			{
				Config: testAccSagemakerNotebookInstanceLifecycleConfigurationConfig_Update(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNotebookInstanceLifecycleConfigurationExists(resourceName, &lifecycleConfig),

					resource.TestCheckResourceAttr(resourceName, "on_create", verify.Base64Encode([]byte("echo bla"))),
					resource.TestCheckResourceAttr(resourceName, "on_start", verify.Base64Encode([]byte("echo blub"))),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckNotebookInstanceLifecycleConfigurationExists(resourceName string, lifecycleConfig *sagemaker.DescribeNotebookInstanceLifecycleConfigOutput) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).SageMakerConn
		output, err := conn.DescribeNotebookInstanceLifecycleConfig(&sagemaker.DescribeNotebookInstanceLifecycleConfigInput{
			NotebookInstanceLifecycleConfigName: aws.String(rs.Primary.ID),
		})

		if err != nil {
			return err
		}

		if output == nil {
			return fmt.Errorf("no SageMaker Notebook Instance Lifecycle Configuration")
		}

		*lifecycleConfig = *output

		return nil
	}
}

func testAccCheckNotebookInstanceLifecycleConfigurationDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_sagemaker_notebook_instance_lifecycle_configuration" {
			continue
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).SageMakerConn
		lifecycleConfig, err := conn.DescribeNotebookInstanceLifecycleConfig(&sagemaker.DescribeNotebookInstanceLifecycleConfigInput{
			NotebookInstanceLifecycleConfigName: aws.String(rs.Primary.ID),
		})

		if err != nil {
			if tfawserr.ErrMessageContains(err, "ValidationException", "") {
				continue
			}
			return err
		}

		if lifecycleConfig != nil && aws.StringValue(lifecycleConfig.NotebookInstanceLifecycleConfigName) == rs.Primary.ID {
			return fmt.Errorf("SageMaker Notebook Instance Lifecycle Configuration %s still exists", rs.Primary.ID)
		}
	}
	return nil
}

func testAccSagemakerNotebookInstanceLifecycleConfigurationConfig_Basic(rName string) string {
	return fmt.Sprintf(`
resource "aws_sagemaker_notebook_instance_lifecycle_configuration" "test" {
  name = %q
}
`, rName)
}

func testAccSagemakerNotebookInstanceLifecycleConfigurationConfig_Update(rName string) string {
	return fmt.Sprintf(`
resource "aws_sagemaker_notebook_instance_lifecycle_configuration" "test" {
  name      = %q
  on_create = base64encode("echo bla")
  on_start  = base64encode("echo blub")
}
`, rName)
}
