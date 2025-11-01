# OpenFeature CLI Schema Definitions

This directory contains schema definitions for the OpenFeature CLI.

## Files

### flag-manifest.json
JSON Schema for the flag manifest file format. This schema defines the structure for feature flag configurations stored locally.

### ~~push-api.yaml~~ (Moved)
The push API specification has been moved to `api/v1/push.yaml` for better organization and to support code generation.

## Push API Implementation

Services that want to receive flag configurations from the OpenFeature CLI push command should implement the API defined in `api/v1/push.yaml`.

Key requirements:
- Accept POST or PUT requests with JSON payload
- Support Bearer token authentication (optional)
- Validate flag data according to the manifest schema
- Return appropriate HTTP status codes

### Example Request
```json
{
  "$schema": "https://raw.githubusercontent.com/open-feature/cli/refs/heads/main/schema/v0/flag-manifest.json",
  "flags": [
    {
      "key": "feature-x-enabled",
      "type": "boolean",
      "description": "Enable feature X",
      "defaultValue": true
    }
  ]
}
```

### Example Success Response
```json
{
  "message": "Successfully pushed 1 flag",
  "count": 1,
  "flags": ["feature-x-enabled"]
}
```

## Usage

To view the OpenAPI specification in a user-friendly format, you can use tools like:
- [Swagger Editor](https://editor.swagger.io/)
- [Redoc](https://github.com/Redocly/redoc)
- [OpenAPI Generator](https://openapi-generator.tech/) to generate server stubs or client SDKs