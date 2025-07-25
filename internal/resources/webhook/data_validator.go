package webhook

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// mainCustomDiffFunc orchestrates all custom diff validations.
func mainCustomDiffFunc(ctx context.Context, diff *schema.ResourceDiff, i interface{}) error {
	if err := validateAuthenticationRequirements(ctx, diff, i); err != nil {
		return err
	}

	if err := validateSmartGroupIDRequirement(ctx, diff, i); err != nil {
		return err
	}

	return nil
}

// validateAuthenticationRequirements checks the conditions related to the 'authentication_type' attribute.
func validateAuthenticationRequirements(_ context.Context, diff *schema.ResourceDiff, _ interface{}) error {
	resourceName := diff.Get("name").(string)
	authType, ok := diff.GetOk("authentication_type")
	if !ok || authType.(string) != "Basic Authentication" {
		return nil
	}

	username, usernameOk := diff.GetOk("username")
	password, passwordOk := diff.GetOk("password")

	if !usernameOk || username == "" {
		return fmt.Errorf("in 'jamfpro_webhook.%s': when 'authentication_type' is set to 'Basic Authentication', 'username' must be provided", resourceName)
	}
	if !passwordOk || password == "" {
		return fmt.Errorf("in 'jamfpro_webhook.%s': when 'authentication_type' is set to 'Basic Authentication', 'password' must be provided", resourceName)
	}

	return nil
}

// validateSmartGroupIDRequirement checks if the specified events require a smart_group_id and validates its presence.
func validateSmartGroupIDRequirement(_ context.Context, diff *schema.ResourceDiff, _ interface{}) error {
	resourceName := diff.Get("name").(string)
	event, ok := diff.GetOk("event")
	if !ok {
		return nil
	}

	// List of events that require a smart_group_id
	requiredEvents := []string{
		"SmartGroupComputerMembershipChange",
		"SmartGroupMobileDeviceMembershipChange",
		"SmartGroupUserMembershipChange",
	}

	// Check if the current event is in the list of required events
	for _, reqEvent := range requiredEvents {
		if event.(string) == reqEvent {
			smartGroupID, smartGroupIDOk := diff.GetOk("smart_group_id")
			if !smartGroupIDOk || smartGroupID == 0 {
				return fmt.Errorf("in 'jamfpro_webhook.%s': when 'event' is set to '%s', 'smart_group_id' must be provided and must be a valid non-zero integer", resourceName, event)
			}
			break
		}
	}

	return nil
}
