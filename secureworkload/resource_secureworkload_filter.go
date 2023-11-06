package secureworkload

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	// client "github.com/secureworkload-exchange/terraform-go-sdk"
)

func resourceSecureWorkloadFilter() *schema.Resource {
	return &schema.Resource{
		Create: resourceSecureWorkloadFilterCreate,
		Update: nil,
		Read:   resourceSecureWorkloadFilterRead,
		Delete: resourceSecureWorkloadFilterDelete,

		SchemaVersion: 1,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Optional:    false,
				ForceNew:    true,
				Description: "User-specified name for the inventory filter.",
			},
			"query": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "JSON object representation of an inventory filter query.",
			},
			"app_scope_id": {
				Type:        schema.TypeString,
				Optional:    false,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the scope associated with the filter.",
			},
			"primary": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Default:     false,
				Description: "(Optional) When true, the filter is restricted to the ownership scope.",
			},
			"public": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Default:     false,
				Description: "(Optional) When true the filter provides a service for its scope. Must also be primary/scope restricted.",
			},
		},
	}
}

var requiredCreateFilterParams = []string{"name", "app_scope_id", "query"}

func resourceSecureWorkloadFilterCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(Client)
	for _, param := range requiredCreateFilterParams {
		if d.Get(param) == "" {
			return fmt.Errorf("%s is required but was not provided", param)
		}
	}
	createFilterParams := CreateFilterRequest{
		Name:       d.Get("name").(string),
		AppScopeId: d.Get("app_scope_id").(string),
		Query:      []byte(d.Get("query").(string)),
		Primary:    d.Get("primary").(bool),
		Public:     d.Get("public").(bool),
	}
	filter, err := client.CreateFilter(createFilterParams)
	if err != nil {
		return err
	}
	d.SetId(filter.Id)
	return nil
}

func resourceSecureWorkloadFilterRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(Client)
	filter, err := client.DescribeFilter(d.Id())
	if err != nil {
		return err
	}
	d.Set("name", filter.Name)
	d.Set("app_scope_id", filter.AppScopeId)
	d.Set("primary", filter.Primary)
	d.Set("public", filter.Public)
	return nil
}

func resourceSecureWorkloadFilterDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(Client)
	return client.DeleteFilter(d.Id())
}
