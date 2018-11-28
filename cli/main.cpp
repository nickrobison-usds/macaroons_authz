#include <iostream>
#include <CLI11.hpp>
#include <termcolor/termcolor.hpp>
#include <cpprest/http_client.h>

using namespace std;
using namespace utility;                    // Common utilities like string conversions
using namespace web;                        // Common features like URIs.
using namespace web::http;                  // Common HTTP functionality
using namespace web::http::client;          // HTTP client features
using namespace concurrency::streams;       // Asynchronous streams

int main(int argc, char **argv) {
    CLI::App app{"CLI client for CMS AuthZ Demo"};

    string filename = "default";

    app.add_option("-f,--file", filename, "Config file.");

    try {
        app.parse(argc, argv);
    } catch (const CLI::ParseError &e) {
        return app.exit(e);
    }
    cout << termcolor::green << "Starting up." << endl;
    // Make sure to reset the terminal color, otherwise the remaining text output is this way.
    cout << termcolor::reset << endl;

    // Try to lookup a given ACO ID

    http_client nameClient(U("http://localhost:8080"));
    uri_builder nameBuilder(U("/api/acos/find"));
    nameBuilder.append_query(U("name"), U("Test ACO 1"));

    string acoID;


    auto nameTask = nameClient.request(methods::GET, nameBuilder.to_string())
            .then([](http_response resp) {
                if (resp.status_code() == status_codes::OK) {
                    return resp.extract_string();
                }
                throw invalid_argument(resp.extract_string().get());
            });

    try {
        acoID = nameTask.get();
    } catch (const exception &e) {
        cout << termcolor::red << e.what() << termcolor::reset << endl;
    }

    cout << "ACO ID: " << acoID << endl;

//    Fetch the user token from the environment and decode it.
    const string token = getenv("TOKEN");
    cout << token << endl;
    
    /*

    http_client client(U("http://localhost:3002"));

    uri_builder builder(U("/api/aco/show/") << aocID);

    auto task = client.request(methods::GET, builder.to_string())
            .then([](http_response response) {
                printf("Received response status code:%u\n", response.status_code());
            });

    task.wait();
     */

    return 0;
}