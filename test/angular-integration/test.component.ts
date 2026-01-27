import { Component, Input } from "@angular/core";
import { NgTemplateOutlet } from "@angular/common";
import {
  GeneratedFeatureFlagDirectives,
  GeneratedFeatureFlagService,
} from "./generated/openfeature.generated";

/**
 * Test component that uses all generated feature flag directives.
 * Each test case is wrapped in a container with a unique class for easy querying.
 *
 * Note: Since we use hostDirectives pattern, inputs must be passed as separate
 * attributes rather than using structural directive microsyntax.
 */
@Component({
  selector: "app-test",
  standalone: true,
  imports: [GeneratedFeatureFlagDirectives, NgTemplateOutlet],
  template: `
    <!-- Boolean flag: enableFeatureA - Basic usage -->
    <div class="case-boolean-basic">
      <div *enableFeatureA class="flag-content">Feature A is enabled</div>
    </div>

    <!-- String flag: greetingMessage - Basic usage -->
    <div class="case-string-basic">
      <div *greetingMessage class="flag-content">String flag rendered</div>
    </div>

    <!-- Number flag: discountPercentage - Basic usage -->
    <div class="case-number-basic">
      <div *discountPercentage class="flag-content">Discount is active</div>
    </div>

    <!-- Number flag: usernameMaxLength - Basic usage -->
    <div class="case-username-basic">
      <div *usernameMaxLength class="flag-content">Username limit active</div>
    </div>

    <!-- Object flag: themeCustomization - Basic usage -->
    <div class="case-object-basic">
      <div *themeCustomization class="flag-content">
        Theme customization active
      </div>
    </div>
  `,
})
export class TestComponent {
  @Input() domain?: string;
}

/**
 * Simple test component for service tests (no template directives).
 */
@Component({
  selector: "app-service-test",
  standalone: true,
  template: `<div>Service Test Component</div>`,
})
export class ServiceTestComponent {
  constructor(public flagService: GeneratedFeatureFlagService) {}
}
