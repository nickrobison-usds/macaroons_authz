declare module "retryer" {
    interface RetryerOptions {
        total?: number;
        timeout?: number;
        onStart?: (attempt: number) => void;
        onError?: (error: Error, attempt: number) => void;

    }
    export default function retry<T>(promise: (args: any[]) => Promise<T>, options?: RetryerOptions): Promise<T>;
}
