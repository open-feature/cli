import { describe, it, expect } from "vitest";
import fs from "node:fs";
import os from "node:os";
import path from "node:path";
import ts from "typescript";
import { fileURLToPath } from "node:url";
import {
  NgtscProgram,
  createCompilerHost,
  readConfiguration,
} from "@angular/compiler-cli";

describe("ngtsc compilation", () => {
  const compileAndAssert = (componentSource: string) => {
    const baseDir = path.resolve(
      path.dirname(fileURLToPath(import.meta.url)),
      "..",
    );
    const tmpDir = fs.mkdtempSync(path.join(baseDir, ".tmp-ngtsc-"));
    const componentPath = path.join(tmpDir, "component.ts");
    const tsconfigPath = path.join(tmpDir, "tsconfig.json");

    try {
      const tsconfig = {
        compilerOptions: {
          target: "ES2022",
          module: "ESNext",
          moduleResolution: "bundler",
          lib: ["ES2022", "DOM"],
          strict: true,
          esModuleInterop: true,
          experimentalDecorators: true,
          emitDecoratorMetadata: true,
          useDefineForClassFields: false,
          skipLibCheck: true,
          baseUrl: baseDir,
          paths: {
            "@generated/*": ["generated/*"],
          },
          types: ["node"],
        },
        angularCompilerOptions: {
          enableI18nLegacyMessageIdFormat: false,
          strictInjectionParameters: true,
          strictInputAccessModifiers: true,
          strictTemplates: true,
        },
        files: [componentPath],
      };

      fs.writeFileSync(componentPath, componentSource, "utf8");
      fs.writeFileSync(tsconfigPath, JSON.stringify(tsconfig, null, 2), "utf8");

      const parsed = readConfiguration(tsconfigPath);

      const host = createCompilerHost({ options: parsed.options });
      const ngProgram = new NgtscProgram(
        parsed.rootNames,
        parsed.options,
        host,
      );

      const diagnostics = [
        ...parsed.errors,
        ...ngProgram.getTsOptionDiagnostics(),
        ...ngProgram.getTsSyntacticDiagnostics(),
        ...ngProgram.getTsSemanticDiagnostics(),
        ...ngProgram.getNgOptionDiagnostics(),
        ...ngProgram.getNgStructuralDiagnostics(),
        ...ngProgram.getNgSemanticDiagnostics(),
      ];

      const errors = diagnostics.filter(
        (diag) => diag.category === ts.DiagnosticCategory.Error,
      );

      if (errors.length > 0) {
        const formatted = ts.formatDiagnosticsWithColorAndContext(errors, {
          getCanonicalFileName: (fileName) => fileName,
          getCurrentDirectory: () => tmpDir,
          getNewLine: () => "\n",
        });
        throw new Error(`ngtsc reported errors:\n${formatted}`);
      }
    } finally {
      fs.rmSync(tmpDir, { recursive: true, force: true });
    }
  };

  it("should compile structural directive usage without requiring a default input", () => {
    const componentSource = `
import { Component } from "@angular/core";
import { EnableFeatureAFeatureFlagDirective } from "@generated/openfeature.generated";

@Component({
  standalone: true,
  selector: "test-host",
  imports: [EnableFeatureAFeatureFlagDirective],
  template: \`
    <ng-template enableFeatureA>
      <div>Feature A enabled</div>
    </ng-template>
  \`,
})
export class RequiredInputHostComponent {}
`;
    compileAndAssert(componentSource);
  });

  it("should compile structural directive usage with else template binding", () => {
    const componentSource = `
import { Component } from "@angular/core";
import { EnableFeatureAFeatureFlagDirective } from "@generated/openfeature.generated";

@Component({
  standalone: true,
  selector: "test-host",
  imports: [EnableFeatureAFeatureFlagDirective],
  template: \`
    <ng-template #elseTemplate>Else</ng-template>

    <ng-template
      enableFeatureA
      [enableFeatureAElse]="elseTemplate">
      <div>Feature A enabled</div>
    </ng-template>
  \`,
})
export class ElseTemplateHostComponent {}
`;
    compileAndAssert(componentSource);
  });

  it("should compile simple microsyntax usage without else templates", () => {
    const componentSource = `
import { Component } from "@angular/core";
import { EnableFeatureAFeatureFlagDirective } from "@generated/openfeature.generated";

@Component({
  standalone: true,
  selector: "test-host",
  imports: [EnableFeatureAFeatureFlagDirective],
  template: \`
    <div *enableFeatureA>Feature A enabled</div>
  \`,
})
export class SimpleMicrosyntaxHostComponent {}
`;
    compileAndAssert(componentSource);
  });

  it("should compile structural directive usage with templates and all options", () => {
    const componentSource = `
import { Component } from "@angular/core";
import { GreetingMessageFeatureFlagDirective } from "@generated/openfeature.generated";

@Component({
  standalone: true,
  selector: "test-host",
  imports: [GreetingMessageFeatureFlagDirective],
  template: \`
    <ng-template #elseTemplate>Else</ng-template>
    <ng-template #initTemplate>Init</ng-template>
    <ng-template #reconcilingTemplate>Reconciling</ng-template>

    <div
      *greetingMessage="let value; let details = evaluationDetails; default: expectedValue; else: elseTemplate; initializing: initTemplate; reconciling: reconcilingTemplate">
      Flag value: {{ value }}
    </div>
  \`,
})
export class AllOptionsHostComponent {
  expectedValue = "hello";
}
`;
    compileAndAssert(componentSource);
  });

  it("should compile structural directive usage on a custom component with inputs", () => {
    const componentSource = `
import { Component, Input } from "@angular/core";
import { GreetingMessageFeatureFlagDirective } from "@generated/openfeature.generated";

@Component({
  standalone: true,
  selector: "custom-widget",
  template: "<span>{{ label }}</span>",
})
export class CustomWidgetComponent {
  @Input() label = "";
}

@Component({
  standalone: true,
  selector: "test-host",
  imports: [CustomWidgetComponent, GreetingMessageFeatureFlagDirective],
  template: 
    \`<custom-widget
      *greetingMessage=\"let value; default: expectedValue\"
      [label]=\"value\">
    </custom-widget>\`,
})
export class CustomComponentHostComponent {
  expectedValue = "hello";
}
`;
    compileAndAssert(componentSource);
  });
});
