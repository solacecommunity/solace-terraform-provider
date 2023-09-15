package resources

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/solacecommunity/solace-go-semp-client/client/all"
	"github.com/solacecommunity/solace-go-semp-client/models"
	"github.com/solacecommunity/solace-terraform-provider/provider"
)

// Main resource definition for Client Usernames
func ResourceAclProfile() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,
		// List of supported configuration fields for your resource
		Schema: schemaAclProfile(),
		// Provider CRUD functions
		CreateContext: createAclProfile,
		ReadContext:   readAclProfile,
		UpdateContext: updateAclProfile,
		DeleteContext: deleteAclProfile,
	}
}

func schemaAclProfile() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		MSG_VPN_NAME: {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		ACL_PROFILE_NAME: {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "The name of the ACL Profile.",
		},
		CLIENT_CONNECT_DEFAULT_ACTION: {
			Type:         schema.TypeString,
			Optional:     true,
			Default:      "allow",
			ValidateFunc: validation.StringInSlice([]string{"allow", "disallow"}, false),
			Description: "The default action to take when a client using the ACL Profile connects to the Message VPN. " +
				"The default value is \"allow\". The allowed values and their meaning are: \"allow\" - Allow client connection " +
				"unless an exception is found for it OR \"disallow\" - Disallow client connection unless an exception is found for it.",
		},
		PUBLISH_TOPIC_DEFAULT_ACTION: {
			Type:         schema.TypeString,
			Optional:     true,
			Default:      "disallow",
			ValidateFunc: validation.StringInSlice([]string{"allow", "disallow"}, false),
			Description: "The default action to take when a client using the ACL Profile publishes to a topic in the Message VPN. " +
				"The default value is \"disallow\".",
		},
		SUBSCRIBE_TOPIC_DEFAULT_ACTION: {
			Type:         schema.TypeString,
			Optional:     true,
			Default:      "disallow",
			ValidateFunc: validation.StringInSlice([]string{"allow", "disallow"}, false),
			Description: "The default action to take when a client using the ACL Profile subscribes to a topic in the Message VPN. " +
				"The default value is \"disallow\".",
		},
		SUBSCRIBE_SHARE_NAME_DEFAULT_ACTION: {
			Type:         schema.TypeString,
			Optional:     true,
			Default:      "disallow",
			ValidateFunc: validation.StringInSlice([]string{"allow", "disallow"}, false),
			Description:  "The default action to take when a client using the ACL Profile subscribes to a share-name subscription in the Message VPN",
		},
		ACL_PUBLISH_EXCEPTION_LIST: {
			Type:     schema.TypeSet,
			Optional: true,
			Elem: &schema.Schema{
				Type:             schema.TypeString,
				ValidateDiagFunc: validateTopicExceptionString, // unfortunately this is ignored by terraform... we'll call this function manually then
			},
			Description: "A Publish Topic Exception is an exception to the default action to take when a client using the ACL " +
				"Profile publishes to a topic in the Message VPN. Exceptions must be expressed as a topic.",
		},
		ACL_SUBSCRIBE_EXCEPTION_LIST: {
			Type:     schema.TypeSet,
			Optional: true,
			Elem: &schema.Schema{
				Type:             schema.TypeString,
				ValidateDiagFunc: validateTopicExceptionString, // unfortunately this is ignored by terraform... we'll call this function manually then
			},
			Description: "A Subscription Topic Exception is an exception to the default action to take when a client using the ACL " +
				"Profile subscribes to a topic in the Message VPN. Exceptions must be expressed as a topic.",
		},
		ACL_SUBSCRIBE_SHARE_NAME_EXCEPTION_LIST: {
			Type:     schema.TypeSet,
			Optional: true,
			Elem: &schema.Schema{
				Type:             schema.TypeString,
				ValidateDiagFunc: validateTopicExceptionString, // unfortunately this is ignored by terraform... we'll call this function manually then
			},
			Description: "A Subscription Topic Exception is an exception to the default action to take when a client using the ACL " +
				"Profile subscribes to a topic in the Message VPN. Exceptions must be expressed as a topic.",
		},
	}
}

// validateTopicString validated a string and returns an error if it is not valid
func validateTopicExceptionString(v interface{}, p cty.Path) diag.Diagnostics {
	topic := v.(string)
	var diags diag.Diagnostics
	const SUMMARY = "wrong topic exptions subscription"

	if err := validateACLTopicException(topic); err != nil {
		diag := diag.Diagnostic{
			Severity: diag.Error,
			Summary:  SUMMARY,
			Detail:   err.Error(),
		}
		return append(diags, diag)
	}

	return diags
}

func getAclProfileModelFromResource(d *schema.ResourceData) *models.MsgVpnACLProfile {
	q := &models.MsgVpnACLProfile{
		MsgVpnName:     d.Get(MSG_VPN_NAME).(string),
		ACLProfileName: d.Get(ACL_PROFILE_NAME).(string),
	}
	if v, s := d.GetOk(CLIENT_CONNECT_DEFAULT_ACTION); s {
		q.ClientConnectDefaultAction = v.(string)
	}
	if v, s := d.GetOk(PUBLISH_TOPIC_DEFAULT_ACTION); s {
		q.PublishTopicDefaultAction = v.(string)
	}
	if v, s := d.GetOk(SUBSCRIBE_TOPIC_DEFAULT_ACTION); s {
		q.SubscribeTopicDefaultAction = v.(string)
	}
	if v, s := d.GetOk(SUBSCRIBE_SHARE_NAME_DEFAULT_ACTION); s {
		q.SubscribeShareNameDefaultAction = v.(string)
	}
	return q
}

func createAclProfile(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	state := meta.(*provider.ProviderState)
	msgvpn := d.Get(MSG_VPN_NAME).(string)

	body := getAclProfileModelFromResource(d)
	params := all.NewCreateMsgVpnACLProfileParams().WithContext(ctx).WithMsgVpnName(msgvpn).WithBody(body)
	_, err := state.Client.All.CreateMsgVpnACLProfile(params, state.Auth)
	if err != nil {
		if semp, err2 := provider.GetSempError(err); err2 == nil {
			if *semp.Status == "ALREADY_EXISTS" {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  fmt.Sprintf("ACL Profile %s already exists. Going to import state from Broker", body.ACLProfileName),
				})
			} else {
				return provider.AppendError(diags, err)
			}
		} else {
			return provider.AppendError(diags, err)
		}
	}

	if d.HasChange(ACL_SUBSCRIBE_EXCEPTION_LIST) {
		err := syncSubscribeTopicExceptions(ctx, d, state, msgvpn, body.ACLProfileName)
		if err != nil {
			if semp, err2 := provider.GetSempError(err); err2 == nil {
				if *semp.Status == "ALREADY_EXISTS" {
					diags = append(diags, diag.Diagnostic{
						Severity: diag.Warning,
						Summary:  fmt.Sprintf("ACL Subscribe Exception already exists for %s. Going to import state from Broker", body.ACLProfileName),
					})
				} else {
					return provider.AppendError(diags, err)
				}
			} else {
				return provider.AppendError(diags, err)
			}
		}
	}

	if d.HasChange(ACL_PUBLISH_EXCEPTION_LIST) {
		err := syncPublishTopicExceptions(ctx, d, state, msgvpn, body.ACLProfileName)
		if err != nil {
			if semp, err2 := provider.GetSempError(err); err2 == nil {
				if *semp.Status == "ALREADY_EXISTS" {
					diags = append(diags, diag.Diagnostic{
						Severity: diag.Warning,
						Summary:  fmt.Sprintf("ACL Publish Exception already exists for %s. Going to import state from Broker", body.ACLProfileName),
					})
				} else {
					return provider.AppendError(diags, err)
				}
			} else {
				return provider.AppendError(diags, err)
			}
		}
	}

	if d.HasChange(ACL_SUBSCRIBE_SHARE_NAME_EXCEPTION_LIST) {
		err := syncSubscribeShareNameExceptions(ctx, d, state, msgvpn, body.ACLProfileName)
		if err != nil {
			if semp, err2 := provider.GetSempError(err); err2 == nil {
				if *semp.Status == "ALREADY_EXISTS" {
					diags = append(diags, diag.Diagnostic{
						Severity: diag.Warning,
						Summary:  fmt.Sprintf("ACL Subscribe Share Name Exception already exists for %s. Going to import state from Broker", body.ACLProfileName),
					})
				} else {
					return provider.AppendError(diags, err)
				}
			} else {
				return provider.AppendError(diags, err)
			}
		}
	}

	d.SetId("sol_" + msgvpn + "_ap_" + body.ACLProfileName)

	// read after create to make sure tf state is in sync with broker
	return append(diags, readAclProfile(ctx, d, meta)...)
}

func readAclProfile(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	state := meta.(*provider.ProviderState)
	msgvpn := d.Get(MSG_VPN_NAME).(string)
	profile := d.Get(ACL_PROFILE_NAME).(string)

	params := all.NewGetMsgVpnACLProfileParamsWithContext(ctx).WithMsgVpnName(msgvpn).WithACLProfileName(profile)
	resp, err := state.Client.All.GetMsgVpnACLProfile(params, state.Auth)
	if err != nil {
		if semp, err2 := provider.GetSempError(err); err2 == nil {
			if *semp.Status == provider.SEMP_STATUS_NOT_FOUND {
				d.SetId("")
				return append(diags, diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  fmt.Sprintf("ACL Profile %s does not exist. Going to remove ressource from tfstate", profile),
				})
			} else {
				return provider.AppendError(diags, err)
			}
		} else {
			return provider.AppendError(diags, err)
		}
	}
	p := resp.Payload.Data
	d.Set(CLIENT_CONNECT_DEFAULT_ACTION, p.ClientConnectDefaultAction)
	d.Set(PUBLISH_TOPIC_DEFAULT_ACTION, p.PublishTopicDefaultAction)
	d.Set(SUBSCRIBE_TOPIC_DEFAULT_ACTION, p.SubscribeTopicDefaultAction)
	d.Set(SUBSCRIBE_SHARE_NAME_DEFAULT_ACTION, p.SubscribeShareNameDefaultAction)

	ptl, err := readPublishTopicExceptions(ctx, state, msgvpn, profile)
	if err != nil {
		return provider.AppendError(diags, err)
	}

	stl, err := readSubscribeTopicExceptions(ctx, state, msgvpn, profile)
	if err != nil {
		return provider.AppendError(diags, err)
	}

	ssnl, err := readSubscribeShareNameExceptions(ctx, state, msgvpn, profile)
	if err != nil {
		return provider.AppendError(diags, err)
	}

	d.Set(ACL_PUBLISH_EXCEPTION_LIST, &ptl)
	d.Set(ACL_SUBSCRIBE_EXCEPTION_LIST, &stl)
	d.Set(ACL_SUBSCRIBE_SHARE_NAME_EXCEPTION_LIST, &ssnl)

	return diags
}

func updateAclProfile(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	state := meta.(*provider.ProviderState)
	msgvpn := d.Get(MSG_VPN_NAME).(string)
	profile := d.Get(ACL_PROFILE_NAME).(string)
	body := getAclProfileModelFromResource(d)

	params := all.NewUpdateMsgVpnACLProfileParamsWithContext(ctx).WithMsgVpnName(msgvpn).WithACLProfileName(profile).WithBody(body)
	_, err := state.Client.All.UpdateMsgVpnACLProfile(params, state.Auth)
	if err != nil {
		return provider.AppendError(diags, err)
	}

	if d.HasChange(ACL_PUBLISH_EXCEPTION_LIST) {
		err := syncPublishTopicExceptions(ctx, d, state, msgvpn, profile)
		if err != nil {
			if semp, err2 := provider.GetSempError(err); err2 == nil {
				if *semp.Status == "ALREADY_EXISTS" {
					diags = append(diags, diag.Diagnostic{
						Severity: diag.Warning,
						Summary:  fmt.Sprintf("ACL Publish Exception already exists for %s. Going to import state from Broker", body.ACLProfileName),
					})
				} else {
					return provider.AppendError(diags, err)
				}
			} else {
				return provider.AppendError(diags, err)
			}
		}
	}

	if d.HasChange(ACL_SUBSCRIBE_EXCEPTION_LIST) {
		err := syncSubscribeTopicExceptions(ctx, d, state, msgvpn, profile)
		if err != nil {
			if semp, err2 := provider.GetSempError(err); err2 == nil {
				if *semp.Status == "ALREADY_EXISTS" {
					diags = append(diags, diag.Diagnostic{
						Severity: diag.Warning,
						Summary:  fmt.Sprintf("ACL Subscribe Exception already exists for %s. Going to import state from Broker", body.ACLProfileName),
					})
				} else {
					return provider.AppendError(diags, err)
				}
			} else {
				return provider.AppendError(diags, err)
			}
		}
	}

	if d.HasChange(ACL_SUBSCRIBE_SHARE_NAME_EXCEPTION_LIST) {
		err := syncSubscribeShareNameExceptions(ctx, d, state, msgvpn, profile)
		if err != nil {
			if semp, err2 := provider.GetSempError(err); err2 == nil {
				if *semp.Status == "ALREADY_EXISTS" {
					diags = append(diags, diag.Diagnostic{
						Severity: diag.Warning,
						Summary:  fmt.Sprintf("ACL Subscribe Share Name Exception already exists for %s. Going to import state from Broker", body.ACLProfileName),
					})
				} else {
					return provider.AppendError(diags, err)
				}
			} else {
				return provider.AppendError(diags, err)
			}
		}
	}

	// read after update to make sure tf state is in sync with broker
	return append(diags, readAclProfile(ctx, d, meta)...)
}

func deleteAclProfile(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	state := meta.(*provider.ProviderState)
	msgvpn := d.Get(MSG_VPN_NAME).(string)
	profile := d.Get(ACL_PROFILE_NAME).(string)

	params := all.NewDeleteMsgVpnACLProfileParamsWithContext(ctx).WithMsgVpnName(msgvpn).WithACLProfileName(profile)
	_, err := state.Client.All.DeleteMsgVpnACLProfile(params, state.Auth)
	if err != nil {
		return provider.AppendError(diags, err)
	}
	return diags
}

func syncSubscribeTopicExceptions(ctx context.Context, d *schema.ResourceData, state *provider.ProviderState, msgvpn string, profile string) error {
	// get a list of existing exceptions from broker
	ref, err := readSubscribeTopicExceptions(ctx, state, msgvpn, profile)
	if err != nil {
		return err
	}
	// get a list of estimated state
	t := d.Get(ACL_SUBSCRIBE_EXCEPTION_LIST).(*schema.Set).List()
	target := iSliceTosSlice(&t)

	// validate topics from tf file
	for _, topic := range target {
		if err := validateACLTopicException(topic); err != nil {
			return fmt.Errorf("the topic ACL topic string %s is invalid because: %w", topic, err)
		}
	}

	// compare both lists
	news, olds := SliceDelta(&ref, &target)
	// add new entries
	for _, s := range news {
		err := createSubscribeTopicException(ctx, state, msgvpn, profile, s)
		if err != nil {
			return err
		}
	}
	// delete old ones
	for _, s := range olds {
		err := deleteSubscribeTopicException(ctx, state, msgvpn, profile, s)
		if err != nil {
			return err
		}
	}
	return nil
}

func readSubscribeTopicExceptions(ctx context.Context, state *provider.ProviderState, msgvpn string, profile string) ([]string, error) {
	count := int64(100000)

	params := all.NewGetMsgVpnACLProfileSubscribeTopicExceptionsParamsWithContext(ctx).
		WithMsgVpnName(msgvpn).WithACLProfileName(profile).WithCount(&count)

	resp, err := state.Client.All.GetMsgVpnACLProfileSubscribeTopicExceptions(params, state.Auth)
	if err != nil {
		return nil, err
	}

	result := make([]string, len(resp.Payload.Data))
	for i, v := range resp.Payload.Data {
		result[i] = v.SubscribeTopicException + "@" + v.SubscribeTopicExceptionSyntax
	}
	return result, nil
}

func createSubscribeTopicException(ctx context.Context, state *provider.ProviderState, msgvpn string, profile string, topicAndSyntax string) error {
	topicAndSyntaxSlice := strings.Split(topicAndSyntax, "@")
	topic := topicAndSyntaxSlice[0]
	syntax := "smf"
	if len(topicAndSyntaxSlice) > 1 {
		syntax = topicAndSyntaxSlice[1]
	}
	body := &models.MsgVpnACLProfileSubscribeTopicException{
		ACLProfileName:                profile,
		MsgVpnName:                    msgvpn,
		SubscribeTopicException:       topic,
		SubscribeTopicExceptionSyntax: syntax,
	}
	params := all.NewCreateMsgVpnACLProfileSubscribeTopicExceptionParamsWithContext(ctx).
		WithMsgVpnName(msgvpn).WithACLProfileName(profile).WithBody(body)
	_, err := state.Client.All.CreateMsgVpnACLProfileSubscribeTopicException(params, state.Auth)
	return err
}

func deleteSubscribeTopicException(ctx context.Context, state *provider.ProviderState, msgvpn string, profile string, topicAndSyntax string) error {
	topicAndSyntaxSlice := strings.Split(topicAndSyntax, "@")
	topic := topicAndSyntaxSlice[0]
	syntax := "smf"
	if len(topicAndSyntaxSlice) > 1 {
		syntax = topicAndSyntaxSlice[1]
	}
	params := all.NewDeleteMsgVpnACLProfileSubscribeTopicExceptionParamsWithContext(ctx).
		WithMsgVpnName(msgvpn).WithACLProfileName(profile).WithSubscribeTopicException(topic).
		WithSubscribeTopicExceptionSyntax(syntax)
	_, err := state.Client.All.DeleteMsgVpnACLProfileSubscribeTopicException(params, state.Auth)
	return err
}

func syncPublishTopicExceptions(ctx context.Context, d *schema.ResourceData, state *provider.ProviderState, msgvpn string, profile string) error {
	// get a list of existing exceptions from broker
	ref, err := readPublishTopicExceptions(ctx, state, msgvpn, profile)
	if err != nil {
		return err
	}

	// get a list of estimated state
	t := d.Get(ACL_PUBLISH_EXCEPTION_LIST).(*schema.Set).List()
	target := iSliceTosSlice(&t)

	// validate topics from tf file
	for _, topic := range target {
		if err := validateACLTopicException(topic); err != nil {
			return fmt.Errorf("the topic ACL topic string %s is invalid because: %w", topic, err)
		}
	}

	// compare both lists
	news, olds := SliceDelta(&ref, &target)
	// add new entries
	for _, s := range news {
		err := createPublishTopicException(ctx, state, msgvpn, profile, s)
		if err != nil {
			return err
		}
	}
	// delete old ones
	for _, s := range olds {
		err := deletePublishTopicException(ctx, state, msgvpn, profile, s)
		if err != nil {
			return err
		}
	}
	return nil
}

func readPublishTopicExceptions(ctx context.Context, state *provider.ProviderState, msgvpn string, profile string) ([]string, error) {
	count := int64(100000)

	params := all.NewGetMsgVpnACLProfilePublishTopicExceptionsParamsWithContext(ctx).
		WithMsgVpnName(msgvpn).WithACLProfileName(profile).WithCount(&count)

	resp, err := state.Client.All.GetMsgVpnACLProfilePublishTopicExceptions(params, state.Auth)
	if err != nil {
		return nil, err
	}

	result := make([]string, len(resp.Payload.Data))
	for i, v := range resp.Payload.Data {
		result[i] = v.PublishTopicException + "@" + v.PublishTopicExceptionSyntax
	}

	return result, nil
}

func createPublishTopicException(ctx context.Context, state *provider.ProviderState, msgvpn string, profile string, topicAndSyntax string) error {
	topicAndSyntaxSlice := strings.Split(topicAndSyntax, "@")
	topic := topicAndSyntaxSlice[0]
	syntax := "smf"
	if len(topicAndSyntaxSlice) > 1 {
		syntax = topicAndSyntaxSlice[1]
	}
	body := &models.MsgVpnACLProfilePublishTopicException{
		ACLProfileName:              profile,
		MsgVpnName:                  msgvpn,
		PublishTopicException:       topic,
		PublishTopicExceptionSyntax: syntax,
	}
	params := all.NewCreateMsgVpnACLProfilePublishTopicExceptionParamsWithContext(ctx).
		WithMsgVpnName(msgvpn).WithACLProfileName(profile).WithBody(body)
	_, err := state.Client.All.CreateMsgVpnACLProfilePublishTopicException(params, state.Auth)
	return err
}

func deletePublishTopicException(ctx context.Context, state *provider.ProviderState, msgvpn string, profile string, topicAndSyntax string) error {
	topicAndSyntaxSlice := strings.Split(topicAndSyntax, "@")
	topic := topicAndSyntaxSlice[0]
	syntax := "smf"
	if len(topicAndSyntaxSlice) > 1 {
		syntax = topicAndSyntaxSlice[1]
	}
	params := all.NewDeleteMsgVpnACLProfilePublishTopicExceptionParamsWithContext(ctx).
		WithMsgVpnName(msgvpn).WithACLProfileName(profile).WithPublishTopicException(topic).
		WithPublishTopicExceptionSyntax(syntax)
	_, err := state.Client.All.DeleteMsgVpnACLProfilePublishTopicException(params, state.Auth)
	return err
}

func syncSubscribeShareNameExceptions(ctx context.Context, d *schema.ResourceData, state *provider.ProviderState, msgvpn string, profile string) error {
	// get a list of existing exceptions from broker
	ref, err := readSubscribeShareNameExceptions(ctx, state, msgvpn, profile)
	if err != nil {
		return err
	}
	// get a list of estimated state
	t := d.Get(ACL_SUBSCRIBE_SHARE_NAME_EXCEPTION_LIST).(*schema.Set).List()
	target := iSliceTosSlice(&t)

	// validate topics from tf file
	for _, topic := range target {
		if err := validateACLTopicException(topic); err != nil {
			return fmt.Errorf("the topic ACL topic string %s is invalid because: %w", topic, err)
		}
	}

	// compare both lists
	news, olds := SliceDelta(&ref, &target)
	// add new entries
	for _, s := range news {
		err := createSubscribeShareNameException(ctx, state, msgvpn, profile, s)
		if err != nil {
			return err
		}
	}
	// delete old ones
	for _, s := range olds {
		err := deleteSubscribeShareNameException(ctx, state, msgvpn, profile, s)
		if err != nil {
			return err
		}
	}
	return nil
}

func readSubscribeShareNameExceptions(ctx context.Context, state *provider.ProviderState, msgvpn string, profile string) ([]string, error) {
	count := int64(100000)
	params := all.NewGetMsgVpnACLProfileSubscribeShareNameExceptionsParamsWithContext(ctx).
		WithMsgVpnName(msgvpn).WithACLProfileName(profile).WithCount(&count)
	resp, err := state.Client.All.GetMsgVpnACLProfileSubscribeShareNameExceptions(params, state.Auth)
	if err != nil {
		return nil, err
	}

	result := make([]string, len(resp.Payload.Data))
	for i, v := range resp.Payload.Data {
		result[i] = v.SubscribeShareNameException + "@" + v.SubscribeShareNameExceptionSyntax
	}
	return result, nil
}

func createSubscribeShareNameException(ctx context.Context, state *provider.ProviderState, msgvpn string, profile string, topicAndSyntax string) error {
	topicAndSyntaxSlice := strings.Split(topicAndSyntax, "@")
	topic := topicAndSyntaxSlice[0]
	syntax := "smf"
	if len(topicAndSyntaxSlice) > 1 {
		syntax = topicAndSyntaxSlice[1]
	}
	body := &models.MsgVpnACLProfileSubscribeShareNameException{
		ACLProfileName:                    profile,
		MsgVpnName:                        msgvpn,
		SubscribeShareNameException:       topic,
		SubscribeShareNameExceptionSyntax: syntax,
	}
	params := all.NewCreateMsgVpnACLProfileSubscribeShareNameExceptionParamsWithContext(ctx).
		WithMsgVpnName(msgvpn).WithACLProfileName(profile).WithBody(body)
	_, err := state.Client.All.CreateMsgVpnACLProfileSubscribeShareNameException(params, state.Auth)
	return err
}

func deleteSubscribeShareNameException(ctx context.Context, state *provider.ProviderState, msgvpn string, profile string, topicAndSyntax string) error {
	topic := strings.Split(topicAndSyntax, "@")[0]
	syntax := strings.Split(topicAndSyntax, "@")[1]
	params := all.NewDeleteMsgVpnACLProfileSubscribeShareNameExceptionParamsWithContext(ctx).
		WithMsgVpnName(msgvpn).WithACLProfileName(profile).WithSubscribeShareNameException(topic).WithSubscribeShareNameExceptionSyntax(syntax)
	_, err := state.Client.All.DeleteMsgVpnACLProfileSubscribeShareNameException(params, state.Auth)
	return err
}
