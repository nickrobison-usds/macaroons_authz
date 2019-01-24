//
// Created by Nicholas Robison on 2019-01-24.
//

#ifndef CMSAUTHZCLI_HELPERS_HPP
#define CMSAUTHZCLI_HELPERS_HPP

#include <cpprest/http_client.h>

using namespace web::http;                  // Common HTTP functionality


class helpers {
public:
    /**
     * Helper function for handling http responses
     * Either returns a {@link std::string} or a {@link web::json}, based on the template parameter
     * Currently, those are the only two types support, but should be expanded in the future
     * @tparam T - Template return type
     * @param resp - {@link http_response}
     * @param error_handler - Function for generating exception
     * @return - {@link pplx::task<T>}
     */
    template<typename T>
    static pplx::task<T>
    handle_response(const http_response &resp, const std::function<std::string(std::string)> &error_handler) {
        // Maybe we're a string, so do that
        if (resp.status_code() != status_codes::OK) {
            throw std::runtime_error(error_handler(resp.reason_phrase()));
        }
        if constexpr (std::is_same<T, std::string>::value) {
            return resp.extract_string();
        } else {
            return resp.extract_json();
        }
    }
};


#endif //CMSAUTHZCLI_HELPERS_HPP
