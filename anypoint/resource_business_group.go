package anypoint

import "github.com/hashicorp/terraform/helper/schema"

func resourceBusinessGroup() *schema.Resource{
	return &schema.Resource{
		Create: resourceBGCreate,
		Read: resourceBGRead,
		Update: resourceBGUpdate,
		Delete: resourceBGDelete,

		Schema: map[string] *schema.Schema{
			"name" : &schema.Schema{
				Type: schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceBGCreate(d *schema.ResourceData, m interface{}) error {
	name := d.Get("name").(string)
	//At this point I should create the BG (if it doesn't exist)
	d.SetId(name)

	return nil
}

func resourceBGRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceBGDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceBGUpdate(d *schema.ResourceData, m interface{}) error {
	return nil
}