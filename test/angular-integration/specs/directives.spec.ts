import { describe, it, expect, beforeEach, afterEach } from "vitest";
import { Component, Input } from "@angular/core";
import { JsonPipe } from "@angular/common";
import { TestBed, ComponentFixture } from "@angular/core/testing";
import { OpenFeature, InMemoryProvider } from "@openfeature/web-sdk";
import { v4 as uuid } from "uuid";
import {
  GeneratedFeatureFlagDirectives,
  EnableFeatureADirective,
  GreetingMessageDirective,
  DiscountPercentageDirective,
  UsernameMaxLengthDirective,
  ThemeCustomizationDirective,
} from "../generated/openfeature.generated";

// Test component for boolean directive with domain input
@Component({
  selector: "test-boolean",
  standalone: true,
  imports: [EnableFeatureADirective],
  template: `
    <div class="container">
      <ng-container *enableFeatureA="let value; domain: domain">
        <div class="flag-content">Feature A is enabled</div>
      </ng-container>
    </div>
  `,
})
class TestBooleanComponent {
  @Input() domain?: string;
}

// Test component for string directive with domain input
@Component({
  selector: "test-string",
  standalone: true,
  imports: [GreetingMessageDirective],
  template: `
    <div class="container">
      <ng-container *greetingMessage="let value; domain: domain">
        <div class="flag-content">Greeting: {{ value }}</div>
      </ng-container>
    </div>
  `,
})
class TestStringComponent {
  @Input() domain?: string;
}

// Test component for number directive with domain input
@Component({
  selector: "test-number",
  standalone: true,
  imports: [DiscountPercentageDirective],
  template: `
    <div class="container">
      <ng-container *discountPercentage="let value; domain: domain">
        <div class="flag-content">Discount: {{ value }}</div>
      </ng-container>
    </div>
  `,
})
class TestNumberComponent {
  @Input() domain?: string;
}

// Test component for username max length directive with domain input
@Component({
  selector: "test-username",
  standalone: true,
  imports: [UsernameMaxLengthDirective],
  template: `
    <div class="container">
      <ng-container *usernameMaxLength="let value; domain: domain">
        <div class="flag-content">Max Length: {{ value }}</div>
      </ng-container>
    </div>
  `,
})
class TestUsernameComponent {
  @Input() domain?: string;
}

// Test component for object directive with domain input
@Component({
  selector: "test-object",
  standalone: true,
  imports: [ThemeCustomizationDirective, JsonPipe],
  template: `
    <div class="container">
      <ng-container *themeCustomization="let value; domain: domain">
        <div class="flag-content">Theme: {{ value | json }}</div>
      </ng-container>
    </div>
  `,
})
class TestObjectComponent {
  @Input() domain?: string;
}

// Test component using all directives via GeneratedFeatureFlagDirectives array
@Component({
  selector: "test-all",
  standalone: true,
  imports: [GeneratedFeatureFlagDirectives],
  template: `
    <div class="boolean-flag">
      <ng-container *enableFeatureA="let v; domain: domain">
        <div class="flag-content">Boolean flag: {{ v }}</div>
      </ng-container>
    </div>
    <div class="string-flag">
      <ng-container *greetingMessage="let v; domain: domain">
        <div class="flag-content">String flag: {{ v }}</div>
      </ng-container>
    </div>
    <div class="number-flag">
      <ng-container *discountPercentage="let v; domain: domain">
        <div class="flag-content">Number flag: {{ v }}</div>
      </ng-container>
    </div>
    <div class="username-flag">
      <ng-container *usernameMaxLength="let v; domain: domain">
        <div class="flag-content">Username flag: {{ v }}</div>
      </ng-container>
    </div>
    <div class="object-flag">
      <ng-container *themeCustomization="let v; domain: domain">
        <div class="flag-content">Object flag</div>
      </ng-container>
    </div>
  `,
})
class TestAllDirectivesComponent {
  @Input() domain?: string;
}

describe("Generated Directives Tests", () => {
  let provider: InMemoryProvider;
  let domain: string;

  beforeEach(async () => {
    // Use unique domain for each test to ensure isolation
    domain = uuid();
    // Reset TestBed to avoid shared state
    TestBed.resetTestingModule();
  });

  afterEach(async () => {
    await OpenFeature.clearProviders();
    TestBed.resetTestingModule();
  });

  describe("EnableFeatureADirective (boolean)", () => {
    it("should render content when flag is true", async () => {
      provider = new InMemoryProvider({
        enableFeatureA: {
          variants: { on: true },
          defaultVariant: "on",
          disabled: false,
        },
      });
      await OpenFeature.setProviderAndWait(domain, provider);

      const fixture = TestBed.configureTestingModule({
        imports: [TestBooleanComponent],
      }).createComponent(TestBooleanComponent);
      fixture.componentRef.setInput("domain", domain);
      fixture.detectChanges();
      await fixture.whenStable();

      const content = fixture.nativeElement.querySelector(".flag-content");
      expect(content).not.toBeNull();
      expect(content.textContent.trim()).toBe("Feature A is enabled");
    });

    it("should NOT render content when flag is false", async () => {
      provider = new InMemoryProvider({
        enableFeatureA: {
          variants: { off: false },
          defaultVariant: "off",
          disabled: false,
        },
      });
      await OpenFeature.setProviderAndWait(domain, provider);

      const fixture = TestBed.configureTestingModule({
        imports: [TestBooleanComponent],
      }).createComponent(TestBooleanComponent);
      fixture.componentRef.setInput("domain", domain);
      fixture.detectChanges();
      await fixture.whenStable();

      const content = fixture.nativeElement.querySelector(".flag-content");
      expect(content).toBeNull();
    });

    it("should use default value (false) when flag is not configured", async () => {
      provider = new InMemoryProvider({});
      await OpenFeature.setProviderAndWait(domain, provider);

      const fixture = TestBed.configureTestingModule({
        imports: [TestBooleanComponent],
      }).createComponent(TestBooleanComponent);
      fixture.componentRef.setInput("domain", domain);
      fixture.detectChanges();
      await fixture.whenStable();

      // Default is false, so content should not render
      const content = fixture.nativeElement.querySelector(".flag-content");
      expect(content).toBeNull();
    });
  });

  describe("GreetingMessageDirective (string)", () => {
    it("should render content and display custom value from provider", async () => {
      provider = new InMemoryProvider({
        greetingMessage: {
          variants: { custom: "Custom greeting from provider!" },
          defaultVariant: "custom",
          disabled: false,
        },
      });
      await OpenFeature.setProviderAndWait(domain, provider);

      const fixture = TestBed.configureTestingModule({
        imports: [TestStringComponent],
      }).createComponent(TestStringComponent);
      fixture.componentRef.setInput("domain", domain);
      fixture.detectChanges();
      await fixture.whenStable();

      const content = fixture.nativeElement.querySelector(".flag-content");
      expect(content).not.toBeNull();
      expect(content.textContent.trim()).toBe(
        "Greeting: Custom greeting from provider!",
      );
    });

    it("should render content with default value when flag not configured", async () => {
      provider = new InMemoryProvider({});
      await OpenFeature.setProviderAndWait(domain, provider);

      const fixture = TestBed.configureTestingModule({
        imports: [TestStringComponent],
      }).createComponent(TestStringComponent);
      fixture.componentRef.setInput("domain", domain);
      fixture.detectChanges();
      await fixture.whenStable();

      // Default is "Hello there!"
      const content = fixture.nativeElement.querySelector(".flag-content");
      expect(content).not.toBeNull();
      expect(content.textContent.trim()).toBe("Greeting: Hello there!");
    });
  });

  describe("DiscountPercentageDirective (number)", () => {
    it("should render content and display custom value from provider", async () => {
      provider = new InMemoryProvider({
        discountPercentage: {
          variants: { custom: 0.42 },
          defaultVariant: "custom",
          disabled: false,
        },
      });
      await OpenFeature.setProviderAndWait(domain, provider);

      const fixture = TestBed.configureTestingModule({
        imports: [TestNumberComponent],
      }).createComponent(TestNumberComponent);
      fixture.componentRef.setInput("domain", domain);
      fixture.detectChanges();
      await fixture.whenStable();

      const content = fixture.nativeElement.querySelector(".flag-content");
      expect(content).not.toBeNull();
      expect(content.textContent.trim()).toBe("Discount: 0.42");
    });

    it("should render content with default value when flag not configured", async () => {
      provider = new InMemoryProvider({});
      await OpenFeature.setProviderAndWait(domain, provider);

      const fixture = TestBed.configureTestingModule({
        imports: [TestNumberComponent],
      }).createComponent(TestNumberComponent);
      fixture.componentRef.setInput("domain", domain);
      fixture.detectChanges();
      await fixture.whenStable();

      // Default is 0.15
      const content = fixture.nativeElement.querySelector(".flag-content");
      expect(content).not.toBeNull();
      expect(content.textContent.trim()).toBe("Discount: 0.15");
    });
  });

  describe("UsernameMaxLengthDirective (number)", () => {
    it("should render content and display custom value from provider", async () => {
      provider = new InMemoryProvider({
        usernameMaxLength: {
          variants: { custom: 100 },
          defaultVariant: "custom",
          disabled: false,
        },
      });
      await OpenFeature.setProviderAndWait(domain, provider);

      const fixture = TestBed.configureTestingModule({
        imports: [TestUsernameComponent],
      }).createComponent(TestUsernameComponent);
      fixture.componentRef.setInput("domain", domain);
      fixture.detectChanges();
      await fixture.whenStable();

      const content = fixture.nativeElement.querySelector(".flag-content");
      expect(content).not.toBeNull();
      expect(content.textContent.trim()).toBe("Max Length: 100");
    });

    it("should render content with default value when flag not configured", async () => {
      provider = new InMemoryProvider({});
      await OpenFeature.setProviderAndWait(domain, provider);

      const fixture = TestBed.configureTestingModule({
        imports: [TestUsernameComponent],
      }).createComponent(TestUsernameComponent);
      fixture.componentRef.setInput("domain", domain);
      fixture.detectChanges();
      await fixture.whenStable();

      // Default is 50
      const content = fixture.nativeElement.querySelector(".flag-content");
      expect(content).not.toBeNull();
      expect(content.textContent.trim()).toBe("Max Length: 50");
    });
  });

  describe("ThemeCustomizationDirective (object)", () => {
    it("should render content and display custom object from provider", async () => {
      provider = new InMemoryProvider({
        themeCustomization: {
          variants: {
            custom: { primaryColor: "#ff0000", secondaryColor: "#00ff00" },
          },
          defaultVariant: "custom",
          disabled: false,
        },
      });
      await OpenFeature.setProviderAndWait(domain, provider);

      const fixture = TestBed.configureTestingModule({
        imports: [TestObjectComponent],
      }).createComponent(TestObjectComponent);
      fixture.componentRef.setInput("domain", domain);
      fixture.detectChanges();
      await fixture.whenStable();

      const content = fixture.nativeElement.querySelector(".flag-content");
      expect(content).not.toBeNull();
      // Verify the custom object is displayed
      expect(content.textContent).toContain("#ff0000");
      expect(content.textContent).toContain("#00ff00");
    });

    it("should render content with default object when flag not configured", async () => {
      provider = new InMemoryProvider({});
      await OpenFeature.setProviderAndWait(domain, provider);

      const fixture = TestBed.configureTestingModule({
        imports: [TestObjectComponent],
      }).createComponent(TestObjectComponent);
      fixture.componentRef.setInput("domain", domain);
      fixture.detectChanges();
      await fixture.whenStable();

      // Default is {"primaryColor":"#007bff","secondaryColor":"#6c757d"}
      const content = fixture.nativeElement.querySelector(".flag-content");
      expect(content).not.toBeNull();
      expect(content.textContent).toContain("#007bff");
      expect(content.textContent).toContain("#6c757d");
    });
  });

  describe("GeneratedFeatureFlagDirectives array", () => {
    it("should render all directive types when flags are enabled/have values", async () => {
      provider = new InMemoryProvider({
        enableFeatureA: {
          variants: { on: true },
          defaultVariant: "on",
          disabled: false,
        },
        greetingMessage: {
          variants: { default: "Hello" },
          defaultVariant: "default",
          disabled: false,
        },
        discountPercentage: {
          variants: { default: 0.1 },
          defaultVariant: "default",
          disabled: false,
        },
        usernameMaxLength: {
          variants: { default: 50 },
          defaultVariant: "default",
          disabled: false,
        },
        themeCustomization: {
          variants: { default: {} },
          defaultVariant: "default",
          disabled: false,
        },
      });
      await OpenFeature.setProviderAndWait(domain, provider);

      const fixture = TestBed.configureTestingModule({
        imports: [TestAllDirectivesComponent],
      }).createComponent(TestAllDirectivesComponent);
      fixture.componentRef.setInput("domain", domain);
      fixture.detectChanges();
      await fixture.whenStable();

      // All flags should render
      expect(
        fixture.nativeElement.querySelector(".boolean-flag .flag-content"),
      ).not.toBeNull();
      expect(
        fixture.nativeElement.querySelector(".string-flag .flag-content"),
      ).not.toBeNull();
      expect(
        fixture.nativeElement.querySelector(".number-flag .flag-content"),
      ).not.toBeNull();
      expect(
        fixture.nativeElement.querySelector(".username-flag .flag-content"),
      ).not.toBeNull();
      expect(
        fixture.nativeElement.querySelector(".object-flag .flag-content"),
      ).not.toBeNull();
    });

    it("should hide boolean directive when flag is false", async () => {
      provider = new InMemoryProvider({
        enableFeatureA: {
          variants: { off: false },
          defaultVariant: "off",
          disabled: false,
        },
        greetingMessage: {
          variants: { default: "Bonjour" },
          defaultVariant: "default",
          disabled: false,
        },
        discountPercentage: {
          variants: { default: 0.5 },
          defaultVariant: "default",
          disabled: false,
        },
        usernameMaxLength: {
          variants: { default: 100 },
          defaultVariant: "default",
          disabled: false,
        },
        themeCustomization: {
          variants: { default: { color: "blue" } },
          defaultVariant: "default",
          disabled: false,
        },
      });
      await OpenFeature.setProviderAndWait(domain, provider);

      const fixture = TestBed.configureTestingModule({
        imports: [TestAllDirectivesComponent],
      }).createComponent(TestAllDirectivesComponent);
      fixture.componentRef.setInput("domain", domain);
      fixture.detectChanges();
      await fixture.whenStable();

      // Boolean is false, so it should not render
      expect(
        fixture.nativeElement.querySelector(".boolean-flag .flag-content"),
      ).toBeNull();

      // Verify the other contents are present with their values
      const stringContent = fixture.nativeElement.querySelector(
        ".string-flag .flag-content",
      );
      expect(stringContent).not.toBeNull();
      expect(stringContent.textContent).toContain("Bonjour");

      const numberContent = fixture.nativeElement.querySelector(
        ".number-flag .flag-content",
      );
      expect(numberContent).not.toBeNull();
      expect(numberContent.textContent).toContain("0.5");

      const usernameContent = fixture.nativeElement.querySelector(
        ".username-flag .flag-content",
      );
      expect(usernameContent).not.toBeNull();
      expect(usernameContent.textContent).toContain("100");

      expect(
        fixture.nativeElement.querySelector(".object-flag .flag-content"),
      ).not.toBeNull();
    });
  });
});
