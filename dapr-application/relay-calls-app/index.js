const { DaprClient, CommunicationProtocolEnum } = require("@dapr/dapr");
const { performance } = require('perf_hooks');
const express = require('express');
const bodyParser = require('body-parser');
const axios = require('axios');
const app = express();
app.use(bodyParser.json());

const daprPort = process.env.DAPR_HTTP_PORT || 3500;
const port = 3000;
const resourceURL = "https://httpbin.org/post";

// Echo the auth headers passed in
/* app.get('/echo', (req, res) => {
    var text = req.query.text;    
    console.log("Echoing: " + text);
    res.send("Access token: " + req.headers["authorization"] + " Text: " + text); 
}); */

// Route to fetch data from the resource server
app.get('/fetch-resource', async (req, res) => {
    try {
        console.log("Accessing resource from server at: " + resourceURL);
        //const token = req.headers["authorization"];
        const response = await axios.post(resourceURL,
            {
                "sample_body":"passing json"
            },
            {
            headers: {
                "X-custom": "chandan", // Include the Bearer token
            },
        });

        console.log("Express print: " + JSON.stringify(response.data, null,4));

        res.status(200).json(response.data);
    } catch (error) {
        console.error('Error fetching resource:', error.response?.data || error.message);
        res.status(500).json({ error: 'Failed to fetch resource' });
    }
});


app.post('/send-data', async (req, res) => {
    try {
        console.log("=====================Non-binding======================== " );
        console.log("Express print request: " + JSON.stringify(req.body, null,4) );
        //const token = req.headers["authorization"];
        const startTime = performance.now(); // Start timing
        const response = await axios.post(resourceURL,
            {
                ...req.body,
                "sample_body":"passing json"
            },
            {
            headers: {
                "X-custom": "chandan", // Include the Bearer token
            },
        });

        console.log("Express print resp from httpbin: " + JSON.stringify(response.data, null,4));
        const endTime = performance.now(); // End timing
        const timeTaken = endTime - startTime; // Calculate time taken
        console.log("=================="+timeTaken.toFixed(2)+"milliseconds=========================== " );
        res.status(200).json(response.data);
    } catch (error) {
        console.error('Error fetching resource:', error.response?.data || error.message);
        res.status(500).json({ error: 'Failed to fetch resource' });
    }
});



app.post('/trigger-binding', async (req, res) => {
    try {
        console.log("=====================binding======================== " );
        console.log("Express print request: " + JSON.stringify(req.body, null,4) );
        const startTime = performance.now(); // Start timing

        const BINDING_NAME = "httpbinfetch-binding";
        const BINDING_OPERATION = "create";
        const client = new DaprClient({
            daprHost:"127.0.0.1",
            daprPort: 3500,
            communicationProtocol: CommunicationProtocolEnum.HTTP,
        });
        //Using Dapr SDK to invoke output binding
        const result = await client.binding.send(BINDING_NAME, BINDING_OPERATION, req.body);

        console.log("Express print resp from httpbin: " + JSON.stringify(result, null,4));

        const endTime = performance.now(); // End timing
        const timeTaken = endTime - startTime; // Calculate time taken

        console.log("=================="+timeTaken.toFixed(2)+"milliseconds=========================== " );
        res.status(200).json(JSON.parse(result.data));
    } catch (error) {
        console.error('Error fetching resource:', error.response?.data || error.message);
        res.status(500).json({ error: 'Failed to fetch resource' });
    }
});
app.listen(port, () => console.log(`Node App listening on port ${port}!`));