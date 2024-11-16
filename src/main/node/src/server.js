const express = require('express');
const axios = require('axios');
require('dotenv').config({ path: '../.env' });

const app = express();
const port = 3001;

app.use(express.json());

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

    const base64EncodedURL = Buffer.from(currentURL).toString('base64');

    if (!team || !project) {
        return res.status(400).json({ message: 'Missing required headers' });
    }

    const siteURL = process.env.SITE_URL;

    const targetUrl = `${siteURL}/api/v1/mock/${team}/${project}?url=${base64EncodedURL}`;

    if (!targetUrl) {
        return res.status(400).json({ message: 'No URL provided' });
    }

    try {
        console.log(`Team: ${team}, Project: ${project}`);
      
        const response = await axios.get(targetUrl, {
            timeout: 10000,
            proxy: false,
        });

        return res.status(response.status).json(response.data);
    } catch (error) {
        return res.status(500).json({ message: 'Error forwarding the request', error: error.message });
    }
});

app.listen(port, () => {
    console.log(`Server is running on http://localhost:${port}`);
});
