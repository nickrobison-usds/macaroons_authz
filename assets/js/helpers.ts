export class JQueryPromise<T> {
    constructor(executor: (resolve: (value?: T | PromiseLike<T>) => void, reject: (reason?: any) => void) => void) {
        let dfd = $.Deferred<T>();
        function fulfilled(value?: T | PromiseLike<T>) {
            let promise = <PromiseLike<T>>value;
            if (value && promise.then) {
                promise.then(fulfilled, rejected);
            }
            else {
                dfd.resolve(<T>value);
            }
        }
        function rejected(reason) {
            let promise = <PromiseLike<T>>reason;
            if (reason && promise.then) {
                promise.then(fulfilled, rejected);
            }
            else {
                dfd.reject(<T>reason);
            }
        }
        executor(fulfilled, rejected);
        return dfd.promise();
    }
}


export class AsyncFetch {

    public static async fetchValuesOnChange<T>(value: string, builder: (value: T) => HTMLOptionElement): Promise<HTMLOptionElement[]> {

        const values = await AsyncFetch.fetchData<T[]>("GET", "/api/" + String(value).toLocaleLowerCase() + "s/list");

        return values.map(builder);
    };

    private static async fetchData<T>(method: string, url: string): Promise<T> {
        const data = <T> await Promise.resolve($.getJSON(url));

        console.debug("Logging result:", data);

        return data;
    }
}
