import { NestFactory } from '@nestjs/core';
import { Module, Injectable } from '@nestjs/common';
import { OpenFeatureModule, OPENFEATURE_CLIENT } from '@openfeature/nestjs-sdk';
import { InMemoryProvider, Client } from '@openfeature/server-sdk';
import * as generated from './generated/openfeature';
import { GeneratedOpenFeatureModule } from './generated/openfeature-module';

// Service that uses generated decorators to test NestJS-specific functionality
@Injectable()
class TestService {
  constructor(
    @generated.EnableFeatureA() private enableFeatureA: boolean,
    @generated.DiscountPercentage() private discountPercentage: number,
    @generated.GreetingMessage() private greetingMessage: string,
    @generated.UsernameMaxLength() private usernameMaxLength: number,
    @generated.ThemeCustomization() private themeCustomization: any,
  ) {}

  getFlags() {
    return {
      enableFeatureA: this.enableFeatureA,
      discountPercentage: this.discountPercentage,
      greetingMessage: this.greetingMessage,
      usernameMaxLength: this.usernameMaxLength,
      themeCustomization: this.themeCustomization,
    };
  }
}

@Module({
  imports: [
    GeneratedOpenFeatureModule.forRoot({
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
  providers: [TestService],
})
class AppModule {}

async function bootstrap() {
  const app = await NestFactory.createApplicationContext(AppModule);

  try {
    const client = app.get<Client>(OPENFEATURE_CLIENT);
    const testService = app.get(TestService);

    // Test NestJS decorators by getting flags from the service
    const flagsFromDecorators = testService.getFlags();
    console.log('Flags from NestJS decorators:');
    console.log('  enableFeatureA:', flagsFromDecorators.enableFeatureA);
    console.log('  discountPercentage:', flagsFromDecorators.discountPercentage.toFixed(2));
    console.log('  greetingMessage:', flagsFromDecorators.greetingMessage);
    console.log('  usernameMaxLength:', flagsFromDecorators.usernameMaxLength);
    console.log('  themeCustomization:', flagsFromDecorators.themeCustomization);

    // Use the generated code for flag evaluations with client
    const enableFeatureA = await generated.EnableFeatureA.value(client, {});
    console.log('\nDirect flag evaluation:');
    console.log('  enableFeatureA:', enableFeatureA);

    const enableFeatureADetails = await generated.EnableFeatureA.valueWithDetails(client, {});
    if (enableFeatureADetails.errorCode) {
      throw new Error('Error evaluating boolean flag');
    }

    const discount = await generated.DiscountPercentage.value(client, {});
    console.log('  Discount Percentage:', discount.toFixed(2));

    const discountDetails = await generated.DiscountPercentage.valueWithDetails(client, {});
    if (discountDetails.errorCode) {
      throw new Error('Failed to get discount');
    }

    const greetingMessage = await generated.GreetingMessage.value(client, {});
    console.log('  greetingMessage:', greetingMessage);

    const greetingDetails = await generated.GreetingMessage.valueWithDetails(client, {});
    if (greetingDetails.errorCode) {
      throw new Error('Error evaluating string flag');
    }

    const usernameMaxLength = await generated.UsernameMaxLength.value(client, {});
    console.log('  usernameMaxLength:', usernameMaxLength);

    const usernameDetails = await generated.UsernameMaxLength.valueWithDetails(client, {});
    if (usernameDetails.errorCode) {
      throw new Error('Error evaluating int flag');
    }

    const themeCustomization = await generated.ThemeCustomization.value(client, {});
    console.log('  themeCustomization:', themeCustomization);

    const themeDetails = await generated.ThemeCustomization.valueWithDetails(client, {});
    if (themeDetails.errorCode) {
      throw new Error('Error evaluating object flag');
    }

    // Test the getKey() method functionality for all flags
    console.log('\nFlag keys:');
    console.log('  enableFeatureA flag key:', generated.EnableFeatureA.getKey());
    console.log('  discountPercentage flag key:', generated.DiscountPercentage.getKey());
    console.log('  greetingMessage flag key:', generated.GreetingMessage.getKey());
    console.log('  usernameMaxLength flag key:', generated.UsernameMaxLength.getKey());
    console.log('  themeCustomization flag key:', generated.ThemeCustomization.getKey());

    console.log('\n✅ Generated NestJS code compiles successfully!');
    console.log('✅ NestJS decorators work correctly!');
    console.log('✅ GeneratedOpenFeatureModule integrates properly!');

    await app.close();
    process.exit(0);
  } catch (error) {
    console.error('Error:', error);
    await app.close();
    process.exit(1);
  }
}

bootstrap();
