package app

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/onelogin/onelogin-go-sdk/pkg/models"
	"github.com/onelogin/onelogin-go-sdk/pkg/oltypes"
	"github.com/onelogin/onelogin-terraform-provider/resources/app/parameters"
	"github.com/onelogin/onelogin-terraform-provider/resources/app/provisioning"
)

// App returns a key/value map of the various fields that make up an App at OneLogin.
func AppSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"visible": &schema.Schema{
			Type:     schema.TypeBool,
			Optional: true,
			Default:  true,
		},
		"description": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"notes": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
			Default:  "",
		},
		"icon_url": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
		"auth_method": &schema.Schema{
			Type:     schema.TypeInt,
			Computed: true,
		},
		"policy_id": &schema.Schema{
			Type:     schema.TypeInt,
			Computed: true,
		},
		"allow_assumed_signin": &schema.Schema{
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},
		"tab_id": &schema.Schema{
			Type:     schema.TypeInt,
			Computed: true,
		},
		"connector_id": &schema.Schema{
			Type:     schema.TypeInt,
			Required: true,
		},
		"created_at": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
		"updated_at": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
		"provisioning": &schema.Schema{
			Type:     schema.TypeSet,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: provisioning.ProvisioningSchema(),
			},
		},
		"parameters": &schema.Schema{
			Type:     schema.TypeSet,
			Optional: true,
			Elem: &schema.Resource{
				Schema: parameters.ParameterSchema(),
			},
		},
	}
}

type SubSchema func() map[string]*schema.Schema

// AddSubSchema attaches a TypeSet of 1 to the parent schema object under the node with the given key.
// SubSchema is generated by a function that takes no input and returns a map[string]*schema.Schema object
func AddSubSchema(schemaKey string, parentSchema *map[string]*schema.Schema, subSchema SubSchema) {
	(*parentSchema)[schemaKey] = &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: subSchema(),
		},
	}
}

// InflateApp takes a pointer to a ResourceData struct and uses it to construct a
// OneLogin App struct to be used in requests to OneLogin.
func InflateApp(s *map[string]interface{}) models.App {

	app := models.App{
		Name:               oltypes.String((*s)["name"].(string)),
		Description:        oltypes.String((*s)["description"].(string)),
		Notes:              oltypes.String((*s)["notes"].(string)),
		ConnectorID:        oltypes.Int32(int32((*s)["connector_id"].(int))),
		Visible:            oltypes.Bool((*s)["visible"].(bool)),
		AllowAssumedSignin: oltypes.Bool((*s)["allow_assumed_signin"].(bool)),
	}

	return app
}
