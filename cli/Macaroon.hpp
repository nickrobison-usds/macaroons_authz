//
// Created by usds on 2018-11-30.
//

#ifndef CMSAUTHZCLI_MACAROON_HPP
#define CMSAUTHZCLI_MACAROON_HPP


#include <string>
#include <utility>
#include <macaroons.h>

struct MacaroonCaveat {
    std::string location;
    std::string identifier;

    MacaroonCaveat(std::string loc, std::string id): location(std::move(loc)), identifier(std::move(id)) {
//        Not used
    }

    bool isLocal() const {
        return this->location.empty();
    }
};

class Macaroon {

private:
    const struct macaroon* m;
    pplx::task<Macaroon> static dischargeCaveat(const MacaroonCaveat &cav);


public:
    Macaroon();
    explicit Macaroon(const macaroon *m);
    const std::string discharge_all_caveats();
    void inspect();
    const std::vector<const MacaroonCaveat> get_third_party_caveats();
    const std::string location();
    const macaroon * M() const;
    web::json::value as_json() const;
    const std::string base64_string() const;

    const static Macaroon importMacaroons(const std::string &string);
};


#endif //CMSAUTHZCLI_MACAROON_HPP
