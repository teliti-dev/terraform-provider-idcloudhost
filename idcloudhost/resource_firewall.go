package idcloudhost

import (
	"context"
	"fmt"
	"strconv"

	idcloudhostAPI "github.com/bapung/idcloudhost-go-client-library/idcloudhost/api"
	idcloudhostFirewall "github.com/bapung/idcloudhost-go-client-library/idcloudhost/firewall"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceFirewall() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceFirewallCreate,
		ReadContext:   resourceFirewallRead,
		UpdateContext: resourceFirewallUpdate,
		DeleteContext: resourceFirewallDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"display_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"billing_account_id": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"rules": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"uuid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"direction": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
								v := val.(string)
								if v != "inbound" && v != "outbound" {
									errs = append(errs, fmt.Errorf("%q must be either 'inbound' or 'outbound', got: %s", key, v))
								}
								return
							},
						},
						"protocol": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
								v := val.(string)
								if v != "tcp" && v != "udp" && v != "icmp" && v != "any" {
									errs = append(errs, fmt.Errorf("%q must be one of 'tcp', 'udp', 'icmp', or 'any', got: %s", key, v))
								}
								return
							},
						},
						"port_start": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"port_end": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"endpoint_spec_type": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
								v := val.(string)
								if v != "any" && v != "ip_prefixes" {
									errs = append(errs, fmt.Errorf("%q must be one of 'any', or 'ip_prefixes', got: %s", key, v))
								}
								return
							},
						},
						"endpoint_spec": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"description": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
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

func resourceFirewallCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*idcloudhostAPI.APIClient)
	var diags diag.Diagnostics

	firewall := &idcloudhostFirewall.Firewall{
		DisplayName:      d.Get("display_name").(string),
		BillingAccountID: d.Get("billing_account_id").(int),
		Description:      d.Get("description").(string),
		Rules:            expandFirewallRules(d.Get("rules").([]interface{})),
	}

	firewallApi := c.Firewall
	if err := firewallApi.CreateFirewall(firewall); err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create firewall",
			Detail:   fmt.Sprint(err),
		})
		return diags
	}

	d.SetId(firewallApi.Firewall.UUID)

	resourceFirewallRead(ctx, d, m)

	return diags
}

func resourceFirewallRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*idcloudhostAPI.APIClient)
	var diags diag.Diagnostics

	uuid := d.Id()

	// Get all firewalls and find the one we're looking for
	firewallApi := c.Firewall
	if err := firewallApi.ListFirewalls(); err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to list firewalls",
			Detail:   fmt.Sprint(err),
		})
		return diags
	}

	// Find the firewall with matching UUID
	var foundFirewall *idcloudhostFirewall.Firewall
	for _, fw := range firewallApi.Firewalls {
		if fw.UUID == uuid {
			foundFirewall = &fw
			break
		}
	}

	if foundFirewall == nil {
		d.SetId("")
		return diags
	}

	if err := setFirewallResource(d, foundFirewall); err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to set firewall resource",
			Detail:   fmt.Sprint(err),
		})
		return diags
	}

	return diags
}

func resourceFirewallUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*idcloudhostAPI.APIClient)
	var diags diag.Diagnostics

	uuid := d.Id()

	if d.HasChange("rules") {
		firewall := &idcloudhostFirewall.Firewall{
			Rules: expandFirewallRules(d.Get("rules").([]interface{})),
		}

		firewallApi := c.Firewall
		if err := firewallApi.UpdateFirewall(uuid, firewall); err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Unable to update firewall",
				Detail:   fmt.Sprint(err),
			})
			return diags
		}
	}

	return resourceFirewallRead(ctx, d, m)
}

func resourceFirewallDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*idcloudhostAPI.APIClient)
	var diags diag.Diagnostics

	uuid := d.Id()

	firewallApi := c.Firewall
	if err := firewallApi.DeleteFirewall(uuid); err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to delete firewall",
			Detail:   fmt.Sprint(err),
		})
		return diags
	}

	d.SetId("")

	return diags
}

func expandFirewallRules(rules []interface{}) []idcloudhostFirewall.FirewallRule {
	if len(rules) == 0 {
		return []idcloudhostFirewall.FirewallRule{}
	}

	firewallRules := make([]idcloudhostFirewall.FirewallRule, len(rules))

	for i, rule := range rules {
		ruleMap := rule.(map[string]interface{})

		firewallRule := idcloudhostFirewall.FirewallRule{
			Direction:        ruleMap["direction"].(string),
			Protocol:         ruleMap["protocol"].(string),
			EndpointSpecType: ruleMap["endpoint_spec_type"].(string),
		}

		if v, ok := ruleMap["uuid"].(string); ok && v != "" {
			firewallRule.UUID = v
		}

		if v, ok := ruleMap["port_start"].(int); ok && v > 0 {
			firewallRule.PortStart = v
		}

		if v, ok := ruleMap["port_end"].(int); ok && v > 0 {
			firewallRule.PortEnd = v
		}

		if v, ok := ruleMap["description"].(string); ok {
			firewallRule.Description = v
		}

		if v, ok := ruleMap["endpoint_spec"].([]interface{}); ok && len(v) > 0 {
			endpointSpec := make([]string, len(v))
			for j, spec := range v {
				endpointSpec[j] = spec.(string)
			}
			firewallRule.EndpointSpec = endpointSpec
		} else {
			firewallRule.EndpointSpec = []string{}
		}

		firewallRules[i] = firewallRule
	}

	return firewallRules
}

func flattenFirewallRules(rules []idcloudhostFirewall.FirewallRule) []interface{} {
	if len(rules) == 0 {
		return []interface{}{}
	}

	flattenedRules := make([]interface{}, len(rules))

	for i, rule := range rules {
		ruleMap := map[string]interface{}{
			"uuid":               rule.UUID,
			"direction":          rule.Direction,
			"protocol":           rule.Protocol,
			"port_start":         rule.PortStart,
			"port_end":           rule.PortEnd,
			"endpoint_spec_type": rule.EndpointSpecType,
			"endpoint_spec":      rule.EndpointSpec,
			"description":        rule.Description,
		}
		flattenedRules[i] = ruleMap
	}

	return flattenedRules
}

func setFirewallResource(d *schema.ResourceData, firewall *idcloudhostFirewall.Firewall) error {
	if err := d.Set("id", strconv.Itoa(firewall.ID)); err != nil {
		return err
	}
	if err := d.Set("display_name", firewall.DisplayName); err != nil {
		return err
	}
	if err := d.Set("uuid", firewall.UUID); err != nil {
		return err
	}
	if err := d.Set("user_id", firewall.UserID); err != nil {
		return err
	}
	if err := d.Set("billing_account_id", firewall.BillingAccountID); err != nil {
		return err
	}
	if err := d.Set("description", firewall.Description); err != nil {
		return err
	}
	if err := d.Set("rules", flattenFirewallRules(firewall.Rules)); err != nil {
		return err
	}
	if err := d.Set("created_at", firewall.CreatedAt); err != nil {
		return err
	}
	if err := d.Set("updated_at", firewall.UpdatedAt); err != nil {
		return err
	}

	return nil
}
