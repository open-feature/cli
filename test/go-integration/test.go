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
	})

	// Set the provider and wait for it to be ready
	err := openfeature.SetProviderAndWait(provider)
	if err != nil {
		return fmt.Errorf("Failed to set provider: %w", err)
	}

	ctx := context.Background()
	evalCtx := openfeature.NewEvaluationContext("someid", map[string]any{})

	// Use the generated code for all flag evaluations
	enableFeatureA, err := generated.EnableFeatureA.Value(ctx, evalCtx)
	if err != nil {
		return fmt.Errorf("Error evaluating boolean flag: %w", err)
	}
	fmt.Printf("enableFeatureA: %v\n", enableFeatureA)

	discount, err := generated.DiscountPercentage.Value(ctx, evalCtx)
	if err != nil {
		return fmt.Errorf("Failed to get discount: %w", err)
	}
	fmt.Printf("Discount Percentage: %.2f\n", discount)

	greetingMessage, err := generated.GreetingMessage.Value(ctx, evalCtx)
	if err != nil {
		return fmt.Errorf("Error evaluating string flag: %w", err)
	}
	fmt.Printf("greetingMessage: %v\n", greetingMessage)

	usernameMaxLength, err := generated.UsernameMaxLength.Value(ctx, evalCtx)
	if err != nil {
		return fmt.Errorf("Error evaluating int flag: %v\n", err)
	}
	fmt.Printf("usernameMaxLength: %v\n", usernameMaxLength)

	fmt.Println("Generated Go code compiles successfully!")

	return nil
}
