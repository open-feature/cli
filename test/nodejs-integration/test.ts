import { OpenFeature } from "@openfeature/server-sdk";
import { existsSync, readdirSync } from 'node:fs';
import { join, dirname } from 'node:path';
import { fileURLToPath } from 'node:url';

const __filename = fileURLToPath(import.meta.url);
const __dirname = dirname(__filename);

// Simple test provider
class TestProvider {
    get name() { return 'test-provider'; }
    get metadata() { return { name: 'test-provider' }; }

    async resolveBooleanEvaluation(flagKey: string, defaultValue: boolean) {
        return { value: flagKey === 'enableFeatureA' ? true : defaultValue, reason: 'STATIC' };
    }

    async resolveStringEvaluation(flagKey: string, defaultValue: string) {
        return { value: flagKey === 'greetingMessage' ? 'Hello from test!' : defaultValue, reason: 'STATIC' };
    }

    async resolveNumberEvaluation(flagKey: string, defaultValue: number) {
        const values: Record<string, number> = { usernameMaxLength: 100, discountPercentage: 0.25 };
        return { value: values[flagKey] ?? defaultValue, reason: 'STATIC' };
    }
    async resolveObjectEvaluation<T>(flagKey: string, defaultValue: T) {
        const values: Record<string, T> = {
            themeCustomization: {
                primaryColor: '#ff0000',
                secondaryColor: '#00ff00'
            } as T
        };
        return { value: values[flagKey] ?? defaultValue, reason: 'STATIC' };
    }
}

async function main() {
    try {
        console.log('🚀 Node.js OpenFeature Integration Test');
        
        // 1. Check generated files
        const generatedDir = join(__dirname, 'generated');
        if (!existsSync(generatedDir)) {
            throw new Error('Generated directory not found');
        }

        const files = readdirSync(generatedDir);
        const clientFile = files.find(file => file.includes('openfeature'));
        if (!clientFile) {
            throw new Error('openfeature.ts not found');
        }
        console.log(`✅ Found: ${clientFile}`);

        
        const clientPath = join(generatedDir, clientFile);
        
        // 3. Setup OpenFeature provider and test
        await OpenFeature.setProvider(new TestProvider());
        
       
        const { getGeneratedClient } = await import(clientPath);
        const client = getGeneratedClient();
        
        console.log('🧪 Testing flags...');
        
        // Test each flag
        const tests = [
            { name: 'enableFeatureA', expected: 'boolean' },
            { name: 'greetingMessage', expected: 'string' },
            { name: 'usernameMaxLength', expected: 'number' },
            { name: 'discountPercentage', expected: 'number' }
        ];

        for (const test of tests) {
            if (client[test.name]) {
                const result = await client[test.name]();
                const type = typeof result;
                if (type === test.expected) {
                    console.log(`✅ ${test.name}: ${result} | type: (${type})`);
                } else {
                    throw new Error(`${test.name} returned ${type}, expected ${test.expected}`);
                }
            } else {
                console.log(`⚠️  ${test.name} method not found`);
                process.exit(1);
            }
        }

        console.log('🎉 All tests passed!');
        process.exit(0);

    } catch (error) {
        console.error('❌ Test failed:', error.message);
        process.exit(1);
    }
}

main();