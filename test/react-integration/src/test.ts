import * as generated from './generated';

async function main() {
  try {
    // Validate that all generated exports exist and have the expected structure
    const flags = [
      { name: 'EnableFeatureA', flag: generated.EnableFeatureA, expectedKey: 'enableFeatureA' },
      { name: 'DiscountPercentage', flag: generated.DiscountPercentage, expectedKey: 'discountPercentage' },
      { name: 'GreetingMessage', flag: generated.GreetingMessage, expectedKey: 'greetingMessage' },
      { name: 'UsernameMaxLength', flag: generated.UsernameMaxLength, expectedKey: 'usernameMaxLength' },
      { name: 'ThemeCustomization', flag: generated.ThemeCustomization, expectedKey: 'themeCustomization' },
    ];

    for (const { name, flag, expectedKey } of flags) {
      // Validate the flag object has the expected properties
      if (typeof flag !== 'object' || flag === null) {
        throw new Error(`${name} is not an object`);
      }

      // Check for getKey method
      if (typeof flag.getKey !== 'function') {
        throw new Error(`${name}.getKey is not a function`);
      }

      const key = flag.getKey();
      console.log(`${name} flag key:`, key);

      if (key !== expectedKey) {
        throw new Error(`${name} has incorrect key. Expected '${expectedKey}', but got '${key}'.`);
      }

      // Check for useFlag hook
      if (typeof flag.useFlag !== 'function') {
        throw new Error(`${name}.useFlag is not a function`);
      }

      // Check for useFlagWithDetails hook
      if (typeof flag.useFlagWithDetails !== 'function') {
        throw new Error(`${name}.useFlagWithDetails is not a function`);
      }
    }

    console.log('All generated React hooks are properly structured!');
    console.log('Generated React code compiles successfully!');
    process.exit(0);
  } catch (error: unknown) {
    const message = error instanceof Error ? error.message : String(error);
    console.error('Error:', message);
    process.exit(1);
  }
}

main();
