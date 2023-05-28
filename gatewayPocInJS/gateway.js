const express = require("express");
const axios = require("axios");
const bodyParser = require("body-parser");

const app = express();
const port = 4000;

app.use(bodyParser.json());

const threshold = 2;

let requestCount = 0;
let startTime = new Date();
// Endpoint that forwards the request

const checkAuthorization = (req, res, next) => {
  const authHeader = req.headers.authorization;

  if (!authHeader) {
    return res.status(401).json({ error: "Authorization header missing" });
  }

  //   // Extract the key from the authorization header
  //   const [, key] = authHeader.split(" ");

  //   // Array of valid keys
  const validKeys = ["key1", "key2", "key3"];

  //   console.log("keys are:");
  //   console.log(authHeader);
  if (!validKeys.includes(authHeader)) {
    return res.status(403).json({ error: "Invalid authorization key" });
  }

  // Proceed to the next middleware or route handler
  next();
};

const srviceMap = {
  service1: {
    url: "http://localhost:3000/modify-message",
  },
  service2: {
    url: "http://localhost:3001/modify-message",
  },
};

app.post("/service-forward", checkAuthorization, async (req, res) => {
  try {
    if (new Date() - startTime > 10000) {
      requestCount = 0;
      startTime = new Date();
    }
    if (requestCount > threshold) {
      res.status(500).json({ error: "Too many frequent requests!!" });
      return;
    }
    requestCount++;
    const serviceUrl = srviceMap[req.query.service].url;
    const response = await axios.post(serviceUrl, req.body);
    console.log("response is:", response);
    res.json(response.data);
  } catch (error) {
    console.log(error.stack);
    res.status(500).json({ error: "Something went wrong" });
  }
});

// Start the server
app.listen(port, () => {
  console.log(`Server listening on port ${port}`);
});
