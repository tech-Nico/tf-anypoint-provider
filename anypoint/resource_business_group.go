package anypoint

import "github.com/hashicorp/terraform/helper/schema"

func resourceBusinessGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceBGCreate,
		Read:   resourceBGRead,
		Update: resourceBGUpdate,
		Delete: resourceBGDelete,
		Exists: resourceBGExists,

		Schema: map[string]*schema.Schema{
			"path_to_bg": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceBGCreate(d *schema.ResourceData, conf interface{}) error {
	name := d.Get("path_to_bg").(string)
	theConf := conf.(Config)
	client := theConf.AnypointClient

	bgId, err := client.Auth.FindBusinessGroup(name)
	if err != nil {
		//I need to create a new BG
	}

	//At this point I should create the BG (if it doesn't exist)
	d.SetId(bgId)

	return nil
}

func resourceBGRead(d *schema.ResourceData, conf interface{}) error {
	apClient := conf.(*Config).AnypointClient
	path := d.Get("path)_to_bg").(string)
	bgID, err := apClient.Auth.FindBusinessGroup(path)

	if err != nil {
		return err
	}

	d.SetId(bgID)

	return nil
}

func resourceBGDelete(d *schema.ResourceData, conf interface{}) error {
	return nil
}

func resourceBGUpdate(d *schema.ResourceData, conf interface{}) error {
	return nil
}

func resourceBGExists(d *schema.ResourceData, conf interface{}) (bool, error) {
	return false, nil
}
