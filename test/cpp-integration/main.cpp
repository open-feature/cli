#include <any>
#include <cmath>
#include <cstdint>
#include <iostream>
#include <map>
#include <memory>
#include <optional>
#include <string>
#include <string_view>
#include <unordered_map>

#include "generated/openfeature.h"
#include "openfeature/evaluation_context.h"
#include "openfeature/flag.h"
#include "openfeature/flag_metadata.h"
#include "openfeature/in_memory_provider.h"
#include "openfeature/openfeature_api.h"
#include "openfeature/value.h"

namespace {

template <typename T>
openfeature::Flag<T> MakeFlag(T value) {
  return openfeature::Flag<T>(
      std::unordered_map<std::string, T>{{"default", std::move(value)}},
      std::optional<std::string>{"default"}, nullptr,
      openfeature::FlagMetadata{}, false);
}

bool Check(bool condition, std::string_view message) {
  if (!condition) {
    std::cerr << "FAIL: " << message << "\n";
    return false;
  }
  return true;
}

}  // namespace

int main() {
  bool ok = true;

  // Verify generated flag key constants.
  ok &= Check(testnamespace::FlagKeys::ENABLE_FEATURE_A == "enableFeatureA",
              "ENABLE_FEATURE_A key mismatch");
  ok &= Check(testnamespace::FlagKeys::USERNAME_MAX_LENGTH ==
                  "usernameMaxLength",
              "USERNAME_MAX_LENGTH key mismatch");
  ok &= Check(testnamespace::FlagKeys::GREETING_MESSAGE == "greetingMessage",
              "GREETING_MESSAGE key mismatch");
  ok &= Check(testnamespace::FlagKeys::DISCOUNT_PERCENTAGE ==
                  "discountPercentage",
              "DISCOUNT_PERCENTAGE key mismatch");
  ok &= Check(testnamespace::FlagKeys::THEME_CUSTOMIZATION ==
                  "themeCustomization",
              "THEME_CUSTOMIZATION key mismatch");

  openfeature::Value theme_value(std::map<std::string, openfeature::Value>{
      {"primaryColor", openfeature::Value(std::string("#112233"))},
      {"secondaryColor", openfeature::Value(std::string("#445566"))},
  });

  std::unordered_map<std::string, std::any> flags = {
      {"enableFeatureA", MakeFlag<bool>(true)},
      {"usernameMaxLength", MakeFlag<int64_t>(100)},
      {"greetingMessage", MakeFlag<std::string>("Welcome!")},
      {"discountPercentage", MakeFlag<double>(0.25)},
      {"themeCustomization", MakeFlag<openfeature::Value>(theme_value)},
  };

  auto& api = openfeature::OpenFeatureAPI::GetInstance();
  api.SetProviderAndWait(
      std::make_shared<openfeature::InMemoryProvider>(std::move(flags)));

  std::unordered_map<std::string, std::any> domain_flags = {
      {"enableFeatureA", MakeFlag<bool>(false)},
  };
  api.SetProviderAndWait(
      "custom-domain",
      std::make_shared<openfeature::InMemoryProvider>(std::move(domain_flags)));

  openfeature::EvaluationContext eval_ctx =
      openfeature::EvaluationContext::Builder()
          .WithTargetingKey("user-123")
          .Build();

  testnamespace::GeneratedClient client = testnamespace::GetClient();

  // Boolean flag
  ok &= Check(client.EnableFeatureA() == true, "EnableFeatureA() != true");
  ok &= Check(client.EnableFeatureA(eval_ctx) == true,
              "EnableFeatureA(ctx) != true");
  auto bool_details = client.EnableFeatureADetails();
  ok &= Check(bool_details.GetValue() == true &&
                  !bool_details.GetErrorCode().has_value(),
              "EnableFeatureADetails() failed");
  auto bool_ctx_details = client.EnableFeatureADetails(eval_ctx);
  ok &= Check(bool_ctx_details.GetValue() == true &&
                  !bool_ctx_details.GetErrorCode().has_value(),
              "EnableFeatureADetails(ctx) failed");

  // Integer flag
  ok &= Check(client.UsernameMaxLength() == 100,
              "UsernameMaxLength() != 100");
  ok &= Check(client.UsernameMaxLength(eval_ctx) == 100,
              "UsernameMaxLength(ctx) != 100");
  auto int_details = client.UsernameMaxLengthDetails();
  ok &= Check(int_details.GetValue() == 100 &&
                  !int_details.GetErrorCode().has_value(),
              "UsernameMaxLengthDetails() failed");
  auto int_ctx_details = client.UsernameMaxLengthDetails(eval_ctx);
  ok &= Check(int_ctx_details.GetValue() == 100 &&
                  !int_ctx_details.GetErrorCode().has_value(),
              "UsernameMaxLengthDetails(ctx) failed");

  // String flag
  ok &= Check(client.GreetingMessage() == "Welcome!",
              "GreetingMessage() != Welcome!");
  ok &= Check(client.GreetingMessage(eval_ctx) == "Welcome!",
              "GreetingMessage(ctx) != Welcome!");
  auto str_details = client.GreetingMessageDetails();
  ok &= Check(str_details.GetValue() == "Welcome!" &&
                  !str_details.GetErrorCode().has_value(),
              "GreetingMessageDetails() failed");
  auto str_ctx_details = client.GreetingMessageDetails(eval_ctx);
  ok &= Check(str_ctx_details.GetValue() == "Welcome!" &&
                  !str_ctx_details.GetErrorCode().has_value(),
              "GreetingMessageDetails(ctx) failed");

  // Float flag
  ok &= Check(std::abs(client.DiscountPercentage() - 0.25) < 1e-9,
              "DiscountPercentage() != 0.25");
  ok &= Check(std::abs(client.DiscountPercentage(eval_ctx) - 0.25) < 1e-9,
              "DiscountPercentage(ctx) != 0.25");
  auto dbl_details = client.DiscountPercentageDetails();
  ok &= Check(std::abs(dbl_details.GetValue() - 0.25) < 1e-9 &&
                  !dbl_details.GetErrorCode().has_value(),
              "DiscountPercentageDetails() failed");
  auto dbl_ctx_details = client.DiscountPercentageDetails(eval_ctx);
  ok &= Check(std::abs(dbl_ctx_details.GetValue() - 0.25) < 1e-9 &&
                  !dbl_ctx_details.GetErrorCode().has_value(),
              "DiscountPercentageDetails(ctx) failed");

  // Object flag
  ok &= Check(client.ThemeCustomization() == theme_value,
              "ThemeCustomization() mismatch");
  ok &= Check(client.ThemeCustomization(eval_ctx) == theme_value,
              "ThemeCustomization(ctx) mismatch");
  auto obj_details = client.ThemeCustomizationDetails();
  ok &= Check(obj_details.GetValue() == theme_value &&
                  !obj_details.GetErrorCode().has_value(),
              "ThemeCustomizationDetails() failed");
  auto obj_ctx_details = client.ThemeCustomizationDetails(eval_ctx);
  ok &= Check(obj_ctx_details.GetValue() == theme_value &&
                  !obj_ctx_details.GetErrorCode().has_value(),
              "ThemeCustomizationDetails(ctx) failed");

  // Domain-scoped client
  testnamespace::GeneratedClient domain_client =
      testnamespace::GetClient("custom-domain");
  ok &= Check(domain_client.EnableFeatureA() == false,
              "domain_client.EnableFeatureA() != false");

  api.Shutdown();

  if (!ok) {
    return 1;
  }

  std::cout << "Generated C++ code compiles and passes all evaluations!\n";
  return 0;
}
