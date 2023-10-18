package panos

import (
	"log"
	"strings"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/objs/schedule"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// Data source (listing).

func dataSourceScheduleObjects() *schema.Resource {
	s := listingSchema()
	s["vsys"] = vsysSchema("vsys1")
	s["device_group"] = deviceGroupSchema()

	return &schema.Resource{
		Read: dataSourceSchedulebjectsRead,

		Schema: s,
	}
}

func dataSourceSchedulebjectsRead(d *schema.ResourceData, meta interface{}) error {
	var err error
	var listing []string
	var id string

	switch con := meta.(type) {
	case *pango.Firewall:
		id = d.Get("vsys").(string)
		listing, err = con.Objects.Schedule.GetList(id)
	case *pango.Panorama:
		id = d.Get("device_group").(string)
		listing, err = con.Objects.Schedule.GetList(id)
	}

	if err != nil {
		return err
	}

	d.SetId(id)
	saveListing(d, listing)
	return nil
}

// Data source.
func dataSourceScheduleObject() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceScheduleObjectRead,

		Schema: scheduleObjectSchema(false, true, nil),
	}
}

func dataSourceScheduleObjectRead(d *schema.ResourceData, meta interface{}) error {
	var err error
	var id string
	var o schedule.Entry
	name := d.Get("name").(string)

	switch con := meta.(type) {
	case *pango.Firewall:
		vsys := d.Get("vsys").(string)
		id = buildScheduleObjectId(vsys, name)
		o, err = con.Objects.Schedule.Get(vsys, name)
	case *pango.Panorama:
		dg := d.Get("device_group").(string)
		id = buildScheduleObjectId(dg, name)
		o, err = con.Objects.Schedule.Get(dg, name)
	}

	if err != nil {
		if isObjectNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	d.SetId(id)
	saveScheduleObject(d, o)

	return nil
}

// Resource.
func resourceScheduleObject() *schema.Resource {
	return &schema.Resource{
		Create: createScheduleObject,
		Read:   readScheduleObject,
		Update: updateScheduleObject,
		Delete: deleteScheduleObject,

		SchemaVersion: 1,
		StateUpgraders: []schema.StateUpgrader{
			{
				Version: 0,
				Type: (&schema.Resource{
					Schema: scheduleObjectSchema(true, true, []string{"device_group"}),
				}).CoreConfigSchema().ImpliedType(),
				Upgrade: scheduleObjectUpgradeV0,
			},
		},

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: scheduleObjectSchema(true, true, nil),
	}
}

func resourcePanoramaScheduleObject() *schema.Resource {
	return &schema.Resource{
		Create: createScheduleObject,
		Read:   readScheduleObject,
		Update: updateScheduleObject,
		Delete: deleteScheduleObject,

		SchemaVersion: 1,
		StateUpgraders: []schema.StateUpgrader{
			{
				Version: 0,
				Type: (&schema.Resource{
					Schema: scheduleObjectSchema(true, true, []string{"vsys"}),
				}).CoreConfigSchema().ImpliedType(),
				Upgrade: scheduleObjectUpgradeV0,
			},
		},

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: scheduleObjectSchema(true, true, nil),
	}
}

func scheduleObjectUpgradeV0(raw map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
	if _, ok := raw["vsys"]; ok {
		raw["device_group"] = "shared"
	}
	if _, ok := raw["device_group"]; ok {
		raw["vsys"] = "vsys1"
	}

	return raw, nil
}

func createScheduleObject(d *schema.ResourceData, meta interface{}) error {
	var err error
	var id string
	o := loadScheduleObject(d)

	switch con := meta.(type) {
	case *pango.Firewall:
		vsys := d.Get("vsys").(string)
		id = buildScheduleObjectId(vsys, o.Name)
		err = con.Objects.Schedule.Set(vsys, o)
	case *pango.Panorama:
		dg := d.Get("device_group").(string)
		id = buildScheduleObjectId(dg, o.Name)
		err = con.Objects.Schedule.Set(dg, o)
	}

	if err != nil {
		return err
	}

	d.SetId(id)
	return readScheduleObject(d, meta)
}

func readScheduleObject(d *schema.ResourceData, meta interface{}) error {
	var err error
	var o schedule.Entry

	switch con := meta.(type) {
	case *pango.Firewall:
		vsys, name := parseScheduleObjectId(d.Id())
		o, err = con.Objects.Schedule.Get(vsys, name)
		d.Set("vsys", vsys)
		d.Set("device_group", "shared")
	case *pango.Panorama:
		dg, name := parseScheduleObjectId(d.Id())
		o, err = con.Objects.Schedule.Get(dg, name)
		d.Set("vsys", "vsys1")
		d.Set("device_group", dg)
	}

	if err != nil {
		if isObjectNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	saveScheduleObject(d, o)
	return nil
}

func updateScheduleObject(d *schema.ResourceData, meta interface{}) error {
	o := loadScheduleObject(d)

	switch con := meta.(type) {
	case *pango.Firewall:
		vsys := d.Get("vsys").(string)
		lo, err := con.Objects.Schedule.Get(vsys, o.Name)
		if err != nil {
			return err
		}
		lo.Copy(o)
		if err = con.Objects.Schedule.Edit(vsys, lo); err != nil {
			return err
		}
	case *pango.Panorama:
		dg := d.Get("device_group").(string)
		lo, err := con.Objects.Schedule.Get(dg, o.Name)
		if err != nil {
			return err
		}
		lo.Copy(o)
		if err = con.Objects.Schedule.Edit(dg, lo); err != nil {
			return err
		}
	}

	return readScheduleObject(d, meta)
}

func deleteScheduleObject(d *schema.ResourceData, meta interface{}) error {
	var err error

	switch con := meta.(type) {
	case *pango.Firewall:
		vsys, name := parseScheduleObjectId(d.Id())
		err = con.Objects.Schedule.Delete(vsys, name)
	case *pango.Panorama:
		dg, name := parseScheduleObjectId(d.Id())
		err = con.Objects.Schedule.Delete(dg, name)
	}

	if err != nil {
		if !isObjectNotFound(err) {
			return err
		}
	}

	d.SetId("")
	return nil
}

// Schema handling.
func scheduleObjectSchema(isResource, forceNew bool, rmKeys []string) map[string]*schema.Schema {
	ans := map[string]*schema.Schema{
		"device_group": deviceGroupSchema(),
		"vsys":         vsysSchema("vsys1"),
		"name": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: forceNew,
		},

		"type": {
			Type:        schema.TypeString,
			Description: "Schedule once or daily or weekly",
			Required:    true,
		},
		"value": {
			Type:        schema.TypeList,
			Required:    true,
			Description: "Time to schedule",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
	}

	if !isResource {
		computed(ans, "", []string{"vsys", "device_group", "name"})
	}

	for _, rmKey := range rmKeys {
		delete(ans, rmKey)
	}

	return ans
}

func loadScheduleObject(d *schema.ResourceData) schedule.Entry {
	return schedule.Entry{
		Name:  d.Get("name").(string),
		Type:  d.Get("type").(string),
		Value: asStringList(d.Get("value").([]interface{})),
	}
}

func saveScheduleObject(d *schema.ResourceData, o schedule.Entry) {
	d.Set("name", o.Name)
	d.Set("type", o.Type)
	if err := d.Set("value", o.Value); err != nil {
		log.Printf("[WARN] Error setting 'Value' field (schedule time) for %q: %s", d.Id(), err)
	}
}

// Id functions.
func parseScheduleObjectId(v string) (string, string) {
	t := strings.Split(v, IdSeparator)
	return t[0], t[1]
}

func buildScheduleObjectId(a, b string) string {
	return strings.Join([]string{a, b}, IdSeparator)
}
