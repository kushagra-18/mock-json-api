const express = require('express');
const axios = require('axios');
const jwt = require('jsonwebtoken');
const { URL } = require('url');
const {executeCode} = require('./faker');

require('dotenv').config({ path: '../.env' });

const cors = require('cors');

const corsOptions = {
    origin: '*',
};


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

    const targetUrl = new URL(`/api/v1/mock/${team}/${project}`, siteURL);
    
    targetUrl.searchParams.append('url', base64EncodedURL);
    targetUrl.searchParams.append('ip', ip);
    targetUrl.searchParams.append('token', token);

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

        const content = data.content ?? null;

        const statusCode = data.status_code ?? null;

        const mockType = data.mock_type ?? null;

        if (!content || !statusCode) {
            return res.status(200).send(data);
        }

        if(mockType === 'FAKER') {
            const result = await executeCode(content);
            if(result.success){
                return res.status(statusCode).json(result.result);
            }else{
                return res.status(statusCode).json(result.errors);
            }
        }else if(mockType === 'XML'){
            res.set('Content-Type', 'application/xml');
            console.log(content);
            return res.status(statusCode).send(content);
        }

        return res.status(statusCode).json(content);
    } catch (error) {

        console.error(error);

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
