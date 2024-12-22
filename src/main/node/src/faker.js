const { NodeVM } = require('vm2');
const {faker, ru} = require('@faker-js/faker');
const path = require('path');

/**
 * Executes user-provided code securely with access only to the `faker` library.
 * @param {string} code - The user-provided JavaScript code.
 * @returns {*} - The result of the executed code.
 * @throws {Error} - If the code is invalid or tries to access restricted resources.
 */

async function runFakerCode(code) {
    const vm = new NodeVM({
        console: 'inherit',
        sandbox: { faker },
        require: {
            external: true,
            root: path.resolve(__dirname),
        }
    });

    try {
        const wrappedCode = `
            (() => {
                ${code}
            })();
        `;
        
        const result = await vm.run(wrappedCode, 'vm.js');
        return result;
    } catch (error) {
        console.error('VM Execution Error:', error);
        throw new Error(`Error executing code: ${error.message}`);
    }
}

module.exports = runFakerCode;