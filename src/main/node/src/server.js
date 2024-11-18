const express = require('express');
const axios = require('axios');
require('dotenv').config({ path: '../.env' });

const cors = require('cors');

const corsOptions = {
    origin: '*',
};


const app = express();
const port = 3001;

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

    const targetUrl = `${siteURL}/api/v1/mock/${team}/${project}?url=${base64EncodedURL}&ip=${ip}`;

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
            throw new Error('Something went wrong');
        }

        return res.status(statusCode).json(jsonData);
    } catch (error) {
        const nativeStatusCode = error.response?.status ?? 500;

        const upstreamData = error.response?.data;

        if (typeof upstreamData === "object") {
            return res.status(nativeStatusCode).json(upstreamData);
        }
        return res.status(nativeStatusCode).send(upstreamData || "An error occurred");
    }

});

app.listen(port, () => {
    console.log(`Server is running on http://localhost:${port}`);
});
