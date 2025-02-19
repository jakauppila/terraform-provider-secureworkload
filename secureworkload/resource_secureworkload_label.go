package secureworkload

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	// client "github.com/secureworkload-exchange/terraform-go-sdk"
	// secureworkload "github.com/secureworkload-exchange/terraform-go-sdk"
)

const (
	TagIdDelimter = ":"
)

func resourceSecureWorkloadLabel() *schema.Resource {
	return &schema.Resource{
		Description: "Resource for creating a new label in Secure Workload\n" +
			"\n" +
			"## Example\n" +
			"An example is shown below: \n" +
			"```hcl\n" +
			"resource \"secureworkload_label\" \"label-1\" {\n" +
			"	 ip = \"1.2.3.4\"\n" +
			"    root_scope_name = \"scope_name\"\n" +
			"    attributes = {\n" +
			"        Environment = \"test\"\n" +
			"        Datacenter = \"aws\"\n" +
			"        app_name = \"app_name\"\n" +
			"    }\n" +	
			"}\n" +
			"```\n" +
			"**Note:** If creating multiple rules during a single `terraform apply`, remember to use `depends_on` to chain the rules so that terraform creates it in the same order that you intended.\n" ,
		Create: resourceSecureWorkloadTagCreate,
		Update: resourceSecureWorkloadTagCreate,
		Read:   resourceSecureWorkloadTagRead,
		Delete: resourceSecureWorkloadTagDelete,

		SchemaVersion: 1,

		Schema: map[string]*schema.Schema{
			"root_scope_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				ForceNew:    true,
				Description: "SecureWorkload root app scope name.",
			},
			"ip": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "IPv4/IPv6 address or subnet.",
			},
			"attributes": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: false,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Key/value map for tagging matching flows and inventory items.",
			},
		},
	}
}

var requiredCreateTagParams = []string{"ip", "attributes"}

func resourceSecureWorkloadTagCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(Client)
	for _, param := range requiredCreateTagParams {
		if d.Get(param) == "" {
			return fmt.Errorf("%s is required but was not provided", param)
		}
	}
	tenantName := d.Get("root_scope_name").(string)
	if tenantName == "" {
		tenantURL := client.Config.APIURL
		// strip protocol and extract the tenant name/subdomain from the url
		// e.g. https://acme.secureworkloadpreview.com => acme
		tenantName = strings.Split(strings.Split(tenantURL, "://")[1], ".")[0]
	}
	attributes := d.Get("attributes").(map[string]interface{})
	createTagParams := CreateTagRequest{
		RootScopeName: tenantName,
		Ip:            d.Get("ip").(string),
		Attributes:    attributes,
	}
	tag, err := client.CreateTag(createTagParams)
	if err != nil {
		return err
	}
	d.SetId(fmt.Sprintf("%s%s%s", createTagParams.RootScopeName, TagIdDelimter, tag.Ip))
	return nil
}

func resourceSecureWorkloadTagRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(Client)
	tagIdComponents := strings.Split(d.Id(), TagIdDelimter)
	describeTagRequest := DescribeTagRequest{
		RootAppScopeName: tagIdComponents[0],
		Ip:               tagIdComponents[1],
	}
	attributes := make(map[string]string)
	err := client.DescribeTag(describeTagRequest, &attributes)
	if err != nil {
		return err
	}
	d.Set("root_scope_name", describeTagRequest.RootAppScopeName)
	d.Set("ip", describeTagRequest.Ip)
	d.Set("attributes", attributes)
	return nil
}

func resourceSecureWorkloadTagDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(Client)
	tagIdComponents := strings.Split(d.Id(), TagIdDelimter)
	deleteTagRequest := DeleteTagRequest{
		RootAppScopeName: tagIdComponents[0],
		Ip:               tagIdComponents[1],
	}
	return client.DeleteTag(deleteTagRequest)
}
