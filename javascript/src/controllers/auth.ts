import { Request, Response } from "express";
import { base64ToBytes, importMacaroon, Macaroon, importMacaroons } from "macaroon";


interface IKeyPair {
    pub: string;
    priv: string;
}

export class AuthController {
    private rootKey = "eeyiIuD5a2yrXjj6BlKctUC7k6qF/H6B";
    private decoder: TextDecoder;

    constructor(privateKeyPath = "../user_keys.json") {
        console.debug("Creating controller");
        this.decoder = new TextDecoder("utf-8");

    }
    public dischargeMacaroon(req: Request, res: Response): void {
        const acoID = req.params["acoID"];
        // Get the macaroon from the reuest and import it.
        const token = AuthController.getRequestMacaroons(req);
        console.log(`Verifying access for ACO ${acoID}\n`);

        // Decode the macaroons from base64 encoding
        const b = base64ToBytes(token);
        const mac = this.importMacaroon(token);
        console.log(mac);

        // Verify the macaroon and any discharges
        const macaroons = AuthController.getMacaroonAndDischarges(mac);

        // Print the caveats
        console.log("Caveats:")
        macaroons[0].caveats.forEach((cav) => {
            console.log(cav);
            console.log("Caveat: ", this.decoder.decode(cav.identifier));
        })

        try {
            macaroons[0].verify(base64ToBytes(this.rootKey), ((cond) => AuthController.verifyACOID(cond, acoID)), macaroons[1]);
        } catch (err) {
            console.error(err);
            res.status(404).send(err.message);
            return;
        }
        console.log("Verified");
        res.status(200).send("Successfully accessed data.");
    }

    /** Imports a macaroon, or a slice of macaroons, from a base64 encoded string
        */
    private importMacaroon(token: string): Macaroon | Macaroon[] {
        const b = base64ToBytes(token);
        const decoded = this.decoder.decode(b);
        if (decoded[0] == "[") {
            console.log("Decoded:", decoded);
            console.log("Importing array of macaroons");
            return importMacaroons(b);
        } else {
            return importMacaroon(b);
        }
    }

    private static isSingleton(mac: Macaroon | Macaroon[]): mac is Macaroon {
        return (<Macaroon>mac).location !== undefined;
    }

    private static verifyACOID(condition: string, acoID: string): string | null {
        // Split the condition based on the first space
        const splits = condition.split("= ");
        if (splits[0] == "aco_id") {
            if (splits[1] == acoID) {
                return null;
            }
            return `This token is only valid for ACO: ${splits[1]}`;
        }
        return null;
    }

    // This expects the cookies to be already parsed and ready to go.
    private static getRequestMacaroons(req: Request): string {
        const rc: { [name: string]: string; } = req.cookies;
        console.log(rc);

        // Iterate through the cookies and find anything name that starts with macaroon-
        let value = "";
        for (const key in rc) {
            if (key.startsWith("macaroon-")) {
                value += rc[key];
            }
        }
        return value;
    }

    private static getMacaroonAndDischarges(mac: Macaroon | Macaroon[]): [Macaroon, Macaroon[]] {
        if (AuthController.isSingleton(mac)) {
            return [mac, []];
        }

        const m = mac[0];
        const discharges = mac.filter((m) => m.location === null);
        return [m, discharges];
    }
}
