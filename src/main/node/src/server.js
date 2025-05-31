const express = require('express');
const axios = require('axios');
const jwt = require('jsonwebtoken');

require('dotenv').config({ path: '../.env' });

const cors = require('cors');

const corsOptions = {
    origin: '*',
};

// Require the new DSL processor
const { processDslString } = require('./faker_dsl_processor');

const app = express();
const port = process.env.PORT || 3001;

app.use(express.json());

app.use(cors(corsOptions));

app.use((req, res, next) => {
    if (req.url === '/favicon.ico') {
        res.status(204).end();
        return;
    }
    next();
});

// New route for DSL processing
app.post('/process-dsl', (req, res) => {
    const { dsl } = req.body;

    if (typeof dsl !== 'string') {
        return res.status(400).json({ message: 'Request body must include a "dsl" string.' });
    }

    try {
        const result = processDslString(dsl);
        
        // If the result is a stringified JSON string (wrapped in quotes),
        // we need to parse it twice to get the actual object
        if (typeof result === 'string') {
            // Check if it's a quoted JSON string like "{"key": "value"}"
            if (result.startsWith('"') && result.endsWith('"')) {
                try {
                    // First parse removes the outer quotes
                    const unquoted = JSON.parse(result);
                    // Second parse gets the actual object
                    const parsedResult = JSON.parse(unquoted);
                    return res.status(200).json(parsedResult);
                } catch (parseError) {
                    // If double parsing fails, try single parse
                    try {
                        const parsedResult = JSON.parse(result);
                        return res.status(200).json(parsedResult);
                    } catch (singleParseError) {
                        // If all parsing fails, return as string
                        return res.status(200).json(result);
                    }
                }
            }
            
            // Check if it's a regular JSON string
            if ((result.startsWith('{') && result.endsWith('}')) || 
                (result.startsWith('[') && result.endsWith(']'))) {
                try {
                    const parsedResult = JSON.parse(result);
                    return res.status(200).json(parsedResult);
                } catch (parseError) {
                    // If parsing fails, treat it as a regular string
                    return res.status(200).json(result);
                }
            }
        }
        
        // For other data types, let express handle the serialization
        res.status(200).json(result);
    } catch (error) {
        // This catch block is for unexpected errors within processDslString itself
        // Errors handled by processDslString by returning "[ERROR: ...]" will be in `result`
        console.error("Error in /process-dsl route:", error);
        res.status(500).json({ message: "Internal server error while processing DSL.", error: error.message });
    }
});

app.get('*', async (req, res) => {

    const team = req.headers['x-header-team'];
    const project = req.headers['x-header-project'];

    const jwtOptions = {
        expiresIn: '20s'
      };

    let currentURL = req.url;

    currentURL = currentURL.replace(/^\//, '');

    let ip = req.headers['x-forwarded-for'] || req.connection.remoteAddress;

    const base64EncodedURL = Buffer.from(currentURL).toString('base64');

    if (ip) {
        ip = Buffer.from(ip).toString('base64');
    }

    if (!team || !project) {
        return res.status(400).json({ message: 'Missing required headers' });
    }

    const siteURL = process.env.SITE_URL; 

    const secretKey = process.env.SECRET_KEY;

    const token = jwt.sign({}, secretKey, jwtOptions);

    const targetUrl = `${siteURL}/api/v1/mock/${team}/${project}?url=${base64EncodedURL}&ip=${ip}&token=${token}`;

    if (!targetUrl) {
        return res.status(400).json({ message: 'No URL provided' });
    }

    try {
        console.log(`Target URL: ${targetUrl}`);

        const response = await axios.get(targetUrl, {
            timeout: 20000,
            proxy: false,
        });

        const { data } = response;

        const jsonData = data.json_data ?? null;

        const statusCode = data.status_code ?? null;

        if (!jsonData && !statusCode) {
            return res.status(200).send(data);
        }

        return res.status(statusCode).json(jsonData);
    } catch (error) {
        const nativeStatusCode = error.response?.status ?? 500;
        const upstreamData = error.response?.data;
        const contentType = error.response?.headers['content-type'];

    
        if (contentType && contentType.includes('text/html')) {
            return res.status(nativeStatusCode).send(upstreamData || "An error occurred");
        }
    
        if (typeof upstreamData === "object") {
            return res.status(nativeStatusCode).json(upstreamData);
        }
        return res.status(nativeStatusCode).send(upstreamData || "An error occurred");
    }

});

app.listen(port, () => {
    console.log(`Server is running on http://localhost:${port}`);
});
