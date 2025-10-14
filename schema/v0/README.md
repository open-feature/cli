# OpenFeature CLI Schema Definitions

This directory contains schema definitions for the OpenFeature CLI.

## Files

### flag-manifest.json
JSON Schema for the flag manifest file format. This schema defines the structure for feature flag configurations stored locally.

## Sync API Implementation

Services that want to sync flag configurations with the OpenFeature CLI should implement the API defined in `api/v1/sync.yaml`.

Key requirements:
- Accept POST or PUT requests with JSON payload
- Support Bearer token authentication (optional)
- Validate flag data according to the manifest schema
- Return appropriate HTTP status codes

### Example Request
```json
{
  "key": "feature-x-enabled",
  "type": "boolean",
  "description": "Enable feature X",
  "defaultValue": true
}
```

### Example Success Response
```json
{
  "flag": {
    "key": "feature-x-enabled",
    "name": "feature-x-enabled",
    "type": "boolean",
    "description": "Enable feature X",
    "defaultValue": true
  },
  "updatedAt": "2025-11-03T20:41:52.000Z"
}
```

## Usage

To view the OpenAPI specification in a user-friendly format, you can use tools like:
- [Swagger Editor](https://editor.swagger.io/)
- [Redoc](https://github.com/Redocly/redoc)
- [OpenAPI Generator](https://openapi-generator.tech/) to generate server stubs or client SDKs