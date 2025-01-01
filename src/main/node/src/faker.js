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
            Handlebars: Handlebars
        });

          const template11 = Handlebars.compile(templateStr);

          console.log('test',template11);

        try {
            const code = `
                const { faker } = require('@faker-js/faker');
                
                // Register custom Handlebars helpers
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

                // Register range helper for loops
                Handlebars.registerHelper('range', function(n, options) {
                    const results = [];
                    for (let i = 0; i < n; i++) {
                        results.push(options.fn({ index: i }));
                    }
                    return results.join('');
                });

                // Register random range helper
                Handlebars.registerHelper('randomRange', function(min, max) {
                    return faker.number.int({ min, max });
                });

                // Register array helper
                Handlebars.registerHelper('array', function(n, options) {
                    return Array.from({ length: n }, (_, index) => 
                        options.fn({ index })
                    );
                });

                // Compile template
                const template = Handlebars.compile(\`${templateStr}\`);
                
            `;

            const script = new vm.Script(code);
            
            const result = await Promise.race([
                new Promise((resolve) => {
                    resolve(script.runInContext(context, {
                        timeout: this.timeout,
                        displayErrors: true,
                        breakOnSigint: true
                    }));
                }),
                new Promise((_, reject) => 
                    setTimeout(() => reject(new Error('Execution timed out')), this.timeout)
                )
            ]);

            return {
                result,
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
        console.log(`Executing template with count: ${count}`);
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