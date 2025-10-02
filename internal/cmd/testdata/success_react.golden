'use client';

import {
  type ReactFlagEvaluationOptions,
  type ReactFlagEvaluationNoSuspenseOptions,
  type FlagQuery,
  useFlag,
  useSuspenseFlag,
  JsonValue
} from "@openfeature/react-sdk";

// Flag key constants for programmatic access
export const FlagKeys = {
  /** Flag key for Discount percentage applied to purchases. */
  DISCOUNT_PERCENTAGE: "discountPercentage",
  /** Flag key for Controls whether Feature A is enabled. */
  ENABLE_FEATURE_A: "enableFeatureA",
  /** Flag key for The message to use for greeting users. */
  GREETING_MESSAGE: "greetingMessage",
  /** Flag key for Allows customization of theme colors. */
  THEME_CUSTOMIZATION: "themeCustomization",
  /** Flag key for Maximum allowed length for usernames. */
  USERNAME_MAX_LENGTH: "usernameMaxLength",
} as const;


/**
* Discount percentage applied to purchases.
* 
* **Details:**
* - flag key: `discountPercentage`
* - default value: `0.15`
* - type: `number`
*/
export const useDiscountPercentage = (options?: ReactFlagEvaluationOptions): FlagQuery<number> => {
  return useFlag("discountPercentage", 0.15, options);
};

/**
* Discount percentage applied to purchases.
* 
* **Details:**
* - flag key: `discountPercentage`
* - default value: `0.15`
* - type: `number`
*
* Equivalent to useFlag with options: `{ suspend: true }`
* @experimental — Suspense is an experimental feature subject to change in future versions.
*/
export const useSuspenseDiscountPercentage = (options?: ReactFlagEvaluationNoSuspenseOptions): FlagQuery<number> => {
  return useSuspenseFlag("discountPercentage", 0.15, options);
};

/**
* Controls whether Feature A is enabled.
* 
* **Details:**
* - flag key: `enableFeatureA`
* - default value: `false`
* - type: `boolean`
*/
export const useEnableFeatureA = (options?: ReactFlagEvaluationOptions): FlagQuery<boolean> => {
  return useFlag("enableFeatureA", false, options);
};

/**
* Controls whether Feature A is enabled.
* 
* **Details:**
* - flag key: `enableFeatureA`
* - default value: `false`
* - type: `boolean`
*
* Equivalent to useFlag with options: `{ suspend: true }`
* @experimental — Suspense is an experimental feature subject to change in future versions.
*/
export const useSuspenseEnableFeatureA = (options?: ReactFlagEvaluationNoSuspenseOptions): FlagQuery<boolean> => {
  return useSuspenseFlag("enableFeatureA", false, options);
};

/**
* The message to use for greeting users.
* 
* **Details:**
* - flag key: `greetingMessage`
* - default value: `Hello there!`
* - type: `string`
*/
export const useGreetingMessage = (options?: ReactFlagEvaluationOptions): FlagQuery<string> => {
  return useFlag("greetingMessage", "Hello there!", options);
};

/**
* The message to use for greeting users.
* 
* **Details:**
* - flag key: `greetingMessage`
* - default value: `Hello there!`
* - type: `string`
*
* Equivalent to useFlag with options: `{ suspend: true }`
* @experimental — Suspense is an experimental feature subject to change in future versions.
*/
export const useSuspenseGreetingMessage = (options?: ReactFlagEvaluationNoSuspenseOptions): FlagQuery<string> => {
  return useSuspenseFlag("greetingMessage", "Hello there!", options);
};

/**
* Allows customization of theme colors.
* 
* **Details:**
* - flag key: `themeCustomization`
* - default value: `{"primaryColor":"#007bff","secondaryColor":"#6c757d"}`
* - type: `JsonValue`
*/
export const useThemeCustomization = (options?: ReactFlagEvaluationOptions): FlagQuery<JsonValue> => {
  return useFlag("themeCustomization", {"primaryColor":"#007bff","secondaryColor":"#6c757d"}, options);
};

/**
* Allows customization of theme colors.
* 
* **Details:**
* - flag key: `themeCustomization`
* - default value: `{"primaryColor":"#007bff","secondaryColor":"#6c757d"}`
* - type: `JsonValue`
*
* Equivalent to useFlag with options: `{ suspend: true }`
* @experimental — Suspense is an experimental feature subject to change in future versions.
*/
export const useSuspenseThemeCustomization = (options?: ReactFlagEvaluationNoSuspenseOptions): FlagQuery<JsonValue> => {
  return useSuspenseFlag("themeCustomization", {"primaryColor":"#007bff","secondaryColor":"#6c757d"}, options);
};

/**
* Maximum allowed length for usernames.
* 
* **Details:**
* - flag key: `usernameMaxLength`
* - default value: `50`
* - type: `number`
*/
export const useUsernameMaxLength = (options?: ReactFlagEvaluationOptions): FlagQuery<number> => {
  return useFlag("usernameMaxLength", 50, options);
};

/**
* Maximum allowed length for usernames.
* 
* **Details:**
* - flag key: `usernameMaxLength`
* - default value: `50`
* - type: `number`
*
* Equivalent to useFlag with options: `{ suspend: true }`
* @experimental — Suspense is an experimental feature subject to change in future versions.
*/
export const useSuspenseUsernameMaxLength = (options?: ReactFlagEvaluationNoSuspenseOptions): FlagQuery<number> => {
  return useSuspenseFlag("usernameMaxLength", 50, options);
};
