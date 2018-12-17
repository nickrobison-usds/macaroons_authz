//
// Created by Nick Robison on 2018-11-30.
//

#include <iostream>
#include <cpprest/http_client.h>
#include "include/bakery/Macaroon.hpp"
#include "../extern/cppcodec/cppcodec/base64_url_unpadded.hpp"
#include "../extern/cppcodec/cppcodec/base64_rfc4648.hpp"

using namespace utility;                    // Common utilities like string conversions
using namespace web;                        // Common features like URIs.
using namespace web::http;                  // Common HTTP functionality
using namespace web::http::client;          // HTTP client features
using namespace concurrency;                // Asynchronous streams
using base64 = cppcodec::base64_url_unpadded;
using base64enc = cppcodec::base64_url;
using base64rfc = cppcodec::base64_rfc4648;

const std::string Macaroon::base64_string(const macaroon_format format) const {
    const size_t sz = macaroon_serialize_size_hint(this->M(), format);

    const std::unique_ptr<char[]> output(new char[sz]);
    macaroon_returncode err;
    const size_t buffer_size = macaroon_serialize(this->M(), format, reinterpret_cast<unsigned char *>(output.get()),
                                                  sz, &err);

    // Binary formats are already base64 encoded
    switch (format) {
        case MACAROON_V2J: {
            return base64enc::encode(output.get(), buffer_size);
        }
        default: {
            return std::string(output.get(), buffer_size);
        }
    }
}

/**
 * Imports macaroon from a base64 encoded string and returns a wrapper class around the base macaroon struct
 * @param token - base64 encoded string to decode and import.
 * @return - Macaroon wrapper class
 */
const Macaroon Macaroon::importMacaroons(const std::string &token) {

    enum macaroon_returncode err;
    macaroon *mac;

    // If it's JSON, we can directly import it
    if (token[0] == '{') {
        mac = macaroon_deserialize(reinterpret_cast<const unsigned char *>(token.data()), token.size(),
                                   &err);
    } else {

        // Determine URL safe encoding
        const auto result = std::find_if(token.begin(), token.end(), [](const char t) {
            return (t == '-' || t == '_');
        });

        std::vector<uint8_t> decoded;
        if (result == token.end()) {
//        not URL safe encoding
            decoded = base64enc::decode(token);
        } else {
            decoded = base64enc::decode(token);
        }


        switch (decoded[0]) {
            // If it's un-encoded JSON, or V2 binary, import the decoded value.
            case '\x02': {
                mac = macaroon_deserialize(reinterpret_cast<const unsigned char *>(decoded.data()), decoded.size(),
                                           &err);
                break;
            }
            default: {
                // If it's V1 binary format, re-encode it as non-url safe
                const auto encoded = base64rfc::encode(decoded);
                // Create the macaroon
                mac = macaroon_deserialize(reinterpret_cast<const unsigned char *>(encoded.data()), encoded.size(),
                                           &err);
            }
        }
    }

    return Macaroon(mac);
}

Macaroon::Macaroon(const macaroon *mac) : m(mac) {
// Not used
}

Macaroon::Macaroon() : m(nullptr) {
//    Not used
}


// Static


/**
 * Inspect the macaroon and print out location, id, and caveats
 */
std::string Macaroon::inspect() {
    // Get the size of the macaroon
    const auto sz = macaroon_inspect_size_hint(this->m);
    // Create a buffer to read into
    const std::unique_ptr<char[]> output(new char[sz]);
    enum macaroon_returncode err;
    macaroon_inspect(this->m, &output[0], sz, &err);

    // Print it
    return std::string(output.get());
}

const std::vector<const MacaroonCaveat> Macaroon::get_third_party_caveats() const {

    const auto num_caveats = macaroon_num_third_party_caveats(m);
    std::vector<const MacaroonCaveat> caveats;

    // Interate through the caveats and build a vector of them.
    for (unsigned int i = 0; i < num_caveats; i++) {
        // ID string
        size_t id_sz;
        std::string id_str;
        const char *id_token = id_str.data();

        // Location
        size_t loc_sz;
        std::string loc_string;
        const char *loc_token = loc_string.data();

        macaroon_third_party_caveat(this->m, i,
                                    reinterpret_cast<const unsigned char **>(&loc_token),
                                    &loc_sz,
                                    reinterpret_cast<const unsigned char **>(&id_token),
                                    &id_sz);

        caveats.emplace_back(std::string(loc_token, loc_sz),
                             std::string(id_token, id_sz));
    }

    return std::as_const(caveats);
}

json::value Macaroon::as_json() const {
    macaroon_returncode err;
    macaroon_format format = MACAROON_V2J;
    const size_t sz = macaroon_serialize_size_hint(this->M(), format);
    const std::unique_ptr<char[]> output(new char[sz]);
    std::string ts;
    const size_t total_size = macaroon_serialize(this->M(), format, reinterpret_cast<unsigned char *>(output.get()), sz,
                                                 &err);
//    return std::string(output.get());
    return json::value::parse(std::string(output.get(), total_size));
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

const macaroon *Macaroon::M() const {
    return m;
}
