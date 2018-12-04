package anypoint

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/tech-nico/terraform-provider-anypoint/anypoint/sdk"
	"log"
)

func resourceBusinessGroup() *schema.Resource {

	return &schema.Resource{
		Create: resourceBGCreate,
		Read:   resourceBGRead,
		Update: resourceBGUpdate,
		Delete: resourceBGDelete,
		//Exists: resourceBGExists,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Description: "The name of the (new) business group",
				Required:    true,
			},
			"parent_path": &schema.Schema{
				Type:        schema.TypeString,
				Description: "The path to parent. Example: Company\\Retail\\APIs",
				Optional: true,
			},
			"parent_org_id": &schema.Schema{
				Type:        schema.TypeString,
				Description: "The ID of the parent org",
				Optional: true,
			},
			"owner_username": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Username of the business group's Owner. Required only if the BG does not exist yet. Defaults to current username.",
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ANYPOINT_USERNAME", nil),
			},
			"can_create_sub_orgs": &schema.Schema{
				Type:        schema.TypeBool,
				Description: "Whether or not the owner of the new org can create sub organizations",
				Optional: true,
			},
			"can_create_environments": &schema.Schema{
				Type:        schema.TypeBool,
				Description: "Whether the org owner can create environments",
				Optional: true,
			},
			"production_vcores": &schema.Schema{
				Type:     schema.TypeFloat,
				Optional: true,
			},
			"sandbox_vcores": &schema.Schema{
				Type:     schema.TypeFloat,
				Optional: true,
			},
			"design_vcores": &schema.Schema{
				Type:     schema.TypeFloat,
				Optional: true,
			},
			"static_ips": &schema.Schema{
				Type:     schema.TypeFloat,
				Optional: true,
			},
			"vpcs": &schema.Schema{
				Type:     schema.TypeFloat,
				Optional: true,
			},
			"load_balancers": &schema.Schema{
				Type:     schema.TypeFloat,
				Optional: true,
			},
			"vpns": &schema.Schema{
				Type:     schema.TypeFloat,
				Optional: true,
			},
		},
	}
}


func resourceBGCreate(d *schema.ResourceData, conf interface{}) error {
	newBgName := d.Get("name").(string)
	parentPath := d.Get("parent_path").(string)
	theConf := conf.(*Config)
	client := theConf.AnypointClient

	parentId, err := client.Auth.FindBusinessGroup(parentPath)
	if err != nil {
		return err
	}

	ents := getEntitlementsFromData(d)

	ownerUsername := conf.(*Config).Username

	if val, isSet := d.GetOk("owner_username"); isSet {
		ownerUsername = val.(string)
	}

	ownerUser, err := client.Auth.FindUserByUsername(parentId, ownerUsername)
	if err != nil {
		return err
	}

	newBG, err := client.Auth.CreateBusinessGroup(ownerUser.Username, parentId, newBgName, ents)

	if err != nil {
		return err
	}

	d.SetId(newBG.ID)

	return nil
}

func resourceBGRead(d *schema.ResourceData, conf interface{}) error {
	apClient := conf.(*Config).AnypointClient
	path := d.Get("parent_path").(string)
	bgID, err := apClient.Auth.FindBusinessGroup(path)

	if err != nil {
		log.Printf("Error while searching for business group %s: %s", path, err)
		return err
	}

	if bgID == "" {
		groups := apClient.Auth.CreateBusinessGroupPath(path)
		if len(groups) > 1 {
			path := ""
			for idx, elem := range groups {
				path = path + elem
				if idx < (len(groups) - 1) {
					path = path + "\\"
				} else {
					break
				}
			}
			if id, err := apClient.Auth.FindBusinessGroup(path); id == "" && err != nil {
				return fmt.Errorf("Parent business group [%s] does not exits", path)
			}
		}

	}

	d.SetId(bgID)

	return nil
}

func getEntitlementsFromData(data *schema.ResourceData) sdk.Entitlements {
	ents := sdk.Entitlements{}

	if val, isSet := data.GetOk("can_create_sub_orgs"); isSet {
		ents.CreateSubOrgs = val.(bool)
	}

	if val, isSet := data.GetOk("can_create_environments"); isSet {
		ents.CreateEnvironments = val.(bool)
	}

	if val, isSet := data.GetOk("proudction_vcores"); isSet {
		ents.ProductionVCores = sdk.EntitlementStatus{Assigned: val.(float64)}
	}

	if val, isSet := data.GetOk("design_vcores"); isSet {
		ents.DesignVCores = sdk.EntitlementStatus{Assigned: val.(float64)}
	}

	if val, isSet := data.GetOk("load_balancers"); isSet {
		ents.LoadBalancer = sdk.EntitlementStatus{Assigned: val.(float64)}
	}

	if val, isSet := data.GetOk("production_vcores"); isSet {
		ents.ProductionVCores = sdk.EntitlementStatus{Assigned: val.(float64)}
	}

	if val, isSet := data.GetOk("static_ips"); isSet {
		ents.StaticIPs = sdk.EntitlementStatus{Assigned: val.(float64)}
	}

	if val, isSet := data.GetOk("vpcs"); isSet {
		ents.VPCs = sdk.EntitlementStatus{Assigned: val.(float64)}
	}

	if val, isSet := data.GetOk("vpns"); isSet {
		ents.VPNs = sdk.EntitlementStatus{Assigned: val.(float64)}
	}

	return ents
}

func resourceBGDelete(d *schema.ResourceData, conf interface{}) error {
	//apClient := conf.(*Config).AnypointClient

	return nil
}

func resourceBGUpdate(d *schema.ResourceData, conf interface{}) error {
	//apClient := conf.(*Config).AnypointClient

	return nil
}

func resourceBGExists(d *schema.ResourceData, conf interface{}) (bool, error) {
	apClient := conf.(*Config).AnypointClient
	bgID := d.Id()
	org, err := apClient.Auth.GetBusinessGroupHierarchy(bgID)

	if err != nil {
		return false, err
	}

	if org.ID != "" {
		return true, nil
	}

	return false, nil
}
