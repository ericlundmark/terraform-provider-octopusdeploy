package octopusdeploy

import (
	"fmt"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataTenant() *schema.Resource {
	return &schema.Resource{
		Read: dataTenantReadByName,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"space_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"cloned_from_tenant_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataTenantReadByName(d *schema.ResourceData, m interface{}) error {
	client := m.(*octopusdeploy.Client)

	tenantName := d.Get("name")
	tenant, err := client.Tenant.GetByName(tenantName.(string))

	if err == octopusdeploy.ErrItemNotFound {
		return nil
	}

	if err != nil {
		return fmt.Errorf("error reading tenant with name %s: %s", tenantName, err.Error())
	}

	d.SetId(tenant.ID)

	d.Set("name", tenant.Name)
	d.Set("space_id", tenant.Description)
	d.Set("cloned_from_tenant_id", tenant.SpaceID)
	d.Set("description", tenant.ClonedFromTenantID)

	return nil
}
