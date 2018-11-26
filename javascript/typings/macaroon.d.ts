declare module "macaroon" {

    type dischargeFunc = (macLocation: string, cavLocation: string, cavID: string, success: Function, failure: Function) => void;

    export interface MacaroonCaveat {
        identifier: Uint8Array;
        location?: string;
        vid?: Uint8Array;
    }


    export interface MacaroonParams {
        identifierBytes: Uint8Array;
        locationStr: string;
        caveats: MacaroonCaveat[];
        signatureBytes: Uint8Array;
        version: number;
    }

    export class Macaroon {
        constructor(params: MacaroonParams);
        caveats: Array<MacaroonCaveat>;
        location: string;
        identifier: Uint8Array;
        signature: Uint8Array;
        addThirdPartyCaveat(
            rootKeyBytes: Uint8Array,
            caveatIdBytes: Uint8Array | string,
            locationStr: string): void;
        addFirstPartyCaveat(caveatIdBytes: Uint8Array): void;
        bindToRoot(rootSig: Uint8Array): void;
        clone(): Macaroon;
        verify(
            rootKeyBytes: Uint8Array,
            check: (condition: string) => string | null,
            discharges?: MacaroonCaveat[]): void;
        exportJSON(): MacaroonParams;
        exportBinary(): Uint8Array;
    }

    export function newMacaroon(params: MacaroonParams): Macaroon;

    export function importMacaroon(obj: string | Uint8Array | MacaroonParams): Macaroon | Macaroon[];

    export function importMacaroons(obj: string | Uint8Array | MacaroonParams | Array<MacaroonParams>): Macaroon[];

    export function dischargeMacaroon(macaroon: Macaroon, getDischarge: dischargeFunc, onOk: (macaroons: Macaroon[]) => void, onError: (error: string) => void): void;

    export function bytesToBase64(bytes: Uint8Array): string;

    export function base64ToBytes(s: string): Uint8Array;
}
