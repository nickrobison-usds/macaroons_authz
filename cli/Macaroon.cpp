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

std::vector<MacaroonCaveat> Macaroon::get_third_party_caveats() {

    const auto num_caveats = macaroon_num_third_party_caveats(m);
    std::vector<MacaroonCaveat> caveats;

    for (int i = 0; i < num_caveats; i ++) {
        // ID string
        size_t id_sz;
        std::string id_str;
        const char* id_token = id_str.data();

        // Location
        size_t loc_sz;
        std::string loc_string;
        const char* loc_token = loc_string.data();

        macaroon_third_party_caveat(this->m, i,
                                    reinterpret_cast<const unsigned char **>(&loc_token),
                                    &loc_sz,
                                    reinterpret_cast<const unsigned char **>(&id_token),
                                    &id_sz);

        caveats.emplace_back(std::string(loc_token, loc_sz),
                std::string(id_token, id_sz));
    }

    return caveats;
}

const std::string Macaroon::location() {

    // Get the size;
    size_t id_sz;
    std::unique_ptr<char[]> te(new char[100]);
    const char *token = te.get();

    macaroon_location(m, reinterpret_cast<const unsigned char **>(&token), &id_sz);
    // This feels redundant, but ok, I guess.
    return std::string(token, id_sz);
}


