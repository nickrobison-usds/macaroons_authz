import express from "express";
import { AuthController, CreateAuthController } from "./controllers/auth";
import cookieParser from "cookie-parser";
// import cookieParser = require("cookie-parser");
import { TokenController } from "./controllers/token";

console.log("Starting API Service");

(async () => {

    const app = express();

    app.use(cookieParser());

    // Express config

    app.set("port", process.env.Port || 3002);

    // Add the controllers and routes
    const ac = await CreateAuthController();
    const tc = new TokenController();

    app.get("/token", (req, res) => tc.getToken(req, res));
    app.get("/:acoID", (req, res) => ac.dischargeMacaroon(req, res));
    // Start it up
    //app.get("/:acoID/token", ())


    app.listen(app.get("port"), () => {
        console.log("App is running at http://localhost:%s in %s mode",
            app.get("port"),
            app.get("env"));

        console.log("Press CTRL-C to stop");
    })
})();
