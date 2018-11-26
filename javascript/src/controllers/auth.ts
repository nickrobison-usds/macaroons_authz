import { Request, Response } from "express";
import Macaroon from "macaroon";

export class AuthController {

    public hello(req: Request, res: Response): void {
        res.json("Hello there!");
    }

    public dischargeMacaroon(req: Request, res: Response): void {

        console.log("Doing the discharge things.");
        // Get the macaroon from the reuest and import it.

        const token = req.query["id64"];

        console.log("token:", token);
        console.log("Macaroon:", Macaroon);
        const b = Macaroon.base64ToBytes(token);

        const mac = Macaroon.importMacaroon(b);
        if (AuthController.isSingleton(mac)) {
            // Print the caveats
            console.log("Caveats:")
            mac.caveats().forEach((cav) => {
                console.log("Caveat: ", cav);
            })
        }
    }

    private static isSingleton(mac: Macaroon.Macaroon | Macaroon.Macaroon[]): mac is Macaroon.Macaroon {
        return (<Macaroon.Macaroon>mac).location2 !== undefined;
    }
}
