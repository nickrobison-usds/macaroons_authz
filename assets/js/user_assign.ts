import { AsyncFetch, IACONamePair, UserAssign } from "./helpers"
$(document).ready(handleForm)


function handleForm($: JQueryStatic): void {
    console.log("Attaching");
    $("#assignEntity").change(changeHandler);
    $(".open-AssignUserModal").click(e => openModalHandler(e));
}

function openModalHandler(e: JQuery.Event<HTMLElement, null>): void {
    console.log("Clicked");
    const userID = $(e.currentTarget).data("id");
    console.debug("Setting userID:", userID);
    $("#userID").val(userID);
}

function changeHandler(value: JQuery.Event<HTMLElement, null>): void {

    const val = this.value;
    console.log("Value:", val)
    AsyncFetch.fetchValuesOnChange<IACONamePair>(val, UserAssign.BuildACOOption)
        .then((newOptions) => {
            const opts = $("#entityOptions").empty();
            opts.append(newOptions);
        })
}



