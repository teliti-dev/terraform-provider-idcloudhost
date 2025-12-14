package idcloudhost

import (
	"context"
	"fmt"

	idcloudhostAPI "github.com/bapung/idcloudhost-go-client-library/idcloudhost/api"
	idcloudhostLoadBalancer "github.com/bapung/idcloudhost-go-client-library/idcloudhost/loadbalancer"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceLoadBalancer() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceLoadBalancerCreate,
		ReadContext:   resourceLoadBalancerRead,
		UpdateContext: resourceLoadBalancerUpdate,
		DeleteContext: resourceLoadBalancerDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"display_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"billing_account_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"network_uuid": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"reserve_public_ip": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				ForceNew: true,
			},
			"targets": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"target_type": {
							Type:     schema.TypeString,
							Required: true,
						},
						"target_uuid": {
							Type:     schema.TypeString,
							Required: true,
						},
						"target_ip_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"created_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"forwarding_rules": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"source_port": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"target_port": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"protocol": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"uuid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"created_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"uuid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"user_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"private_address": {
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

func resourceLoadBalancerCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*idcloudhostAPI.APIClient)
	var diags diag.Diagnostics

	request := &idcloudhostLoadBalancer.CreateLoadBalancerRequest{
		DisplayName:      d.Get("display_name").(string),
		BillingAccountID: d.Get("billing_account_id").(int),
		NetworkUUID:      d.Get("network_uuid").(string),
		ReservePublicIP:  d.Get("reserve_public_ip").(bool),
	}

	// Add targets if provided
	if v, ok := d.GetOk("targets"); ok {
		request.Targets = expandLoadBalancerTargets(v.([]interface{}))
	}

	// Add forwarding rules if provided
	if v, ok := d.GetOk("forwarding_rules"); ok {
		request.Rules = expandLoadBalancerRules(v.([]interface{}))
	}

	lbApi := c.LoadBalancer
	if err := lbApi.CreateLoadBalancer(request); err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create load balancer",
			Detail:   fmt.Sprint(err),
		})
		return diags
	}

	d.SetId(lbApi.LoadBalancer.UUID)

	resourceLoadBalancerRead(ctx, d, m)

	return diags
}

func resourceLoadBalancerRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*idcloudhostAPI.APIClient)
	var diags diag.Diagnostics

	uuid := d.Id()

	lbApi := c.LoadBalancer
	if err := lbApi.GetLoadBalancer(uuid); err != nil {
		// If load balancer not found, remove from state
		d.SetId("")
		return diags
	}

	if err := setLoadBalancerResource(d, lbApi.LoadBalancer); err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to set load balancer resource",
			Detail:   fmt.Sprint(err),
		})
		return diags
	}

	return diags
}

func resourceLoadBalancerUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*idcloudhostAPI.APIClient)
	var diags diag.Diagnostics

	uuid := d.Id()
	lbApi := c.LoadBalancer

	// Update display name
	if d.HasChange("display_name") {
		displayName := d.Get("display_name").(string)
		if err := lbApi.RenameLoadBalancer(uuid, displayName); err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Unable to rename load balancer",
				Detail:   fmt.Sprint(err),
			})
			return diags
		}
	}

	// Update billing account
	if d.HasChange("billing_account_id") {
		billingAccountID := d.Get("billing_account_id").(int)
		if err := lbApi.ChangeBillingAccount(uuid, billingAccountID); err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Unable to change billing account",
				Detail:   fmt.Sprint(err),
			})
			return diags
		}
	}

	// Handle targets changes
	if d.HasChange("targets") {
		old, new := d.GetChange("targets")
		oldTargets := expandLoadBalancerTargets(old.([]interface{}))
		newTargets := expandLoadBalancerTargets(new.([]interface{}))

		// Get current state to have target UUIDs
		if err := lbApi.GetLoadBalancer(uuid); err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Unable to get load balancer state",
				Detail:   fmt.Sprint(err),
			})
			return diags
		}

		// Build a map of existing targets by their UUID
		existingTargetsMap := make(map[string]idcloudhostLoadBalancer.Target)
		for _, target := range lbApi.LoadBalancer.Targets {
			key := fmt.Sprintf("%s-%s", target.TargetType, target.TargetUUID)
			existingTargetsMap[key] = target
		}

		// Build maps for comparison
		oldMap := make(map[string]bool)
		newMap := make(map[string]bool)

		for _, t := range oldTargets {
			key := fmt.Sprintf("%s-%s", t.TargetType, t.TargetUUID)
			oldMap[key] = true
		}

		for _, t := range newTargets {
			key := fmt.Sprintf("%s-%s", t.TargetType, t.TargetUUID)
			newMap[key] = true
		}

		// Remove targets that are in old but not in new
		for key := range oldMap {
			if !newMap[key] {
				if target, exists := existingTargetsMap[key]; exists {
					if err := lbApi.RemoveTarget(uuid, target.TargetUUID); err != nil {
						diags = append(diags, diag.Diagnostic{
							Severity: diag.Warning,
							Summary:  "Unable to remove target",
							Detail:   fmt.Sprint(err),
						})
					}
				}
			}
		}

		// Add targets that are in new but not in old
		for key := range newMap {
			if !oldMap[key] {
				// Find the target in newTargets
				for _, t := range newTargets {
					targetKey := fmt.Sprintf("%s-%s", t.TargetType, t.TargetUUID)
					if targetKey == key {
						if err := lbApi.AddTarget(uuid, &t); err != nil {
							diags = append(diags, diag.Diagnostic{
								Severity: diag.Error,
								Summary:  "Unable to add target",
								Detail:   fmt.Sprint(err),
							})
							return diags
						}
						break
					}
				}
			}
		}
	}

	// Handle forwarding rules changes
	if d.HasChange("forwarding_rules") {
		old, new := d.GetChange("forwarding_rules")
		oldRules := expandLoadBalancerRules(old.([]interface{}))
		newRules := expandLoadBalancerRules(new.([]interface{}))

		// Get current state to have rule UUIDs
		if err := lbApi.GetLoadBalancer(uuid); err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Unable to get load balancer state",
				Detail:   fmt.Sprint(err),
			})
			return diags
		}

		// Build a map of existing rules by their ports
		existingRulesMap := make(map[string]idcloudhostLoadBalancer.ForwardingRule)
		for _, rule := range lbApi.LoadBalancer.ForwardingRules {
			key := fmt.Sprintf("%d-%d", rule.SourcePort, rule.TargetPort)
			existingRulesMap[key] = rule
		}

		// Build maps for comparison
		oldMap := make(map[string]bool)
		newMap := make(map[string]bool)

		for _, r := range oldRules {
			key := fmt.Sprintf("%d-%d", r.SourcePort, r.TargetPort)
			oldMap[key] = true
		}

		for _, r := range newRules {
			key := fmt.Sprintf("%d-%d", r.SourcePort, r.TargetPort)
			newMap[key] = true
		}

		// Remove rules that are in old but not in new
		for key := range oldMap {
			if !newMap[key] {
				if rule, exists := existingRulesMap[key]; exists {
					if err := lbApi.RemoveRule(uuid, rule.UUID); err != nil {
						diags = append(diags, diag.Diagnostic{
							Severity: diag.Warning,
							Summary:  "Unable to remove forwarding rule",
							Detail:   fmt.Sprint(err),
						})
					}
				}
			}
		}

		// Add rules that are in new but not in old
		for key := range newMap {
			if !oldMap[key] {
				// Find the rule in newRules
				for _, r := range newRules {
					ruleKey := fmt.Sprintf("%d-%d", r.SourcePort, r.TargetPort)
					if ruleKey == key {
						if err := lbApi.AddRule(uuid, &r); err != nil {
							diags = append(diags, diag.Diagnostic{
								Severity: diag.Error,
								Summary:  "Unable to add forwarding rule",
								Detail:   fmt.Sprint(err),
							})
							return diags
						}
						break
					}
				}
			}
		}
	}

	return resourceLoadBalancerRead(ctx, d, m)
}

func resourceLoadBalancerDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*idcloudhostAPI.APIClient)
	var diags diag.Diagnostics

	uuid := d.Id()

	lbApi := c.LoadBalancer
	if err := lbApi.DeleteLoadBalancer(uuid); err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to delete load balancer",
			Detail:   fmt.Sprint(err),
		})
		return diags
	}

	d.SetId("")

	return diags
}

func expandLoadBalancerTargets(targets []interface{}) []idcloudhostLoadBalancer.CreateTargetRequest {
	if len(targets) == 0 {
		return []idcloudhostLoadBalancer.CreateTargetRequest{}
	}

	expandedTargets := make([]idcloudhostLoadBalancer.CreateTargetRequest, len(targets))

	for i, target := range targets {
		targetMap := target.(map[string]interface{})

		expandedTargets[i] = idcloudhostLoadBalancer.CreateTargetRequest{
			TargetType: targetMap["target_type"].(string),
			TargetUUID: targetMap["target_uuid"].(string),
		}
	}

	return expandedTargets
}

func expandLoadBalancerRules(rules []interface{}) []idcloudhostLoadBalancer.CreateRuleRequest {
	if len(rules) == 0 {
		return []idcloudhostLoadBalancer.CreateRuleRequest{}
	}

	expandedRules := make([]idcloudhostLoadBalancer.CreateRuleRequest, len(rules))

	for i, rule := range rules {
		ruleMap := rule.(map[string]interface{})

		expandedRules[i] = idcloudhostLoadBalancer.CreateRuleRequest{
			SourcePort: ruleMap["source_port"].(int),
			TargetPort: ruleMap["target_port"].(int),
		}
	}

	return expandedRules
}

func flattenLoadBalancerTargets(targets []idcloudhostLoadBalancer.Target) []interface{} {
	if len(targets) == 0 {
		return []interface{}{}
	}

	flattenedTargets := make([]interface{}, len(targets))

	for i, target := range targets {
		targetMap := map[string]interface{}{
			"target_type":       target.TargetType,
			"target_uuid":       target.TargetUUID,
			"target_ip_address": target.TargetIPAddress,
			"created_at":        target.CreatedAt,
		}
		flattenedTargets[i] = targetMap
	}

	return flattenedTargets
}

func flattenLoadBalancerRules(rules []idcloudhostLoadBalancer.ForwardingRule) []interface{} {
	if len(rules) == 0 {
		return []interface{}{}
	}

	flattenedRules := make([]interface{}, len(rules))

	for i, rule := range rules {
		ruleMap := map[string]interface{}{
			"source_port": rule.SourcePort,
			"target_port": rule.TargetPort,
			"protocol":    rule.Protocol,
			"uuid":        rule.UUID,
			"created_at":  rule.CreatedAt,
		}
		flattenedRules[i] = ruleMap
	}

	return flattenedRules
}

func setLoadBalancerResource(d *schema.ResourceData, lb *idcloudhostLoadBalancer.LoadBalancer) error {
	if err := d.Set("uuid", lb.UUID); err != nil {
		return err
	}
	if err := d.Set("display_name", lb.DisplayName); err != nil {
		return err
	}
	if err := d.Set("user_id", lb.UserID); err != nil {
		return err
	}
	if err := d.Set("billing_account_id", lb.BillingAccountID); err != nil {
		return err
	}
	if err := d.Set("network_uuid", lb.NetworkUUID); err != nil {
		return err
	}
	if err := d.Set("private_address", lb.PrivateAddress); err != nil {
		return err
	}
	if err := d.Set("targets", flattenLoadBalancerTargets(lb.Targets)); err != nil {
		return err
	}
	if err := d.Set("forwarding_rules", flattenLoadBalancerRules(lb.ForwardingRules)); err != nil {
		return err
	}
	if err := d.Set("created_at", lb.CreatedAt); err != nil {
		return err
	}
	if err := d.Set("updated_at", lb.UpdatedAt); err != nil {
		return err
	}

	return nil
}
