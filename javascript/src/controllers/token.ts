import { Request, Response } from "express";
import { Macaroon, base64ToBytes, newMacaroon, bytesToBase64 } from "macaroon";
import fetch from "node-fetch";
import { box, randomBytes } from "tweetnacl";
import { decodeUTF8, encodeUTF8, decodeBase64, encodeBase64 } from "tweetnacl-util";
import * as varint from "varint";
import { TextDecoder, TextEncoder } from "util";


interface IJWKSResponse {
    k: string;
}

export class TokenController {

    private encoder: TextEncoder;

    constructor() {
        // Not used
        this.encoder = new TextEncoder();
    }

    public getToken(req: Request, res: Response): void {

        const id = req.query["user_id"];
        console.debug("User ID: ", id);
        // Get the JWKS
        this.getPublicKey()
            .then((thirdPartyPublicKey) => {

                const keyPair = box.keyPair();

                const rootKeyStr = "this is a test key, it should be long enough.";
                const nonceStr = "this is a test nonce,...";
                const nonce = this.encoder.encode(nonceStr);
                const rootKey = decodeUTF8(rootKeyStr);

                console.debug("Encoding with root key: ", rootKey);

                //const nonce = randomBytes(24);

                const mac = newMacaroon({
                    identifier: "test identifier",
                    location: "http://localhost:8080",
                    rootKey: rootKey,
                    version: 2
                });

                const thirdPartyKeyBytes = base64ToBytes(thirdPartyPublicKey);
                console.debug("Key bytes: ", thirdPartyKeyBytes);

                console.debug("Boxing");

                const msg = "user_id= " + id;

                const rootKeyLength = varint.encode(rootKey.length);
                const decMSG = decodeUTF8(msg)


                const myPub = keyPair.publicKey

                // Sealed part.

                const fullmessage = new Uint8Array(
                    1
                    + rootKeyLength.length
                    + rootKey.length
                    + decMSG.length);

                fullmessage.set([2], 0);
                fullmessage.set(rootKeyLength, 1);
                fullmessage.set(rootKey, 1 + rootKeyLength.length);
                fullmessage.set(decMSG, 1 + rootKeyLength.length + rootKey.length);

                console.debug("Secret part:", fullmessage);

                // Seal the full thing?

                const shared = box.before(thirdPartyKeyBytes, keyPair.secretKey);
                console.debug("Shared key");
                console.debug(shared);

                const sealed = box(fullmessage, nonce, thirdPartyKeyBytes, keyPair.secretKey)

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
                withheader.set(thirdPartyKeyBytes.slice(0, 4), 1);
                withheader.set(myPub, 5);
                withheader.set(nonce, 5 + 32);
                withheader.set(sealedBack, 5 + 32 + 24);

                console.debug("Headers:", new TextDecoder("utf-8").decode(withheader));

                mac.addThirdPartyCaveat(rootKey, withheader, "http://localhost:8080/api/users/verify");

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
