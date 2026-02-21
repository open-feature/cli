import * as generated from './generated';

async function main() {
  try {
    // Validate that all generated exports exist and have the expected structure
    const flags = [
      { name: 'EnableFeatureA', flag: generated.EnableFeatureA },
      { name: 'DiscountPercentage', flag: generated.DiscountPercentage },
      { name: 'GreetingMessage', flag: generated.GreetingMessage },
      { name: 'UsernameMaxLength', flag: generated.UsernameMaxLength },
      { name: 'ThemeCustomization', flag: generated.ThemeCustomization },
    ];

    for (const { name, flag } of flags) {
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

      if (typeof key !== 'string' || key.length === 0) {
        throw new Error(`${name}.getKey() did not return a valid string`);
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

    // Verify expected flag keys
    if (generated.EnableFeatureA.getKey() !== 'enableFeatureA') {
      throw new Error('EnableFeatureA has incorrect key');
    }
    if (generated.DiscountPercentage.getKey() !== 'discountPercentage') {
      throw new Error('DiscountPercentage has incorrect key');
    }
    if (generated.GreetingMessage.getKey() !== 'greetingMessage') {
      throw new Error('GreetingMessage has incorrect key');
    }
    if (generated.UsernameMaxLength.getKey() !== 'usernameMaxLength') {
      throw new Error('UsernameMaxLength has incorrect key');
    }
    if (generated.ThemeCustomization.getKey() !== 'themeCustomization') {
      throw new Error('ThemeCustomization has incorrect key');
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
