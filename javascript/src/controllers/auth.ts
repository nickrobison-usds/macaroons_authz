import { Request, Response } from "express";
import Macaroon from "js-macaroon";

export class AuthController {

    public hello(req: Request, res: Response): void {
        res.json("Hello there!");
    }

    public dischargeMacaroon(req: Request, res: Response): void {

        console.log("Doing the discharge things.");
        // Get the macaroon from the reuest and import it.

        const b = Macaroon.base64ToBytes(req.params["id64"]);
        const mac = Macaroon.importMacaroon(b);
        if (AuthController.isSingleton(mac)) {
            // Print the caveats
            console.log("Caveats:")
            mac.caveats().forEach((cav) => {
                console.log("Caveat: ", cav);
            })
        }

    }

    private static isSingleton(mac: Macaroon.Macaroon.Macaroon | Macaroon.Macaroon.Macaroon[]): mac is Macaroon.Macaroon.Macaroon {
        return (<Macaroon.Macaroon.Macaroon>mac).location2 !== undefined;
    }
}
