package resources

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/solacecommunity/solace-terraform-provider/provider"
)

// Main resource definition for Jndi Queue entities
func ResourceCliUser() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,
		// List of supported configuration fields for your resource
		Schema: schemaCliUser(),
		// Provider CRUD functions
		CreateContext: createCliUser,
		ReadContext:   readCliUser,
		UpdateContext: updateCliUser,
		DeleteContext: deleteCliUser,
	}
}

func schemaCliUser() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		CLI_USERNAME: {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		CLI_PASSWORD: {
			Type:      schema.TypeString,
			Required:  true,
			Sensitive: true,
		},
		CLI_USERTYPE: {
			Type:     schema.TypeString,
			Computed: true,
		},
		CLI_GLOBAL_ACCESS_LEVEL: {
			Type:         schema.TypeString,
			Default:      "none",
			Optional:     true,
			ValidateFunc: validation.StringInSlice([]string{"none", "admin", "read-only", "read-write"}, false),
		},
		CLI_MSGVPN_DEFAULT_ACCESS_LEVEL: {
			Type:         schema.TypeString,
			Default:      "none",
			Optional:     true,
			ValidateFunc: validation.StringInSlice([]string{"none", "read-only", "read-write"}, false),
		},
		CLI_MSGVPN_ACCESS_LEVEL_EXCEPTION: {
			Type:     schema.TypeMap,
			Optional: true,
			Elem: &schema.Schema{
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"none", "read-only", "read-write"}, false),
			},
		},
	}
}

func createCliUser(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := meta.(*provider.ProviderState).SempV1Client
	username := d.Get(CLI_USERNAME).(string)
	password := d.Get(CLI_PASSWORD).(string)

	err := client.CreateCliUser(ctx, username, password)
	if err != nil {
		return provider.AppendError(diags, fmt.Errorf("cannot create cli user %s: %s", username, err))
	}

	err = client.SetCliUserGlobalAccessLevel(ctx, username, d.Get(CLI_GLOBAL_ACCESS_LEVEL).(string))
	if err != nil {
		return provider.AppendError(diags, fmt.Errorf("cannot set global access-level for user %s: %s", username, err))
	}
	err = client.SetCliUserMsgVpnDefaultAccessLevel(ctx, username, d.Get(CLI_MSGVPN_DEFAULT_ACCESS_LEVEL).(string))
	if err != nil {
		return provider.AppendError(diags, fmt.Errorf("cannot set msgvpn default access-level for user %s: %s", username, err))
	}

	if exceptions, ok := d.Get(CLI_MSGVPN_ACCESS_LEVEL_EXCEPTION).(map[string]interface{}); ok {
		for vpn, accessLevel := range exceptions {
			err = client.CreateCliUserVpnAccessLevelException(ctx, username, vpn, accessLevel.(string))
			if err != nil {
				return provider.AppendError(diags, fmt.Errorf("cannot create msgvpn access-level exception for user %s and vpn %s: %s", username, vpn, err))
			}
		}
	} else {
		return provider.AppendError(diags, fmt.Errorf("cannot read CLI_MSGVPN_ACCESS_LEVEL_EXCEPTIONs"))
	}

	d.SetId("sol_cliu_" + username)
	return append(diags, readCliUser(ctx, d, meta)...)
}

func readCliUser(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := meta.(*provider.ProviderState).SempV1Client
	username := d.Get(CLI_USERNAME).(string)

	result, err := client.GetCliUser(ctx, username)
	if err != nil {
		return provider.AppendError(diags, err)
	}

	if len(result.Rpc.Show.Username.Users.User) == 1 {
		user := result.Rpc.Show.Username.Users.User[0]
		d.Set(CLI_USERNAME, user.Name)
		d.Set(CLI_USERTYPE, user.UserType)
		d.Set(CLI_GLOBAL_ACCESS_LEVEL, user.GlobalAccessLevel)
		d.Set(CLI_MSGVPN_DEFAULT_ACCESS_LEVEL, user.DefaultVpnAccessLevel)

		exceptions := map[string]string{}
		for _, ex := range user.VpnAccessLevelException {
			exceptions[ex.VpnName] = ex.AccessLevel
		}
		d.Set(CLI_MSGVPN_ACCESS_LEVEL_EXCEPTION, exceptions)

	} else {
		return provider.AppendError(diags, fmt.Errorf("query requires exactly one returned username"))
	}

	return diags
}

func updateCliUser(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := meta.(*provider.ProviderState).SempV1Client
	username := d.Get(CLI_USERNAME).(string)
	password := d.Get(CLI_PASSWORD).(string)

	err := client.SetCliUserPassword(ctx, username, password)
	if err != nil {
		return provider.AppendError(diags, fmt.Errorf("cannot set password for user %s: %s", username, err))
	}
	err = client.SetCliUserGlobalAccessLevel(ctx, username, d.Get(CLI_GLOBAL_ACCESS_LEVEL).(string))
	if err != nil {
		return provider.AppendError(diags, fmt.Errorf("cannot set global access-level for user %s: %s", username, err))
	}
	err = client.SetCliUserMsgVpnDefaultAccessLevel(ctx, username, d.Get(CLI_MSGVPN_DEFAULT_ACCESS_LEVEL).(string))
	if err != nil {
		return provider.AppendError(diags, fmt.Errorf("cannot set msgvpn default access-level for user %s: %s", username, err))
	}

	if d.HasChange(CLI_MSGVPN_ACCESS_LEVEL_EXCEPTION) {
		target_ex := d.Get(CLI_MSGVPN_ACCESS_LEVEL_EXCEPTION).(map[string]interface{})
		target := iSMapTosSMap(&target_ex)
		user, err := client.GetCliUser(ctx, username)
		if err != nil {
			return provider.AppendError(diags, fmt.Errorf("cannot read current msgvpn exceptions user %s: %s", username, err))
		}
		current_ex := user.Exceptions()

		// Get the delta of specified exceptions and currently assigned exceptions (c=toCreate, u=toUpdatem, d=toDelete)
		c, u, d := StringMapDelta(target, current_ex)
		for k, v := range c {
			err = client.CreateCliUserVpnAccessLevelException(ctx, username, k, v)
			if err != nil {
				return provider.AppendError(diags, fmt.Errorf("cannot create msgvpn exception for user %s, msgvpn %s and access level %s: %s", username, k, v, err))
			}
		}
		for k, v := range u {
			err = client.SetCliUserVpnAccessLevelException(ctx, username, k, v)
			if err != nil {
				return provider.AppendError(diags, fmt.Errorf("cannot update msgvpn exception for user %s, msgvpn %s and access level %s: %s", username, k, v, err))
			}
		}
		for k := range d {
			err = client.DeleteCliUserVpnAccessLevelException(ctx, username, k)
			if err != nil {
				return provider.AppendError(diags, fmt.Errorf("cannot delete msgvpn exception for user %s and msgvpn %s: %s", username, k, err))
			}
		}
	}

	return append(diags, readCliUser(ctx, d, meta)...)
}

func deleteCliUser(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := meta.(*provider.ProviderState).SempV1Client
	username := d.Get(CLI_USERNAME).(string)

	err := client.DeleteCliUser(ctx, username)
	if err != nil {
		return provider.AppendError(diags, err)
	}
	d.SetId("")
	return diags
}
