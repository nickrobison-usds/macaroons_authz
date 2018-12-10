#include <iostream>
#include <CLI11.hpp>
#include <cpprest/http_client.h>
#include <fmt/format.h>
#include <spdlog/spdlog.h>
#include <spdlog/sinks/stdout_color_sinks.h>
#include "Macaroon.hpp"

using namespace std;
using namespace utility;                    // Common utilities like string conversions
using namespace web;                        // Common features like URIs.
using namespace web::http;                  // Common HTTP functionality
using namespace web::http::client;          // HTTP client features
using namespace concurrency::streams;       // Asynchronous streams

int main(int argc, char **argv) {
    // Setup the logger
    const auto console = spdlog::stdout_color_st("console");

    CLI::App app{"CLI client for CMS AuthZ Demo"};

    string filename = "default";

    app.add_option("-f,--file", filename, "Config file.");

    try {
        app.parse(argc, argv);
    } catch (const CLI::ParseError &e) {
        return app.exit(e);
    }

    console->info("Starting up demo client");

    const auto aco_name = "Test ACO 1";

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
    string userID;
    const string user_name = "Test User 1";

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
        userID = userTask.get();
    } catch (const exception &e) {
        console->critical("Unable to get User ID: {:s}", e.what());
        return 1;
    }

    console->info("Looking up Macaroon for '{:s}' associated with '{:s}'", aco_name, user_name);

    // Try to find the ACO token associated with the user
    string token;
    std::string token_query = fmt::format("api/users/token/{}/ACO/{}", userID, acoID);
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

    auto mac = Macaroon::importMacaroons(token);
    // Debug
    console->debug("Inspected macaroon: {:s}", mac.inspect());


    // Try to bind macaroons
    console->info("Discharging third party caveats");
    auto bound_mac = mac.discharge_all_caveats();

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
                    if (response.status_code() == status_codes::OK) {
                        return response.extract_string();
                    }
            });

    const string resp = task.get();

    console->info(resp);

    return 0;
}