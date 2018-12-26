import { Request, Response } from "express";
import { base64ToBytes, importMacaroon, Macaroon, importMacaroons } from "macaroon";
import { Client, ClientConfig } from "pg";
import retry from "retryer";
import { TextEncoder, TextDecoder } from "util";

interface IKeyPair {
    pub: string;
    priv: string;
}

interface IVaultKeyResponse {
    keys: { [id: string]: string };
    name: string;
}

interface IVaultResponse {
    request_id: string;
    lease_id?: string;
    renewable: boolean;
    lease_duration: number;
    data: IVaultKeyResponse;
    wrap_info?: string;
    warnings?: string;
    auth?: string;
}

export class AuthController {
    private rootKey: string;
    private decoder: TextDecoder;

    constructor(key: string, privateKeyPath = "../user_keys.json") {
        console.debug("Creating controller with key: ", key);
        this.decoder = new TextDecoder("utf-8");
        this.rootKey = key;
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



        try {
            console.log("Decrypting with key:", this.rootKey);
            macaroons[0].verify(base64ToBytes(this.rootKey), ((cond) => AuthController.verifyACOID(cond, acoID)), macaroons[1]);
        } catch (err) {
            console.error(err);
            res.status(401).send(err.message);
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
        console.log("Decoded:", decoded);
        if (decoded[0] == "[") {
            console.log("Importing array of macaroons");
            // Check for JSON
            var toImport: string;
            if (decoded[1] == "{") {
                const parsed = JSON.parse(decoded);
                console.log("Parsed:", parsed);
                toImport = parsed;
            } else {
                toImport = decoded;
            }
            return importMacaroons(toImport);

        } else {
            // Check for JSON
            var toImport: string;
            if (decoded[0] == "{") {
                toImport = JSON.parse(decoded);
            } else {
                toImport = decoded;
            }
            return importMacaroon(toImport);
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

export async function CreateAuthController(): Promise<AuthController> {

    const connectionString = process.env.DATABASE_URL ? process.env.DATABASE_URL : "postgres://raac:@127.0.0.1:5432/macaroons_authz_development?sslmode=disable";
    const client = await retry(() => connectToDB({
        connectionString: connectionString,
    }), {
            timeout: 10000,
        });

    // We cheat and just grab the first key
    const res = await client.query("SELECT encode(rootkey::bytea, 'base64') as key FROM public.root_keys ORDER BY id ASC LIMIT 1");
    console.log("Result: ", res);
    await client.end();

    return new AuthController(res.rows[0]["key"]);
}

async function connectToDB(options: ClientConfig): Promise<Client> {
    const client = new Client(options);
    await client.connect();
    return client;
}
