// Main module for processing Faker DSL strings.
const { faker } = require('@faker-js/faker');

/**
 * Evaluates a single Faker.js directive.
 * Examples:
 *  - "name.firstName"
 *  - "number.int({\"min\": 1, \"max\": 99})"
 *  - "locales.en_US.address.zipCode"
 *
 * @param {string} directiveContent - The content of the directive (inside {{...}}).
 * @returns {*} The result of the Faker.js call.
 * @throws {Error} If the directive is invalid, arguments are malformed, or Faker.js call fails.
 */
function evaluateFakerDirective(directiveContent) {
  // Regex to parse directives like "module.method(args)" or "module.property"
  // It captures:
  // 1. The full path to the method/property (e.g., "name.firstName", "number.int")
  // 2. Optionally, the arguments string if it's a method call (e.g., "{\"min\": 10}")
  const directiveRegex = /^([a-zA-Z0-9_.]+)(?:\((.*)\))?$/;
  const match = directiveContent.match(directiveRegex);

  if (!match) {
    throw new Error(`Invalid directive format: ${directiveContent}`);
  }

  const fullPath = match[1];
  const argsString = match[2]; // undefined if no parentheses, or empty string if "()"

  const pathParts = fullPath.split('.');
  let current = faker;
  let targetObject = faker; // To keep track of the object before the final property/method

  for (let i = 0; i < pathParts.length; i++) {
    const part = pathParts[i];
    if (current && typeof current === 'object' && part in current) {
      if (i < pathParts.length - 1) {
        targetObject = current[part];
      }
      current = current[part];
    } else {
      throw new Error(`Invalid module or path: "${fullPath}" (failed at "${part}")`);
    }
  }

  if (typeof current === 'function') {
    let parsedArgs = [];
    if (argsString !== undefined) { // Arguments were provided in parentheses
      try {
        // Allow empty argsString for functions called with no arguments e.g. {{date.past()}}
        let rawParsedArgs;
        if (argsString.trim() === '') {
          rawParsedArgs = [];
        } else {
          // Try to fix common escaping issues before parsing
          let cleanedArgsString = argsString;
          
          // Fix double-escaped quotes like \\\"string\\\" -> \"string\"
          cleanedArgsString = cleanedArgsString.replace(/\\\\\"/g, '"');
          
          // Fix incorrectly escaped quotes in arrays like [\"item\"] -> ["item"]
          cleanedArgsString = cleanedArgsString.replace(/\[\\"/g, '["').replace(/\\"]/g, '"]').replace(/\\",\s*\\"/g, '", "');
          
          rawParsedArgs = JSON.parse(cleanedArgsString);
        }
        
        // Ensure parsedArgs is an array for .apply()
        if (Array.isArray(rawParsedArgs)) {
          // If the function expects the array as a single argument (like arrayElement),
          // we need to pass the array as one argument, not spread its elements
          if (fullPath === 'helpers.arrayElement' || fullPath === 'helpers.arrayElements' || fullPath === 'helpers.weightedArrayElement') {
            parsedArgs = [rawParsedArgs];
          } else {
            parsedArgs = rawParsedArgs;
          }
        } else {
          parsedArgs = [rawParsedArgs];
        }
      } catch (e) {
        throw new Error(`Failed to parse JSON arguments for "${fullPath}": ${e.message}. Args: ${argsString}`);
      }
    }
    try {
      return current.apply(targetObject, parsedArgs);
    } catch (e) {
      throw new Error(`Error executing faker function "${fullPath}": ${e.message}`);
    }
  } else { // Property
    if (argsString !== undefined) {
      // Arguments provided (e.g., {{name.firstName()}} which is not a function in faker.js)
      throw new Error(`Property "${fullPath}" does not accept arguments. Found: (${argsString})`);
    }
    return current;
  }
}

/**
 * Processes a string containing Faker DSL directives.
 * Directives are in the format {{module.method(JSONArguments)}} or {{module.property}}.
 * Can also handle an array multiplier like {{directive}}*N.
 *
 * @param {string} dslString - The input string with DSL directives.
 * @returns {string|*} - The processed string, or the raw result if the input is a single directive.
 */
function processDslString(dslString) {
  // Regular expression to match {{JSON.stringify(INNER_DSL)}}
  // Captures:
  // 1. The INNER_DSL part
  const stringifyPattern = /^{{JSON\.stringify\((.+)\)}}$/;
  const stringifyMatch = dslString.match(stringifyPattern);

  if (stringifyMatch) {
    const innerDsl = stringifyMatch[1];
    // Recursively process the inner DSL
    const resultObject = processDslString(innerDsl);
    // Then stringify the result of the inner processing
    return JSON.stringify(resultObject);
  }

  // Regex to find directives:
  // - {{directiveContent}}
  // - {{directiveContent}} *N (captures N)
  // It captures:
  // 1. The directive content (e.g., "name.firstName", "number.int({\"min\":1})")
  // 2. Optionally, the multiplier N
  const directiveRegex = /\{\{(.+?)\}\}(?:\s*\*(\d+))?/g;

  let match;
  let lastIndex = 0;
  let resultParts = [];
  let isSingleDirective = true; // Assume it's a single directive until proven otherwise

  // Check if the entire string is a single directive for raw output
  // Need to do this before iterative replacement
  const firstMatch = directiveRegex.exec(dslString);
  if (firstMatch && firstMatch.index === 0 && firstMatch[0].length === dslString.length) {
    const directiveContent = firstMatch[1];
    const multiplier = firstMatch[2] ? parseInt(firstMatch[2], 10) : 1;

    try {
      if (multiplier > 1) {
        const resultsArray = [];
        for (let i = 0; i < multiplier; i++) {
          resultsArray.push(evaluateFakerDirective(directiveContent));
        }
        return resultsArray; // Return raw array
      } else {
        return evaluateFakerDirective(directiveContent); // Return raw result
      }
    } catch (e) {
      // If the single directive fails, return its error stringified if it's for JSON.stringify scenario,
      // or raw error string if it's a direct single directive.
      // This case (single directive) should return raw error.
      return `[ERROR: ${e.message}]`;
    }
  }
  // Reset regex for iterative processing if not a single directive
  directiveRegex.lastIndex = 0;


  while ((match = directiveRegex.exec(dslString)) !== null) {
    isSingleDirective = false; // It's part of a template
    // Append text before the directive
    if (match.index > lastIndex) {
      resultParts.push(dslString.substring(lastIndex, match.index));
    }

    const directiveContent = match[1];
    const multiplier = match[2] ? parseInt(match[2], 10) : 1;
    let evaluatedValue;

    try {
      if (multiplier > 1) {
        const resultsArray = [];
        for (let i = 0; i < multiplier; i++) {
          resultsArray.push(evaluateFakerDirective(directiveContent));
        }
        evaluatedValue = resultsArray; // This will be stringified later
      } else {
        evaluatedValue = evaluateFakerDirective(directiveContent);
      }

      // If part of a template, non-string results should be stringified.
      // Strings are appended directly.
      if (typeof evaluatedValue === 'string') {
        resultParts.push(evaluatedValue);
      } else if (evaluatedValue instanceof Date) {
        // Handle Date objects specially - convert to ISO string without quotes
        resultParts.push(evaluatedValue.toISOString());
      } else {
        resultParts.push(JSON.stringify(evaluatedValue));
      }
    } catch (e) {
      resultParts.push(`[ERROR: ${e.message}]`);
    }
    lastIndex = directiveRegex.lastIndex;
  }

  // If isSingleDirective is still true here, it means no directives were found at all.
  // (e.g. input is "hello world", not "hello {{name.firstName}}")
  // The firstMatch logic above handles cases where one directive spans the whole string.
  if (isSingleDirective && firstMatch === null) {
     return dslString; // No directives found, return original string
  }

  // Append any remaining text after the last directive
  if (lastIndex < dslString.length) {
    resultParts.push(dslString.substring(lastIndex));
  }

  return resultParts.join('');
}

module.exports = { processDslString };
