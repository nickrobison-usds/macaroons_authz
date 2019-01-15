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

    MacaroonCaveat(std::string loc, std::string id) : location(std::move(loc)), identifier(std::move(id)) {
//        Not used
    }

    bool isLocal() const {
        return this->location.empty();
    }
};

class Macaroon {

private:
    std::shared_ptr<const macaroon> m;


public:
    Macaroon();

    explicit Macaroon(const macaroon* m);

    std::string inspect();

    const std::vector<const MacaroonCaveat> get_third_party_caveats() const;

    const std::string location();

    const macaroon *M() const;

    web::json::value as_json() const;

    const std::string serialize(macaroon_format format = MACAROON_V2J) const;

    const static Macaroon importMacaroons(const std::string &string);
};


#endif //CMSAUTHZCLI_MACAROON_HPP
