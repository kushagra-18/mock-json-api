# Faker DSL Processor Documentation

A powerful Node.js service that processes Faker.js Domain Specific Language (DSL) strings to generate realistic mock data. This service provides a simple HTTP API endpoint that can evaluate Faker.js directives embedded in templates.

## Table of Contents

- [Overview](#overview)
- [Installation](#installation)
- [Quick Start](#quick-start)
- [API Reference](#api-reference)
- [DSL Syntax](#dsl-syntax)
- [Available Faker Modules](#available-faker-modules)
- [Examples](#examples)
- [Error Handling](#error-handling)
- [Best Practices](#best-practices)
- [Troubleshooting](#troubleshooting)

## Overview

The Faker DSL Processor allows you to:
- Generate realistic mock data using Faker.js
- Process template strings with embedded Faker directives
- Create complex JSON objects with nested fake data
- Handle arrays and multipliers for bulk data generation
- Convert data to various formats (JSON strings, objects, arrays)

## Installation

### Prerequisites
- Node.js (v14 or higher)
- npm or yarn

### Setup

1. **Install dependencies:**
   ```bash
   cd src/main/node
   npm install
   ```

2. **Start the server:**
   ```bash
   npm start
   # or
   node src/server.js
   ```

3. **Server will run on:**
   ```
   http://localhost:3001
   ```

## Quick Start

### Basic Example

Send a POST request to `/process-dsl`:

```bash
curl -X POST http://localhost:3001/process-dsl \
  -H "Content-Type: application/json" \
  -d '{
    "dsl": "Hello {{person.fullName()}}, your email is {{internet.email()}}"
  }'
```

**Response:**
```json
"Hello John Doe, your email is john.doe@example.com"
```

## API Reference

### POST `/process-dsl`

Processes a DSL string and returns the generated mock data.

**Request Body:**
```json
{
  "dsl": "string" // Required: The DSL template to process
}
```

**Response:**
- **200 OK**: Returns the processed result
- **400 Bad Request**: Invalid request format
- **500 Internal Server Error**: Processing error

**Response Types:**
- `string`: For template strings with embedded directives
- `object`: For single directive results that are objects
- `array`: For multiplied directives or array results

## DSL Syntax

### Basic Directive Format

```
{{module.method(arguments)}}
```

### Examples:

1. **Simple method call:**
   ```
   {{person.firstName()}}
   ```

2. **Method with arguments:**
   ```
   {{number.int({"min": 1, "max": 100})}}
   ```

3. **Property access:**
   ```
   {{date.now()}}
   ```

4. **Array multiplier:**
   ```
   {{person.firstName()}} *5
   ```

5. **JSON.stringify wrapper:**
   ```
   {{JSON.stringify({"name": "{{person.fullName()}}", "age": {{number.int({"min": 18, "max": 80})}}})}}
   ```

### Special Features

#### JSON.stringify Wrapper
Use `JSON.stringify()` to convert objects to JSON strings:

```json
{
  "dsl": "{{JSON.stringify({\"user\": \"{{person.fullName()}}\", \"id\": {{number.int({\"min\": 1, \"max\": 1000})}}}})"
}
```

#### Array Multipliers
Generate multiple items using the `*N` syntax:

```json
{
  "dsl": "{{person.firstName()}} *3"
}
```

Returns: `["John", "Jane", "Bob"]`

## Available Faker Modules

### Core Modules

| Module | Description | Example Methods |
|--------|-------------|-----------------|
| `person` | Person-related data | `fullName()`, `firstName()`, `lastName()` |
| `internet` | Internet/web data | `email()`, `userName()`, `url()`, `domainName()` |
| `location` | Geographic data | `city()`, `country()`, `latitude()`, `longitude()` |
| `date` | Date/time data | `recent()`, `past()`, `future()`, `now()` |
| `number` | Numeric data | `int(options)`, `float(options)` |
| `string` | String data | `uuid()`, `alphanumeric(length)` |
| `lorem` | Lorem ipsum text | `paragraph()`, `sentence()`, `word()` |
| `datatype` | Data types | `boolean()`, `json()`, `array()` |
| `image` | Image URLs | `url()`, `avatar()`, `placeholder()` |
| `helpers` | Utility functions | `arrayElement(array)`, `randomize()` |

### Method Examples

#### Person Module
```javascript
// Basic person data
{{person.fullName()}}          // "John Doe"
{{person.firstName()}}         // "John"
{{person.lastName()}}          // "Doe"
{{person.prefix()}}            // "Mr."
{{person.suffix()}}            // "Jr."
```

#### Internet Module
```javascript
// Internet-related data
{{internet.email()}}           // "john.doe@email.com"
{{internet.userName()}}        // "john_doe123"
{{internet.url()}}             // "https://example.com"
{{internet.domainName()}}      // "example.com"
{{internet.password()}}        // "xK9#mN2$"
```

#### Location Module
```javascript
// Geographic data
{{location.city()}}            // "New York"
{{location.country()}}         // "United States"
{{location.latitude()}}        // "40.7128"
{{location.longitude()}}       // "-74.0060"
{{location.streetAddress()}}   // "123 Main St"
{{location.zipCode()}}         // "10001"
```

#### Number Module
```javascript
// Numeric data with options
{{number.int({"min": 1, "max": 100})}}        // 42
{{number.float({"min": 0, "max": 1, "precision": 2})}}  // 0.73
```

#### Date Module
```javascript
// Date/time data
{{date.recent()}}              // "2025-05-31T10:30:00.000Z"
{{date.past()}}                // "2024-03-15T14:20:00.000Z"
{{date.future()}}              // "2026-08-22T09:45:00.000Z"
{{date.now()}}                 // Current timestamp
```

#### Helpers Module
```javascript
// Array selection
{{helpers.arrayElement(["red", "green", "blue"])}}     // "red"

// Multiple array elements
{{helpers.arrayElements(["a", "b", "c", "d"], {"count": 2})}}  // ["a", "c"]
```

## Examples

### 1. Simple User Profile

**Request:**
```json
{
  "dsl": "{{JSON.stringify({\"id\": {{number.int({\"min\": 1, \"max\": 1000})}}, \"name\": \"{{person.fullName()}}\", \"email\": \"{{internet.email()}}\", \"age\": {{number.int({\"min\": 18, \"max\": 80})}}}})"
}
```

**Response:**
```json
"{\"id\": 42, \"name\": \"John Doe\", \"email\": \"john.doe@email.com\", \"age\": 28}"
```

### 2. Complex Social Media Post

**Request:**
```json
{
  "dsl": "{{JSON.stringify({\"post_id\": \"{{string.uuid()}}\", \"author\": {\"username\": \"{{internet.userName()}}\", \"display_name\": \"{{person.fullName()}}\", \"avatar\": \"{{image.avatar()}}\"}, \"content\": \"{{lorem.paragraph()}} #{{word.noun()}} #{{word.adjective()}}\", \"location\": {\"city\": \"{{location.city()}}\", \"country\": \"{{location.country()}}\", \"coordinates\": {\"lat\": \"{{location.latitude()}}\", \"lng\": \"{{location.longitude()}}\"}}, \"created_at\": \"{{date.recent()}}\", \"edited\": {{datatype.boolean()}}})"
}
```

### 3. Array of Users

**Request:**
```json
{
  "dsl": "{{JSON.stringify([{\"name\": \"{{person.fullName()}}\", \"email\": \"{{internet.email()}}\"}])}} *5"
}
```

### 4. E-commerce Product

**Request:**
```json
{
  "dsl": "{{JSON.stringify({\"product_id\": \"{{string.uuid()}}\", \"name\": \"{{commerce.productName()}}\", \"price\": {{number.float({\"min\": 10, \"max\": 500, \"precision\": 2})}}, \"description\": \"{{lorem.sentence()}}\", \"category\": \"{{helpers.arrayElement([\\\"electronics\\\", \\\"clothing\\\", \\\"books\\\", \\\"home\\\"])}\", \"in_stock\": {{datatype.boolean()}}, \"image_url\": \"{{image.url()}}\"})"
}
```

## Error Handling

### Common Error Types

1. **Invalid Module/Method:**
   ```
   [ERROR: Invalid module or path: "invalid.method" (failed at "invalid")]
   ```

2. **Invalid Arguments:**
   ```
   [ERROR: Failed to parse JSON arguments for "number.int": Unexpected token]
   ```

3. **Function Execution Error:**
   ```
   [ERROR: Error executing faker function "helpers.arrayElement": Expected array argument]
   ```

### Error Response Format

When processing fails, errors are embedded in the response as `[ERROR: message]` strings:

```json
{
  "name": "John Doe",
  "age": "[ERROR: Invalid module or path: \"invalid.method\"]"
}
```

## Best Practices

### 1. Proper JSON Escaping

When using arrays or objects in arguments, ensure proper escaping:

**✅ Correct:**
```json
{
  "dsl": "{{helpers.arrayElement([\"red\", \"green\", \"blue\"])}}"
}
```

**❌ Incorrect:**
```json
{
  "dsl": "{{helpers.arrayElement([red, green, blue])}}"
}
```

### 2. Use Appropriate Modules

Use the correct module for your data type:

**✅ Correct:**
```javascript
{{datatype.boolean()}}     // For boolean values
{{location.city()}}        // For city names
{{image.url()}}           // For image URLs
```

**❌ Incorrect:**
```javascript
{{helpers.boolean()}}      // doesn't exist
{{address.city()}}         // old API
{{image.imageUrl()}}       // old API
```

### 3. Function vs Property

Always use parentheses for function calls:

**✅ Correct:**
```javascript
{{person.fullName()}}
{{date.recent()}}
```

**❌ Incorrect:**
```javascript
{{person.fullName}}
{{date.recent}}
```

### 4. Complex Objects

For complex nested objects, use `JSON.stringify()`:

```json
{
  "dsl": "{{JSON.stringify({\"user\": {\"profile\": {\"name\": \"{{person.fullName()}}\", \"settings\": {\"theme\": \"{{helpers.arrayElement([\\\"dark\\\", \\\"light\\\"])}}\"}}}}})"
}
```

## Troubleshooting

### Common Issues and Solutions

1. **Module not found errors:**
   - Check the [Available Faker Modules](#available-faker-modules) section
   - Ensure you're using the correct module name (e.g., `location` not `address`)

2. **JSON parsing errors:**
   - Verify proper escaping of quotes in arguments
   - Use online JSON validators to check syntax

3. **Unexpected string responses:**
   - Check if you need to use `JSON.stringify()` for object conversion
   - Verify the template structure

4. **Array element selection not working:**
   - Ensure array arguments are properly formatted: `["item1", "item2"]`
   - Use double escaping for quotes: `[\"item1\", \"item2\"]`

### Testing Your DSL

Use simple examples first:

```bash
# Test basic functionality
curl -X POST http://localhost:3001/process-dsl \
  -H "Content-Type: application/json" \
  -d '{"dsl": "{{person.fullName()}}"}'

# Test with arguments
curl -X POST http://localhost:3001/process-dsl \
  -H "Content-Type: application/json" \
  -d '{"dsl": "{{number.int({\"min\": 1, \"max\": 10})}}"}'
```

### Debugging Tips

1. **Start Simple:** Begin with basic directives and gradually add complexity
2. **Check Logs:** Monitor server console for detailed error messages
3. **Validate JSON:** Use tools like JSONLint to verify your JSON structure
4. **Test Incrementally:** Build complex templates step by step

---

## Version Information

- **Faker.js Version:** 8.4.1
- **Node.js:** Compatible with v14+
- **Express:** 4.21.1

For more information about Faker.js methods and options, visit the [official Faker.js documentation](https://fakerjs.dev/).
