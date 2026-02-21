package dev.openfeature;

import dev.openfeature.generated.*;
import dev.openfeature.sdk.*;
import dev.openfeature.sdk.providers.memory.Flag;
import dev.openfeature.sdk.providers.memory.InMemoryProvider;

import java.util.HashMap;
import java.util.Map;

public class Main {
    public static void main(String[] args) {
        try {
            run();
            System.out.println("Generated Java code compiles successfully!");
        } catch (Exception e) {
            System.err.println("Error: " + e.getMessage());
            e.printStackTrace();
            System.exit(1);
        }
    }

    private static void run() throws Exception {
        // Set up the in-memory provider with test flags
        Map<String, Flag<?>> flags = new HashMap<>();

        flags.put("discountPercentage", Flag.builder()
            .variant("default", 0.15)
            .defaultVariant("default")
            .build());

        flags.put("enableFeatureA", Flag.builder()
            .variant("default", false)
            .defaultVariant("default")
            .build());

        flags.put("greetingMessage", Flag.builder()
            .variant("default", "Hello there!")
            .defaultVariant("default")
            .build());

        flags.put("usernameMaxLength", Flag.builder()
            .variant("default", 50)
            .defaultVariant("default")
            .build());

        Map<String, Object> themeConfig = new HashMap<>();
        themeConfig.put("primaryColor", "#007bff");
        themeConfig.put("secondaryColor", "#6c757d");

        flags.put("themeCustomization", Flag.builder()
            .variant("default", new Value(themeConfig))
            .defaultVariant("default")
            .build());

        InMemoryProvider provider = new InMemoryProvider(flags);

        // Set the provider
        OpenFeatureAPI.getInstance().setProviderAndWait(provider);

        Client client = OpenFeatureAPI.getInstance().getClient();
        MutableContext evalContext = new MutableContext();

        // Use the generated code for all flag evaluations
        Boolean enableFeatureA = EnableFeatureA.value(client, evalContext);
        System.out.println("enableFeatureA: " + enableFeatureA);
        FlagEvaluationDetails<Boolean> enableFeatureADetails = EnableFeatureA.valueWithDetails(client, evalContext);
        if (enableFeatureADetails.getErrorCode() != null) {
            throw new Exception("Error evaluating boolean flag");
        }

        Double discount = DiscountPercentage.value(client, evalContext);
        System.out.printf("Discount Percentage: %.2f%n", discount);
        FlagEvaluationDetails<Double> discountDetails = DiscountPercentage.valueWithDetails(client, evalContext);
        if (discountDetails.getErrorCode() != null) {
            throw new Exception("Failed to get discount");
        }

        String greetingMessage = GreetingMessage.value(client, evalContext);
        System.out.println("greetingMessage: " + greetingMessage);
        FlagEvaluationDetails<String> greetingDetails = GreetingMessage.valueWithDetails(client, evalContext);
        if (greetingDetails.getErrorCode() != null) {
            throw new Exception("Error evaluating string flag");
        }

        Integer usernameMaxLength = UsernameMaxLength.value(client, evalContext);
        System.out.println("usernameMaxLength: " + usernameMaxLength);
        FlagEvaluationDetails<Integer> usernameDetails = UsernameMaxLength.valueWithDetails(client, evalContext);
        if (usernameDetails.getErrorCode() != null) {
            throw new Exception("Error evaluating int flag");
        }

        Value themeCustomization = ThemeCustomization.value(client, evalContext);
        FlagEvaluationDetails<Value> themeDetails = ThemeCustomization.valueWithDetails(client, evalContext);
        if (themeDetails.getErrorCode() != null) {
            throw new Exception("Error evaluating object flag");
        }
        System.out.println("themeCustomization: " + themeCustomization);

        // Test the getKey() method functionality for all flags
        System.out.println("enableFeatureA flag key: " + EnableFeatureA.getKey());
        System.out.println("discountPercentage flag key: " + DiscountPercentage.getKey());
        System.out.println("greetingMessage flag key: " + GreetingMessage.getKey());
        System.out.println("usernameMaxLength flag key: " + UsernameMaxLength.getKey());
        System.out.println("themeCustomization flag key: " + ThemeCustomization.getKey());
    }
}
