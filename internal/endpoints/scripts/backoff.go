package scripts

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// APICallFunc is a generic function type for API calls that can handle different types of IDs.
type APICallFunc func(interface{}) (interface{}, error)

// APICallFuncInt is specifically for API calls that require an integer ID.
type APICallFuncInt func(int) (interface{}, error)

// APICallFuncString is specifically for API calls that require a string ID.
type APICallFuncString func(string) (interface{}, error)

// ByResourceIntID is a wrapper function that facilitates retrying API Read calls which require an integer ID.
// It adapts an API call expecting an integer ID to the generic retry mechanism provided by RetryAPIReadCall.
//
// Parameters:
//   - ctx: A context.Context instance that carries deadlines, cancellation signals, and other request-scoped values across API boundaries and between processes.
//   - d: The Terraform resource data schema instance, providing access to the operations timeout settings and the resource's state.
//   - resourceID: The unique integer identifier of the resource to be fetched.
//   - apiCall: The specific API call function that accepts an integer ID and returns the resource along with any error encountered during the fetch operation.
//
// Returns:
//   - interface{}: The resource fetched by the API call if successful. This will need to be type-asserted to the specific resource type expected by the caller.
//   - diag.Diagnostics: A collection of diagnostic information including any errors encountered during the operation or warnings related to the resource's state.
//
// Note: If the resource cannot be found or if an error occurs, appropriate diagnostics are returned to Terraform, potentially marking the resource for deletion from the state if not found.
func ByResourceIntID(ctx context.Context, d *schema.ResourceData, resourceID int, apiCall APICallFuncInt) (interface{}, diag.Diagnostics) {
	genericAPICall := func(id interface{}) (interface{}, error) {
		intID, ok := id.(int)
		if !ok {
			return nil, fmt.Errorf("expected int ID, got %T", id)
		}
		return apiCall(intID)
	}
	return WaitForResourceToBeAvailable(ctx, d, resourceID, genericAPICall)
}

// ByResourceStringID is a wrapper function that facilitates retrying API Read calls which require a string ID.
// It adapts an API call expecting a string ID to the generic retry mechanism provided by RetryAPIReadCall.
//
// Parameters:
//   - ctx: A context.Context instance that carries deadlines, cancellation signals, and other request-scoped values across API boundaries and between processes.
//   - d: The Terraform resource data schema instance, providing access to the operations timeout settings and the resource's state.
//   - resourceID: The unique string identifier of the resource to be fetched.
//   - apiCall: The specific API call function that accepts a string ID and returns the resource along with any error encountered during the fetch operation.
//
// Returns:
//   - interface{}: The resource fetched by the API call if successful. This will need to be type-asserted to the specific resource type expected by the caller.
//   - diag.Diagnostics: A collection of diagnostic information including any errors encountered during the operation or warnings related to the resource's state.
//
// Note: If the resource cannot be found or if an error occurs, appropriate diagnostics are returned to Terraform, potentially marking the resource for deletion from the state if not found.
func ByResourceStringID(ctx context.Context, d *schema.ResourceData, resourceID string, apiCall APICallFuncString) (interface{}, diag.Diagnostics) {
	genericAPICall := func(id interface{}) (interface{}, error) {
		strID, ok := id.(string)
		if !ok {
			return nil, fmt.Errorf("expected string ID, got %T", id)
		}
		return apiCall(strID)
	}
	return WaitForResourceToBeAvailable(ctx, d, resourceID, genericAPICall)
}

// WaitForResourceToBeAvailable employs a retry mechanism with exponential backoff and jitter to wait for a resource to become available. This function is particularly useful in scenarios where a resource creation is asynchronous and may not be immediately available after a create API call.
//
// The function uses an APICallFunc to repeatedly check for the existence of the resource, retrying in the face of "resource not found" errors, which are common immediately after resource creation. Other types of errors are not retried and lead to an immediate return.
//
// Exponential backoff helps in efficiently spacing out retry attempts to reduce load on the server and minimize the chance of failures due to rate limiting or server overload. Jitter is added to the backoff duration to prevent retry storms in scenarios with many concurrent operations.
//
// The retry process respects the provided context's deadline, ensuring that the function does not exceed the overall timeout specified for the resource creation operation in Terraform. This approach ensures robustness in transient network issues or temporary server-side unavailability.
//
// Parameters:
//   - ctx: The context governing the retry operation, carrying timeout and cancellation signals.
//   - d: The Terraform resource data schema instance, providing access to the resource's operational timeout settings.
//   - resourceID: The unique identifier of the resource being waited on.
//   - checkResourceExists: A function conforming to the APICallFunc type that attempts to fetch the resource by its ID, returning the resource or an error.
//
// Returns:
//   - interface{}: The successfully fetched resource if available, needing type assertion to the expected resource type by the caller.
//   - diag.Diagnostics: Diagnostic information including any errors encountered during the wait operation, or warnings related to the resource's availability state.
func WaitForResourceToBeAvailable(ctx context.Context, d *schema.ResourceData, resourceID interface{}, checkResourceExists APICallFunc) (interface{}, diag.Diagnostics) {
	var diags diag.Diagnostics
	var lastError error
	var resource interface{}

	initialBackoff := 1 * time.Second
	maxBackoff := 30 * time.Second
	backoffFactor := 2.0
	jitterFactor := 0.5

	currentBackoff := initialBackoff

	err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		var apiErr error
		resource, apiErr = checkResourceExists(resourceID)
		if apiErr != nil {
			lastError = apiErr

			// Check specifically for "resource not found" errors to retry
			if strings.Contains(apiErr.Error(), "404") || strings.Contains(apiErr.Error(), "410") {
				// Apply exponential backoff with jitter
				time.Sleep(currentBackoff + time.Duration(rand.Float64()*jitterFactor*float64(currentBackoff)))
				currentBackoff = time.Duration(float64(currentBackoff) * backoffFactor) // Corrected line
				if currentBackoff > maxBackoff {
					currentBackoff = maxBackoff
				}
				return retry.RetryableError(apiErr)
			}

			// For other types of errors, do not retry and return the error
			return retry.NonRetryableError(apiErr)
		}

		// If no error, the resource exists, stop retrying
		lastError = nil
		return nil
	})

	// If an error occurred during retries (other than the resource not found),
	// add it to diagnostics
	if err != nil {
		diags = append(diags, diag.FromErr(fmt.Errorf("error waiting for resource with ID '%v' to become available: %v", resourceID, lastError))...)
		return nil, diags // Return nil as the resource and the diagnostics
	}

	// Return the successfully fetched resource and any diagnostics
	return resource, diags
}
