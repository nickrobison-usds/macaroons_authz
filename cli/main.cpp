#include <iostream>
#include <CLI/Error.hpp>
#include <CLI/App.hpp>

int main(int argc, char **argv) {
//    std::cout << "Works." << std::endl;
    CLI::App app {"CLI client for CMS AuthZ Demo"};

    std::string filename = "default";

    app.add_option("-f,--file", filename, "Config file.");

    try {
        app.parse(argc, argv);
    } catch(const CLI::ParseError &e) {
        return app.exit(e);
    }
}