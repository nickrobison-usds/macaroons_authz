//
// Created by usds on 2018-12-16.
//

#ifndef CMSAUTHZCLI_CLIENT_HPP
#define CMSAUTHZCLI_CLIENT_HPP

#include <bakery/Macaroon.hpp>

class Client {

public:
    Client() = default;

    const std::string dischargeMacaroon(Macaroon m, macaroon_format format = MACAROON_V2J) const;

private:
    pplx::task<std::vector<pplx::task<Macaroon>>> dischargeCaveat(const MacaroonCaveat &cav) const;
};

#endif //CMSAUTHZCLI_CLIENT_HPP
