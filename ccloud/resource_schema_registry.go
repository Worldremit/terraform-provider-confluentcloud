package ccloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ccloud "github.com/worldremit/go-client-confluent-cloud/confluentcloud"
)

func schemaRegistryResource() *schema.Resource {
	return &schema.Resource{
		CreateContext: schemaRegistryCreate,
		ReadContext:   schemaRegistryRead,
		DeleteContext: schemaRegistryDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"environment_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Environment ID",
			},
			"region": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "where",
			},
			"service_provider": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Cloud provider",
			},
			"endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func schemaRegistryCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*ccloud.Client)

	environment := d.Get("environment_id").(string)
	region := d.Get("region").(string)
	serviceProvider := d.Get("service_provider").(string)

	log.Printf("[INFO] Creating Schema Registry %s", environment)

	reg, err := c.CreateSchemaRegistry(environment, region, serviceProvider)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(reg.ID)
	err = d.Set("endpoint", reg.Endpoint)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func schemaRegistryRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*ccloud.Client)

	environment := d.Get("environment_id").(string)
	log.Printf("[INFO] Reading Schema Registry %s", environment)

	env, err := c.GetSchemaRegistry(environment)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set("environment_id", environment)
	if err != nil {
		err = d.Set("endpoint", env.Endpoint)
	}

	return diag.FromErr(err)
}

func schemaRegistryDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[WARN] Schema registry cannot be deleted: %s", d.Id())
	return nil
}
