package main

import (
	"context"
	"fmt"
	"os"

	generated "github.com/open-feature/cli/test/go-integration/openfeature"
	"github.com/open-feature/go-sdk/openfeature"
	"github.com/open-feature/go-sdk/openfeature/memprovider"
)

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run() error {
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
		"themeCustomization": {
			State:          memprovider.Enabled,
			DefaultVariant: "default",
			Variants: map[string]any{
				"default": map[string]any{
					"primaryColor":   "#007bff",
					"secondaryColor": "#6c757d",
				},
			},
		},
	})

	// Set the provider and wait for it to be ready
	err := openfeature.SetProviderAndWait(provider)
	if err != nil {
		return fmt.Errorf("failed to set provider: %w", err)
	}

	ctx := context.Background()
	evalCtx := openfeature.NewEvaluationContext("someid", map[string]any{})

	// Use the generated code for all flag evaluations
	enableFeatureA := generated.EnableFeatureA.Value(ctx, evalCtx)
	fmt.Printf("enableFeatureA: %v\n", enableFeatureA)
	_, err = generated.EnableFeatureA.ValueWithDetails(ctx, evalCtx)
	if err != nil {
		return fmt.Errorf("Error evaluating boolean flag: %w", err)
	}

	discount := generated.DiscountPercentage.Value(ctx, evalCtx)
	fmt.Printf("Discount Percentage: %.2f\n", discount)
	_, err = generated.DiscountPercentage.ValueWithDetails(ctx, evalCtx)
	if err != nil {
		return fmt.Errorf("Failed to get discount: %w", err)
	}

	greetingMessage := generated.GreetingMessage.Value(ctx, evalCtx)
	fmt.Printf("greetingMessage: %v\n", greetingMessage)
	_, err = generated.GreetingMessage.ValueWithDetails(ctx, evalCtx)
	if err != nil {
		return fmt.Errorf("Error evaluating string flag: %w", err)
	}

	usernameMaxLength := generated.UsernameMaxLength.Value(ctx, evalCtx)
	fmt.Printf("usernameMaxLength: %v\n", usernameMaxLength)
	_, err = generated.UsernameMaxLength.ValueWithDetails(ctx, evalCtx)
	if err != nil {
		return fmt.Errorf("Error evaluating int flag: %v\n", err)
	}

	themeCustomization := generated.ThemeCustomization.Value(ctx, evalCtx)
	_, err = generated.ThemeCustomization.ValueWithDetails(ctx, evalCtx)
	if err != nil {
		return fmt.Errorf("Error evaluating int flag: %v\n", err)
	}
	fmt.Printf("themeCustomization: %v\n", themeCustomization)

	// Test the String() method functionality for all flags
	fmt.Printf("enableFeatureA flag key: %s\n", generated.EnableFeatureA.String())
	fmt.Printf("discountPercentage flag key: %s\n", generated.DiscountPercentage.String())
	fmt.Printf("greetingMessage flag key: %s\n", generated.GreetingMessage.String())
	fmt.Printf("usernameMaxLength flag key: %s\n", generated.UsernameMaxLength.String())
	fmt.Printf("themeCustomization flag key: %s\n", generated.ThemeCustomization.String())

	fmt.Println("Generated Go code compiles successfully!")

	return nil
}
