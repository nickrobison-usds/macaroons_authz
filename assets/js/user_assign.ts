class JQueryPromise<T> {
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

interface IACONamePair {
    ID: string
    Name: string
}

$(document).ready(handleForm)


function handleForm($: JQueryStatic): void {
    console.log("Attaching");
    $("#assignEntity").change(fetchValuesOnChange);
    $(".open-AssignUserModal").click(e => openModalHandler(e));
}

function openModalHandler(e: JQueryEventObject): void {
    console.log("Clicked");
    const userID = $(e.currentTarget).data("id");
    console.debug("Setting userID:", userID);
    $("#userID").val(userID);
}


function fetchValuesOnChange(): void {
    console.log("Changed");
    console.log(this.value);

    const data = fetchData<IACONamePair[]>("GET", "/api/acos/list");
    data.then((d) => {
        console.log("Data:", d);

        const opts = $("#entityOptions").empty();

        const newOptions = d.map((option) => {
            return new Option(option.Name, option.ID, false, false);
        });

        opts.append(newOptions);
     })
}

async function fetchData<T>(method: string, url: string): JQueryPromise<T> {
    const data = <T>await $.getJSON(url);

    console.debug("Logging result:", data);

    return data;
}
