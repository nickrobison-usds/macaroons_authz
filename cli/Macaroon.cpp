//
// Created by Nick Robison on 2018-11-30.
//

#include <iostream>
#include "Macaroon.hpp"
#include "extern/cppcodec/cppcodec/base64_url.hpp"

/**
 * Imports macaroon from a base64 encoded string and returns a wrapper class around the base macaroon struct
 * @param token - base64 encoded string to decode and import.
 * @return - Macaroon wrapper class
 */
const Macaroon Macaroon::importMacaroons(const std::string &token) {
    using base64 = cppcodec::base64_url;
    // Decode the macaroon from base64 string
    const auto decoded = base64::decode(token);
    // Create the macaroon
    enum macaroon_returncode err;
    const auto mac = macaroon_deserialize(decoded.data(), decoded.size(), &err);
    return Macaroon(mac);
}

Macaroon::Macaroon(const macaroon *mac) : m(mac) {

}

const macaroon *Macaroon::M() {
    return m;
}

/**
 * Inspect the macaroon and print out location, id, and caveats
 */
void Macaroon::inspect() {
    // Get the size of the macaroon
    const auto sz = macaroon_inspect_size_hint(this->m);
    // Create a buffer to read into
    const std::unique_ptr<char[]> output(new char[sz]);
    enum macaroon_returncode err;
    macaroon_inspect(this->m, &output[0], sz, &err);

    // Print it
    std::cout << "Inspected macaroon: " << output.get() << std::endl;
}

const std::string Macaroon::location() {

    // Get the size;
    size_t id_sz;
    std::unique_ptr<char[]> te(new char[100]);
    const char* token = te.get();

    macaroon_location(m, reinterpret_cast<const unsigned char **>(&token), &id_sz);
    // This feels redundant, but ok, I guess.
    return std::string(token, id_sz);
}
