package idcloudhost

import (
	"context"
	"fmt"

	idcloudhostAPI "github.com/bapung/idcloudhost-go-client-library/idcloudhost/api"
	idcloudhostObjectStorage "github.com/bapung/idcloudhost-go-client-library/idcloudhost/objectstorage"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceObjectStorage() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceObjectStorageCreate,
		ReadContext:   resourceObjectStorageRead,
		DeleteContext: resourceObjectStorageDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"billing_account_id": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"user_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"size_bytes": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"num_objects": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"owner": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_suspended": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"acl": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"modified_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceObjectStorageCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*idcloudhostAPI.APIClient)
	var diags diag.Diagnostics

	name := d.Get("name").(string)
	billingAccountID := d.Get("billing_account_id").(int)

	osApi := c.ObjectStorage
	if err := osApi.CreateBucket(name, billingAccountID); err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create object storage bucket",
			Detail:   fmt.Sprint(err),
		})
		return diags
	}

	d.SetId(osApi.Bucket.Name)

	resourceObjectStorageRead(ctx, d, m)

	return diags
}

func resourceObjectStorageRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*idcloudhostAPI.APIClient)
	var diags diag.Diagnostics

	name := d.Id()

	osApi := c.ObjectStorage
	if err := osApi.ListBuckets(); err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to list buckets",
			Detail:   fmt.Sprint(err),
		})
		return diags
	}

	// Find the bucket with matching name
	var foundBucket *idcloudhostObjectStorage.Bucket
	for _, bucket := range osApi.Buckets {
		if bucket.Name == name {
			foundBucket = &bucket
			break
		}
	}

	if foundBucket == nil {
		d.SetId("")
		return diags
	}

	if err := setObjectStorageResource(d, foundBucket); err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to set object storage resource",
			Detail:   fmt.Sprint(err),
		})
		return diags
	}

	return diags
}

func resourceObjectStorageDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*idcloudhostAPI.APIClient)
	var diags diag.Diagnostics

	name := d.Id()

	osApi := c.ObjectStorage
	if err := osApi.DeleteBucket(name); err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to delete object storage bucket",
			Detail:   fmt.Sprint(err),
		})
		return diags
	}

	d.SetId("")

	return diags
}

func setObjectStorageResource(d *schema.ResourceData, bucket *idcloudhostObjectStorage.Bucket) error {
	if err := d.Set("name", bucket.Name); err != nil {
		return err
	}
	if err := d.Set("billing_account_id", bucket.BillingAccount); err != nil {
		return err
	}
	if err := d.Set("user_id", bucket.UserID); err != nil {
		return err
	}
	if err := d.Set("size_bytes", bucket.SizeBytes); err != nil {
		return err
	}
	if err := d.Set("num_objects", bucket.NumObjects); err != nil {
		return err
	}
	if err := d.Set("owner", bucket.Owner); err != nil {
		return err
	}
	if err := d.Set("is_suspended", bucket.IsSuspended); err != nil {
		return err
	}
	if err := d.Set("acl", bucket.ACL); err != nil {
		return err
	}
	if err := d.Set("created_at", bucket.CreatedAt); err != nil {
		return err
	}
	if err := d.Set("modified_at", bucket.ModifiedAt); err != nil {
		return err
	}

	return nil
}
