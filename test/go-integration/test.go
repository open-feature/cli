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
	if enableFeatureA != false {
		return fmt.Errorf("error evaluating %s flag. want: %v, got: %v", generated.EnableFeatureA, false, enableFeatureA)
	}
	_, err = generated.EnableFeatureA.ValueWithDetails(ctx, evalCtx)
	if err != nil {
		return fmt.Errorf("Error evaluating boolean flag: %w", err)
	}
	fmt.Printf("enableFeatureA: %v\n", enableFeatureA)

	discount := generated.DiscountPercentage.Value(ctx, evalCtx)
	if discount != 0.15 {
		return fmt.Errorf("error evaluating %s flag. want: %v, got: %v", generated.DiscountPercentage, 0.15, discount)
	}
	_, err = generated.DiscountPercentage.ValueWithDetails(ctx, evalCtx)
	if err != nil {
		return fmt.Errorf("Failed to get discount: %w", err)
	}
	fmt.Printf("Discount Percentage: %.2f\n", discount)

	greetingMessage := generated.GreetingMessage.Value(ctx, evalCtx)
	if greetingMessage != "Hello there!" {
		return fmt.Errorf("error evaluating %s flag. want: %v, got: %v", generated.GreetingMessage, "Hello there!", greetingMessage)
	}
	_, err = generated.GreetingMessage.ValueWithDetails(ctx, evalCtx)
	if err != nil {
		return fmt.Errorf("Error evaluating string flag: %w", err)
	}
	fmt.Printf("greetingMessage: %v\n", greetingMessage)

	usernameMaxLength := generated.UsernameMaxLength.Value(ctx, evalCtx)
	if usernameMaxLength != 50 {
		return fmt.Errorf("error evaluating %s flag. want: %v, got: %v", generated.UsernameMaxLength, 50, usernameMaxLength)
	}
	_, err = generated.UsernameMaxLength.ValueWithDetails(ctx, evalCtx)
	if err != nil {
		return fmt.Errorf("Error evaluating int flag: %v\n", err)
	}
	fmt.Printf("usernameMaxLength: %v\n", usernameMaxLength)

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
