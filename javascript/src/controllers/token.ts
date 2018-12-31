import { Request, Response } from "express";
import { Macaroon, base64ToBytes, newMacaroon, bytesToBase64 } from "macaroon";
import fetch from "node-fetch";
import { box, randomBytes } from "tweetnacl";
import { decodeUTF8, decodeBase64, encodeBase64 } from "tweetnacl-util";
import * as varint from "varint";
import { TextDecoder } from "util";


interface IJWKSResponse {
    k: string;
}

export class TokenController {

    constructor() {
        // Not used
    }

    public getToken(req: Request, res: Response): void {

        // Get the JWKS
        this.getPublicKey()
            .then((key) => {

                const keyPair = box.keyPair();

                const rootKey = "this is a test key, it should be long enough.";

                const nonce = randomBytes(24);

                const mac = newMacaroon({
                    identifier: "test identifier",
                    location: "http://localhost:8080",
                    rootKey: nonce,
                    version: 2
                });

                const keyBytes = base64ToBytes(key);
                console.debug("Key bytes: ", keyBytes);

                console.debug("Boxing");

                const msg = "This is a test message"

                const rootKeyLength = varint.encode(nonce.length);
                const decMSG = decodeUTF8(msg)


                const myPub = keyPair.publicKey

                // Sealed part.

                const fullmessage = new Uint8Array(
                    1
                    + rootKeyLength.length
                    + nonce.length
                    + decMSG.length);

                fullmessage.set([2], 0);
                fullmessage.set(rootKeyLength, 1);
                fullmessage.set(nonce, 1 + rootKeyLength.length);
                fullmessage.set(decMSG, 1 + rootKeyLength.length + nonce.length);

                console.debug("Box parts:");
                console.debug("Key:", [2]);
                console.debug(myPub.slice(0, 4));
                console.debug(myPub);
                console.debug("Nonce: ", nonce)
                console.debug(fullmessage);

                console.debug("Secret part:", fullmessage);

                // Seal the full thing?

                const shared = box.before(keyBytes, keyPair.secretKey);
                console.debug("Shared key");
                console.debug(shared);

                const sealed = box(fullmessage, nonce, keyBytes, keyPair.secretKey)

                const sealedB64 = bytesToBase64(sealed);
                const sealedBack = base64ToBytes(sealedB64);

                console.debug("Sealed:", sealed);
                console.debug("Sealed B64:", sealedB64);
                console.debug("SealedBack", sealedBack);

                // Now that we have the sealed message, we need to add more to the header?
                const withheader = new Uint8Array(
                    1
                    + 4
                    + 32
                    + 24
                    + sealedBack.length);

                withheader.set([2], 0);
                withheader.set(keyBytes.slice(0, 4), 1);
                withheader.set(myPub, 5);
                withheader.set(nonce, 5 + 32);
                withheader.set(sealedBack, 5 + 32 + 24);

                console.debug("Headers:", new TextDecoder("utf-8").decode(withheader));

                mac.addThirdPartyCaveat(nonce, withheader, "http://localhost:8080/api/users/verify");

                const macJSON = mac.exportJSON();

                console.debug("Exporting as JSON: ", macJSON);

                res.status(200).send(mac.exportJSON());
            });
    }

    private async getPublicKey(): Promise<string> {

        const resp = await fetch("http://localhost:8080/api/users/.well-known/jwks.json");


        const jwks: IJWKSResponse = await resp.json();
        console.debug("JWKS:", jwks);

        //return jwks.k;
        return "+VN3lXu8QBuH561ueQcN0vo7LCDtTdH8jQWtz2VaSRs=";
    }
}
