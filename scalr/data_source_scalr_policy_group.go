package scalr

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/scalr/go-scalr"
)

func dataSourceScalrPolicyGroup() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceScalrPolicyGroupRead,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:         schema.TypeString,
				Optional:     true,
				AtLeastOneOf: []string{"name"},
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"error_message": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"opa_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vcs_repo": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"identifier": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"branch": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"path": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"account_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				DefaultFunc: scalrAccountIDDefaultFunc,
			},
			"vcs_provider_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"policies": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"enforced_level": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"environments": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceScalrPolicyGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	scalrClient := meta.(*scalr.Client)

	// required fields
	pgID := d.Get("id").(string)
	name := d.Get("name").(string)
	accountID := d.Get("account_id").(string)

	options := scalr.PolicyGroupListOptions{
		Account: accountID,
		Include: "policies",
	}

	if pgID != "" {
		options.PolicyGroup = pgID
	}

	if name != "" {
		options.Name = name
	}

	log.Printf("[DEBUG] Read configuration of policy group with ID '%s', name '%s' and account_id '%s'", pgID, name, accountID)

	pgl, err := scalrClient.PolicyGroups.List(ctx, options)
	if err != nil {
		return diag.Errorf("error retrieving policy group: %v", err)
	}

	if pgl.TotalCount == 0 {
		return diag.Errorf("policy group with ID '%s', name '%s' and account_id '%s' not found", pgID, name, accountID)
	}

	pg := pgl.Items[0]

	// Update the configuration.
	_ = d.Set("name", pg.Name)
	_ = d.Set("status", pg.Status)
	_ = d.Set("error_message", pg.ErrorMessage)
	_ = d.Set("opa_version", pg.OpaVersion)

	if pg.VcsProvider != nil {
		_ = d.Set("vcs_provider_id", pg.VcsProvider.ID)
	}

	var vcsRepo []interface{}
	if pg.VCSRepo != nil {
		vcsConfig := map[string]interface{}{
			"identifier": pg.VCSRepo.Identifier,
			"branch":     pg.VCSRepo.Branch,
			"path":       pg.VCSRepo.Path,
		}
		vcsRepo = append(vcsRepo, vcsConfig)
	}
	_ = d.Set("vcs_repo", vcsRepo)

	var policies []map[string]interface{}
	if len(pg.Policies) != 0 {
		for _, policy := range pg.Policies {
			policies = append(policies, map[string]interface{}{
				"name":           policy.Name,
				"enabled":        policy.Enabled,
				"enforced_level": policy.EnforcementLevel,
			})
		}
	}
	_ = d.Set("policies", policies)

	var envs []string
	if len(pg.Environments) != 0 {
		for _, env := range pg.Environments {
			envs = append(envs, env.ID)
		}
	}
	_ = d.Set("environments", envs)

	d.SetId(pg.ID)

	return nil
}
