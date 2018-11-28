#include <iostream>
#include <CLI11.hpp>
#include <termcolor/termcolor.hpp>

int main(int argc, char **argv) {
    CLI::App app {"CLI client for CMS AuthZ Demo"};

    std::string filename = "default";

    app.add_option("-f,--file", filename, "Config file.");

    try {
        app.parse(argc, argv);
    } catch(const CLI::ParseError &e) {
        return app.exit(e);
    }
    std::cout << termcolor::yellow << "Works!" << std::endl;
    return 0;
}