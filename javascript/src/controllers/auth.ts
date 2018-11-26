import { Request, Response } from "express";
import { base64ToBytes, importMacaroon, Macaroon, bytesToBase64 } from "macaroon";
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
    private rootKey = "eeYrGr8rnzSnPNvrwV4gIFAkqxCG8dBA";
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
    public hello(req: Request, res: Response): void {
        res.json("Hello there!");
    }

    public dischargeMacaroon(req: Request, res: Response): void {
        // Get the macaroon from the reuest and import it.
        const token = req.query["id64"];

        console.log("token:", token);
        const b = base64ToBytes(token);
        // Decrypt the token and import it.
        /*
                eccrypto.decrypt(this.privBytes, b)
                    .then((plaintext: string) => console.log("Decrypted:", plaintext))
                    .catch((e: any) => console.log("Error:", e));*/

        const mac = importMacaroon(token);
        if (AuthController.isSingleton(mac)) {
            // Print the caveats
            console.log("Caveats:")
            mac.caveats.forEach((cav) => {
                console.log(cav);
                console.log("Caveat: ", this.decoder.decode(cav.identifier));
            })

            try {
                mac.verify(base64ToBytes(this.rootKey), ((cond) => {
                    console.log("In check:", cond);
                    return null;
                }))
            } catch (err) {
                console.log("Checked with error.")
                console.log(err);
                res.status(404).send("Nope, not authed.");
            }
            console.log("Verified");
        }
    }

    private static isSingleton(mac: Macaroon | Macaroon[]): mac is Macaroon {
        return (<Macaroon>mac).location !== undefined;
    }
}
