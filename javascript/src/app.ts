import express from "express";
import { AuthController, CreateAuthController } from "./controllers/auth";
import cookieParser = require("cookie-parser");

console.log("Starting API Service");

(async () => {

    const app = express();

    app.use(cookieParser());

    // Express config

    app.set("port", process.env.Port || 3002);

    // Add the controllers and routes
    const ac = await CreateAuthController();

    app.get("/:acoID", (req, res) => ac.dischargeMacaroon(req, res));
    // Start it up

    app.listen(app.get("port"), () => {
        console.log("App is running at http://localhost:%s in %s mode",
            app.get("port"),
            app.get("env"));

        console.log("Press CTRL-C to stop");
    })
})();
