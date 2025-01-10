const vm = require('vm');
const util = require('util');
const Handlebars = require('handlebars');

class CodeSandbox {
    constructor(options = {}) {
        this.timeout = options.timeout || 1000;
        this.memory = options.memory || 50;
        this.allowedModules = new Set(options.allowedModules || ['@faker-js/faker', 'handlebars']);
        this.contextDefaults = {
            console: {
                log: (...args) => this.logs.push(util.format(...args)),
                error: (...args) => this.errors.push(util.format(...args))
            },
            setTimeout: undefined,
            setInterval: undefined,
            setImmediate: undefined,
            process: undefined,
            Buffer: undefined
        };
    }

    async runCode(templateStr, count = 1) {
        this.logs = [];
        this.errors = [];

        const context = vm.createContext({
            ...this.contextDefaults,
            require: (moduleName) => {
                if (this.allowedModules.has(moduleName)) {
                    return require(moduleName);
                }
                throw new Error(`Module '${moduleName}' is not allowed`);
            },
            Handlebars: Handlebars,
            resultValue: null
        });

        try {
            const code = `
                const { faker } = require('@faker-js/faker');
                
                Handlebars.registerHelper('faker', function(property) {
                    const props = property.split('.');
                    let current = faker;
                    for (const prop of props) {
                        if (typeof current[prop] === 'function') {
                            return current[prop]();
                        }
                        current = current[prop];
                    }
                    return '';
                });

                Handlebars.registerHelper('randomRange', function(min, max) {
                    return faker.number.int({ min, max });
                });

                // Improved array helper with better handling of index and separators
                Handlebars.registerHelper('array', function(n, options) {
                    return Array.from({ length: n }, (_, index) => 
                        options.fn({
                            index,
                            '@index': index,
                            '@first': index === 0,
                            '@last': index === n - 1
                        })
                    ).join('');
                });

                const template = Handlebars.compile(\`${fixedTemplate}\`);
                resultValue = template();
            `;

            const script = new vm.Script(code);

            await Promise.race([
                new Promise((resolve) => {
                    script.runInContext(context, {
                        timeout: this.timeout,
                        displayErrors: true,
                        breakOnSigint: true
                    });
                    resolve();
                }),
                new Promise((_, reject) =>
                    setTimeout(() => reject(new Error('Execution timed out')), this.timeout)
                )
            ]);

            // Parse and reformat the JSON for proper formatting
            const parsed = JSON.parse(context.resultValue);
            const formatted = JSON.stringify(parsed, null, 2);

            return {
                result: formatted,
                logs: this.logs,
                errors: this.errors,
                success: true
            };

        } catch (error) {
            return {
                result: null,
                logs: this.logs,
                errors: [...this.errors, error.message],
                success: false
            };
        }
    }
}

const sandbox = new CodeSandbox({
    timeout: 2000,
    memory: 100
});

async function executeCode(template, count = 1) {
    try {
        const result = await sandbox.runCode(template, count);
        return result;
    } catch (error) {
        return {
            result: null,
            logs: [],
            errors: [error.message],
            success: false
        };
    }
}

module.exports = {
    executeCode
};