//
// Created by Nick Robison on 2018-11-30.
//

#include <iostream>
#include <cpprest/http_client.h>
#include "Macaroon.hpp"
#include "extern/cppcodec/cppcodec/base64_url_unpadded.hpp"
#include "extern/cppcodec/cppcodec/base64_rfc4648.hpp"

using namespace utility;                    // Common utilities like string conversions
using namespace web;                        // Common features like URIs.
using namespace web::http;                  // Common HTTP functionality
using namespace web::http::client;          // HTTP client features
using namespace concurrency;                // Asynchronous streams
using base64 = cppcodec::base64_url_unpadded;
using base64enc = cppcodec::base64_url ;

const std::string Macaroon::base64_string() {
    const macaroon_format format = MACAROON_V2;
    const size_t sz = macaroon_serialize_size_hint(this->M(), format);

    const std::unique_ptr<char[]> output(new char[sz]);
    macaroon_returncode err;
    macaroon_serialize(this->M(), format, reinterpret_cast<unsigned char *>(output.get()), sz, &err);

    return base64enc::encode(output.get(), sz);
}

/**
 * Imports macaroon from a base64 encoded string and returns a wrapper class around the base macaroon struct
 * @param token - base64 encoded string to decode and import.
 * @return - Macaroon wrapper class
 */
const Macaroon Macaroon::importMacaroons(const std::string &token) {

    enum macaroon_returncode err;

    // If it's JSON, just import it, otherwise assume base64 and decode
    if (token[0] == '{') {
        const auto mac = macaroon_deserialize(reinterpret_cast<const unsigned char *>(token.data()), token.size(), &err);
        return Macaroon(mac);
    }

    // Decode the macaroon from base64 string
    const auto decoded = base64::decode(token);
    // Create the macaroon
    const auto mac = macaroon_deserialize(decoded.data(), decoded.size(), &err);
    return Macaroon(mac);
}

Macaroon::Macaroon(const macaroon *mac) : m(mac) {
// Not used
}

Macaroon::Macaroon() : m(nullptr) {
//    Not used
}

const Macaroon Macaroon::discharge_all_caveats() {
    // Get all the caveats
    const auto caveats = this->get_third_party_caveats();
/*
    auto test = [](const MacaroonCaveat &m) { std::cout << m.location << std::endl; };

    auto t2 = std::transform(std::remove_if(caveats,
                                            [](const MacaroonCaveat &m) {
                                                return m.isLocal();
                                            }),
                             test);
/*/

    std::vector<pplx::task<Macaroon>> discharges;
    discharges.reserve(caveats.size());

    std::for_each(caveats.begin(), caveats.end(), [&discharges](const MacaroonCaveat &cav) {
        auto test = Macaroon::dischargeCaveat(cav);
        discharges.push_back(test);
    });

    // Bind it all
    const auto discharged = pplx::when_all(discharges.begin(), discharges.end()).get();

    const auto bound = std::accumulate(discharged.begin(), discharged.end(), this->M(),
                                       [](const macaroon *acc, const Macaroon &val) {
                                           macaroon_returncode err;
                                           return macaroon_prepare_for_request(acc, val.M(), &err);
                                       });
    return Macaroon(bound);
}


// Static
pplx::task<Macaroon> Macaroon::dischargeCaveat(const MacaroonCaveat &cav) {
    // Encode the caveat caveat ID as base64
    const auto encoded = base64::encode(cav.identifier);

    // Create the URL client
    std::cout << "Location: " << cav.location << std::endl;
    http_client client(U(cav.location));
    uri_builder builder(U("/discharge"));
    builder.append_query(U("id64"), U(encoded));
//    Create the request
    http_request req(methods::POST);
    req.set_request_uri(builder.to_uri());

    // JSON encode the token
    const json::value id64 = json::value::string(encoded);
    json::value obj = json::value::object();
    obj["id64"] = id64;
    req.set_body(obj);

    // Make the call
    return client.request(req)
            .then([](http_response resp) {
                return resp.extract_json();
            })
//            Build the macaroons
            .then([](json::value json) {
                std::cout << json << std::endl;
                auto j_mac = json[U("Macaroon")];
                // Add v2
                // De-url encode
                std::string i64_string = j_mac[U("i64")].as_string();
                std::string s64_string = j_mac[U("s64")].as_string();

//                std::replace(i64_string.begin(), i64_string.end(), '-', '+');
//                std::replace(i64_string.begin(), i64_string.end(), '_', '/');
//                std::replace(s64_string.begin(), s64_string.end(), '-', '+');
//                std::replace(s64_string.begin(), s64_string.end(), '_', '/');
                const auto i64_dec = base64::decode(i64_string);
                const auto s64_dec = base64::decode(s64_string);
                auto new_mac = json::value::object();
                new_mac[U("i64")] = json::value(base64enc::encode(i64_dec.data(), i64_dec.size()));
                new_mac[U("s64")] = json::value(base64enc::encode(s64_dec.data(), s64_dec.size()));
                // Both the version and an empty array of caveats needs to be present
                new_mac[U("c")] = json::value::array();
                new_mac[U("v")] = json::value(2);
//                j_mac[U("i64")] = json::value(std::re)
//                const std::string mac_string = json["Macaroon"]["i64"].as_string();
                // Decode it and parse it
                const auto mac_string = new_mac.serialize();
                return Macaroon::importMacaroons(mac_string);
            });
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

const std::vector<const MacaroonCaveat> Macaroon::get_third_party_caveats() {

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
