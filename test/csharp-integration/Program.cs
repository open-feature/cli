using System;
using System.Threading.Tasks;
using Microsoft.Extensions.DependencyInjection;
using OpenFeature;
using OpenFeature.Model;
using TestNamespace;

// This program just validates that the generated OpenFeature C# client code compiles
// We don't need to run the code since the goal is to test compilation only
namespace CompileTest
{
    class Program
    {
        static void Main(string[] args)
        {
            Console.WriteLine("Testing compilation of generated OpenFeature client...");
            
            // Test DI initialization
            var services = new ServiceCollection();
            // Register OpenFeature services manually for the test
            services.AddSingleton(_ => Api.Instance);
            services.AddSingleton<IFeatureClient>(_ => Api.Instance.GetClient());
            services.AddSingleton<GeneratedClient>();
            var serviceProvider = services.BuildServiceProvider();
            
            // Test client retrieval from DI
            var client = serviceProvider.GetRequiredService<GeneratedClient>();
            
            // Also test the traditional factory method
            var clientFromFactory = GeneratedClient.CreateClient();
            
            // Success!
            Console.WriteLine("Generated C# code compiles successfully!");
        }
    }
}