using System;
using System.Threading.Tasks;
using OpenFeature;
using OpenFeature.Model;

// This program just validates that the generated OpenFeature C# client code compiles
// We don't need to run the code since the goal is to test compilation only
namespace CompileTest
{
    class Program
    {
        static void Main(string[] args)
        {
            Console.WriteLine("Testing compilation of generated OpenFeature client...");
            
            // Success!
            Console.WriteLine("Generated C# code compiles successfully!");
        }
    }
}