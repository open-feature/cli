package main

import (
	"context"
	"fmt"
	"github.com/open-feature/go-sdk/openfeature"
	"github.com/open-feature/go-sdk/openfeature/memprovider"
)

// This program validates that the generated OpenFeature Go client code compiles
// We don't need to run the code since the goal is to test compilation only
func main() {
	// Set up the in-memory provider with test flags
	provider := memprovider.NewInMemoryProvider(map[string]memprovider.InMemoryFlag{
		"discountPercentage": {
			State:          memprovider.Enabled,
			DefaultVariant: "default",
			Variants: map[string]any{
				"default": 0.15,
			},
		},
		"enableFeatureA": {
			State:          memprovider.Enabled,
			DefaultVariant: "default",
			Variants: map[string]any{
				"default": false,
			},
		},
		"greetingMessage": {
			State:          memprovider.Enabled,
			DefaultVariant: "default",
			Variants: map[string]any{
				"default": "Hello there!",
			},
		},
		"usernameMaxLength": {
			State:          memprovider.Enabled,
			DefaultVariant: "default",
			Variants: map[string]any{
				"default": 50,
			},
		},
	})

	// Set the provider and wait for it to be ready
	err := openfeature.SetProviderAndWait(provider)
	if err != nil {
		fmt.Printf("Failed to set provider: %v\n", err)
		return
	}

	// Create a new client
	client := openfeature.NewClient("test-app")

	// Test flag evaluations
	ctx := context.Background()

	// Test boolean flag
	enableFeatureA, err := client.BooleanValue(ctx, "enableFeatureA", true, openfeature.EvaluationContext{})
	if err != nil {
		fmt.Printf("Error evaluating boolean flag: %v\n", err)
		return
	}
	fmt.Printf("enableFeatureA: %v\n", enableFeatureA)

	// Test float flag
	discountPercentage, err := client.FloatValue(ctx, "discountPercentage", 0.0, openfeature.EvaluationContext{})
	if err != nil {
		fmt.Printf("Error evaluating float flag: %v\n", err)
		return
	}
	fmt.Printf("discountPercentage: %v\n", discountPercentage)

	// Test string flag
	greetingMessage, err := client.StringValue(ctx, "greetingMessage", "", openfeature.EvaluationContext{})
	if err != nil {
		fmt.Printf("Error evaluating string flag: %v\n", err)
		return
	}
	fmt.Printf("greetingMessage: %v\n", greetingMessage)

	// Test integer flag
	usernameMaxLength, err := client.IntValue(ctx, "usernameMaxLength", 0, openfeature.EvaluationContext{})
	if err != nil {
		fmt.Printf("Error evaluating int flag: %v\n", err)
		return
	}
	fmt.Printf("usernameMaxLength: %v\n", usernameMaxLength)

	fmt.Println("Generated Go code compiles successfully!")
} 