import { Request, Response } from "express";
import { base64ToBytes, importMacaroons, Macaroon } from "macaroon";
import { readFileSync } from "fs";
import { resolve } from "path";
import { privateDecrypt } from "crypto";

const nacl: any = require("tweetnacl");
const pemtools: any = require("pemtools");

interface IKeyPair {
    pub: string;
    priv: string;
}

export class AuthController {
    private keys: IKeyPair;
    private privPemKey: string;
    private pubPemKey: string;
    private rootKey = "eeyiIuD5a2yrXjj6BlKctUC7k6qF/H6B";
    private rkEncode: Uint8Array;
    private privBytes: Buffer;
    private decoder: TextDecoder;

    constructor(privateKeyPath = "../user_keys.json") {
        console.debug("Creating controller");
        // Read the file
        const absPath = resolve(privateKeyPath);
        const json = JSON.parse(readFileSync(absPath, "utf8"));
        this.keys = {
            pub: json["public"],
            priv: json["private"]
        }

        this.privPemKey = "-----BEGIN EC PRIVATE KEY-----\n" + this.keys.priv + "\n-----END EC PRIVATE KEY-----";
        this.pubPemKey = "";
        this.rkEncode = base64ToBytes(this.rootKey);
        this.privBytes = Buffer.from(this.keys.priv, "utf8");
        console.log("Size:", this.privBytes.length);

        // Read real file?
        const file = readFileSync(resolve("../vendor__ca-key.pem"), "utf8");
        this.privPemKey = file;
        console.log(file);
        const pem = pemtools(file);
        console.log("Pem:", pemtools(file));
        this.decoder = new TextDecoder("utf-8");

    }
    public dischargeMacaroon(req: Request, res: Response): void {
        const acoID = req.params["acoID"];
        // Get the macaroon from the reuest and import it.
        const token = AuthController.getRequestMacaroons(req);
        console.log(`Verifying access for ACO ${acoID}\n`);

        // Decode the macaroons from base64 encoding
        const b = base64ToBytes(token);
        // Parse it as JSON and add the missing version key.
        // Not sure why it's missing, it might be an issue with the go library, or my code.
        let decoded: Array<any> = JSON.parse(this.decoder.decode(b));
        decoded = decoded.map((mac: any) => {
            mac["v"] = 2;
            return mac;
        })
        console.log(decoded);

        const mac = importMacaroons(decoded);

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
            console.error("Checked with error.")
            console.error(err);
            res.status(404).send("Nope, not authed.");
            return;
        }
        console.log("Verified");
        res.status(200).send("Successfully accessed data.");
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
            return "ACO ID does not match";
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
