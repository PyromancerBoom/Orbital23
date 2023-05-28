const express = require("express");
const bodyParser = require("body-parser");

const app = express();
const port = 3001;

// Middleware to parse JSON bodies
app.use(bodyParser.json());

// Endpoint that modifies the "message" field
app.post("/modify-message", (req, res) => {
  // Assuming the incoming JSON has a "message" field
  if (req.body.message) {
    req.body.message = "Server2: " + req.body.message;
    res.json(req.body);
  } else {
    res.status(400).json({ error: "Invalid JSON payload" });
  }
});

// Start the server
app.listen(port, () => {
  console.log(`Server listening on port ${port}`);
});
