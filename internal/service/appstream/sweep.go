package appstream

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/appstream"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/sweep"
)

func init() {
	resource.AddTestSweepers("aws_appstream_image_builder", &resource.Sweeper{
		Name: "aws_appstream_image_builder",
		F:    sweepImageBuilder,
		Dependencies: []string{
			"aws_vpc",
			"aws_subnet",
		},
	})
}

func sweepImageBuilder(region string) error {
	client, err := sweep.SharedRegionalSweepClient(region)
	if err != nil {
		return fmt.Errorf("error getting client: %w", err)
	}

	conn := client.(*conns.AWSClient).AppStreamConn
	sweepResources := make([]*sweep.SweepResource, 0)
	var errs *multierror.Error

	input := &appstream.DescribeImageBuildersInput{}

	err = DescribeImageBuildersPagesWithContext(context.TODO(), conn, input, func(page *appstream.DescribeImageBuildersOutput, lastPage bool) bool {
		if page == nil {
			return !lastPage
		}

		for _, imageBuilder := range page.ImageBuilders {
			if imageBuilder == nil {
				continue
			}

			id := aws.StringValue(imageBuilder.Name)

			r := ResourceImageBuilder()
			d := r.Data(nil)
			d.SetId(id)

			sweepResources = append(sweepResources, sweep.NewSweepResource(r, d, client))
		}

		return !lastPage
	})

	if err != nil {
		errs = multierror.Append(errs, fmt.Errorf("error listing AppStream Image Builders: %w", err))
	}

	if err = sweep.SweepOrchestrator(sweepResources); err != nil {
		errs = multierror.Append(errs, fmt.Errorf("error sweeping AppStream Image Builders for %s: %w", region, err))
	}

	if sweep.SkipSweepError(err) {
		log.Printf("[WARN] Skipping AppStream Image Builders sweep for %s: %s", region, err)
		return nil // In case we have completed some pages, but had errors
	}

	return errs.ErrorOrNil()
}
