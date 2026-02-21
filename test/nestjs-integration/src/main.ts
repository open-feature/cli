import { NestFactory } from '@nestjs/core';
import { Module } from '@nestjs/common';
import { OpenFeatureModule } from '@openfeature/nestjs-sdk';
import { InMemoryProvider } from '@openfeature/server-sdk';
import * as generated from './generated';

@Module({
  imports: [
    OpenFeatureModule.forRoot({
      provider: new InMemoryProvider({
        discountPercentage: {
          disabled: false,
          variants: {
            default: 0.15,
          },
          defaultVariant: 'default',
        },
        enableFeatureA: {
          disabled: false,
          variants: {
            default: false,
          },
          defaultVariant: 'default',
        },
        greetingMessage: {
          disabled: false,
          variants: {
            default: 'Hello there!',
          },
          defaultVariant: 'default',
        },
        usernameMaxLength: {
          disabled: false,
          variants: {
            default: 50,
          },
          defaultVariant: 'default',
        },
        themeCustomization: {
          disabled: false,
          variants: {
            default: {
              primaryColor: '#007bff',
              secondaryColor: '#6c757d',
            },
          },
          defaultVariant: 'default',
        },
      }),
    }),
  ],
})
class AppModule {}

async function bootstrap() {
  const app = await NestFactory.createApplicationContext(AppModule);

  try {
    const client = app.get('OPENFEATURE_CLIENT');

    // Use the generated code for all flag evaluations
    const enableFeatureA = await generated.EnableFeatureA.value(client, {});
    console.log('enableFeatureA:', enableFeatureA);

    const enableFeatureADetails = await generated.EnableFeatureA.valueWithDetails(client, {});
    if (enableFeatureADetails.errorCode) {
      throw new Error('Error evaluating boolean flag');
    }

    const discount = await generated.DiscountPercentage.value(client, {});
    console.log('Discount Percentage:', discount.toFixed(2));

    const discountDetails = await generated.DiscountPercentage.valueWithDetails(client, {});
    if (discountDetails.errorCode) {
      throw new Error('Failed to get discount');
    }

    const greetingMessage = await generated.GreetingMessage.value(client, {});
    console.log('greetingMessage:', greetingMessage);

    const greetingDetails = await generated.GreetingMessage.valueWithDetails(client, {});
    if (greetingDetails.errorCode) {
      throw new Error('Error evaluating string flag');
    }

    const usernameMaxLength = await generated.UsernameMaxLength.value(client, {});
    console.log('usernameMaxLength:', usernameMaxLength);

    const usernameDetails = await generated.UsernameMaxLength.valueWithDetails(client, {});
    if (usernameDetails.errorCode) {
      throw new Error('Error evaluating int flag');
    }

    const themeCustomization = await generated.ThemeCustomization.value(client, {});
    console.log('themeCustomization:', themeCustomization);

    const themeDetails = await generated.ThemeCustomization.valueWithDetails(client, {});
    if (themeDetails.errorCode) {
      throw new Error('Error evaluating object flag');
    }

    // Test the getKey() method functionality for all flags
    console.log('enableFeatureA flag key:', generated.EnableFeatureA.getKey());
    console.log('discountPercentage flag key:', generated.DiscountPercentage.getKey());
    console.log('greetingMessage flag key:', generated.GreetingMessage.getKey());
    console.log('usernameMaxLength flag key:', generated.UsernameMaxLength.getKey());
    console.log('themeCustomization flag key:', generated.ThemeCustomization.getKey());

    console.log('Generated NestJS code compiles successfully!');

    await app.close();
    process.exit(0);
  } catch (error) {
    console.error('Error:', error);
    await app.close();
    process.exit(1);
  }
}

bootstrap();
