//
// Created by Nicholas Robison on 2019-01-18.
//

#include <cpprest/http_client.h>
#include "UserInterceptor.hpp"

http_request UserInterceptor::intercept(http_request &request, const std::__1::string &location) const {
    // When discharging as a Vendor, we need to specify our userID, otherwise the application gets confused
    if (location.rfind("http://localhost:8080/api/acos") == 0) {
        uri_builder builder(request.absolute_uri());
        builder.append_query("user_id", this->userID);
        request.set_request_uri(builder.to_uri());
    }
    return request;
}