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
        const cook = AuthController.getRequestMacaroons(req);
        console.log("Cookie:", req.cookies);
        let token = req.cookies["macaroon-4809536c6ce8caa9dc97f074780ac3da29d3cb2d0ea64c8bf87cb43322832279"];


        console.log("ACO ID:", acoID);
        console.log("token:", token);
        const b = base64ToBytes(token);

        let decoded: Array<any> = JSON.parse(this.decoder.decode(b));
        // Add version 2 manually. Not sure why it's missing.
        decoded = decoded.map((mac: any) => {
            mac["v"] = 2;
            return mac;
        })
        console.log(decoded);

        const mac = importMacaroons(decoded);

        console.log(mac);
        if (AuthController.isSingleton(mac)) {
            // Print the caveats
            console.log("Caveats:")
            mac.caveats.forEach((cav) => {
                console.log(cav);
                console.log("Caveat: ", this.decoder.decode(cav.identifier));
            })

            try {

                mac.verify(base64ToBytes(this.rootKey), ((cond) => AuthController.verifyACOID(cond, acoID)));
            } catch (err) {
                console.error("Checked with error.")
                console.error(err);
                res.status(404).send("Nope, not authed.");
                return;
            }
        } else {
            // If it's an array, check to see if we have any discharges.
            const discharges = mac
                .filter((m) => m.location === null);
            console.log("Have discharges: ", discharges.length);
            for (const d of discharges) {
                console.log("Discharge:", this.decoder.decode(d.identifier));
            }
            try {
                mac[0].verify(base64ToBytes(this.rootKey), ((cond) => AuthController.verifyACOID(cond, acoID)), discharges);
            } catch (err) {
                console.error("Checked with error.");
                console.error(err);
                res.status(404).send("Nope, not authed.");
                return;
            }
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

    private static getRequestMacaroons(req: Request): string {
        const rc = req.headers.cookie;
        console.log(rc);
        let value = "";
        if (rc) {
            const splitCookies = AuthController.splitStringMaybeArray(rc, " ");
            console.log("Split:", splitCookies);
            for (const cook of splitCookies) {
                if (cook.startsWith("macaroon-")) {
                    value += cook.replace("macaroon-", "");
                }
            }
        }

        return "";
    }

    private static splitStringMaybeArray(str: string | string[], delim: string): string[] {
        if (AuthController.isStringArray(str)) {
            return str;
        }
        return str.split(" ");

    }

    private static isStringArray(str: string | string[]): str is string[] {
        return (<string[]>str).length !== undefined;
    }
}
