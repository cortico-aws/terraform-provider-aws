package memorydb_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/service/memorydb"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tfmemorydb "github.com/hashicorp/terraform-provider-aws/internal/service/memorydb"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

func TestAccMemoryDBCluster_basic(t *testing.T) {
	rName := "tf-test-" + sdkacctest.RandString(8)
	resourceName := "aws_memorydb_cluster.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t); testAccPreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, memorydb.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccClusterConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClusterExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "acl_name", "open-access"),
					acctest.CheckResourceAttrRegionalARN(resourceName, "arn", "memorydb", "cluster/"+rName),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "node_type", "db.t4g.small"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.Test", "test"),
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

func TestAccMemoryDBCluster_disappears(t *testing.T) {
	rName := "tf-test-" + sdkacctest.RandString(8)
	resourceName := "aws_memorydb_cluster.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t); testAccPreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, memorydb.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccClusterConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClusterExists(resourceName),
					acctest.CheckResourceDisappears(acctest.Provider, tfmemorydb.ResourceCluster(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccCheckClusterDestroy(s *terraform.State) error {
	conn := acctest.Provider.Meta().(*conns.AWSClient).MemoryDBConn

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_memorydb_cluster" {
			continue
		}

		_, err := tfmemorydb.FindClusterByName(context.Background(), conn, rs.Primary.Attributes["name"])

		if tfresource.NotFound(err) {
			continue
		}

		if err != nil {
			return err
		}

		return fmt.Errorf("MemoryDB Cluster %s still exists", rs.Primary.ID)
	}

	return nil
}

func testAccCheckClusterExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No MemoryDB Cluster ID is set")
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).MemoryDBConn

		_, err := tfmemorydb.FindClusterByName(context.Background(), conn, rs.Primary.Attributes["name"])

		if err != nil {
			return err
		}

		return nil
	}
}

func testAccClusterConfig_basic(rName string) string {
	return acctest.ConfigCompose(
		//testAccClusterConfigBase(rName),
		fmt.Sprintf(`
resource "aws_memorydb_cluster" "test" {
  acl_name  = "open-access"
  name      = %[1]q
  node_type = "db.t4g.small"

  tags = {
    Test = "test"
  }
}
`, rName),
	)
}
