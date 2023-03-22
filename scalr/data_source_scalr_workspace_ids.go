package scalr

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/scalr/go-scalr"
)

func dataSourceScalrWorkspaceIDs() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceScalrWorkspaceIDsRead,

		Schema: map[string]*schema.Schema{
			"environment_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"names": {
				Type:         schema.TypeList,
				Elem:         &schema.Schema{Type: schema.TypeString},
				Optional:     true,
				AtLeastOneOf: []string{"tag_ids"},
			},
			"tag_ids": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
			"exact_match": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"ids": {
				Type:     schema.TypeMap,
				Computed: true,
			},
		},
	}
}

func dataSourceScalrWorkspaceIDsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	scalrClient := meta.(*scalr.Client)

	// Get the environment_id.
	environmentID := d.Get("environment_id").(string)
	exact := d.Get("exact_match").(bool)
	var id string
	// Create a map to store workspace IDs
	ids := make(map[string]string, 0)
	options := scalr.WorkspaceListOptions{Environment: &environmentID}

	// Create a map with all the names we are looking for.
	names := make(map[string]bool)
	if namesI, ok := d.GetOk("names"); ok {
		for _, name := range namesI.([]interface{}) {
			id += name.(string)
			names[name.(string)] = true
		}
	}

	if tagIDsI, ok := d.GetOk("tag_ids"); ok {
		tagIDs := make([]string, 0)
		for _, t := range tagIDsI.(*schema.Set).List() {
			id += t.(string)
			tagIDs = append(tagIDs, t.(string))
		}
		if len(tagIDs) > 0 {
			options.Tag = scalr.String("in:" + strings.Join(tagIDs, ","))
		}
	}

	for {
		wl, err := scalrClient.Workspaces.List(ctx, options)
		if err != nil {
			return diag.Errorf("Error retrieving workspaces: %v", err)
		}

		for _, w := range wl.Items {
			if len(names) > 0 {
				if names["*"] || (exact && names[w.Name]) || (!exact && matchesPattern(w.Name, names)) {
					ids[w.Name] = w.ID
				}
			} else {
				ids[w.Name] = w.ID
			}
		}

		// Exit the loop when we've seen all pages.
		if wl.CurrentPage >= wl.TotalPages {
			break
		}

		// Update the page number to get the next page.
		options.PageNumber = wl.NextPage
	}

	_ = d.Set("ids", ids)
	d.SetId(fmt.Sprintf("%s/%d", environmentID, schema.HashString(id)))

	return nil
}
