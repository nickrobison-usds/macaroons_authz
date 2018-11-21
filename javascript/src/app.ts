import express from "express";
import { AuthController } from "./controllers/auth";

console.log("Starting User Service");

const app = express();

// Express config

app.set("port", process.env.Port || 3002);

// Add the controllers and routes
const ac = new AuthController();

app.get("/", ac.hello);
app.post("/users/verify/discharge", ac.dischargeMacaroon);
// Start it up

app.listen(app.get("port"), () => {
    console.log("App is running at http://localhost:%s in %s mode",
        app.get("port"),
        app.get("env"));

    console.log("Press CTRL-C to stop");
})
