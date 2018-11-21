declare module "js-macaroon" {

    function newMaracoon(): Macaroon;

    interface Macaroon {
        identifier: string | Uint8Array;
        location: string | null | undefined;
        rootKey: string | Uint8Array;
        version: number;
    }
}
