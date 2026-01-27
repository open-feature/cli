import { describe, it, expect, beforeEach } from "vitest";
import { TestBed } from "@angular/core/testing";
import { firstValueFrom } from "rxjs";
import { OpenFeature, InMemoryProvider } from "@openfeature/web-sdk";
import { GeneratedFeatureFlagService } from "../generated/openfeature.generated";
import { v4 as uuid } from "uuid";

describe("GeneratedFeatureFlagService", () => {
  let service: GeneratedFeatureFlagService;
  let domain: string;
  let provider: InMemoryProvider;

  beforeEach(async () => {
    domain = uuid();
    provider = new InMemoryProvider({
      enableFeatureA: {
        variants: { on: true, off: false },
        defaultVariant: "on",
        disabled: false,
      },
      greetingMessage: {
        variants: { default: "Hello from provider!" },
        defaultVariant: "default",
        disabled: false,
      },
      discountPercentage: {
        variants: { default: 0.25 },
        defaultVariant: "default",
        disabled: false,
      },
      usernameMaxLength: {
        variants: { default: 100 },
        defaultVariant: "default",
        disabled: false,
      },
      themeCustomization: {
        variants: {
          default: {
            primaryColor: "#ff0000",
            secondaryColor: "#00ff00",
          },
        },
        defaultVariant: "default",
        disabled: false,
      },
    });

    await OpenFeature.setProviderAndWait(domain, provider);

    TestBed.configureTestingModule({
      providers: [GeneratedFeatureFlagService],
    });

    service = TestBed.inject(GeneratedFeatureFlagService);
  });

  describe("getEnableFeatureADetails", () => {
    it("should return boolean evaluation details", async () => {
      const details = await firstValueFrom(
        service.getEnableFeatureADetails(domain),
      );

      expect(details).toBeDefined();
      expect(typeof details.value).toBe("boolean");
      expect(details.flagKey).toBe("enableFeatureA");
    });

    it("should return true when provider returns true", async () => {
      const details = await firstValueFrom(
        service.getEnableFeatureADetails(domain),
      );

      expect(details.value).toBe(true);
    });

    it("should return default value when flag is not configured", async () => {
      const emptyDomain = uuid();
      const emptyProvider = new InMemoryProvider({});
      await OpenFeature.setProviderAndWait(emptyDomain, emptyProvider);

      const details = await firstValueFrom(
        service.getEnableFeatureADetails(emptyDomain),
      );

      // Default value from manifest is false
      expect(details.value).toBe(false);
    });
  });

  describe("getGreetingMessageDetails", () => {
    it("should return string evaluation details", async () => {
      const details = await firstValueFrom(
        service.getGreetingMessageDetails(domain),
      );

      expect(details).toBeDefined();
      expect(typeof details.value).toBe("string");
      expect(details.flagKey).toBe("greetingMessage");
    });

    it("should return provider value", async () => {
      const details = await firstValueFrom(
        service.getGreetingMessageDetails(domain),
      );

      expect(details.value).toBe("Hello from provider!");
    });

    it("should return default value when flag is not configured", async () => {
      const emptyDomain = uuid();
      const emptyProvider = new InMemoryProvider({});
      await OpenFeature.setProviderAndWait(emptyDomain, emptyProvider);

      const details = await firstValueFrom(
        service.getGreetingMessageDetails(emptyDomain),
      );

      // Default value from manifest
      expect(details.value).toBe("Hello there!");
    });
  });

  describe("getDiscountPercentageDetails", () => {
    it("should return number evaluation details", async () => {
      const details = await firstValueFrom(
        service.getDiscountPercentageDetails(domain),
      );

      expect(details).toBeDefined();
      expect(typeof details.value).toBe("number");
      expect(details.flagKey).toBe("discountPercentage");
    });

    it("should return provider value", async () => {
      const details = await firstValueFrom(
        service.getDiscountPercentageDetails(domain),
      );

      expect(details.value).toBe(0.25);
    });

    it("should return default value when flag is not configured", async () => {
      const emptyDomain = uuid();
      const emptyProvider = new InMemoryProvider({});
      await OpenFeature.setProviderAndWait(emptyDomain, emptyProvider);

      const details = await firstValueFrom(
        service.getDiscountPercentageDetails(emptyDomain),
      );

      // Default value from manifest
      expect(details.value).toBe(0.15);
    });
  });

  describe("getUsernameMaxLengthDetails", () => {
    it("should return number evaluation details", async () => {
      const details = await firstValueFrom(
        service.getUsernameMaxLengthDetails(domain),
      );

      expect(details).toBeDefined();
      expect(typeof details.value).toBe("number");
      expect(details.flagKey).toBe("usernameMaxLength");
    });

    it("should return provider value", async () => {
      const details = await firstValueFrom(
        service.getUsernameMaxLengthDetails(domain),
      );

      expect(details.value).toBe(100);
    });
  });

  describe("getThemeCustomizationDetails", () => {
    it("should return object evaluation details", async () => {
      const details = await firstValueFrom(
        service.getThemeCustomizationDetails(domain),
      );

      expect(details).toBeDefined();
      expect(typeof details.value).toBe("object");
      expect(details.flagKey).toBe("themeCustomization");
    });

    it("should return provider value", async () => {
      const details = await firstValueFrom(
        service.getThemeCustomizationDetails(domain),
      );

      expect(details.value).toEqual({
        primaryColor: "#ff0000",
        secondaryColor: "#00ff00",
      });
    });

    it("should return default value when flag is not configured", async () => {
      const emptyDomain = uuid();
      const emptyProvider = new InMemoryProvider({});
      await OpenFeature.setProviderAndWait(emptyDomain, emptyProvider);

      const details = await firstValueFrom(
        service.getThemeCustomizationDetails(emptyDomain),
      );

      // Default value from manifest
      expect(details.value).toEqual({
        primaryColor: "#007bff",
        secondaryColor: "#6c757d",
      });
    });
  });
});
