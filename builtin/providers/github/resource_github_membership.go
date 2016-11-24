package github

import (
	"github.com/google/go-github/github"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceGithubMembership() *schema.Resource {

	return &schema.Resource{
		Create: resourceGithubMembershipCreate,
		Read:   resourceGithubMembershipRead,
		Update: resourceGithubMembershipUpdate,
		Delete: resourceGithubMembershipDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"username": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"role": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateValueFunc([]string{"member", "admin"}),
				Default:      "member",
			},
		},
	}
}

func resourceGithubMembershipCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Clients).OrgClient
	n := d.Get("username").(string)
	r := d.Get("role").(string)

	membership, _, err := client.Organizations.EditOrgMembership(n, meta.(*Clients).OrgName,
		&github.Membership{Role: &r})
	if err != nil {
		return err
	}

	d.SetId(buildTwoPartID(membership.Organization.Login, membership.User.Login))

	return resourceGithubMembershipRead(d, meta)
}

func resourceGithubMembershipRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Clients).OrgClient
	_, n := parseTwoPartID(d.Id())

	membership, _, err := client.Organizations.GetOrgMembership(n, meta.(*Clients).OrgName)
	if err != nil {
		d.SetId("")
		return nil
	}

	d.Set("username", membership.User.Login)
	d.Set("role", membership.Role)
	return nil
}

func resourceGithubMembershipUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Clients).OrgClient
	n := d.Get("username").(string)
	r := d.Get("role").(string)

	membership, _, err := client.Organizations.EditOrgMembership(n, meta.(*Clients).OrgName, &github.Membership{
		Role: &r,
	})
	if err != nil {
		return err
	}
	d.SetId(buildTwoPartID(membership.Organization.Login, membership.User.Login))

	return nil
}

func resourceGithubMembershipDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Clients).OrgClient
	n := d.Get("username").(string)

	_, err := client.Organizations.RemoveOrgMembership(n, meta.(*Clients).OrgName)

	return err
}
