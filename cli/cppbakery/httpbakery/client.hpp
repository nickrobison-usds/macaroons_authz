//
// Created by usds on 2018-12-16.
//

#ifndef CMSAUTHZCLI_CLIENT_HPP
#define CMSAUTHZCLI_CLIENT_HPP

#include <bakery/Macaroon.hpp>
#include <cpprest/http_client.h>
#include <fmt/format.h>
#include "../extern/cppcodec/cppcodec/base64_url_unpadded.hpp"
#include "../extern/cppcodec/cppcodec/base64_url.hpp"
#include "../extern/cppcodec/cppcodec/base64_rfc4648.hpp"
#include "interceptor.hpp"
#include "helpers.hpp"

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

using namespace web::http;                  // Common HTTP functionality


template<class Logger>
class Client : public Logger {

public:
    Client(): logger(nullptr) {
        // Not used
    };

    Client(Logger logger): logger(std::unique_ptr<Logger>(logger)) {
        // Not used
    };

    void addInterceptor(const std::string &location, const Interceptor *interceptor) {
        Logger::debug(fmt::format("Registering interceptor for location {}", location));
        this->interceptors.push_back(interceptor);
    };
    const std::string dischargeMacaroon(Macaroon m, macaroon_format format = MACAROON_V2J) const {
        // Get all the caveats
        const auto caveats = m.get_third_party_caveats();

        std::vector<pplx::task<std::vector<Macaroon>>> discharges;
//    discharges.reserve(caveats.size());

        std::for_each(caveats.begin(), caveats.end(), [this, &discharges](const MacaroonCaveat &cav) {
//        const auto discharged =
//                .then([](std::vector<pplx::task<Macaroon>> tasks) {
//                    return pplx::when_all(tasks.begin(), tasks.end());
//                });
            discharges.emplace_back(this->dischargeCaveat(cav));
        });

//    std::transform(caveats.begin(), caveats.end(), discharges.begin(), [this](const MacaroonCaveat &cav) {
//        // Discharge the caveat (and any nested caveats) and then transform them into a vector of waiting tasks
//        const auto done = this->dischargeCaveat(cav)
//                .then([](std::vector<pplx::task<Macaroon>> tasks) {
//                    return pplx::when_all(tasks.begin(), tasks.end());
//                });
//        return done;
//    });
        std::vector<Macaroon> discharged = pplx::when_all(discharges.begin(), discharges.end()).get();

        // Bind everything
        // Create the json value
        std::vector<std::string> discharged_macs;
        discharged_macs.reserve(discharged.size() + 1);
        discharged_macs.emplace_back(m.serialize(format));

        std::for_each(discharged.begin(), discharged.end(), [&m, &discharged_macs, format](const Macaroon &dm) {
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
        std::copy(discharged_macs.begin(), discharged_macs.end() - 1, std::ostream_iterator<std::string>(output, ", "));
        // Copy the last value and the ending array block
        output << discharged_macs.back() << "]";

        const std::string outs = output.str();

        return base64enc::encode(outs);
    };

private:

    std::vector<const Interceptor*> interceptors;
    std::unique_ptr<Logger> logger;

    pplx::task<std::vector<Macaroon>> dischargeCaveat(const MacaroonCaveat &cav) const {
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

        // Intercept!
        const auto loc = cav.location;

        const auto intercepted_req = std::reduce(this->interceptors.begin(), this->interceptors.end(), req,
                                                 [&loc](http_request acc, const Interceptor* interceptor) {
                                                     return interceptor->intercept(acc, loc);
                                                 });

        // Make the call
        return client.request(intercepted_req)
                .then([&cav](http_response resp) {
                    return helpers::handle_response<json::value>(resp, [&cav](const std::string& reason) {
                        return fmt::format("Unable to discharge Macaroon from: {}. {}", cav.location, reason);
                    });
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
                    // It's possible that a macaroon has no caveats, so we can't just add it to the new mac.
                    if (j_mac.has_field(U("c"))) {
                        new_mac[U("c")] = j_mac[U("c")];
                    }
                    new_mac[U("v")] = json::value(2);
//                j_mac[U("i64")] = json::value(std::re)
//                const std::string mac_string = json["Macaroon"]["i64"].as_string();
                    // Decode it and parse it
                    const auto mac_string = new_mac.serialize();
                    return Macaroon::importMacaroons(mac_string);
                })
                        // Discharge any additional caveats
                .then([this](const Macaroon &mac) {
                    const auto caveats = mac.get_third_party_caveats();
                    if (caveats.empty()) {
                        return std::vector({mac});
                    }

                    std::vector<pplx::task<std::vector<Macaroon>>> discharged_caveats;
                    discharged_caveats.reserve(caveats.size() + 1);
                    // Add the original mac, so its discharges come second.
                    const auto original_mac = pplx::task_from_result(std::vector<Macaroon>({mac}));
                    discharged_caveats.push_back(original_mac);

                    std::for_each(caveats.begin(), caveats.end(), [this, &discharged_caveats](const MacaroonCaveat &cav) {
                        discharged_caveats.emplace_back(this->dischargeCaveat(cav));
                    });
                    return pplx::when_all(discharged_caveats.begin(), discharged_caveats.end()).get();
                });


//    std::vector<pplx::task<const Macaroon>> discharges;

        // Transform the task of
    }
};

#endif //CMSAUTHZCLI_CLIENT_HPP
