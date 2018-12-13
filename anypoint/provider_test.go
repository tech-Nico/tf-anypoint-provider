package anypoint

import (
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"log"
	"os"
	"testing"
)

var testAccProvider *schema.Provider
var testAccTemplateProvider *schema.Provider
var testAccProviders map[string]terraform.ResourceProvider
var testAccProviderFactories func(providers *[]*schema.Provider) map[string]terraform.ResourceProviderFactory

func init() {

	testAccProvider = Provider()
	testAccProviders = map[string]terraform.ResourceProvider{
		"ap": testAccProvider,
	}

	testAccProviderFactories = func(providers *[]*schema.Provider) map[string]terraform.ResourceProviderFactory {
		return map[string]terraform.ResourceProviderFactory{
			"ap": func() (terraform.ResourceProvider, error) {
				p := Provider()
				*providers = append(*providers, p)
				return p, nil
			},
		}
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func testAccPreCheck(t *testing.T) {

	if v := os.Getenv("ANYPOINT_URL"); v == "" {
		t.Fatalf("ANYPOINT_URL must be set for acceptance tests")
	}

	if v := os.Getenv("ANYPOINT_USERNAME"); v == "" {
		t.Fatalf("ANYPOINT_USERNAME must be set for acceptance tests")
	}

	if v := os.Getenv("ANYPOINT_PASSWORD"); v == "" {
		t.Fatalf("ANYPOINT_PASSWORD must be set for acceptance tests")
	}

	err := testAccProvider.Configure(terraform.NewResourceConfig(nil))

	if err != nil {
		t.Fatal(err)
	}
}

func testAccCheckWithProviders(f func(*terraform.State, *schema.Provider) error, providers *[]*schema.Provider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		numberOfProviders := len(*providers)
		for i, provider := range *providers {
			if provider.Meta() == nil {
				log.Printf("[DEBUG] Skipping empty provider %d (total: %d)", i, numberOfProviders)
				continue
			}
			log.Printf("[DEBUG] Calling check with provider %d (total: %d)", i, numberOfProviders)
			if err := f(s, provider); err != nil {
				return err
			}
		}
		return nil
	}
}

func TestProvider_impl(t *testing.T) {
	var _ terraform.ResourceProvider = Provider()
}
