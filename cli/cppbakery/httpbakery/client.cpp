
#include <cpprest/http_client.h>
#include <fmt/format.h>
#include "include/httpbakery/client.hpp"
#include "../extern/cppcodec/cppcodec/base64_url_unpadded.hpp"
#include "../extern/cppcodec/cppcodec/base64_url.hpp"
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
    const std::vector<Macaroon> discharged = pplx::when_all(discharges.begin(), discharges.end()).get();

    // Bind everything
    // Create the json value
    std::vector<std::string> discharged_macs;
    discharged_macs.emplace_back(m.serialize(format));
    std::for_each(discharged.begin(), discharged.end(), [&discharged_macs, &m, format](const Macaroon &dm) {
        macaroon_returncode err;
        const macaroon *mm = macaroon_prepare_for_request(m.M(), dm.M(), &err);
        const Macaroon m2 = Macaroon(mm);
        discharged_macs.emplace_back(m2.serialize(format));
    });

    // We manually build the discharged array, to avoid double quoting everything
    // This probably wouldn't be an issue with another JSON library, but the cpprestsdk doesn't seem to have an intuitive way of handling this.
    std::ostringstream output;
    output << "[";
    // Copy the all but the last discharge into the array
    std::copy(discharged_macs.begin(), discharged_macs.end()-1, std::ostream_iterator<std::string>(output, ", "));
    // Copy the last value and the ending array block
    output << discharged_macs.back() << "]";

    return base64enc::encode(output.str());
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
            .then([&cav](http_response resp) {
                if (resp.status_code() != status_codes::OK) {
                    const auto json_error = resp.extract_json().get();
                    const std::string error_msg = json_error.at("error").as_string();
                    const std::string msg = fmt::format("Unable to discharge Macaroon from: {}. {}", cav.location, error_msg);
                    throw std::runtime_error(msg);
                }
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
