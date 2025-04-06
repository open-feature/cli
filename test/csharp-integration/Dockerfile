FROM mcr.microsoft.com/dotnet/sdk:8.0

WORKDIR /app

# Copy necessary files
COPY expected/OpenFeature.cs /app/
COPY CompileTest.csproj /app/
COPY Program.cs /app/

# Restore dependencies
RUN dotnet restore

# Build the project
RUN dotnet build

# The image will be used to validate C# compilation only
ENTRYPOINT ["dotnet", "run"]