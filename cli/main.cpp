#include <iostream>
#include <CLI11.hpp>
#include <cpprest/http_client.h>
#include <fmt/format.h>
#include <spdlog/spdlog.h>
#include <spdlog/sinks/stdout_color_sinks.h>
#include <httpbakery/client.hpp>
#include <httpbakery/interceptor.hpp>
#include <bakery/Macaroon.hpp>

using namespace std;
using namespace utility;                    // Common utilities like string conversions
using namespace web;                        // Common features like URIs.
using namespace web::http;                  // Common HTTP functionality
using namespace web::http::client;          // HTTP client features
using namespace concurrency::streams;       // Asynchronous streams

struct UserInterceptor : public Interceptor {
    http_request intercept(http_request request) const {
        cout << "Intercepting! from the tester!" << endl;
        return request;
    }
};

int main(int argc, char **argv) {
    // Setup the logger
    const auto console = spdlog::stdout_color_st("console");

    CLI::App app{"CLI client for CMS AuthZ Demo"};

    std::optional<string> aco_name_opt;
    std::optional<string> user_name_opt;
    std::optional<string> vendor_name_opt;
    bool gather_discharges = true;
    bool java_service = false;

    app.add_flag("--no-discharge", [&gather_discharges](size_t count) {
        if (count > 0) {
            gather_discharges = false;
        }
    }, "Disable gathering required discharges");
    app.add_flag("--java", [&java_service](size_t count) {
        if (count > 0) {
            java_service = true;
        }
    });
    app.add_option("user", user_name_opt, "User to perform queries as");
    app.add_option("aco", aco_name_opt, "ACO to query against");
    app.add_option("--vendor", vendor_name_opt, "Vendor name to lookup");

    try {
        app.parse(argc, argv);
    } catch (const CLI::ParseError &e) {
        return app.exit(e);
    }

    // Validate args
    string user_name;
    string aco_name;
    // std::optional doesn't work on MacOS <10.14, so we'll need to come up with a workaround.
    // This is fine for now, but clunky.
    if (!user_name_opt) {
        console->critical("Must provide a username to query as.");
        return 1;
    } else {
        user_name = *user_name_opt;
    }

    if (!aco_name_opt) {
        console->critical("Must provide an ACO Name to query against.");
        return 1;
    } else {
        aco_name = *aco_name_opt;
    }

    console->info("Starting up demo client");

    // Try to lookup a given ACO ID

    console->info("Looking up ID for ACO '{:s}'", aco_name);

    http_client nameClient(U("http://localhost:8080"));
    uri_builder nameBuilder(U("/api/acos/find"));
    nameBuilder.append_query(U("name"), U(aco_name));

    string acoID;

    http_request name_req(methods::GET);
    name_req.set_request_uri(nameBuilder.to_uri());

    auto nameTask = nameClient.request(name_req)
            .then([](http_response resp) {
                if (resp.status_code() == status_codes::OK) {
                    return resp.extract_string();
                }
                throw invalid_argument(resp.extract_string().get());
            });
    try {
        acoID = nameTask.get();
    } catch (const exception &e) {
        console->critical("Error getting aco ID: {:s}", e.what());
        return 1;
    }

    // And a User ID
    string user_id;

    console->info("Looking up ID for user '{:s}'", user_name);

    uri_builder userBuilder(U("/api/users/find"));
    userBuilder.append_query(U("name"), U(user_name));
    http_request user_req(methods::GET);
    user_req.set_request_uri(userBuilder.to_uri());

    auto userTask = nameClient.request(user_req)
            .then([](http_response resp) {
                if (resp.status_code() == status_codes::OK) {
                    return resp.extract_string();
                }
                throw invalid_argument(resp.extract_string().get());
            });
    try {
        user_id = userTask.get();
    } catch (const exception &e) {
        console->critical("Unable to get User ID: {:s}", e.what());
        return 1;
    }

    // And a vendor ID, if one was given
    optional<string> vendor_id;

    if (vendor_name_opt) {
        console->info("Looking up ID for vendor '{:s}'", *vendor_name_opt);

        uri_builder vendorBuilder(U("/api/vendors/find"));
        vendorBuilder.append_query(U("name"), U(*vendor_name_opt));
        http_request vendor_req(methods::GET);
        vendor_req.set_request_uri(vendorBuilder.to_uri());

        auto vendorTask = nameClient.request(vendor_req)
                .then([](http_response resp) {
                    if (resp.status_code() == status_codes::OK) {
                        return resp.extract_string();
                    }
                    throw invalid_argument(resp.extract_string().get());
                });
        try {
            vendor_id.emplace(vendorTask.get());
        } catch (const exception &e) {
            console->critical("Unable to get Vendor ID: {:s}", e.what());
            return 1;
        }
    }

    string token;

    // Java service request

    if (java_service) {
        console->info("Making request to Java service");

        http_client standaloneClient(U(fmt::format("http://localhost:3002/{}", acoID)));
        uri_builder standaloneBuilder(U("/token"));
        standaloneBuilder.append_query(U("user_id"), user_id);

        // Add the vendor id, if we're making a request on their behalf
        if (vendor_id) {
            standaloneBuilder.append_query(U("vendor_id"), *vendor_id);
        }

        http_request standalone_req(methods::GET);
        standalone_req.set_request_uri(standaloneBuilder.to_uri());

        const auto client_task = standaloneClient.request(standalone_req)
        .then([](http_response resp) {
            if (resp.status_code() == status_codes::OK) {
                return resp.extract_string();
            }
            throw invalid_argument(resp.extract_string().get());
        });

        try {
            token = client_task.get();
        } catch (const exception &e) {
            console->critical("Error making service Request: {:s}", e.what());
            return 1;
        }
    }
    else {
        console->info("Getting token from ACO manager");


        console->info("Looking up Macaroon for '{:s}' associated with '{:s}'", aco_name, user_name);

        // Try to find the ACO token associated with the user
        std::string token_query = fmt::format("api/users/token/{}/ACO/{}", user_id, acoID);
        uri_builder tokenBuilder(U(token_query));
        http_request token_req(methods::GET);
        token_req.set_request_uri(tokenBuilder.to_uri());

        auto tokenTask = nameClient.request(token_req)
                .then([](http_response resp) {
                    if (resp.status_code() == status_codes::OK) {
                        return resp.extract_string();
                    }
                    throw invalid_argument(resp.extract_string().get());
                });
        try {
            token = tokenTask.get();
        } catch (const exception &e) {
            console->critical("Unable to get user token; {:s}", e.what());
        }
    }

    auto mac = Macaroon::importMacaroons(token);
    // Debug
//    console->debug("Inspected macaroon: {:s}", mac.inspect());

    string bound_mac;


    // Try to bind macaroons
    if (gather_discharges) {
        console->info("Discharging third party caveats");
        Client mac_client;
        const UserInterceptor tic;
        mac_client.addInterceptor("http://local.test", tic);
        bound_mac = mac_client.dischargeMacaroon(mac);
//        bound_mac = mac.discharge_all_caveats();
//bound_mac = "REMOVE ME!!!";
    } else {
        console->info("Not discharging caveats");
        bound_mac = mac.serialize(MACAROON_V2J);
    }

    // Now make the actual request for the ACO data

    console->info("Making request to endpoint.");

    http_client client(U("http://localhost:3002"));
    uri_builder builder(U("/" + acoID));

//    Attach the macaroon as a cookie
    http_request req(methods::GET);
    name_req.headers().add(U("Cookie"), U(fmt::format("macaroon-1={:s};", bound_mac)));
    name_req.set_request_uri(builder.to_uri());

    auto task = client.request(name_req)
            .then([](http_response response) {
                return response.extract_string();
            });

    const string resp = task.get();
    console->info(resp);
    return 0;
}
