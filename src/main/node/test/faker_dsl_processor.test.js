// src/main/node/test/faker_dsl_processor.test.js
const assert = require('assert');
const { processDslString } = require('../src/faker_dsl_processor'); // Adjust path if needed

// Helper to run tests
function test(description, fn) {
    try {
        fn();
        console.log(`\u001b[32m✓\u001b[0m ${description}`);
    } catch (error) {
        console.error(`\u001b[31m✗\u001b[0m ${description}`);
        console.error(error);
    }
}

// --- Test Cases ---

test('should process a simple Faker.js directive: {{name.firstName}}', () => {
    const result = processDslString('{{name.firstName}}');
    assert.strictEqual(typeof result, 'string');
    assert.ok(result.length > 0, 'First name should not be empty');
});

test('should process a simple Faker.js directive with arguments: {{lorem.words(3)}}', () => {
    const result = processDslString('{{lorem.words(3)}}');
    assert.strictEqual(typeof result, 'string');
    assert.strictEqual(result.split(' ').length, 3, 'Should generate 3 words');
});

test('should handle repetition: {{internet.email}}*3', () => {
    const result = processDslString('{{internet.email}}*3');
    assert.ok(Array.isArray(result), 'Result should be an array');
    assert.strictEqual(result.length, 3, 'Array should have 3 emails');
    result.forEach(email => {
        assert.strictEqual(typeof email, 'string');
        assert.ok(email.includes('@'), 'Each item should be an email');
    });
});

test('should handle repetition with count 1: {{number.int}}*1', () => {
    // As per current implementation, *1 means an array of 1.
    // If {{number.int}}*1 was intended to be same as {{number.int}}, behavior of processDslString needs change.
    // Current code implies *N always yields array.
    const result = processDslString('{{number.int}}*1');
    assert.ok(Array.isArray(result), 'Result should be an array even for *1');
    assert.strictEqual(result.length, 1, 'Array should have 1 number');
    assert.strictEqual(typeof result[0], 'number');
});

test('should process a directive with JSON object arguments: {{number.int({"min": 100, "max": 101})}}', () => {
    const result = processDslString('{{number.int({"min": 100, "max": 101})}}');
    assert.strictEqual(typeof result, 'number');
    assert.ok(result === 100 || result === 101, 'Number should be between 100 and 101');
});

test('should process a directive with JSON array arguments: {{helpers.arrayElement(["a", "b", "c"])}}', () => {
    const result = processDslString('{{helpers.arrayElement(["a", "b", "c"])}}');
    assert.strictEqual(typeof result, 'string');
    assert.ok(['a', 'b', 'c'].includes(result), 'Result should be one of the array elements');
});

test('should return raw number for {{number.int}}', () => {
    const result = processDslString('{{number.int}}');
    assert.strictEqual(typeof result, 'number');
});

test('should return raw boolean for {{datatype.boolean}}', () => {
    const result = processDslString('{{datatype.boolean}}');
    assert.strictEqual(typeof result, 'boolean');
});

test('should embed processed value in a string template: Hello {{name.firstName}}!', () => {
    const result = processDslString('Hello {{name.firstName}}!');
    assert.strictEqual(typeof result, 'string');
    assert.ok(result.startsWith('Hello '));
    assert.ok(result.endsWith('!'));
    assert.ok(result.length > 'Hello !'.length);
});

test('should embed multiple processed values: {{name.firstName}} {{name.lastName}}', () => {
    const result = processDslString('{{name.firstName}} {{name.lastName}}');
    assert.strictEqual(typeof result, 'string');
    const parts = result.split(' ');
    assert.strictEqual(parts.length, 2);
    assert.ok(parts[0].length > 0);
    assert.ok(parts[1].length > 0);
});

test('should embed repeated values in a string: Emails: {{internet.email}}*2', () => {
    const result = processDslString('Emails: {{internet.email}}*2');
    assert.strictEqual(typeof result, 'string');
    assert.ok(result.startsWith('Emails: [')); // Array becomes JSON string
    assert.ok(result.endsWith(']'));
    const jsonPart = result.substring('Emails: '.length);
    const arr = JSON.parse(jsonPart);
    assert.ok(Array.isArray(arr));
    assert.strictEqual(arr.length, 2);
    assert.ok(arr[0].includes('@'));
});

test('should handle invalid Faker module/path: {{nonExistent.module}}', () => {
    const result = processDslString('{{nonExistent.module}}');
    assert.strictEqual(typeof result, 'string');
    // Error message comes from evaluateFakerDirective
    assert.ok(result.includes('[ERROR: Invalid module or path: "nonExistent.module" (failed at "nonExistent")]'));
});

test('should handle invalid Faker method in existing module: {{name.nonExistentMethod}}', () => {
    const result = processDslString('{{name.nonExistentMethod}}');
    assert.strictEqual(typeof result, 'string');
    assert.ok(result.includes('[ERROR: Invalid module or path: "name.nonExistentMethod" (failed at "nonExistentMethod")]'));
});

test('should handle invalid arguments (non-JSON): {{lorem.words(badArg)}}', () => {
    const result = processDslString('{{lorem.words(badArg)}}'); // badArg is not valid JSON
    assert.strictEqual(typeof result, 'string');
    assert.ok(result.includes('[ERROR: Failed to parse JSON arguments for "lorem.words"'));
});

test('should handle arguments for property: {{name.firstName("arg")}}', () => {
    const result = processDslString('{{name.firstName("arg")}}');
    assert.strictEqual(typeof result, 'string');
    assert.ok(result.includes('[ERROR: Property "name.firstName" does not accept arguments. Found: ("arg")]'));
});

test('should process top-level JSON.stringify: {{JSON.stringify({"name": "{{name.firstName}}"})}}', () => {
    const result = processDslString('{{JSON.stringify({"name": "{{name.firstName}}"})}}');
    assert.strictEqual(typeof result, 'string');
    const parsed = JSON.parse(result);
    assert.strictEqual(typeof parsed, 'object');
    assert.ok(parsed.hasOwnProperty('name'));
    assert.strictEqual(typeof parsed.name, 'string');
    assert.ok(parsed.name.length > 0);
});

test('should process JSON.stringify with repetition: {{JSON.stringify({"emails": "{{internet.email}}*2"})}}', () => {
    const result = processDslString('{{JSON.stringify({"emails": "{{internet.email}}*2"})}}');
    assert.strictEqual(typeof result, 'string');
    const parsed = JSON.parse(result);
    assert.strictEqual(typeof parsed, 'object');
    assert.ok(parsed.hasOwnProperty('emails'));
    assert.ok(Array.isArray(parsed.emails));
    assert.strictEqual(parsed.emails.length, 2);
    assert.ok(parsed.emails[0].includes('@'));
});

test('should process JSON.stringify with raw types: {{JSON.stringify({"age": "{{number.int}}", "active": "{{datatype.boolean}}"})}}', () => {
    const result = processDslString('{{JSON.stringify({"age": "{{number.int}}", "active": "{{datatype.boolean}}"})}}');
    assert.strictEqual(typeof result, 'string');
    const parsed = JSON.parse(result);
    assert.strictEqual(typeof parsed, 'object');
    assert.ok(parsed.hasOwnProperty('age'));
    assert.ok(parsed.hasOwnProperty('active'));
    assert.strictEqual(typeof parsed.age, 'number');
    assert.strictEqual(typeof parsed.active, 'boolean');
});

test('should handle nested JSON.stringify (complex): {{JSON.stringify({"user": "{{JSON.stringify({ "name": "{{name.firstName}}" })}}" })}}', () => {
    const result = processDslString('{{JSON.stringify({"user": "{{JSON.stringify({ "name": "{{name.firstName}}" })}}" })}}');
    assert.strictEqual(typeof result, 'string');
    const outerParsed = JSON.parse(result);
    assert.strictEqual(typeof outerParsed, 'object');
    assert.ok(outerParsed.hasOwnProperty('user'));
    assert.strictEqual(typeof outerParsed.user, 'string'); // Inner stringify means user is a string
    const innerParsed = JSON.parse(outerParsed.user);
    assert.strictEqual(typeof innerParsed, 'object');
    assert.ok(innerParsed.hasOwnProperty('name'));
    assert.strictEqual(typeof innerParsed.name, 'string');
});

test('should handle plain string without directives', () => {
    const result = processDslString('Hello world');
    assert.strictEqual(result, 'Hello world');
});

test('should handle empty string', () => {
    const result = processDslString('');
    assert.strictEqual(result, '');
});

test('should handle directive with sub-module path: {{animal.dog}}', () => {
    // This test assumes 'animal.dog' resolves to an object, which then gets stringified.
    // Faker.js typically has functions at the end of paths, or simple properties.
    // If animal.dog is an object of functions, this tests if it's stringified.
    const result = processDslString('{{animal.dog}}');
    assert.strictEqual(typeof result, 'object'); // animal.dog is an object of functions
});

test('should handle functions that return objects: {{helpers.contextualCard}}', () => {
    const result = processDslString('{{helpers.contextualCard}}');
    assert.strictEqual(typeof result, 'object');
    assert.ok(result.hasOwnProperty('name') && result.hasOwnProperty('username') && result.hasOwnProperty('avatar'));
});

test('should handle functions that return objects within a template string: User: {{helpers.contextualCard}}', () => {
    const result = processDslString('User: {{helpers.contextualCard}}');
    assert.strictEqual(typeof result, 'string');
    assert.ok(result.startsWith('User: {'));
    assert.ok(result.endsWith('}'));
    const jsonPart = result.substring('User: '.length);
    const parsed = JSON.parse(jsonPart);
    assert.strictEqual(typeof parsed, 'object');
    assert.ok(parsed.hasOwnProperty('name') && parsed.hasOwnProperty('username') && parsed.hasOwnProperty('avatar'));
});


// To run these tests, you would typically save this file and run `node src/main/node/test/faker_dsl_processor.test.js`
// from the root of the repository or `src/main/node` directory.
// The output will show ✓ for pass and ✗ for fail.
// Add more tests as needed for edge cases or more complex structures.
console.log("\n--- All tests run ---");
console.log("Please review the output above. Failures (✗) will include error details.");
console.log("If all tests show ✓, then the basic functionality is working as expected.");

// Note: Some Faker functions might be locale-dependent or highly random.
// Tests for exact string matches might be flaky. Focus on type and structure.
// For example, instead of assert.strictEqual(result, "Specific Name"),
// use assert.strictEqual(typeof result, "string") and assert.ok(result.length > 0).
// The `faker.seed()` method could be used if deterministic output is absolutely necessary for some tests,
// but that would require modifying processDslString or having a global setup for tests.
// For now, we accept the inherent randomness of Faker.js for most string outputs.

// Example of a test that might be flaky if not handled carefully:
// test('specific date (potentially flaky without seeding)', () => {
//    const result = processDslString('{{date.month}}');
//    // This is okay as it checks type and that it's one of the valid months.
//    assert.strictEqual(typeof result, 'string');
//    assert.ok(faker.definitions.date.month.includes(result));
// });

// The provided tests are designed to be robust against typical Faker randomness.
// Errors like '[ERROR: Faker module 'nonExistent' not found]'
// were updated to match the actual error messages from the implementation.
// '[ERROR: Faker path 'nonExistentMethod' not found in module 'name']'
// '[ERROR: Property 'name.firstName' does not accept arguments. Found: ("arg")]'
// '[ERROR: Failed to parse JSON arguments for "lorem.words"]'
// These are based on the current error messages in `evaluateFakerDirective`.
// If those messages change, these tests would need updating.

// The test for '{{animal.dog}}' was updated. 'animal.dog' is an object containing functions,
// so the raw result is an object. If it were embedded in a string, it would be JSON.stringified.

// Added tests for functions returning objects (e.g. helpers.contextualCard)
// both as raw output and within a template string.
// End of tests
