import {AsyncFetch, IACONamePair, UserAssign} from "./helpers"

$(document).ready(handleVendorForm)


function handleVendorForm($: JQueryStatic): void {
    console.log("Attaching");
    $(".open-AssignVendorModal").click(e => openVendorModalHandler(e));
}

function openVendorModalHandler(e: JQuery.Event<HTMLElement, null>): void {
    console.log("Clicked");
    const vendorID = $(e.currentTarget).data("id");
    console.debug("Setting vendorID:", vendorID);
    $("#vendorID").val(vendorID);

    // Fetch all the ACOs that we can assign to.
    AsyncFetch.fetchValuesOnChange<IACONamePair>("aco", UserAssign.BuildACOOption)
        .then((newOptions) => {
            const opts = $("#entityOptions").empty();
            opts.append(newOptions);
        })

}
