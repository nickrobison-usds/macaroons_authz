
#include <cpprest/http_client.h>
#include "include/httpbakery/client.hpp"
#include "../extern/cppcodec/cppcodec/base64_url_unpadded.hpp"
#include "../extern/cppcodec/cppcodec/base64_rfc4648.hpp"

//
// Created by usds on 2018-12-16.
//
using namespace utility;                    // Common utilities like string conversions
using namespace web;                        // Common features like URIs.
using namespace web::http;                  // Common HTTP functionality
using namespace web::http::client;          // HTTP client features
using namespace concurrency;                // Asynchronous streams
using base64 = cppcodec::base64_url_unpadded;
using base64enc = cppcodec::base64_url;
using base64rfc = cppcodec::base64_rfc4648;

Client::Client() {
//        not used, yet
}

const std::string Client::dischargeMacaroon(const Macaroon m, const macaroon_format format) const {
    // Get all the caveats
    const auto caveats = m.get_third_party_caveats();

    std::vector<pplx::task<Macaroon>> discharges;
    discharges.reserve(caveats.size());

    std::for_each(caveats.begin(), caveats.end(), [&discharges, this](const MacaroonCaveat &cav) {
        auto test = this->dischargeCaveat(cav);
        discharges.push_back(test);
    });
    const auto discharged = pplx::when_all(discharges.begin(), discharges.end()).get();

    // Bind everything
    // Create the json value
    std::vector<json::value> discharged_macs;
    discharged_macs.emplace_back(m.base64_string(format));
    std::for_each(discharged.begin(), discharged.end(), [&discharged_macs, this, format](const Macaroon &mac) {
        macaroon_returncode err;
        const macaroon *mm = macaroon_prepare_for_request(mac.M(), mac.M(), &err);
        const Macaroon m2 = Macaroon(mm);
        discharged_macs.emplace_back(json::value::string(m2.base64_string(format)));
    });

    json::value val_array = json::value::array(discharged_macs);
    const std::string serialized = val_array.serialize();
    return base64enc::encode(serialized);
}

pplx::task<Macaroon> Client::dischargeCaveat(const MacaroonCaveat &cav) const {
    // Encode the caveat caveat ID as base64
    const auto encoded = base64::encode(cav.identifier);

    // Create the URL client
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

