import { describe, it, expect } from "vitest";
import fs from "node:fs";
import os from "node:os";
import path from "node:path";
import ts from "typescript";
import { fileURLToPath } from "node:url";
import { NgtscProgram, createCompilerHost } from "@angular/compiler-cli";

describe("ngtsc compilation", () => {
  it("should compile structural directive usage without requiring a default input", () => {
    const baseDir = path.resolve(
      path.dirname(fileURLToPath(import.meta.url)),
      "..",
    );
    const tmpDir = fs.mkdtempSync(path.join(baseDir, ".tmp-ngtsc-"));
    const componentPath = path.join(tmpDir, "component.ts");
    const tsconfigPath = path.join(tmpDir, "tsconfig.json");

    try {
      const componentSource = `
import { Component } from "@angular/core";
import { EnableFeatureADirective } from "@generated/openfeature.generated";

@Component({
  standalone: true,
  selector: "test-host",
  imports: [EnableFeatureADirective],
  template: \`<div *enableFeatureA></div>\`,
})
export class RequiredInputHostComponent {}
`;

      const tsconfig = {
        compilerOptions: {
          target: "ES2022",
          module: "ESNext",
          moduleResolution: "bundler",
          lib: ["ES2022", "DOM"],
          strict: true,
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
          strictTemplates: true,
        },
        files: [componentPath],
      };

      fs.writeFileSync(componentPath, componentSource, "utf8");
      fs.writeFileSync(tsconfigPath, JSON.stringify(tsconfig, null, 2), "utf8");

      const config = ts.readConfigFile(tsconfigPath, ts.sys.readFile);
      const parsed = ts.parseJsonConfigFileContent(
        config.config,
        ts.sys,
        tmpDir,
        undefined,
        tsconfigPath,
      );

      const host = createCompilerHost({ options: parsed.options });
      const ngProgram = new NgtscProgram(parsed.fileNames, parsed.options, host);

      const diagnostics = [
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
  });
});
