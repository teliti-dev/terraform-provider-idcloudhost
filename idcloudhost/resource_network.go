package idcloudhost

import (
	"context"
	"fmt"
	"strconv"

	idcloudhostAPI "github.com/bapung/idcloudhost-go-client-library/idcloudhost/api"
	idcloudhostNetwork "github.com/bapung/idcloudhost-go-client-library/idcloudhost/network"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceNetwork() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNetworkCreate,
		ReadContext:   resourceNetworkRead,
		UpdateContext: resourceNetworkUpdate,
		DeleteContext: resourceNetworkDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"uuid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"user_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"default": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceNetworkCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*idcloudhostAPI.APIClient)
	var diags diag.Diagnostics

	name := d.Get("name").(string)

	networkApi := c.Network
	if err := networkApi.CreateNetwork(name); err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create network",
			Detail:   fmt.Sprint(err),
		})
		return diags
	}

	d.SetId(networkApi.Network.UUID)

	// Set as default if requested
	if d.Get("default").(bool) {
		if err := networkApi.SetAsDefault(networkApi.Network.UUID); err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "Unable to set network as default",
				Detail:   fmt.Sprint(err),
			})
		}
	}

	resourceNetworkRead(ctx, d, m)

	return diags
}

func resourceNetworkRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*idcloudhostAPI.APIClient)
	var diags diag.Diagnostics

	uuid := d.Id()

	networkApi := c.Network
	if err := networkApi.GetNetwork(uuid); err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to get network",
			Detail:   fmt.Sprint(err),
		})
		return diags
	}

	if err := setNetworkResource(d, networkApi.Network); err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to set network resource",
			Detail:   fmt.Sprint(err),
		})
		return diags
	}

	return diags
}

func resourceNetworkUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*idcloudhostAPI.APIClient)
	var diags diag.Diagnostics

	uuid := d.Id()
	networkApi := c.Network

	if d.HasChange("name") {
		name := d.Get("name").(string)
		if err := networkApi.UpdateNetwork(uuid, name); err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Unable to update network",
				Detail:   fmt.Sprint(err),
			})
			return diags
		}
	}

	if d.HasChange("default") {
		if d.Get("default").(bool) {
			if err := networkApi.SetAsDefault(uuid); err != nil {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "Unable to set network as default",
					Detail:   fmt.Sprint(err),
				})
				return diags
			}
		}
	}

	return resourceNetworkRead(ctx, d, m)
}

func resourceNetworkDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*idcloudhostAPI.APIClient)
	var diags diag.Diagnostics

	uuid := d.Id()

	networkApi := c.Network
	if err := networkApi.DeleteNetwork(uuid); err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to delete network",
			Detail:   fmt.Sprint(err),
		})
		return diags
	}

	d.SetId("")

	return diags
}

func setNetworkResource(d *schema.ResourceData, network *idcloudhostNetwork.Network) error {
	if err := d.Set("id", strconv.Itoa(network.ID)); err != nil {
		return err
	}
	if err := d.Set("name", network.Name); err != nil {
		return err
	}
	if err := d.Set("uuid", network.UUID); err != nil {
		return err
	}
	if err := d.Set("user_id", network.UserID); err != nil {
		return err
	}
	if err := d.Set("default", network.Default); err != nil {
		return err
	}
	if err := d.Set("description", network.Description); err != nil {
		return err
	}
	if err := d.Set("created_at", network.CreatedAt); err != nil {
		return err
	}
	if err := d.Set("updated_at", network.UpdatedAt); err != nil {
		return err
	}

	return nil
}
