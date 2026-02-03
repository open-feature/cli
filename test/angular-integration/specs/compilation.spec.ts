import { describe, it, expect } from "vitest";

/**
 * Compilation tests verify that the generated code imports correctly
 * and all types are properly defined.
 */
describe("Compilation Tests", () => {
  it("should import FlagKeys constant", async () => {
    const { FlagKeys } = await import("../generated/openfeature.generated");

    expect(FlagKeys).toBeDefined();
    expect(typeof FlagKeys).toBe("object");
  });

  it("should have all expected flag keys", async () => {
    const { FlagKeys } = await import("../generated/openfeature.generated");

    expect(FlagKeys.DISCOUNT_PERCENTAGE).toBe("discountPercentage");
    expect(FlagKeys.ENABLE_FEATURE_A).toBe("enableFeatureA");
    expect(FlagKeys.GREETING_MESSAGE).toBe("greetingMessage");
    expect(FlagKeys.THEME_CUSTOMIZATION).toBe("themeCustomization");
    expect(FlagKeys.USERNAME_MAX_LENGTH).toBe("usernameMaxLength");
  });

  it("should import GeneratedFeatureFlagService", async () => {
    const { GeneratedFeatureFlagService } =
      await import("../generated/openfeature.generated");

    expect(GeneratedFeatureFlagService).toBeDefined();
    expect(typeof GeneratedFeatureFlagService).toBe("function");
  });

  it("should import all generated directives", async () => {
    const {
      DiscountPercentageDirective,
      EnableFeatureADirective,
      GreetingMessageDirective,
      ThemeCustomizationDirective,
      UsernameMaxLengthDirective,
    } = await import("../generated/openfeature.generated");

    expect(DiscountPercentageDirective).toBeDefined();
    expect(EnableFeatureADirective).toBeDefined();
    expect(GreetingMessageDirective).toBeDefined();
    expect(ThemeCustomizationDirective).toBeDefined();
    expect(UsernameMaxLengthDirective).toBeDefined();
  });

  it("should import GeneratedFeatureFlagDirectives array", async () => {
    const { GeneratedFeatureFlagDirectives } =
      await import("../generated/openfeature.generated");

    expect(GeneratedFeatureFlagDirectives).toBeDefined();
    expect(Array.isArray(GeneratedFeatureFlagDirectives)).toBe(true);
    expect(GeneratedFeatureFlagDirectives.length).toBe(5);
  });

  it("should export FlagKey type (via type inference)", async () => {
    const { FlagKeys } = await import("../generated/openfeature.generated");

    // Type inference test - if this compiles, the type is exported correctly
    const key: keyof typeof FlagKeys = "ENABLE_FEATURE_A";
    expect(FlagKeys[key]).toBe("enableFeatureA");
  });
});
