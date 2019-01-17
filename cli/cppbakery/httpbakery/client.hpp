//
// Created by usds on 2018-12-16.
//

#ifndef CMSAUTHZCLI_CLIENT_HPP
#define CMSAUTHZCLI_CLIENT_HPP

#include "bakery/Macaroon.hpp"
#include <cpprest/http_client.h>
#include "interceptor.hpp"

using namespace web::http;                  // Common HTTP functionality

class Client {

private:
    pplx::task<std::vector<Macaroon>> dischargeCaveat(const MacaroonCaveat &cav) const;

    std::map<std::string, std::unique_ptr<const Interceptor>> interceptors;

public:
    Client() = default;

    void addInterceptor(const std::string &location, const Interceptor &interceptor);
    const std::string dischargeMacaroon(Macaroon m, macaroon_format format = MACAROON_V2J) const;
    const http_request interceptRequest(const http_request request) const;
};

#endif //CMSAUTHZCLI_CLIENT_HPP
