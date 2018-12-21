import { Request, Response } from "express";
import { Macaroon, base64ToBytes, newMacaroon, bytesToBase64 } from "macaroon";
import fetch from "node-fetch";
import { box, randomBytes } from "tweetnacl";
import { decodeUTF8, decodeBase64, encodeBase64 } from "tweetnacl-util";
import * as varint from "varint";


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

                console.debug("Has key:", key);

                console.debug("Box:", box);

                const keyPair = box.keyPair();

                console.debug("KeyPair:", keyPair);

                const rootKey = "this is a test key, it should be long enough.";

                const mac = newMacaroon({
                    identifier: "test identifier",
                    location: "http://localhost:8080",
                    rootKey: rootKey,
                    version: 2
                });

                const keyBytes = base64ToBytes(key);
                const nonce = randomBytes(24);

                console.debug("Key bytes: ", keyBytes);

                console.debug("Boxing");

                const msg = "This is a test message"

                const keyLen = varint.encode(nonce.length);
                const decMSG = decodeUTF8(msg)


                // Seal it?

                // Add everything that we need

                const fullmessage = new Uint8Array(
                    1
                    + keyLen.length
                    + nonce.length
                    + decMSG.length
                );


                fullmessage.set([2], 0);
                fullmessage.set(keyLen, 1);
                fullmessage.set(nonce, 1 + keyLen.length);
                fullmessage.set(decMSG, 1 + keyLen.length + nonce.length);

                console.debug("Box parts:");
                console.debug("Key:", [2]);
                console.debug(keyBytes.slice(0, 4));
                console.debug(keyBytes);
                console.debug("Nonce: ", nonce)
                console.debug(fullmessage);

                console.debug("Secret part:", fullmessage);

                // Seal the full thing?

                const sealed = box(fullmessage, nonce, keyBytes, keyPair.secretKey)

                const sealedB64 = bytesToBase64(sealed);
                const sealedBack = base64ToBytes(sealedB64);

                // Now that we have the sealed message, we need to add more to the header?
                const withheader = new Uint8Array(
                    1
                    + 4
                    + 32
                    + 24
                    + sealedBack.length);

                withheader.set([2], 0);
                withheader.set(keyBytes.slice(0, 4), 1);
                withheader.set(keyBytes, 5);
                withheader.set(nonce, 5 + 32);
                withheader.set(sealedBack, 5 + 32 + 24);

                console.debug("Headers:", new TextDecoder("utf-8").decode(withheader));



                mac.addThirdPartyCaveat(nonce, withheader, "http://localhost:8080/api/users/verify");

                const macJSON = mac.exportJSON();

                console.debug("Exporting as JSON: ", macJSON);

                res.status(200).send(bytesToBase64(mac.exportBinary()));
            });
    }

    private async getPublicKey(): Promise<string> {

        const resp = await fetch("http://localhost:8080/api/users/.well-known/jwks.json");


        const jwks: IJWKSResponse = await resp.json();
        console.debug("JWKS:", jwks);

        return jwks.k;
    }
}
