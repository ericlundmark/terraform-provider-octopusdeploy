package octopusdeploy

import (
	"fmt"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceTenant() *schema.Resource {
	return &schema.Resource{
		Create: resourceTenantCreate,
		Read:   resourceTenantRead,
		Update: resourceTenantUpdate,
		Delete: resourceTenantDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceTenantRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*octopusdeploy.Client)

	tenantID := d.Id()
	tenant, err := client.Tenant.Get(tenantID)

	if err == octopusdeploy.ErrItemNotFound {
		d.SetId("")
		return nil
	}

	if err != nil {
		return fmt.Errorf("error reading tenant %s: %s", tenantID, err.Error())
	}

	d.Set("name", tenant.Name)
	d.Set("description", tenant.Description)

	return nil
}

func buildTenantResource(d *schema.ResourceData) *octopusdeploy.Tenant {
	tenantName := d.Get("name").(string)

	var tenantDesc string

	tenantDescInterface, ok := d.GetOk("description")
	if ok {
		tenantDesc = tenantDescInterface.(string)
	}

	var tenant = octopusdeploy.NewTenant(tenantName, tenantDesc)

	return tenant
}

func resourceTenantCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*octopusdeploy.Client)

	newTenant := buildTenantResource(d)
	tenant, err := client.Tenant.Add(newTenant)

	if err != nil {
		return fmt.Errorf("error creating tenant %s: %s", newTenant.Name, err.Error())
	}

	d.SetId(tenant.ID)

	return nil
}

func resourceTenantUpdate(d *schema.ResourceData, m interface{}) error {
	tenant := buildTenantResource(d)
	tenant.ID = d.Id() // set project struct ID so octopus knows which project to update

	client := m.(*octopusdeploy.Client)

	updatedTenant, err := client.Tenant.Update(tenant)

	if err != nil {
		return fmt.Errorf("error updating tenant id %s: %s", d.Id(), err.Error())
	}

	d.SetId(updatedTenant.ID)
	return nil
}

func resourceTenantDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*octopusdeploy.Client)

	tenantID := d.Id()

	err := client.Tenant.Delete(tenantID)

	if err != nil {
		return fmt.Errorf("error deleting tenant id %s: %s", tenantID, err.Error())
	}

	d.SetId("")
	return nil
}
