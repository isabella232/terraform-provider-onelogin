package onelogin

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/onelogin/onelogin-go-sdk/pkg/client"
	"github.com/onelogin/onelogin-terraform-provider/resources/app"
	"github.com/onelogin/onelogin-terraform-provider/resources/app/configuration"
	"github.com/onelogin/onelogin-terraform-provider/resources/app/parameters"
	"github.com/onelogin/onelogin-terraform-provider/resources/app/provisioning"
	"github.com/onelogin/onelogin-terraform-provider/resources/app/sso"
)

func OneloginSAMLApps() *schema.Resource {
	appSchema := app.AppSchema()
	app.AddSubSchema("configuration", &appSchema, configuration.SAMLConfigurationSchema)
	app.AddSubSchema("sso", &appSchema, sso.SAMLSSOSchema)

	return &schema.Resource{
		Create: samlAppCreate,
		Read:   samlAppRead,
		Update: samlAppUpdate,
		Delete: samlAppDelete,
		Schema: appSchema,
	}
}

// samlAppCreate takes a pointer to the ResourceData Struct and a HTTP client and
// makes the POST request to OneLogin to create an samlApp with its sub-resources
func samlAppCreate(d *schema.ResourceData, m interface{}) error {
	samlApp := app.Inflate(map[string]interface{}{
		"name":                 d.Get("name"),
		"description":          d.Get("description"),
		"notes":                d.Get("notes"),
		"connector_id":         d.Get("connector_id"),
		"visible":              d.Get("visible"),
		"allow_assumed_signin": d.Get("allow_assumed_signin"),
		"parameters":           d.Get("parameters"),
		"provisioning":         d.Get("provisioning"),
		"configuration":        d.Get("configuration"),
	})
	client := m.(*client.APIClient)
	resp, appResp, err := client.Services.AppsV2.CreateApp(&samlApp)
	if err != nil {
		log.Printf("[ERROR] There was a problem creating the app!")
		log.Println(err)
		return err
	}
	log.Printf("[CREATED] Created app with %d", *(appResp.ID))
	log.Println(resp)
	d.SetId(fmt.Sprintf("%d", *(appResp.ID)))
	return samlAppRead(d, m)
}

// samlAppRead takes a pointer to the ResourceData Struct and a HTTP client and
// makes the GET request to OneLogin to read an samlApp with its sub-resources
func samlAppRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*client.APIClient)
	aid, _ := strconv.Atoi(d.Id())
	resp, app, err := client.Services.AppsV2.GetAppByID(int32(aid))
	if err != nil {
		log.Printf("[ERROR] There was a problem reading the app!")
		log.Println(err)
		return err
	}
	if app == nil {
		d.SetId("")
		return nil
	}
	log.Printf("[READ] Reading app with %d", *(app.ID))
	log.Println(resp)

	d.Set("name", app.Name)
	d.Set("visible", app.Visible)
	d.Set("description", app.Description)
	d.Set("notes", app.Notes)
	d.Set("icon_url", app.IconURL)
	d.Set("auth_method", app.AuthMethod)
	d.Set("policy_id", app.PolicyID)
	d.Set("allow_assumed_signin", app.AllowAssumedSignin)
	d.Set("tab_id", app.TabID)
	d.Set("connector_id", app.ConnectorID)
	d.Set("created_at", app.CreatedAt.String())
	d.Set("updated_at", app.UpdatedAt.String())
	d.Set("parameters", parameters.Flatten(app.Parameters))
	d.Set("provisioning", provisioning.Flatten(*app.Provisioning))
	d.Set("configuration", configuration.FlattenSAMLConfig(*app.Configuration))

	return nil
}

// samlAppUpdate takes a pointer to the ResourceData Struct and a HTTP client and
// makes the PUT request to OneLogin to update an samlApp and its sub-resources
func samlAppUpdate(d *schema.ResourceData, m interface{}) error {
	samlApp := app.Inflate(map[string]interface{}{
		"name":                 d.Get("name"),
		"description":          d.Get("description"),
		"notes":                d.Get("notes"),
		"connector_id":         d.Get("connector_id"),
		"visible":              d.Get("visible"),
		"allow_assumed_signin": d.Get("allow_assumed_signin"),
		"parameters":           d.Get("parameters"),
		"provisioning":         d.Get("provisioning"),
		"configuration":        d.Get("configuration"),
	})

	aid, _ := strconv.Atoi(d.Id())
	client := m.(*client.APIClient)

	resp, appResp, err := client.Services.AppsV2.UpdateAppByID(int32(aid), &samlApp)
	if err != nil {
		log.Printf("[ERROR] There was a problem updating the app!")
		log.Println(err)
		return err
	}
	if appResp == nil { // app must be deleted in api so remove from tf state
		d.SetId("")
		return nil
	}
	log.Printf("[UPDATED] Updated app with %d", *(appResp.ID))
	log.Println(resp)
	d.SetId(fmt.Sprintf("%d", *(appResp.ID)))
	return samlAppRead(d, m)
}

// samlAppDelete takes a pointer to the ResourceData Struct and a HTTP client and
// makes the DELETE request to OneLogin to delete an samlApp and its sub-resources
func samlAppDelete(d *schema.ResourceData, m interface{}) error {
	aid, _ := strconv.Atoi(d.Id())

	client := m.(*client.APIClient)
	resp, err := client.Services.AppsV2.DeleteApp(int32(aid))
	if err != nil {
		log.Printf("[ERROR] There was a problem creating the samlApp!")
		log.Println(err)
	} else {
		log.Printf("[DELETED] Deleted samlApp with %d", aid)
		log.Println(resp)
		d.SetId("")
	}

	return nil
}